module.exports = {
  name: 'play',
  description: 'Play music from YouTube',
  async execute(message, args, client, PREFIX) {
    if (!message.member.voice.channel) {
      return message.reply('❌ You need to be in a voice channel!');
    }

    if (!args.length) {
      return message.reply(`❌ Please provide a song name or URL!\nUsage: \`${PREFIX}play <song name or URL>\``);
    }

    const query = args.join(' ');

    try {
      const searchResult = await client.player.search(query, {
        requestedBy: message.author
      });

      if (!searchResult || !searchResult.tracks.length) {
        return message.reply('❌ No results found!');
      }

      const queue = await client.player.nodes.create(message.guild, {
        metadata: {
          channel: message.channel,
          client: message.guild.members.me,
          requestedBy: message.author
        },
        selfDeaf: true,
        volume: 50,
        leaveOnEmpty: true,
        leaveOnEnd: false,
        leaveOnStop: true
      });

      try {
        if (!queue.connection) await queue.connect(message.member.voice.channel);
      } catch {
        queue.delete();
        return message.reply('❌ Could not join voice channel!');
      }

      await message.reply(`🔍 Searching for **${query}**...`);

      searchResult.playlist ? queue.addTrack(searchResult.tracks) : queue.addTrack(searchResult.tracks[0]);

      if (!queue.isPlaying()) await queue.node.play();

    } catch (error) {
      console.error('Play error:', error);
      message.reply('❌ An error occurred while trying to play music!');
    }
  },
};
