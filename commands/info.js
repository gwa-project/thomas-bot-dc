module.exports = {
  name: 'info',
  description: 'Show bot information',
  async execute(message, args, client, PREFIX) {
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
  },
};
