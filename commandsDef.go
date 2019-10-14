package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Description struct {
	description    string
	argDescription []string
}

type ReponseActionType uint8

const (
	ResponsePM    ReponseActionType = 0
	ResponseReply ReponseActionType = 1
)

//Function struct handles the varying parts that go into the various bot commands
//Their nickname, their description, their actual definition etc
type Function struct {
	name        string
	description *Description
	def         func(arg []string, message *discordgo.MessageCreate) (res string, action ReponseActionType)
	minArgsLen  int
	maxArgsLen  int
}

//NewFunction -> Constructor for the Function struct
func NewFunction(name string, def func(arg []string, message *discordgo.MessageCreate) (response string, action ReponseActionType), minArgsLen int, maxArgsLen int, description *Description) *Function {
	if len(description.argDescription) != maxArgsLen {
		panic("not every arg is descripted in " + name)
	}
	return &Function{
		name,
		description,
		def,
		minArgsLen,
		maxArgsLen,
	}
}

func help(arg []string, message *discordgo.MessageCreate) (string, ReponseActionType) {
	var generateHelperArgString = func(amount int) string {
		re := ""
		for i := 0; i < amount; i++ {
			re += " <arg" + strconv.Itoa(i) + ">"
		}
		return re
	}

	var buffer bytes.Buffer

	buffer.WriteString("```")

	for _, v := range funcMap {
		buffer.WriteString(fmt.Sprintf("Command: %s\n\nDescription: %s\n", v.name+generateHelperArgString(v.maxArgsLen), v.description.description))
		for i, argDesc := range v.description.argDescription {
			isOptionalString := ""
			if i >= v.minArgsLen {
				isOptionalString = "*"
			}
			buffer.WriteString(fmt.Sprintf("Argument%s %s: %s\n", isOptionalString, strconv.Itoa(i), argDesc))
		}
		buffer.WriteString("\n\n\n")
	}
	buffer.WriteString(fmt.Sprintf("* argument is optional\n\n\nPrefix is currently: %s", *getGuildPrefix(message.GuildID)))

	buffer.WriteString("```")

	return buffer.String(), ResponsePM
}

func setPrefix(arg []string, message *discordgo.MessageCreate) (string, ReponseActionType) {
	user, err := discord.GuildMember(message.GuildID, message.Author.ID)
	logErr(err)
	for i := 0; i < len(user.Roles); i++ {
		role := strings.ToLower(user.Roles[i])
		if role == "admin" || role == "admins" || role == "mods" || role == "mod" {
			if arg[0] == "" || arg[0] == " " {
				return "invalid prefixes", ResponseReply
			}
			setPrefixForGuild(message.GuildID, arg[0])
			return fmt.Sprintf("Command prefix changed to %s", arg[0]), ResponseReply
		}
	}
	return "you dont have the rights to do this", ResponseReply
}

func getNick(arg []string, message *discordgo.MessageCreate) (string, ReponseActionType) {
	nick := getUser(message.Author.ID).nickname
	if nick.Valid {
		return nick.String, ResponseReply
	}
	return "your nickname is not set", ResponseReply
}

func setNick(arg []string, message *discordgo.MessageCreate) (string, ReponseActionType) {
	// TODO check to make sure arg[0] is valid and good and
	// has a nice cup of coofie and all that user input sanitization
	//logName(message.Author.ID, arg[0])
	setUsersNickname(&User{
		discordId: message.Author.ID,
		nickname: sql.NullString{
			String: arg[0],
			Valid:  true,
		},
	})
	return fmt.Sprintf("Set your nickname to be %s :thumbsup:", arg[0]), ResponseReply
}

func clearNick(arg []string, message *discordgo.MessageCreate) (string, ReponseActionType) {
	setUsersNickname(&User{
		discordId: message.Author.ID,
		nickname: sql.NullString{
			String: "",
			Valid:  false,
		},
	})
	return "your nickname was cleared", ResponseReply
}

func pointsCommand(arg []string, message *discordgo.MessageCreate) (string, ReponseActionType) {
	if len(arg) == 0 {
		points, nicknames := getUsersPointsGiven(message.Author.ID)
		re := "You have given:\n"
		for i := 0; i < len(points); i++ {
			var nick string
			if nicknames[i].Valid {
				nick = nicknames[i].String
			} else {
				member, err := discord.GuildMember(message.GuildID, points[i].receiver)
				if err != nil {
					continue
				}
				if member.Nick == "" {
					nick = member.User.Username
				} else {
					nick = member.Nick
				}
			}
			re += "\t" + nick + ", " + strconv.FormatInt(points[i].amount, 10) + " of your points\n"
		}
		return re, ResponseReply
	} else if len(arg) == 1 {
		points, _ := getUsersPointsReceived(message.Author.ID)
		re := getPrintableName(message.Author.ID, message.GuildID) + " has got\n"
		for i := 0; i < len(points); i++ {
			re += "\t" + strconv.FormatInt(points[i].amount, 10) + " " + getPrintableName(points[i].giver, message.GuildID) + " points\n"
		}
		return re, ResponseReply
	} else {
		return "Too many arguments", ResponseReply
	}
}

func givePoints(arg []string, message *discordgo.MessageCreate) (string, ReponseActionType) {
	recipientID, ok := parseUserIDFromAt(arg[0], message.GuildID)
	if !ok {
		return fmt.Sprintf("Recipient not defined, what is a \"%s\" :thinking:", arg[0]), ResponseReply
	}
	//recipient, err := discord.GuildMember(message.GuildID, recipientID)

	amount, err := strconv.ParseInt(arg[1], 10, 64)
	if amount > 100 {
		return "Limit is 100 points per command", ResponseReply
	} else if amount < -100 {
		return "Limit is -100 points per command", ResponseReply
	}
	if err != nil {
		return fmt.Sprintf("%s is not a number :thumbsdown:", arg[1]), ResponseReply
	}
	giveUserPoints(recipientID, message.Author.ID, amount)

	return getPrintableName(message.Author.ID, message.GuildID) + " gave " + arg[1] + " points to " + getPrintableName(recipientID, message.GuildID), ResponseReply
}
