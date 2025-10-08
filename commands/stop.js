module.exports = {
  name: 'stop',
  description: 'Stop the music and clear the queue',
  execute(message, args, client) {
    const queue = client.distube.getQueue(message);
    if (!queue) return message.reply('❌ Nothing is playing!');

    client.distube.stop(message);
    message.reply('⏹ Stopped the music and cleared the queue!');
  },
};
