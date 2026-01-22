# Quick Reference Guide

## What You Need to Do

### 1. Get Your Discord Bot Token

1. Visit https://discord.com/developers/applications
2. Click "New Application" and give it a name
3. Go to "Bot" section in the left sidebar
4. Click "Add Bot"
5. Click "Reset Token" to get your bot token
6. **IMPORTANT:** Enable "MESSAGE CONTENT INTENT" under Privileged Gateway Intents
7. Save the token somewhere safe

### 2. Invite Bot to Your Server

1. Go to "OAuth2" → "URL Generator" in the left sidebar
2. Select scopes:
   - ✅ `bot`
3. Select bot permissions:
   - ✅ Send Messages
   - ✅ Embed Links
   - ✅ Read Message History
   - ✅ Read Messages/View Channels
4. Copy the generated URL and open it in your browser
5. Select your server and authorize

### 3. Set Up the Bot

```bash
# Copy the example config file
cp .env.example .env

# Edit .env and add your token
# Replace YOUR_ACTUAL_BOT_TOKEN_FROM_DISCORD_DEVELOPER_PORTAL with your real token

# Install dependencies
go mod download

# Run the bot
export $(cat .env | xargs) && go run main.go
```

### 4. Test in Discord

Once the bot is running, go to your Discord server and try:

- `!help` - See all commands
- `!ping` - Test bot response
- `!info` - Bot information
- `!server` - Server details
- `!userinfo` - Your user info

## Build for Production

```bash
# Build binary
go build -o discord-bot

# Run binary (Linux/Mac)
export $(cat .env | xargs) && ./discord-bot

# Run binary (Windows PowerShell)
Get-Content .env | ForEach-Object { $var = $_.Split('='); [Environment]::SetEnvironmentVariable($var[0], $var[1], 'Process') }
.\discord-bot.exe
```

## Troubleshooting

**Bot doesn't respond:**
- Check MESSAGE CONTENT INTENT is enabled
- Verify bot has permissions in the channel
- Check the token is correct in .env

**Bot is offline:**
- Make sure the bot is running
- Check your internet connection
- Verify the token is valid

**Environment variable error:**
- Make sure .env file exists
- Check environment variables are loaded
- Try setting token directly: `DISCORD_BOT_TOKEN=your_token go run main.go`

## Next Steps

You can extend the bot by:
- Adding more commands in `main.go`
- Implementing slash commands
- Adding a database
- Creating moderation features
- Adding music playback
- Integrating external APIs

For detailed instructions, see SETUP.md
