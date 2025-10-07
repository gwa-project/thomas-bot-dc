module.exports = {
  name: 'ping',
  description: 'Check bot latency',
  async execute(message, args, client) {
    const sent = await message.reply('🏓 Pinging...');
    const latency = sent.createdTimestamp - message.createdTimestamp;
    sent.edit(`🏓 Pong!\n⏱️ Latency: ${latency}ms\n💓 API Latency: ${Math.round(client.ws.ping)}ms`);
  },
};
