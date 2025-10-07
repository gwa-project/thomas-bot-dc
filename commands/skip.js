const { getQueue } = require('../utils/musicManager');

module.exports = {
  name: 'skip',
  description: 'Skip the current song',
  async execute(message, args, client, PREFIX) {
    const queue = getQueue(message.guild.id);

    if (!queue || !queue.isPlaying) {
      return message.reply('❌ Nothing is playing right now!');
    }

    if (!message.member.voice.channel) {
      return message.reply('❌ You need to be in a voice channel!');
    }

    if (queue.skip()) {
      message.reply('⏭️ Skipped the current song');
    } else {
      message.reply('❌ Failed to skip');
    }
  },
};
