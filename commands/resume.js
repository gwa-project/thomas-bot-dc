module.exports = {
  name: 'resume',
  description: 'Resume the paused song',
  execute(message, args, client) {
    const queue = client.distube.getQueue(message);
    if (!queue) return message.reply('❌ Nothing is playing!');
    if (!queue.paused) return message.reply('▶️ The song is not paused!');

    client.distube.resume(message);
    message.reply('▶️ Resumed the song!');
  },
};
