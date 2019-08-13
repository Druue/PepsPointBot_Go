package main

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

type Function struct {
	name           string
	description    string
	argDescription []string
	def            func(arg []string, message *discordgo.MessageCreate) (res string)
	minArgsLen     int
	maxArgsLen     int
}

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
	return fmt.Sprintf("Set %s's name to be %s :thumbsup:", message.Author.ID, arg[0])
}

func getName(arg []string, message *discordgo.MessageCreate) string {
	return getNameOr(message.Author.ID, message.Author.Username)
}

func givePoints(arg []string, message *discordgo.MessageCreate) string {
	recipientID, ok := parseUserIDFromAt(arg[0])
	if !ok {
		return fmt.Sprintf("Recipient not defined, what is a \"%s\" :thinking:", arg[0])
	}
	fmt.Println(discord)
	recipient, err := discord.GuildMember(message.GuildID, recipientID)
	if err != nil {
		panic(err)
	}
	amount, err := strconv.ParseInt(arg[1], 10, 64)
	if err != nil {
		return fmt.Sprintf("%s is not a number :thumbsdown:", arg[1])
	}
	//logTransaction(message.Author.ID, recipient, int(amount))
	var recipientBackupName string
	if recipient.Nick == "" {
		recipientBackupName = recipient.User.Username
	} else {
		recipientBackupName = recipient.Nick
	}
	sender, err := discord.GuildMember(message.GuildID, message.Author.ID)
	if err != nil {
		return fmt.Sprintf("%s is not a number :thumbsdown:", arg[1])
	}
	var senderBackupName string
	if sender.Nick == "" {
		senderBackupName = sender.User.Username
	} else {
		senderBackupName = sender.Nick
	}
	return fmt.Sprintf("%s has given %d points to %s :thumbsup:",
		getNameOr(sender.User.ID, senderBackupName), amount, getNameOr(recipientID, recipientBackupName))
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
