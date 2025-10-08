module.exports = {
  name: 'nowplaying',
  description: 'Show the currently playing song',
  execute(message, args, client) {
    const queue = client.distube.getQueue(message);
    if (!queue) return message.reply('❌ Nothing is playing!');

    const song = queue.songs[0];

    message.reply({
      embeds: [{
        color: 0x00ff00,
        title: '🎵 Now Playing',
        description: `[${song.name}](${song.url})`,
        fields: [
          {
            name: 'Duration',
            value: song.formattedDuration,
            inline: true
          },
          {
            name: 'Requested by',
            value: song.user.tag,
            inline: true
          },
          {
            name: 'Volume',
            value: `${queue.volume}%`,
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
