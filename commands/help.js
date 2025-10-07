module.exports = {
  name: 'help',
  description: 'Show command list',
  async execute(message, args, client, PREFIX, commands) {
    const commandList = commands.map(cmd => ({
      name: `${PREFIX}${cmd.name}`,
      value: cmd.description,
      inline: true
    }));

    message.reply({
      embeds: [{
        color: 0x0099ff,
        title: '📚 Thomas Bot - Command List',
        description: `Prefix: \`${PREFIX}\``,
        fields: commandList,
        timestamp: new Date(),
        footer: {
          text: 'Thomas Bot v1.0'
        }
      }]
    });
  },
};
