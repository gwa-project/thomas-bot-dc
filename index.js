require('dotenv').config();
const { Client, GatewayIntentBits, Collection } = require('discord.js');
const express = require('express');

// Bot Configuration
const PREFIX = '!T';
const TOKEN = process.env.DISCORD_TOKEN;
const CLIENT_ID = process.env.CLIENT_ID;
const PORT = process.env.PORT || 8080;

// Initialize Express server for Cloud Run health checks
const app = express();

// Health check endpoint
app.get('/', (req, res) => {
  res.json({
    status: 'online',
    bot: client.user ? client.user.tag : 'Starting...',
    uptime: process.uptime(),
    guilds: client.guilds ? client.guilds.cache.size : 0,
    prefix: PREFIX
  });
});

app.get('/health', (req, res) => {
  res.json({ status: 'healthy' });
});

// Start Express server
app.listen(PORT, () => {
  console.log(`🌐 HTTP server listening on port ${PORT}`);
});

// Initialize Discord Client
const client = new Client({
  intents: [
    GatewayIntentBits.Guilds,
    GatewayIntentBits.GuildMessages,
    GatewayIntentBits.MessageContent,
    GatewayIntentBits.GuildMembers,
  ],
});

// Bot Ready Event
client.once('ready', () => {
  console.log('========================================');
  console.log(`✅ Bot is online as ${client.user.tag}`);
  console.log(`📝 Prefix: ${PREFIX}`);
  console.log(`🆔 Client ID: ${CLIENT_ID}`);
  console.log(`🌐 Serving ${client.guilds.cache.size} servers`);
  console.log('========================================');

  // Set bot status
  client.user.setActivity(`${PREFIX}help`, { type: 'LISTENING' });
});

// Message Handler
client.on('messageCreate', async (message) => {
  // Ignore bot messages
  if (message.author.bot) return;

  // Check if message starts with prefix
  if (!message.content.startsWith(PREFIX)) return;

  // Parse command and arguments
  const args = message.content.slice(PREFIX.length).trim().split(/ +/);
  const command = args.shift().toLowerCase();

  try {
    // Command: ping
    if (command === 'ping') {
      const sent = await message.reply('🏓 Pinging...');
      const latency = sent.createdTimestamp - message.createdTimestamp;
      sent.edit(`🏓 Pong!\n⏱️ Latency: ${latency}ms\n💓 API Latency: ${Math.round(client.ws.ping)}ms`);
    }

    // Command: help
    else if (command === 'help') {
      message.reply({
        embeds: [{
          color: 0x0099ff,
          title: '📚 Thomas Bot - Command List',
          description: `Prefix: \`${PREFIX}\``,
          fields: [
            {
              name: `${PREFIX}ping`,
              value: 'Check bot latency',
              inline: true
            },
            {
              name: `${PREFIX}help`,
              value: 'Show this help message',
              inline: true
            },
            {
              name: `${PREFIX}info`,
              value: 'Show bot information',
              inline: true
            },
            {
              name: `${PREFIX}server`,
              value: 'Show server information',
              inline: true
            },
          ],
          timestamp: new Date(),
          footer: {
            text: 'Thomas Bot v1.0'
          }
        }]
      });
    }

    // Command: info
    else if (command === 'info') {
      message.reply({
        embeds: [{
          color: 0x00ff00,
          title: '🤖 Bot Information',
          fields: [
            {
              name: 'Bot Name',
              value: client.user.tag,
              inline: true
            },
            {
              name: 'Bot ID',
              value: client.user.id,
              inline: true
            },
            {
              name: 'Prefix',
              value: `\`${PREFIX}\``,
              inline: true
            },
            {
              name: 'Servers',
              value: `${client.guilds.cache.size}`,
              inline: true
            },
            {
              name: 'Node.js Version',
              value: process.version,
              inline: true
            },
            {
              name: 'Discord.js Version',
              value: require('discord.js').version,
              inline: true
            }
          ],
          timestamp: new Date()
        }]
      });
    }

    // Command: server
    else if (command === 'server') {
      if (!message.guild) {
        return message.reply('❌ This command can only be used in a server!');
      }

      message.reply({
        embeds: [{
          color: 0xff9900,
          title: '🏰 Server Information',
          thumbnail: {
            url: message.guild.iconURL()
          },
          fields: [
            {
              name: 'Server Name',
              value: message.guild.name,
              inline: true
            },
            {
              name: 'Server ID',
              value: message.guild.id,
              inline: true
            },
            {
              name: 'Owner',
              value: `<@${message.guild.ownerId}>`,
              inline: true
            },
            {
              name: 'Members',
              value: `${message.guild.memberCount}`,
              inline: true
            },
            {
              name: 'Created At',
              value: message.guild.createdAt.toLocaleDateString(),
              inline: true
            }
          ],
          timestamp: new Date()
        }]
      });
    }

    // Unknown command
    else {
      message.reply(`❌ Unknown command. Use \`${PREFIX}help\` to see available commands.`);
    }

  } catch (error) {
    console.error('Command execution error:', error);
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
