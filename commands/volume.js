module.exports = {
  name: 'volume',
  description: 'Set the volume',
  execute(message, args, client, PREFIX) {
    const queue = client.distube.getQueue(message);
    if (!queue) return message.reply('❌ Nothing is playing!');

    const volume = parseInt(args[0]);

    if (isNaN(volume) || volume < 0 || volume > 100) {
      return message.reply(`❌ Please provide a volume between 0 and 100!\nUsage: \`${PREFIX}volume <0-100>\``);
    }

    client.distube.setVolume(message, volume);
    message.reply(`🔊 Volume set to **${volume}%**`);
  },
};
