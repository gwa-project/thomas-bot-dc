module.exports = {
  name: 'queue',
  description: 'Show the music queue',
  execute(message, args, client) {
    const queue = client.player.nodes.get(message.guild.id);
    if (!queue || !queue.tracks.data.length) {
      return message.reply('❌ Queue is empty!');
    }

    const tracks = queue.tracks.data.slice(0, 10);
    const queueString = tracks.map((track, i) =>
      `**${i + 1}.** ${track.title} - \`${track.duration}\``
    ).join('\n');

    message.reply({
      embeds: [{
        color: 0x0099ff,
        title: '📜 Music Queue',
        description: queueString,
        fields: [
          { name: 'Now Playing', value: queue.currentTrack ? queue.currentTrack.title : 'None' },
          { name: 'Total Songs', value: `${queue.tracks.data.length}`, inline: true }
        ]
      }]
    });
  },
};
