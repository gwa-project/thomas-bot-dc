const ytdl = require('ytdl-core');
const YouTube = require('youtube-sr').default;
const { getQueue } = require('../utils/musicManager');

module.exports = {
  name: 'play',
  description: 'Play music from YouTube',
  async execute(message, args, client, PREFIX) {
    // Check if user is in voice channel
    if (!message.member.voice.channel) {
      return message.reply('❌ You need to be in a voice channel to play music!');
    }

    // Check if query provided
    if (!args.length) {
      return message.reply(`❌ Please provide a song name or URL!\nUsage: \`${PREFIX}play <song name or URL>\``);
    }

    const query = args.join(' ');
    const voiceChannel = message.member.voice.channel;

    try {
      const searchMsg = await message.channel.send('🔍 Searching...');

      let video;

      // Check if it's a YouTube URL
      if (ytdl.validateURL(query)) {
        const info = await ytdl.getBasicInfo(query);
        video = {
          title: info.videoDetails.title,
          url: info.videoDetails.video_url,
          duration: parseInt(info.videoDetails.lengthSeconds),
          thumbnail: info.videoDetails.thumbnails[0]?.url || ''
        };
      } else {
        // Search YouTube using youtube-sr
        const searchResults = await YouTube.searchOne(query);

        if (!searchResults) {
          await searchMsg.delete();
          return message.reply('❌ No results found!');
        }

        video = {
          title: searchResults.title,
          url: searchResults.url,
          duration: searchResults.duration / 1000, // Convert ms to seconds
          thumbnail: searchResults.thumbnail?.url || ''
        };
      }

      // Create song object
      const song = {
        title: video.title,
        url: video.url,
        duration: formatDuration(video.duration),
        thumbnail: video.thumbnail,
        requester: message.author.tag
      };

      await searchMsg.delete();

      // Get queue
      const queue = getQueue(message.guild.id);
      queue.addSong(song);

      // If not playing, start playing
      if (!queue.isPlaying) {
        queue.play(voiceChannel, message.channel);
      } else {
        message.reply({
          embeds: [{
            color: 0x0099ff,
            title: '➕ Added to Queue',
            description: `[${song.title}](${song.url})`,
            fields: [
              {
                name: 'Duration',
                value: song.duration,
                inline: true
              },
              {
                name: 'Position',
                value: `${queue.songs.length}`,
                inline: true
              }
            ],
            thumbnail: {
              url: song.thumbnail
            }
          }]
        });
      }

    } catch (error) {
      console.error('Play command error:', error);
      message.reply('❌ An error occurred while trying to play music!');
    }
  },
};

function formatDuration(seconds) {
  const hours = Math.floor(seconds / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  const secs = Math.floor(seconds % 60);

  if (hours > 0) {
    return `${hours}:${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
  }
  return `${minutes}:${secs.toString().padStart(2, '0')}`;
}
