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
	DB *sql.DB
)

//Some core basics to get going
func main() {
	discord, err := discordgo.New("Bot " + getToken())
	errCheck("Error creating discord session", err)

	/*
		DB, err := openDBConnection("CONNECTION STRING - NYI")
		errCheck("Error estabilishing database session", err)
		defer DB.Close()
	*/
	funcName := "help"
	funcMap[funcName] = NewFunction(funcName, getDescription(funcName), help)

	funcName = "set-prefix"
	funcMap[funcName] = NewFunction(funcName, getDescription(funcName), setPrefix)

	funcName = "set-name"
	funcMap[funcName] = NewFunction(funcName, getDescription(funcName), setName)

	funcName = "get-name"
	funcMap[funcName] = NewFunction(funcName, getDescription(funcName), getName)

	funcName = "give"
	funcMap[funcName] = NewFunction(funcName, getDescription(funcName), givePoints)

	discord.AddHandler(commandHandler)
	discord.AddHandler(func(discord *discordgo.Session, ready *discordgo.Ready) {
		err = discord.UpdateStatus(0, "A Friendly bot!")
		errCheck("Error attempting to set status", err)

		servers := discord.State.Guilds
		fmt.Printf("BOT has started on %d servers", len(servers))
	})
	discord.AddHandler(func(s *discordgo.Session, event *discordgo.GuildCreate) {
		if event.Guild.Unavailable {
			return
		}

		for _, channel := range event.Guild.Channels {
			if channel.ID == event.Guild.ID {
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
			_, err := discord.ChannelMessageSend(message.ChannelID, fun.def(args, message, 0, 0))
			errCheck("Oepsie woepsie, er was een stukkiewukkie in 't command handler", err)
		} else {
			discord.ChannelMessageSend(message.ChannelID, "Invalid command")
		}
	}
}
