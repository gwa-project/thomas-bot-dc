module.exports = {
  name: 'pause',
  description: 'Pause the current song',
  execute(message, args, client) {
    const queue = client.distube.getQueue(message);
    if (!queue) return message.reply('❌ Nothing is playing!');
    if (queue.paused) return message.reply('⏸ The song is already paused!');

    client.distube.pause(message);
    message.reply('⏸ Paused the song!');
  },
};
