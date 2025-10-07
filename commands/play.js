const play = require('play-dl');
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
      message.channel.send('🔍 Searching...');

      // Search YouTube
      let songInfo;

      if (play.yt_validate(query) === 'video') {
        // Direct YouTube URL
        songInfo = await play.video_info(query);
      } else {
        // Search by name
        const searched = await play.search(query, {
          limit: 1
        });

        if (!searched || searched.length === 0) {
          return message.reply('❌ No results found!');
        }

        songInfo = searched[0];
      }

      // Create song object
      const song = {
        title: songInfo.title,
        url: songInfo.url,
        duration: formatDuration(songInfo.durationInSec),
        thumbnail: songInfo.thumbnails[0]?.url || '',
        requester: message.author.tag
      };

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
  const secs = seconds % 60;

  if (hours > 0) {
    return `${hours}:${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
  }
  return `${minutes}:${secs.toString().padStart(2, '0')}`;
}
