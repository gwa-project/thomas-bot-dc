module.exports = {
  name: 'stop',
  description: 'Stop the music and clear queue',
  execute(message, args, client) {
    const queue = client.player.nodes.get(message.guild.id);
    if (!queue) {
      return message.reply('❌ Nothing is playing!');
    }

    queue.delete();
    message.reply('⏹ Stopped the music!');
  },
};
