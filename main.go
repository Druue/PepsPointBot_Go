package main

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Secret struct {
	DISCORD_TOKEN string
	DB_USER       string
	DB_PASSWORD   string
	DB_NAME       string
	DB_HOST       string
	DB_PORT       string
}

var (
	funcMap = make(map[string]*Function)
	//DB varaible to handle sql connection
	DB      *sql.DB
	discord *discordgo.Session
	ready   = false
)

//Some core basics to get going
func main() {
	localDiscord, err := discordgo.New("Bot " + SECRET.DISCORD_TOKEN)
	discord = localDiscord
	errCheck("Error creating discord session", err)
	openDBConnection()
	errCheck("Error establishing database session", err)
	defer DB.Close()

	funcName := "help"
	funcMap[funcName] = NewFunction(funcName, help, 0, 0, &Description{
		description:    "Returns the list of commands and their descriptions",
		argDescription: []string{},
	})

	funcName = "prefix"
	funcMap[funcName] = NewFunction(funcName, setPrefix, 1, 1, &Description{
		description:    "Updates the prefix that the bot uses to identify commands",
		argDescription: []string{"Your new prefix"},
	})

	funcName = "setnick"
	funcMap[funcName] = NewFunction(funcName, setNick, 1, 1, &Description{
		description:    "Sets your own nickname, which the bot uses when printing how many points people have",
		argDescription: []string{"Your new nickname"},
	})

	funcName = "getnick"
	funcMap[funcName] = NewFunction(funcName, getNick, 0, 0, &Description{
		description:    "Returns the nickname this bot uses to refer to you",
		argDescription: []string{},
	})

	funcName = "clearnick"
	funcMap[funcName] = NewFunction(funcName, clearNick, 0, 0, &Description{
		description:    "Resets your nickname for the bot",
		argDescription: []string{},
	})

	funcName = "givepoints"
	funcMap[funcName] = NewFunction(funcName, givePoints, 2, 2, &Description{
		description:    "Gives a user an amount of your points",
		argDescription: []string{"The user in question", "The amount of points (must be an integer)"},
	})

	funcName = "points"
	funcMap[funcName] = NewFunction(funcName, pointsCommand, 0, 1, &Description{
		description:    "Prints the amount of points",
		argDescription: []string{"Returns the amount of points of the pinged user, if this is not set, it will return all points you have given"},
	})

	discord.AddHandler(func(discord *discordgo.Session, message *discordgo.MessageCreate) {
		user := message.Author
		if user.Bot {
			return
		}
		prefix := getGuildPrefix(message.GuildID)
		if prefix == nil {
			fmt.Println("the frick, this shouldnt be nil")
			return
		}

		if len(message.Content) > len(*prefix) && message.Content[0:len(*prefix)] == *prefix {
			rawMessage := strings.Split(message.Content[1:], " ")
			funcName := rawMessage[0]
			args := rawMessage[1:]
			fun, ok := funcMap[funcName]
			if ok {
				if fun.minArgsLen > len(args) {
					_, err := discord.ChannelMessageSend(message.ChannelID, "too few arguments")
					errCheck("Oepsie woepsie, er was een stukkiewukkie in 't command handler", err)
					return
				}
				if fun.maxArgsLen < len(args) {
					_, err := discord.ChannelMessageSend(message.ChannelID, "too many arguments")
					errCheck("Oepsie woepsie, er was een stukkiewukkie in 't command handler", err)
					return
				}
				res, action := fun.def(args, message)
				switch action {
				case ResponseReply:
					_, err := discord.ChannelMessageSend(message.ChannelID, res)
					errCheck("Oepsie woepsie, er was een stukkiewukkie in 't command handler", err)
					break
				case ResponsePM:
					channel, err := discord.UserChannelCreate(message.Author.ID)
					errCheck("Oepsie woepsie, er was een stukkiewukkie in 't command handler", err)
					_, err = discord.ChannelMessageSend(channel.ID, res)
					errCheck("Oepsie woepsie, er was een stukkiewukkie in 't command handler", err)
					break
				}

			} else {
				discord.ChannelMessageSend(message.ChannelID, "Invalid command")
			}
		}
	})
	discord.AddHandler(func(discord *discordgo.Session, ready *discordgo.Ready) {
		err = discord.UpdateStatus(0, "gender non-conformity")
		errCheck("Error attempting to set status", err)
		var guilds []string
		for _, s := range ready.Guilds {
			guilds = append(guilds, s.ID)
		}
		startupAddAllGuilds(guilds)
	})
	discord.AddHandler(func(s *discordgo.Session, event *discordgo.GuildCreate) {
		if event.Guild.Unavailable {
			fmt.Printf("\nHi, the bot was here")
			return
		}
	})
	discord.AddHandler(func(s *discordgo.Session, chunk *discordgo.GuildMembersChunk) {
		var users []string
		for i := 0; i < len(chunk.Members); i++ {
			users = append(users, chunk.Members[i].User.ID)
		}
		startupAddAllUsers(users)
	})
	waitForMemberFetch(discord, func(discord *discordgo.Session) {
		ready = true
	})

	err = discord.Open()
	defer discord.Close()

	<-make(chan struct{})
}
