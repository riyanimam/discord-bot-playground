# Discord Bot Setup Guide

This guide will walk you through setting up and running your Discord bot written in Go.

## Prerequisites

Before you begin, make sure you have the following installed:

1. **Go (Golang)** - Version 1.21 or higher
   - Download from: https://golang.org/dl/
   - Verify installation: `go version`

2. **Git** (optional, for cloning)
   - Download from: https://git-scm.com/downloads

## Step 1: Create a Discord Application and Bot

1. **Go to the Discord Developer Portal**
   - Visit: https://discord.com/developers/applications
   - Log in with your Discord account

2. **Create a New Application**
   - Click the "New Application" button (top right)
   - Give your application a name (e.g., "My Go Bot")
   - Click "Create"

3. **Create a Bot User**
   - In the left sidebar, click on "Bot"
   - Click "Add Bot" and confirm
   - You should see your bot's username and avatar

4. **Get Your Bot Token**
   - Under the bot's username, click "Reset Token" (or "Copy" if you see it)
   - **IMPORTANT**: Copy this token and save it somewhere safe
   - **WARNING**: Never share this token publicly or commit it to Git!

5. **Configure Bot Permissions**
   - Scroll down to "Privileged Gateway Intents"
   - Enable "MESSAGE CONTENT INTENT" (required for reading messages)
   - Enable "SERVER MEMBERS INTENT" (optional, for member information)
   - Click "Save Changes"

## Step 2: Invite Your Bot to Your Server

1. **Generate an Invite Link**
   - In the left sidebar, click "OAuth2" → "URL Generator"
   - Under "SCOPES", select:
     - ✅ `bot`
   - Under "BOT PERMISSIONS", select:
     - ✅ `Send Messages`
     - ✅ `Embed Links`
     - ✅ `Read Message History`
     - ✅ `Read Messages/View Channels`
     - ✅ `Use Slash Commands` (optional, for future features)
   
2. **Copy and Use the Invite Link**
   - Copy the generated URL at the bottom
   - Paste it in your browser
   - Select the server you want to add the bot to
   - Click "Authorize"

## Step 3: Configure the Bot

1. **Set Up Environment Variables**
   ```bash
   # Copy the example environment file
   cp .env.example .env
   ```

2. **Edit the `.env` file**
   - Open `.env` in a text editor
   - Replace `your_discord_bot_token_here` with your actual bot token
   - Optionally, change the `BOT_PREFIX` (default is `!`)
   
   Example:
   ```
   DISCORD_BOT_TOKEN=YOUR_ACTUAL_BOT_TOKEN_FROM_DISCORD_DEVELOPER_PORTAL
   BOT_PREFIX=!
   ```

## Step 4: Install Dependencies

Run the following command to download the required Go modules:

```bash
go mod download
```

## Step 5: Run the Bot

### Option 1: Run Directly

```bash
# Set environment variables and run
export $(cat .env | xargs) && go run main.go
```

### Option 2: Build and Run

```bash
# Build the binary
go build -o discord-bot

# Run the binary (Linux/Mac)
export $(cat .env | xargs) && ./discord-bot

# Run the binary (Windows - PowerShell)
Get-Content .env | ForEach-Object { $var = $_.Split('='); [Environment]::SetEnvironmentVariable($var[0], $var[1], 'Process') }
.\discord-bot.exe
```

### Option 3: Using Environment Variables Directly (Linux/Mac)

```bash
DISCORD_BOT_TOKEN=your_token_here BOT_PREFIX=! go run main.go
```

## Step 6: Test Your Bot

Once the bot is running, you should see:
```
Bot is ready! Logged in as: YourBotName#1234
Bot is in X guilds
Bot is now running. Press CTRL-C to exit.
```

Now go to your Discord server and try these commands:

- `!help` - See all available commands
- `!ping` - Check bot latency
- `!info` - Display bot information
- `!server` - Show server information
- `!userinfo` - Show your user information
- `!userinfo @username` - Show another user's information

## Available Commands

| Command | Description |
|---------|-------------|
| `!ping` | Check the bot's latency |
| `!help` | Display all available commands |
| `!info` | Show information about the bot |
| `!server` | Display current server information |
| `!userinfo [@user]` | Show user information |

## Troubleshooting

### Bot doesn't respond to commands
- Make sure "MESSAGE CONTENT INTENT" is enabled in the Discord Developer Portal
- Check that the bot has permission to read and send messages in the channel
- Verify your bot token is correct in the `.env` file
- Make sure you're using the correct prefix (default is `!`)

### "DISCORD_BOT_TOKEN environment variable is required" error
- Ensure your `.env` file exists and contains your bot token
- Make sure you're loading the environment variables correctly
- Try running with the token directly: `DISCORD_BOT_TOKEN=your_token go run main.go`

### Bot is offline in Discord
- Check that the bot is running (you should see "Bot is now running" message)
- Verify your bot token is valid
- Check your internet connection

## Running in Production

For production deployment, consider:

1. **Using a process manager** (e.g., systemd, PM2, or Docker)
2. **Setting environment variables** at the system level instead of using `.env`
3. **Monitoring and logging** the bot's output
4. **Auto-restart** on crashes
5. **Running on a VPS or cloud service** (AWS, DigitalOcean, Heroku, etc.)

### Example systemd Service (Linux)

Create `/etc/systemd/system/discord-bot.service`:

```ini
[Unit]
Description=Discord Bot
After=network.target

[Service]
Type=simple
User=your-user
WorkingDirectory=/path/to/discord-bot-playground
Environment="DISCORD_BOT_TOKEN=your_token_here"
Environment="BOT_PREFIX=!"
ExecStart=/usr/local/go/bin/go run main.go
Restart=always

[Install]
WantedBy=multi-user.target
```

Then:
```bash
sudo systemctl daemon-reload
sudo systemctl enable discord-bot
sudo systemctl start discord-bot
sudo systemctl status discord-bot
```

## Security Best Practices

1. **Never commit your `.env` file** - It's already in `.gitignore`
2. **Keep your bot token secret** - Treat it like a password
3. **Regenerate your token** if it's ever exposed
4. **Use environment variables** for sensitive data
5. **Regularly update dependencies** for security patches

## Next Steps

You can extend this bot by:
- Adding more commands in the `messageCreate` function
- Implementing slash commands
- Adding a database for persistent storage
- Creating custom embeds and reactions
- Integrating with external APIs
- Adding moderation features
- Implementing music playback

## Resources

- [DiscordGo Documentation](https://pkg.go.dev/github.com/bwmarrin/discordgo)
- [Discord Developer Documentation](https://discord.com/developers/docs)
- [Go Documentation](https://golang.org/doc/)
- [Discord API Server](https://discord.gg/discord-api)

## Support

If you encounter issues:
1. Check the troubleshooting section above
2. Review the Discord Developer Portal settings
3. Check the bot logs for error messages
4. Open an issue on the GitHub repository

---

**Happy botting! 🤖**
