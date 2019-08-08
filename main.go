package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var commandPrefix = "?"
var funcMap = make(map[string]func(arg []string, message *discordgo.MessageCreate) (response string))

//Some core basics to get going
func main() {
	discord, err := discordgo.New("Bot " + getToken())
	errCheck("error creating discord session", err)

	funcMap["set-name"] = func(arg []string, message *discordgo.MessageCreate) string {
		//TODO check to make sure arg[0] is valid and good and has a nice cup of coofie and all that user input sanitization
		if len(arg) != 1 {
			return "wrong number of arguments"
		}
		setName(message.Author.ID, arg[0])
		return ":thumbsup:"
	}

	funcMap["get-name"] = func(arg []string, message *discordgo.MessageCreate) string {
		return getNameOr(message.Author.ID, message.Author.Username)
	}

	funcMap["give"] = func(arg []string, message *discordgo.MessageCreate) string {
		if len(arg) != 2 {
			return "wrong number of arguments"
		}
		recipient, ok := parseUserIDFromAt(arg[0])
		if !ok {
			return ":thumbsdown:, could not parse first argument"
		}
		amount, err := strconv.ParseInt(arg[1], 10, 64)
		if err != nil {
			return ":thumbsdown:, second argument is not a number"
		}
		addPoint(message.Author.ID, recipient, int(amount))
		return ":thumbsup:"
	}

	discord.AddHandler(commandHandler)
	discord.AddHandler(func(discord *discordgo.Session, ready *discordgo.Ready) {
		err = discord.UpdateStatus(0, "A Friendly bot!")
		if err != nil {
			fmt.Println("Error attempting to set status")
		}

		servers := discord.State.Guilds
		fmt.Printf("BOT has started on %d servers", len(servers))
	})

	err = discord.Open()
	defer discord.Close()

	<-make(chan struct{})
}

func errCheck(msg string, err error) {
	if err != nil {
		fmt.Printf("%s: %+v", msg, err)
		panic(err)
	}
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
			errCheck("", err)
		} else {
			discord.ChannelMessageSend(message.ChannelID, "function couldnt be found")
		}
	}
}
