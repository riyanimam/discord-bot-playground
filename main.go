package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	botToken  string
	botPrefix string
)

func init() {
	// Load environment variables
	botToken = os.Getenv("DISCORD_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("DISCORD_BOT_TOKEN environment variable is required")
	}

	botPrefix = os.Getenv("BOT_PREFIX")
	if botPrefix == "" {
		botPrefix = "!" // Default prefix
	}
}

func main() {
	// Create a new Discord session using the provided bot token
	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatal("Error creating Discord session:", err)
	}

	// Register the messageCreate func as a callback for MessageCreate events
	dg.AddHandler(messageCreate)

	// Register the ready handler
	dg.AddHandler(ready)

	// We need information about guilds (servers) and messages
	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds

	// Open a websocket connection to Discord
	err = dg.Open()
	if err != nil {
		log.Fatal("Error opening connection:", err)
	}
	defer dg.Close()

	fmt.Println("Bot is now running. Press CTRL-C to exit.")

	// Wait here until CTRL-C or other term signal is received
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	fmt.Println("\nGracefully shutting down...")
}

// ready is called when the bot successfully connects to Discord
func ready(s *discordgo.Session, event *discordgo.Ready) {
	log.Printf("Bot is ready! Logged in as: %s#%s", event.User.Username, event.User.Discriminator)
	log.Printf("Bot is in %d guilds", len(event.Guilds))

	// Set the bot's status
	err := s.UpdateGameStatus(0, fmt.Sprintf("%shelp for commands", botPrefix))
	if err != nil {
		log.Println("Error setting game status:", err)
	}
}

// messageCreate is called whenever a new message is created
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Check if the message starts with the bot prefix
	if !strings.HasPrefix(m.Content, botPrefix) {
		return
	}

	// Remove the prefix and split into command and arguments
	content := strings.TrimPrefix(m.Content, botPrefix)
	args := strings.Fields(content)

	if len(args) == 0 {
		return
	}

	command := strings.ToLower(args[0])

	// Handle commands
	switch command {
	case "ping":
		handlePing(s, m)
	case "help":
		handleHelp(s, m)
	case "info":
		handleInfo(s, m)
	case "server":
		handleServerInfo(s, m)
	case "userinfo":
		handleUserInfo(s, m, args)
	default:
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Unknown command: `%s`. Use `%shelp` to see available commands.", command, botPrefix))
	}
}

// handlePing responds with "Pong!" and the latency
func handlePing(s *discordgo.Session, m *discordgo.MessageCreate) {
	start := time.Now()
	msg, err := s.ChannelMessageSend(m.ChannelID, "Pinging...")
	if err != nil {
		return
	}

	latency := time.Since(start)
	s.ChannelMessageEdit(m.ChannelID, msg.ID, fmt.Sprintf("🏓 Pong! Latency: %dms", latency.Milliseconds()))
}

// handleHelp sends a help message with all available commands
func handleHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "📚 Bot Commands",
		Description: fmt.Sprintf("Here are all the commands you can use (prefix: `%s`):", botPrefix),
		Color:       0x00ff00,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   fmt.Sprintf("%sping", botPrefix),
				Value:  "Check the bot's latency",
				Inline: false,
			},
			{
				Name:   fmt.Sprintf("%shelp", botPrefix),
				Value:  "Show this help message",
				Inline: false,
			},
			{
				Name:   fmt.Sprintf("%sinfo", botPrefix),
				Value:  "Display information about the bot",
				Inline: false,
			},
			{
				Name:   fmt.Sprintf("%sserver", botPrefix),
				Value:  "Show information about the current server",
				Inline: false,
			},
			{
				Name:   fmt.Sprintf("%suserinfo [@user]", botPrefix),
				Value:  "Show information about yourself or a mentioned user",
				Inline: false,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

// handleInfo sends information about the bot
func handleInfo(s *discordgo.Session, m *discordgo.MessageCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "ℹ️ Bot Information",
		Description: "A simple Discord bot written in Go!",
		Color:       0x0099ff,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Language",
				Value:  "Go (Golang)",
				Inline: true,
			},
			{
				Name:   "Library",
				Value:  "DiscordGo",
				Inline: true,
			},
			{
				Name:   "Prefix",
				Value:  fmt.Sprintf("`%s`", botPrefix),
				Inline: true,
			},
			{
				Name:   "Repository",
				Value:  "[GitHub](https://github.com/riyanimam/discord-bot-playground)",
				Inline: false,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Made with ❤️ using Go",
		},
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

// handleServerInfo sends information about the current server
func handleServerInfo(s *discordgo.Session, m *discordgo.MessageCreate) {
	guild, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error retrieving server information.")
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("🏰 %s", guild.Name),
		Description: "Server Information",
		Color:       0xffa500,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: guild.IconURL("256"),
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Server ID",
				Value:  guild.ID,
				Inline: true,
			},
			{
				Name:   "Owner",
				Value:  fmt.Sprintf("<@%s>", guild.OwnerID),
				Inline: true,
			},
			{
				Name:   "Members",
				Value:  fmt.Sprintf("%d", guild.MemberCount),
				Inline: true,
			},
			{
				Name:   "Channels",
				Value:  fmt.Sprintf("%d", len(guild.Channels)),
				Inline: true,
			},
			{
				Name:   "Roles",
				Value:  fmt.Sprintf("%d", len(guild.Roles)),
				Inline: true,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

// handleUserInfo sends information about a user
func handleUserInfo(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	var targetUser *discordgo.User

	// If a user is mentioned, get their info; otherwise, use the message author
	if len(m.Mentions) > 0 {
		targetUser = m.Mentions[0]
	} else {
		targetUser = m.Author
	}

	// Get member information
	member, err := s.GuildMember(m.GuildID, targetUser.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error retrieving user information.")
		return
	}

	nickname := "None"
	if member.Nick != "" {
		nickname = member.Nick
	}

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("👤 %s#%s", targetUser.Username, targetUser.Discriminator),
		Description: "User Information",
		Color:       0x9b59b6,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: targetUser.AvatarURL("256"),
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "User ID",
				Value:  targetUser.ID,
				Inline: true,
			},
			{
				Name:   "Nickname",
				Value:  nickname,
				Inline: true,
			},
			{
				Name:   "Roles",
				Value:  fmt.Sprintf("%d", len(member.Roles)),
				Inline: true,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}
