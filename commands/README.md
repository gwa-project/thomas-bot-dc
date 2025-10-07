# 📚 Commands Directory

Folder ini berisi semua command Discord bot. Setiap command memiliki file terpisah untuk kemudahan maintenance.

## 📁 Struktur Command

Setiap file command harus mengikuti struktur ini:

```javascript
module.exports = {
  name: 'commandname',           // Nama command (lowercase)
  description: 'Command description',  // Deskripsi command
  async execute(message, args, client, PREFIX, commands) {
    // Logic command di sini
    message.reply('Response');
  },
};
```

## 🎯 Parameter `execute()`:

| Parameter | Type | Description |
|-----------|------|-------------|
| `message` | Message | Discord message object |
| `args` | Array | Arguments setelah command name |
| `client` | Client | Discord client instance |
| `PREFIX` | String | Bot prefix (`!T`) |
| `commands` | Array | List semua commands (untuk help command) |

## 📝 Commands Tersedia:

| File | Command | Description |
|------|---------|-------------|
| [ping.js](ping.js) | `!T ping` | Check bot latency |
| [help.js](help.js) | `!T help` | Show command list |
| [info.js](info.js) | `!T info` | Show bot information |
| [server.js](server.js) | `!T server` | Show server information |

## ➕ Cara Menambah Command Baru:

1. Buat file baru di folder `commands/`, contoh: `test.js`

```javascript
module.exports = {
  name: 'test',
  description: 'Test command',
  async execute(message, args, client, PREFIX) {
    message.reply('Test berhasil! 🎉');
  },
};
```

2. Save file
3. Restart bot (atau deploy ke GCP)
4. Command otomatis ter-load!

**Tidak perlu edit `index.js`** - command handler akan auto-detect semua file `.js` di folder ini.

## 🔧 Tips:

- **File name** bebas, yang penting adalah property `name` di dalam module
- Gunakan **async/await** untuk operations yang perlu wait (seperti API calls)
- Tambahkan **error handling** di dalam execute function jika perlu
- Untuk command dengan permission check, tambahkan di execute function

## 📚 Contoh Command dengan Arguments:

```javascript
module.exports = {
  name: 'say',
  description: 'Bot akan mengulang text yang Anda ketik',
  async execute(message, args, client, PREFIX) {
    if (!args.length) {
      return message.reply(`❌ Usage: \`${PREFIX}say <text>\``);
    }

    const text = args.join(' ');
    message.channel.send(text);
  },
};
```

Usage: `!T say Hello World!` → Bot akan kirim "Hello World!"

---

🎉 **Happy coding!**
