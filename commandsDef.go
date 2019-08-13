package main

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

type Function struct {
	name        string
	description string
	def         func(arg []string, message *discordgo.MessageCreate, minArgsLen int, maxArgsLen int) (res string)
}

func NewFunction(name string, description string,
	def func(arg []string, message *discordgo.MessageCreate,
		minArgsLen int, maxArgsLen int) (response string)) *Function {
	return &Function{
		name,
		description,
		def,
	}
}

func help(arg []string, message *discordgo.MessageCreate, minArgsLen int, maxArgsLen int) string {
	return helpCommands()
}

func setPrefix(arg []string, message *discordgo.MessageCreate, minArgsLen int, maxArgsLen int) string {
	if len(arg) != 1 {
		return "Invalid number of arguments"
	}
	commandPrefix = arg[0]

	return fmt.Sprintf("Command prefix changed to %s", commandPrefix)
}

func setName(arg []string, message *discordgo.MessageCreate, minArgsLen int, maxArgsLen int) string {
	// TODO check to make sure arg[0] is valid and good and
	// has a nice cup of coofie and all that user input sanitization
	if len(arg) != 1 {
		return "Invalid number of arguments"
	}
	//logName(message.Author.ID, arg[0])
	return fmt.Sprintf("Set %s's name to be %s :thumbsup:", message.Author.ID, arg[0])
}

func getName(arg []string, message *discordgo.MessageCreate, minArgsLen int, maxArgsLen int) string {
	return getNameOr(message.Author.ID, message.Author.Username)
}

func givePoints(arg []string, message *discordgo.MessageCreate, minArgsLen int, maxArgsLen int) string {
	if len(arg) != 2 {
		return "Invalid number of arguments!"
	}
	recipient, ok := parseUserIDFromAt(arg[0])
	if !ok {
		return fmt.Sprintf("Recipient not defined, what is a \"%s\" :thinking:", arg[0])
	}
	amount, err := strconv.ParseInt(arg[1], 10, 64)
	if err != nil {
		return fmt.Sprintf("%s is not a number :thumbsdown:", arg[1])
	}
	//logTransaction(message.Author.ID, recipient, int(amount))
	return fmt.Sprintf("%s has given %d points to %s :thumbsup:",
		message.Author.ID, amount, recipient)
}

func helpCommands() string {
	var buffer bytes.Buffer

	for _, v := range funcMap {
		commandString := fmt.Sprintf("Command: %s\nHow to use: %s%s\n", v.name, commandPrefix, v.description)
		buffer.WriteString(commandString)
	}

	return buffer.String()
}
