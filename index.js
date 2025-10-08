require('dotenv').config();
const { Client, GatewayIntentBits, Collection, ActivityType } = require('discord.js');
const { DisTube } = require('distube');
const express = require('express');
const fs = require('fs');
const path = require('path');

// Bot Configuration
const PREFIX = '!T';
const TOKEN = process.env.DISCORD_TOKEN;
const CLIENT_ID = process.env.CLIENT_ID;
const PORT = process.env.PORT || 8080;

// Initialize Discord Client
const client = new Client({
  intents: [
    GatewayIntentBits.Guilds,
    GatewayIntentBits.GuildMessages,
    GatewayIntentBits.MessageContent,
    GatewayIntentBits.GuildMembers,
    GatewayIntentBits.GuildVoiceStates,
  ],
});

// Initialize DisTube
client.distube = new DisTube(client, {
  emitNewSongOnly: true,
  leaveOnEmpty: true,
  leaveOnFinish: false,
  leaveOnStop: true,
  savePreviousSongs: true,
  searchSongs: 5,
  nsfw: false,
  emptyCooldown: 60
});

// Load Commands
client.commands = new Collection();
const commandsPath = path.join(__dirname, 'commands');
const commandFiles = fs.readdirSync(commandsPath).filter(file => file.endsWith('.js'));

for (const file of commandFiles) {
  const filePath = path.join(commandsPath, file);
  const command = require(filePath);
  client.commands.set(command.name, command);
  console.log(`📦 Loaded command: ${command.name}`);
}

// Initialize Express server for Cloud Run health checks
const app = express();

// Health check endpoint
app.get('/', (req, res) => {
  res.json({
    status: 'online',
    bot: client.user ? client.user.tag : 'Starting...',
    uptime: process.uptime(),
    guilds: client.guilds ? client.guilds.cache.size : 0,
    prefix: PREFIX,
    commands: client.commands.size
  });
});

app.get('/health', (req, res) => {
  res.json({ status: 'healthy' });
});

// Start Express server
app.listen(PORT, () => {
  console.log(`🌐 HTTP server listening on port ${PORT}`);
});

// Bot Ready Event
client.once('clientReady', () => {
  console.log('========================================');
  console.log(`✅ Bot is online as ${client.user.tag}`);
  console.log(`📝 Prefix: ${PREFIX}`);
  console.log(`🆔 Client ID: ${CLIENT_ID}`);
  console.log(`🌐 Serving ${client.guilds.cache.size} servers`);
  console.log(`🔧 Loaded ${client.commands.size} commands`);
  console.log('========================================');

  // Set bot status: Idle + Listening to new music
  client.user.setPresence({
    status: 'idle',
    activities: [{
      name: 'new music',
      type: ActivityType.Listening
    }]
  });

  console.log('🎭 Status set to: Idle | Listening to new music');
});

// DisTube Events
client.distube.on('playSong', (queue, song) => {
  queue.textChannel.send({
    embeds: [{
      color: 0x00ff00,
      title: '🎵 Now Playing',
      description: `[${song.name}](${song.url})`,
      fields: [
        { name: 'Duration', value: song.formattedDuration, inline: true },
        { name: 'Requested by', value: song.user.tag, inline: true }
      ],
      thumbnail: { url: song.thumbnail }
    }]
  });
});

client.distube.on('addSong', (queue, song) => {
  queue.textChannel.send({
    embeds: [{
      color: 0x0099ff,
      title: '➕ Added to Queue',
      description: `[${song.name}](${song.url})`,
      fields: [
        { name: 'Duration', value: song.formattedDuration, inline: true },
        { name: 'Position', value: `${queue.songs.length}`, inline: true }
      ],
      thumbnail: { url: song.thumbnail }
    }]
  });
});

client.distube.on('error', (channel, error) => {
  console.error('DisTube error:', error);
  if (channel) channel.send(`❌ An error occurred: ${error.message.slice(0, 100)}`);
});

// Message Handler
client.on('messageCreate', async (message) => {
  // Ignore bot messages
  if (message.author.bot) return;

  // Check if message starts with prefix
  if (!message.content.startsWith(PREFIX)) return;

  // Parse command and arguments
  const args = message.content.slice(PREFIX.length).trim().split(/ +/);
  const commandName = args.shift().toLowerCase();

  // Get command from collection
  const command = client.commands.get(commandName);

  if (!command) {
    return message.reply(`❌ Unknown command. Use \`${PREFIX}help\` to see available commands.`);
  }

  // Execute command
  try {
    await command.execute(message, args, client, PREFIX, Array.from(client.commands.values()));
  } catch (error) {
    console.error(`Error executing command ${commandName}:`, error);
    message.reply('❌ An error occurred while executing the command.');
  }
});

// Error Handling
client.on('error', (error) => {
  console.error('Discord client error:', error);
});

process.on('unhandledRejection', (error) => {
  console.error('Unhandled promise rejection:', error);
});

// Login to Discord
client.login(TOKEN).catch((error) => {
  console.error('❌ Failed to login:', error);
  process.exit(1);
});
