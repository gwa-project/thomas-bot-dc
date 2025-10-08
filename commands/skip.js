module.exports = {
  name: 'skip',
  description: 'Skip the current song',
  execute(message, args, client) {
    const queue = client.distube.getQueue(message);
    if (!queue) return message.reply('❌ Nothing is playing!');

    try {
      client.distube.skip(message);
      message.reply('⏭ Skipped the song!');
    } catch (error) {
      message.reply('❌ There is no next song in queue!');
    }
  },
};
