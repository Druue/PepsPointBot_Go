package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	funcMap       = make(map[string]func(arg []string, message *discordgo.MessageCreate) (response string))
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

	funcMap["set-name"] = func(arg []string, message *discordgo.MessageCreate) string {
		// TODO check to make sure arg[0] is valid and good and
		// has a nice cup of coofie and all that user input sanitization
		if len(arg) != 1 {
			return "Invalid number of arguments"
		}
		setName(message.Author.ID, arg[0])
		return ":thumbsup:"
	}

	funcMap["get-name"] = func(arg []string, message *discordgo.MessageCreate) string {
		return getNameOr(message.Author.ID, message.Author.Username)
	}

	funcMap["give"] = func(arg []string, message *discordgo.MessageCreate) string {
		if len(arg) != 2 {
			return "Invalid number of arguments!"
		}
		recipient, ok := parseUserIDFromAt(arg[0])
		if !ok {
			return fmt.Sprintf("Recipient not defined, what is a %s :thinking:", arg[0])
		}
		amount, err := strconv.ParseInt(arg[1], 10, 64)
		if err != nil {
			return fmt.Sprintf("%s is not a number :thumbsdown:", arg[1])
		}
		//logTransaction(message.Author.ID, recipient, int(amount))
		return fmt.Sprintf("%s has given %d points to %s :thumbsup:",
			message.Author.ID, amount, recipient)
	}

	discord.AddHandler(commandHandler)
	discord.AddHandler(func(discord *discordgo.Session, ready *discordgo.Ready) {
		err = discord.UpdateStatus(0, "A Friendly bot!")
		errCheck("Error attempting to set status", err)

		servers := discord.State.Guilds
		fmt.Printf("BOT has started on %d servers", len(servers))
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
			_, err := discord.ChannelMessageSend(message.ChannelID, fun(args, message))
			errCheck("Oepsie woepsie, er was een stukkiewukkie in 't command handler", err)
		} else {
			discord.ChannelMessageSend(message.ChannelID, "Invalid command")
		}
	}
}
