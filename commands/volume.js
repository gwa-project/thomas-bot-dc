const { getQueue } = require('../utils/musicManager');

module.exports = {
  name: 'volume',
  description: 'Set music volume (0-100)',
  async execute(message, args, client, PREFIX) {
    const queue = getQueue(message.guild.id);

    if (!queue || !queue.isPlaying) {
      return message.reply('❌ Nothing is playing right now!');
    }

    if (!message.member.voice.channel) {
      return message.reply('❌ You need to be in a voice channel!');
    }

    if (!args.length) {
      const currentVolume = Math.round(queue.volume * 100);
      return message.reply(`🔊 Current volume: **${currentVolume}%**\nUsage: \`${PREFIX}volume <0-100>\``);
    }

    const volume = parseInt(args[0]);

    if (isNaN(volume) || volume < 0 || volume > 100) {
      return message.reply('❌ Volume must be a number between 0 and 100!');
    }

    const newVolume = queue.setVolume(volume);
    message.reply(`🔊 Volume set to **${Math.round(newVolume)}%**`);
  },
};
