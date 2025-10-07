const { ActivityType } = require('discord.js');

module.exports = {
  name: 'status',
  description: 'Change bot status',
  async execute(message, args, client, PREFIX) {
    if (!args.length) {
      return message.reply({
        embeds: [{
          color: 0x0099ff,
          title: '🎭 Status Command',
          description: `Change bot's activity status`,
          fields: [
            {
              name: 'Usage',
              value: `\`${PREFIX}status <type> <text>\``,
              inline: false
            },
            {
              name: 'Available Types',
              value: '`playing`, `listening`, `watching`, `competing`, `idle`',
              inline: false
            },
            {
              name: 'Examples',
              value: `\`${PREFIX}status listening new music\`\n\`${PREFIX}status playing Minecraft\`\n\`${PREFIX}status watching YouTube\`\n\`${PREFIX}status idle\``,
              inline: false
            }
          ]
        }]
      });
    }

    const type = args[0].toLowerCase();
    const text = args.slice(1).join(' ');

    try {
      if (type === 'idle') {
        client.user.setPresence({
          status: 'idle',
          activities: []
        });
        return message.reply('✅ Status changed to **Idle**');
      }

      if (!text) {
        return message.reply(`❌ Please provide text for the status!\nExample: \`${PREFIX}status listening new music\``);
      }

      let activityType;
      switch (type) {
        case 'playing':
          activityType = ActivityType.Playing;
          break;
        case 'listening':
          activityType = ActivityType.Listening;
          break;
        case 'watching':
          activityType = ActivityType.Watching;
          break;
        case 'competing':
          activityType = ActivityType.Competing;
          break;
        default:
          return message.reply(`❌ Invalid status type!\nAvailable: \`playing\`, \`listening\`, \`watching\`, \`competing\`, \`idle\``);
      }

      client.user.setPresence({
        status: 'online',
        activities: [{
          name: text,
          type: activityType
        }]
      });

      message.reply(`✅ Status changed to **${type.charAt(0).toUpperCase() + type.slice(1)}** ${text}`);
    } catch (error) {
      console.error('Error setting status:', error);
      message.reply('❌ Failed to change status');
    }
  },
};
