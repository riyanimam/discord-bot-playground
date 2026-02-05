# Discord Bot Playground 🤖

A simple yet feature-rich Discord bot written in **Go (Golang)** using the [DiscordGo](https://github.com/bwmarrin/discordgo) library.

## Features ✨

- 🏓 **Ping command** - Check bot latency
- 📚 **Help command** - Interactive help menu
- ℹ️ **Info command** - Bot information
- 🏰 **Server info** - Display server details
- 👤 **User info** - Show user information
- 🖼️ **Avatar command** - Display user avatars in full resolution
- 🎱 **Magic 8-Ball** - Ask questions and get mystical answers
- 🎲 **Dice roller** - Roll dice with standard notation (e.g., 2d6, 1d20)
- 🪙 **Coin flip** - Flip a coin for quick decisions
- 📊 **Poll creator** - Create interactive polls with emoji reactions
- 🎨 **Rich embeds** - Beautiful embedded messages
- ⚡ **Fast and efficient** - Written in Go
- 🔧 **Easy to configure** - Simple environment variable setup

## Quick Start 🚀

### Prerequisites

- Go 1.21 or higher installed
- A Discord account and server
- A Discord bot token (see [SETUP.md](SETUP.md) for detailed instructions)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/riyanimam/discord-bot-playground.git
   cd discord-bot-playground
   ```

2. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env and add your Discord bot token
   ```

3. **Install dependencies**
   ```bash
   go mod download
   ```

4. **Run the bot**
   ```bash
   export $(cat .env | xargs) && go run main.go
   ```

## Detailed Setup Instructions 📖

For complete setup instructions including:
- Creating a Discord application
- Getting your bot token
- Inviting the bot to your server
- Running in production
- Troubleshooting

Please see **[SETUP.md](SETUP.md)** for the full guide.

## Commands 💬

All commands use the `!` prefix by default (configurable):

### Basic Commands

| Command | Description |
|---------|-------------|
| `!ping` | Check the bot's latency and responsiveness |
| `!help` | Display all available commands with descriptions |
| `!info` | Show information about the bot |

### Server & User Information

| Command | Description |
|---------|-------------|
| `!server` | Display information about the current server |
| `!userinfo [@user]` | Show information about yourself or a mentioned user |
| `!avatar [@user]` | Display user's avatar in full resolution with download link |

### Fun Commands

| Command | Description |
|---------|-------------|
| `!8ball <question>` | Ask the magic 8-ball a question and get a mystical answer |
| `!roll [dice notation]` | Roll dice using standard notation (e.g., `!roll 2d6`, `!roll 1d20`) |
| `!coinflip` or `!flip` | Flip a coin - heads or tails |

### Utility Commands

| Command | Description |
|---------|-------------|
| `!poll <question> \| <option1> \| <option2> \| ...` | Create a poll with up to 10 options. The bot will add number reactions for voting |

### Command Examples

```
!8ball Will I win the game?
!roll 2d6
!roll 1d20
!coinflip
!avatar @username
!poll What should we do? | Play games | Watch a movie | Go outside
```

## Project Structure 📁

```
discord-bot-playground/
├── main.go           # Main bot application
├── go.mod            # Go module dependencies
├── go.sum            # Go module checksums
├── .env.example      # Example environment variables
├── .gitignore        # Git ignore rules
├── README.md         # This file
└── SETUP.md          # Detailed setup guide
```

## Configuration ⚙️

The bot uses environment variables for configuration:

- `DISCORD_BOT_TOKEN` - Your Discord bot token (required)
- `BOT_PREFIX` - Command prefix (optional, defaults to `!`)

## Development 🛠️

### Building

```bash
go build -o discord-bot
```

### Running Tests

```bash
go test ./...
```

### Code Structure

The bot is organized with:
- **Event handlers** - Handle Discord events (messages, ready state)
- **Command handlers** - Process individual commands
- **Embedded responses** - Rich Discord embeds for better UX

## Tech Stack 💻

- **Language**: Go (Golang) 1.21+
- **Discord Library**: [DiscordGo](https://github.com/bwmarrin/discordgo) v0.27.1
- **Architecture**: Event-driven bot with message handlers

## Contributing 🤝

Contributions are welcome! Feel free to:
- Add new commands
- Improve existing features
- Fix bugs
- Enhance documentation

## Security 🔒

- Never commit your `.env` file
- Keep your bot token secret
- The `.env` file is already in `.gitignore`
- Regenerate your token if it's ever exposed

## License 📄

This project is open source and available under the MIT License.

## Resources 📚

- [DiscordGo Documentation](https://pkg.go.dev/github.com/bwmarrin/discordgo)
- [Discord Developer Portal](https://discord.com/developers/applications)
- [Discord API Documentation](https://discord.com/developers/docs)
- [Go Documentation](https://golang.org/doc/)

## Support 💬

If you need help:
1. Check [SETUP.md](SETUP.md) for detailed instructions
2. Review the troubleshooting section
3. Open an issue on GitHub

---

Made with ❤️ using Go