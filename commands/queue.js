const { getQueue } = require('../utils/musicManager');

module.exports = {
  name: 'queue',
  description: 'Show the music queue',
  async execute(message, args, client, PREFIX) {
    const queue = getQueue(message.guild.id);

    if (!queue || queue.songs.length === 0) {
      return message.reply('❌ Queue is empty!');
    }

    // Show only first 10 songs
    const songs = queue.songs.slice(0, 10);
    const queueList = songs.map((song, index) => {
      if (index === 0) {
        return `**Now Playing:**\n🎵 [${song.title}](${song.url}) - ${song.duration}`;
      }
      return `**${index}.** [${song.title}](${song.url}) - ${song.duration}`;
    }).join('\n\n');

    const remaining = queue.songs.length > 10 ? `\n\n...and ${queue.songs.length - 10} more songs` : '';

    message.reply({
      embeds: [{
        color: 0x0099ff,
        title: '📜 Music Queue',
        description: queueList + remaining,
        footer: {
          text: `Total: ${queue.songs.length} song(s) in queue`
        }
      }]
    });
  },
};
