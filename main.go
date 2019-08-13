package main

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	funcMap       = make(map[string]*Function)
	commandPrefix = "?"
	//DB varaible to handle sql connection
	DB      *sql.DB
	discord *discordgo.Session
)

//Some core basics to get going
func main() {
	localDiscord, err := discordgo.New("Bot " + getToken())
	discord = localDiscord
	errCheck("Error creating discord session", err)

	DB, err := openDBConnection("Data Source=./PepPointsDBTest.db;Version=3") //??
	errCheck("Error estabilishing database session", err)
	defer DB.Close()

	funcName := "help"
	funcMap[funcName] = NewFunction(funcName, help, 0, 0)

	funcName = "set-prefix"
	funcMap[funcName] = NewFunction(funcName, setPrefix, 1, 1)

	funcName = "set-name"
	funcMap[funcName] = NewFunction(funcName, setName, 1, 1)

	funcName = "get-name"
	funcMap[funcName] = NewFunction(funcName, getName, 0, 0)

	funcName = "give"
	funcMap[funcName] = NewFunction(funcName, givePoints, 2, 2)

	discord.AddHandler(commandHandler)
	discord.AddHandler(func(discord *discordgo.Session, ready *discordgo.Ready) {
		err = discord.UpdateStatus(0, "A Friendly bot!")
		errCheck("Error attempting to set status", err)

		servers := discord.State.Guilds
		fmt.Printf("BOT has started on %d servers", len(servers))
	})
	discord.AddHandler(func(s *discordgo.Session, event *discordgo.GuildCreate) {
		if event.Guild.Unavailable {
			fmt.Printf("\nHi, the bot was here")
			return
		}

		for _, channel := range event.Guild.Channels {
			if channel.Name == "general" {
				_, _ = s.ChannelMessageSend(channel.ID, helpCommands())
				return
			}
		}
	})

	err = discord.Open()
	defer discord.Close()

	<-make(chan struct{})
}

func commandHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
	user := message.Author
	if user.Bot {
		return
	}

	if string(message.Content[0]) == commandPrefix {
		rawMessage := strings.Split(string(message.Content[1:]), " ")
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
			_, err := discord.ChannelMessageSend(message.ChannelID, fun.def(args, message))
			errCheck("Oepsie woepsie, er was een stukkiewukkie in 't command handler", err)
		} else {
			discord.ChannelMessageSend(message.ChannelID, "Invalid command")
		}
	}
}
