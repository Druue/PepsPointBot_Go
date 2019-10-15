package main

import (
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

func noRows(err error) {
	if err == sql.ErrNoRows {
		fmt.Println("No rows to return!")
	}
}

func logErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func errCheck(msg string, err error) {
	if err != nil {
		fmt.Printf("%s: %+v", msg, err)
		panic(err)
	}
}

func parseUserIDFromAt(user string, guildID string) (string, bool) {
	if !(len(user) > 3 && user[0:2] == "<@" && user[len(user)-1:] == ">") {
		dbUser := getUserFromNickname(user)
		if dbUser.Valid {
			return dbUser.String, true
		}
		guild, err := discord.Guild(guildID)
		errCheck("error \"in parse user id from at\", when finding guilds", err)
		for _, m := range guild.Members {
			if m.Nick == user || m.User.Username == user {
				return m.User.ID, true
			}
		}
		return "", false
	}
	first := user[2 : len(user)-1]
	if first[0:1] == "!" {
		return first[1:], true
	}
	return first, true
}

func getPrintableName(discordId string, guildId string) string {
	dbUser := getUser(discordId)
	if dbUser != nil {
		dbNick := dbUser.nickname
		if dbNick.Valid {
			return dbNick.String
		}
	}
	guildMember, err := discord.GuildMember(guildId, discordId)
	if guildMember == nil {
		usr, err := discord.User(discordId)
		errCheck("user doesnt exist or something", err)
		return usr.Username
	}
	errCheck("getting the guild member in getting printable name failed", err)
	if guildMember.Nick != "" {
		return guildMember.Nick
	}
	return guildMember.User.Username
}

func waitForMemberFetch(discord *discordgo.Session, cb func(discord *discordgo.Session)) {
	discord.AddHandler(func(discord *discordgo.Session, ready *discordgo.Ready) {
		for i := 0; i < len(discord.State.Guilds); i++ {
			err := discord.RequestGuildMembers(discord.State.Guilds[i].ID, "", 2147483647)
			errCheck("request guild members failed for guild "+discord.State.Guilds[i].Name, err)
		}
	})
	discord.AddHandler(func(discord *discordgo.Session, chunk *discordgo.GuildMembersChunk) {
		done := true
		for i := 0; i < len(discord.State.Guilds); i++ {
			done = done || len(discord.State.Guilds[i].Members) == discord.State.Guilds[i].MemberCount
		}
		if done {
			cb(discord)
		}
	})
}

func commandLineArgSplit(str string) (string, []string) {
	commandNameIndex := strings.Index(str, " ")
	if commandNameIndex == -1 {
		return str, []string{}
	}
	commandName := str[:commandNameIndex]
	argsStr := str[commandNameIndex+1:]
	prevChar := ""
	strBuilder := ""
	ignoreSpace := false
	var args []string
	for i := 0; i < len(argsStr); i++ {
		fmt.Println()
		char := string(argsStr[i])
		if prevChar == "\\" {
			strBuilder += char
			prevChar = char
			continue
		}
		if char == "\\" {
			prevChar = char
			continue
		}
		if char == "\"" {
			ignoreSpace = !ignoreSpace
			prevChar = char
			continue
		}
		if char == " " && !ignoreSpace {
			if prevChar != " " {
				args = append(args, strBuilder)
				strBuilder = ""
			}
			prevChar = char
			continue
		}
		strBuilder += char
		prevChar = char
	}
	args = append(args, strBuilder)
	return commandName, args
}
