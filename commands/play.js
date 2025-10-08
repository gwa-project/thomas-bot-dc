module.exports = {
  name: 'play',
  description: 'Play music from YouTube',
  async execute(message, args, client, PREFIX) {
    // Check if user is in voice channel
    const voiceChannel = message.member.voice.channel;
    if (!voiceChannel) {
      return message.reply('❌ You need to be in a voice channel to play music!');
    }

    // Check if query provided
    if (!args.length) {
      return message.reply(`❌ Please provide a song name or URL!\nUsage: \`${PREFIX}play <song name or URL>\``);
    }

    const query = args.join(' ');

    try {
      await client.distube.play(voiceChannel, query, {
        textChannel: message.channel,
        member: message.member
      });
    } catch (error) {
      console.error('Play command error:', error);
      message.reply('❌ An error occurred while trying to play music!');
    }
  },
};
