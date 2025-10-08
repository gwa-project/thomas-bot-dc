const {
  joinVoiceChannel,
  createAudioPlayer,
  createAudioResource,
  AudioPlayerStatus,
  VoiceConnectionStatus,
  entersState,
  StreamType
} = require('@discordjs/voice');
const ytdl = require('ytdl-core');

// Map to store queue per guild
const queues = new Map();

class MusicQueue {
  constructor(guildId) {
    this.guildId = guildId;
    this.songs = [];
    this.connection = null;
    this.player = createAudioPlayer();
    this.isPlaying = false;
    this.volume = 0.5;
    this.currentSong = null;

    // Player event handlers
    this.player.on(AudioPlayerStatus.Idle, () => {
      this.playNext();
    });

    this.player.on('error', error => {
      console.error('Audio player error:', error);
      this.playNext();
    });
  }

  addSong(song) {
    this.songs.push(song);
  }

  async play(voiceChannel, textChannel) {
    if (this.songs.length === 0) {
      if (this.connection) {
        this.connection.destroy();
        this.connection = null;
      }
      return textChannel.send('❌ Queue is empty!');
    }

    // Create voice connection if not exists
    if (!this.connection) {
      try {
        this.connection = joinVoiceChannel({
          channelId: voiceChannel.id,
          guildId: voiceChannel.guild.id,
          adapterCreator: voiceChannel.guild.voiceAdapterCreator,
          selfDeaf: true
        });

        // Wait for connection to be ready
        await entersState(this.connection, VoiceConnectionStatus.Ready, 30000);

        // Subscribe player to connection
        this.connection.subscribe(this.player);
      } catch (error) {
        console.error('Failed to join voice channel:', error);
        this.connection = null;
        return textChannel.send('❌ Failed to join voice channel!');
      }
    }

    // Play current song
    this.currentSong = this.songs[0];
    this.isPlaying = true;

    try {
      const stream = ytdl(this.currentSong.url, {
        filter: 'audioonly',
        quality: 'highestaudio',
        highWaterMark: 1 << 25
      });

      const resource = createAudioResource(stream, {
        inputType: StreamType.Arbitrary,
        inlineVolume: true
      });

      resource.volume.setVolume(this.volume);
      this.player.play(resource);

      textChannel.send({
        embeds: [{
          color: 0x00ff00,
          title: '🎵 Now Playing',
          description: `[${this.currentSong.title}](${this.currentSong.url})`,
          fields: [
            {
              name: 'Duration',
              value: this.currentSong.duration,
              inline: true
            },
            {
              name: 'Requested by',
              value: this.currentSong.requester,
              inline: true
            }
          ],
          thumbnail: {
            url: this.currentSong.thumbnail
          }
        }]
      });
    } catch (error) {
      console.error('Error playing song:', error);
      textChannel.send('❌ Error playing this song, skipping...');
      this.playNext();
    }
  }

  playNext() {
    this.songs.shift(); // Remove current song
    this.currentSong = null;
    this.isPlaying = false;

    if (this.songs.length > 0) {
      this.play(null, null); // Will use existing connection
    } else {
      if (this.connection) {
        this.connection.destroy();
        this.connection = null;
      }
    }
  }

  pause() {
    if (this.player) {
      this.player.pause();
      return true;
    }
    return false;
  }

  resume() {
    if (this.player) {
      this.player.unpause();
      return true;
    }
    return false;
  }

  skip() {
    if (this.player) {
      this.player.stop();
      return true;
    }
    return false;
  }

  stop() {
    this.songs = [];
    this.currentSong = null;
    this.isPlaying = false;

    if (this.player) {
      this.player.stop();
    }

    if (this.connection) {
      this.connection.destroy();
      this.connection = null;
    }

    return true;
  }

  setVolume(volume) {
    this.volume = Math.max(0, Math.min(1, volume / 100));
    return this.volume * 100;
  }

  getQueue() {
    return this.songs;
  }

  getCurrentSong() {
    return this.currentSong;
  }
}

// Get or create queue for guild
function getQueue(guildId) {
  if (!queues.has(guildId)) {
    queues.set(guildId, new MusicQueue(guildId));
  }
  return queues.get(guildId);
}

// Delete queue
function deleteQueue(guildId) {
  const queue = queues.get(guildId);
  if (queue) {
    queue.stop();
    queues.delete(guildId);
  }
}

module.exports = {
  getQueue,
  deleteQueue
};
