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
	funcMap       = make(map[string]*Function)
	commandPrefix = "?"
	//DB varaible to handle sql connection
	DB      *sql.DB
	discord *discordgo.Session
)

//Some core basics to get going
func main() {
	localDiscord, err := discordgo.New("Bot " + SECRET.DISCORD_TOKEN)
	discord = localDiscord
	errCheck("Error creating discord session", err)
	openDBConnection()
	errCheck("Error estabilishing database session", err)
	defer DB.Close()

	funcName := "help"
	funcMap[funcName] = NewFunction(funcName, help, 0, 0)

	funcName = "set-prefix"
	funcMap[funcName] = NewFunction(funcName, setPrefix, 1, 1)

	funcName = "set-nickname"
	funcMap[funcName] = NewFunction(funcName, setName, 1, 1)

	funcName = "get-nickname"
	funcMap[funcName] = NewFunction(funcName, getName, 0, 0)

	funcName = "give"
	funcMap[funcName] = NewFunction(funcName, givePoints, 2, 2)

	discord.AddHandler(commandHandler)
	discord.AddHandler(func(discord *discordgo.Session, ready *discordgo.Ready) {
		err = discord.UpdateStatus(0, "A Friendly bot!")
		errCheck("Error attempting to set status", err)

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
	go waitForMemberFetch(discord, func(discord *discordgo.Session) {
		var allUsers []string
		servers := discord.State.Guilds
		for i := 0; i < len(servers); i++ {
			for j := 0; j < len(servers[i].Members); j++ {
				allUsers = append(allUsers, servers[i].Members[j].User.ID)
			}
		}
		startupAddAllUsers(allUsers)
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
