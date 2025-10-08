module.exports = {
  name: 'skip',
  description: 'Skip the current song',
  execute(message, args, client) {
    const queue = client.player.nodes.get(message.guild.id);
    if (!queue || !queue.isPlaying()) {
      return message.reply('❌ Nothing is playing!');
    }

    queue.node.skip();
    message.reply('⏭ Skipped the song!');
  },
};
