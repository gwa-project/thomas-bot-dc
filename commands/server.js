module.exports = {
  name: 'server',
  description: 'Show server information',
  async execute(message, args, client) {
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
  },
};
