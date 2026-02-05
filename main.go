package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"regexp"
	"strconv"
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

	// We need information about guilds (servers), messages, and message content
	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds | discordgo.IntentsMessageContent

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
	log.Printf("Bot is ready! Logged in as: %s", event.User.String())
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
	case "8ball":
		handle8Ball(s, m, args)
	case "roll":
		handleRoll(s, m, args)
	case "coinflip", "flip":
		handleCoinFlip(s, m)
	case "avatar":
		handleAvatar(s, m, args)
	case "poll":
		handlePoll(s, m, args)
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
				Name:   "**🔧 Basic Commands**",
				Value:  "\u200b",
				Inline: false,
			},
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
				Name:   "**📊 Server & User Info**",
				Value:  "\u200b",
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
			{
				Name:   fmt.Sprintf("%savatar [@user]", botPrefix),
				Value:  "Display user's avatar in full resolution",
				Inline: false,
			},
			{
				Name:   "**🎮 Fun Commands**",
				Value:  "\u200b",
				Inline: false,
			},
			{
				Name:   fmt.Sprintf("%s8ball <question>", botPrefix),
				Value:  "Ask the magic 8-ball a question",
				Inline: false,
			},
			{
				Name:   fmt.Sprintf("%sroll [dice notation]", botPrefix),
				Value:  "Roll dice (e.g., `1d6`, `2d20`, or just `roll` for 1d6)",
				Inline: false,
			},
			{
				Name:   fmt.Sprintf("%scoinflip", botPrefix),
				Value:  "Flip a coin (heads or tails)",
				Inline: false,
			},
			{
				Name:   "**🛠️ Utility Commands**",
				Value:  "\u200b",
				Inline: false,
			},
			{
				Name:   fmt.Sprintf("%spoll <question> | <option1> | <option2> | ...", botPrefix),
				Value:  "Create a poll with up to 10 options",
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
		Title:       fmt.Sprintf("👤 %s", targetUser.String()),
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

// handle8Ball responds with a magic 8-ball answer
func handle8Ball(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Please ask a question! Usage: `%s8ball <your question>`", botPrefix))
		return
	}

	responses := []string{
		"It is certain.",
		"It is decidedly so.",
		"Without a doubt.",
		"Yes definitely.",
		"You may rely on it.",
		"As I see it, yes.",
		"Most likely.",
		"Outlook good.",
		"Yes.",
		"Signs point to yes.",
		"Reply hazy, try again.",
		"Ask again later.",
		"Better not tell you now.",
		"Cannot predict now.",
		"Concentrate and ask again.",
		"Don't count on it.",
		"My reply is no.",
		"My sources say no.",
		"Outlook not so good.",
		"Very doubtful.",
	}

	question := strings.Join(args[1:], " ")
	answer := responses[rand.Intn(len(responses))]

	embed := &discordgo.MessageEmbed{
		Title:       "🎱 Magic 8-Ball",
		Description: fmt.Sprintf("**Question:** %s\n\n**Answer:** %s", question, answer),
		Color:       0x8b00ff,
		Timestamp:   time.Now().Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Asked by %s", m.Author.Username),
		},
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

// handleRoll rolls dice using standard dice notation (e.g., 2d6, 1d20)
func handleRoll(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	diceNotation := "1d6" // Default
	if len(args) > 1 {
		diceNotation = strings.ToLower(args[1])
	}

	// Parse dice notation (e.g., "2d6" means 2 dice with 6 sides each)
	re := regexp.MustCompile(`^(\d+)d(\d+)$`)
	matches := re.FindStringSubmatch(diceNotation)

	if matches == nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Invalid dice notation! Use format like `1d6`, `2d20`, etc. Usage: `%sroll [XdY]`", botPrefix))
		return
	}

	numDice, _ := strconv.Atoi(matches[1])
	numSides, _ := strconv.Atoi(matches[2])

	// Sanity checks
	if numDice < 1 || numDice > 100 {
		s.ChannelMessageSend(m.ChannelID, "Please roll between 1 and 100 dice!")
		return
	}
	if numSides < 2 || numSides > 1000 {
		s.ChannelMessageSend(m.ChannelID, "Please use dice with 2 to 1000 sides!")
		return
	}

	// Roll the dice
	rolls := make([]int, numDice)
	total := 0
	for i := 0; i < numDice; i++ {
		rolls[i] = rand.Intn(numSides) + 1
		total += rolls[i]
	}

	// Format the results
	rollsStr := ""
	if numDice <= 20 {
		// Show individual rolls if not too many
		rollsStrs := make([]string, numDice)
		for i, roll := range rolls {
			rollsStrs[i] = fmt.Sprintf("%d", roll)
		}
		rollsStr = fmt.Sprintf("\n**Rolls:** %s", strings.Join(rollsStrs, ", "))
	}

	embed := &discordgo.MessageEmbed{
		Title:       "🎲 Dice Roll",
		Description: fmt.Sprintf("**Dice:** %s%s\n**Total:** %d", diceNotation, rollsStr, total),
		Color:       0xff6347,
		Timestamp:   time.Now().Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Rolled by %s", m.Author.Username),
		},
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

// handleCoinFlip flips a coin
func handleCoinFlip(s *discordgo.Session, m *discordgo.MessageCreate) {
	result := "Heads"
	emoji := "🪙"
	if rand.Intn(2) == 1 {
		result = "Tails"
	}

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("%s Coin Flip", emoji),
		Description: fmt.Sprintf("The coin landed on: **%s**!", result),
		Color:       0xffd700,
		Timestamp:   time.Now().Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Flipped by %s", m.Author.Username),
		},
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

// handleAvatar displays a user's avatar in full resolution
func handleAvatar(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	var targetUser *discordgo.User

	// If a user is mentioned, get their avatar; otherwise, use the message author
	if len(m.Mentions) > 0 {
		targetUser = m.Mentions[0]
	} else {
		targetUser = m.Author
	}

	avatarURL := targetUser.AvatarURL("1024")

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("🖼️ %s's Avatar", targetUser.Username),
		Color:       0x00bfff,
		Image:       &discordgo.MessageEmbedImage{URL: avatarURL},
		Description: fmt.Sprintf("[Download Link](%s)", avatarURL),
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

// handlePoll creates a simple poll with reactions
func handlePoll(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Please provide a question and options! Usage: `%spoll <question> | <option1> | <option2> | ...`", botPrefix))
		return
	}

	// Join all args and split by pipe
	content := strings.Join(args[1:], " ")
	parts := strings.Split(content, "|")

	if len(parts) < 3 {
		s.ChannelMessageSend(m.ChannelID, "Please provide at least a question and 2 options separated by `|`")
		return
	}

	question := strings.TrimSpace(parts[0])
	options := make([]string, 0)
	for i := 1; i < len(parts); i++ {
		option := strings.TrimSpace(parts[i])
		if option != "" {
			options = append(options, option)
		}
	}

	if len(options) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Please provide at least 2 options!")
		return
	}

	if len(options) > 10 {
		s.ChannelMessageSend(m.ChannelID, "Maximum 10 options allowed!")
		return
	}

	// Number emojis for reactions
	numberEmojis := []string{"1️⃣", "2️⃣", "3️⃣", "4️⃣", "5️⃣", "6️⃣", "7️⃣", "8️⃣", "9️⃣", "🔟"}

	// Build poll description
	pollDescription := ""
	for i, option := range options {
		pollDescription += fmt.Sprintf("%s %s\n", numberEmojis[i], option)
	}

	embed := &discordgo.MessageEmbed{
		Title:       "📊 " + question,
		Description: pollDescription,
		Color:       0x3498db,
		Timestamp:   time.Now().Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Poll created by %s", m.Author.Username),
		},
	}

	msg, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Failed to create poll. Please try again.")
		return
	}

	// Add reactions to the poll
	for i := 0; i < len(options); i++ {
		s.MessageReactionAdd(m.ChannelID, msg.ID, numberEmojis[i])
	}
}
