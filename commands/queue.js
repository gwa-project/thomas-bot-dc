module.exports = {
  name: 'queue',
  description: 'Show the music queue',
  execute(message, args, client) {
    const queue = client.distube.getQueue(message);
    if (!queue) return message.reply('❌ Nothing is playing!');

    const queueString = queue.songs.map((song, id) =>
      `**${id + 1}.** ${song.name} - \`${song.formattedDuration}\``
    ).slice(0, 10).join('\n');

    message.reply({
      embeds: [{
        color: 0x0099ff,
        title: '📜 Music Queue',
        description: queueString || 'No songs in queue',
        fields: [
          {
            name: 'Now Playing',
            value: `${queue.songs[0].name} - \`${queue.songs[0].formattedDuration}\``
          },
          {
            name: 'Total Songs',
            value: `${queue.songs.length}`,
            inline: true
          },
          {
            name: 'Duration',
            value: `${queue.formattedDuration}`,
            inline: true
          }
        ]
      }]
    });
  },
};
