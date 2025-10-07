const { getQueue } = require('../utils/musicManager');

module.exports = {
  name: 'nowplaying',
  description: 'Show currently playing song',
  async execute(message, args, client, PREFIX) {
    const queue = getQueue(message.guild.id);

    if (!queue || !queue.isPlaying || !queue.currentSong) {
      return message.reply('❌ Nothing is playing right now!');
    }

    const song = queue.currentSong;

    message.reply({
      embeds: [{
        color: 0x00ff00,
        title: '🎵 Now Playing',
        description: `[${song.title}](${song.url})`,
        fields: [
          {
            name: 'Duration',
            value: song.duration,
            inline: true
          },
          {
            name: 'Requested by',
            value: song.requester,
            inline: true
          },
          {
            name: 'Queue',
            value: `${queue.songs.length - 1} song(s) remaining`,
            inline: true
          }
        ],
        thumbnail: {
          url: song.thumbnail
        }
      }]
    });
  },
};
