package main

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

//Function struct handles the varying parts that go into the various bot commands
//Their nickname, their description, their actual definition etc
type Function struct {
	name           string
	description    string
	argDescription []string
	def            func(arg []string, message *discordgo.MessageCreate) (res string)
	minArgsLen     int
	maxArgsLen     int
}

//NewFunction -> Constructor for the Function struct
func NewFunction(name string, def func(arg []string, message *discordgo.MessageCreate) (response string), minArgsLen int, maxArgsLen int) *Function {
	description, argDescription := getDescription(name)
	if len(argDescription) != maxArgsLen {
		panic("not every arg is descripted in " + name)
	}
	return &Function{
		name,
		description,
		argDescription,
		def,
		minArgsLen,
		maxArgsLen,
	}
}

func help(arg []string, message *discordgo.MessageCreate) string {
	return helpCommands()
}

func setPrefix(arg []string, message *discordgo.MessageCreate) string {
	commandPrefix = arg[0]
	return fmt.Sprintf("Command prefix changed to %s", commandPrefix)
}

func setName(arg []string, message *discordgo.MessageCreate) string {
	// TODO check to make sure arg[0] is valid and good and
	// has a nice cup of coofie and all that user input sanitization
	//logName(message.Author.ID, arg[0])
	setUsersNickname(&User{
		discordId: message.Author.ID,
		nickname:  arg[0],
	})
	return fmt.Sprintf("Set %s's nickname to be %s :thumbsup:", message.Author.ID, arg[0])
}

func getName(arg []string, message *discordgo.MessageCreate) string {
	return getUsersNicknameOr(message.Author.ID, message.Author.Username)
}

func givePoints(arg []string, message *discordgo.MessageCreate) string {
	recipientID, ok := parseUserIDFromAt(arg[0])
	if !ok {
		return fmt.Sprintf("Recipient not defined, what is a \"%s\" :thinking:", arg[0])
	}
	//recipient, err := discord.GuildMember(message.GuildID, recipientID)

	amount, err := strconv.ParseInt(arg[1], 10, 64)
	if err != nil {
		return fmt.Sprintf("%s is not a number :thumbsdown:", arg[1])
	}
	giveUserPoints(recipientID, message.Author.ID, amount)
	return ":thumbsup:"
}

func getPointsGiven(arg []string, message *discordgo.MessageCreate) string {
	points, nicknames := getUsersPointsGiven(message.Author.ID)
	if len(arg) == 1 {
		//return points given to individual person
		return "no implemented yet"
	}
	re := "you have:"
	for i := 0; i < len(points); i++ {
		re += "\t" + nicknames[i] + " " + strconv.FormatInt(points[i].amount, 10) + " points\n"
	}
	return re
}

func helpCommands() string {
	var buffer bytes.Buffer

	buffer.WriteString("```")

	for _, v := range funcMap {
		buffer.WriteString(fmt.Sprintf("Command: %s\n\ndescription: %s\n", v.name, v.description))
		for i, argDesc := range v.argDescription {
			isOptionalString := ""
			if i > v.minArgsLen {
				isOptionalString = "*"
			}
			buffer.WriteString(fmt.Sprintf("Argument%s %s: %s\n", isOptionalString, strconv.Itoa(i), argDesc))
		}
		buffer.WriteString("\n\n\n")
	}
	buffer.WriteString(fmt.Sprintf("* argument is optional\n\n\nPrefix is currently: %s", commandPrefix))

	buffer.WriteString("```")

	return buffer.String()
}
