package main

import (
	"database/sql"
	"fmt"

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
	funcMap = make(map[string]*Function)
	//DB varaible to handle sql connection
	DB      *sql.DB
	discord *discordgo.Session
	ready   = false
)

//Some core basics to get going
func main() {
	localDiscord, err := discordgo.New("Bot " + SECRET.DISCORD_TOKEN)
	discord = localDiscord
	errCheck("Error creating discord session", err)
	openDBConnection()
	errCheck("Error establishing database session", err)
	defer DB.Close()

	// Prints out a list of all the currently available bot commands,
	// descriptions thereof, and an example of how to use them.
	// help
	funcName := "help"
	funcMap[funcName] = NewFunction(funcName, help, 0, 0, &Description{
		description:    "Returns the list of commands and their descriptions",
		argDescription: []string{},
	})

	// Sets the prefix the bot responds to when users issue commands.
	// prefix <arg0>
	// arg0: prefix to use (Char)
	funcName = "setPrefix"
	funcMap[funcName] = NewFunction(funcName, setPrefix, 1, 1, &Description{
		description:    "Updates the prefix that the bot uses to identify commands",
		argDescription: []string{"Your new prefix"},
	})

	// Sets the user's nickname to the provided string
	// setnick <arg0>
	// arg0: provided nickname to switch to (String)
	funcName = "setNick"
	funcMap[funcName] = NewFunction(funcName, setNick, 1, 1, &Description{
		description:    "Sets your own nickname, which the bot uses when printing how many points people have",
		argDescription: []string{"Your new nickname"},
	})

	// Returns the currently stored nickname of the user who issued the command
	funcName = "getNick"
	funcMap[funcName] = NewFunction(funcName, getNick, 0, 0, &Description{
		description:    "Returns the nickname this bot uses to refer to you",
		argDescription: []string{},
	})

	// Clears the nickname the bot currently has stored associated
	// to the user who issued the command
	funcName = "resetNick"
	funcMap[funcName] = NewFunction(funcName, clearNick, 0, 0, &Description{
		description:    "Resets your nickname for the bot",
		argDescription: []string{},
	})

	// Gives a target user some number of points from the issuing user
	// user1: givePoints @user2 20
	// -> Gives user2 20 points of type user1
	funcName = "givePoints"
	funcMap[funcName] = NewFunction(funcName, givePoints, 2, 2, &Description{
		description:    "Gives a user an amount of your points",
		argDescription: []string{"The user in question", "The amount of points (must be an integer)"},
	})

	// Prints out of a list of either your own, or a target user's points
	// and who they were issued by
	// user1: points -> prints out user1's point tallies
	// user1: points user2 -> prints out user2's point tallies
	funcName = "points"
	funcMap[funcName] = NewFunction(funcName, pointsCommand, 0, 1, &Description{
		description:    "Prints the amount of points",
		argDescription: []string{"Returns the amount of points of the pinged user, if this is not set, it will return all points you have given"},
	})

	// Returns the discord tag of the supplied target user
	// whois user -> returns usertag#xyzw for user
	funcName = "whoIs"
	funcMap[funcName] = NewFunction(funcName, whoIsCommand, 1, 1, &Description{
		description:    "Describes who a person is",
		argDescription: []string{"the person you're interesting in knowing "},
	})

	discord.AddHandler(func(discord *discordgo.Session, message *discordgo.MessageCreate) {
		user := message.Author
		if user.Bot {
			return
		}
		if !ready {
			_, err := discord.ChannelMessageSend(message.ChannelID, "sorry, im not ready yet, database sync havnt been done yet")
			errCheck("An error happened when telling users we arent ready yet", err)
		}
		prefix := getGuildPrefix(message.GuildID)
		if prefix == nil {
			fmt.Println("the frick, this shouldnt be nil")
			return
		}

		if len(message.Content) > len(*prefix) && message.Content[0:len(*prefix)] == *prefix {
			funcName, args := commandLineArgSplit(message.Content[1:])
			fun, ok := funcMap[funcName]
			if ok {
				if fun.minArgsLen > len(args) {
					_, err := discord.ChannelMessageSend(message.ChannelID, "too few arguments")
					errCheck("Oepsie woepsie, er was een stukkiewukkie in 't command handler", err)
					return
				}
				if fun.maxArgsLen < len(args) {
					_, err := discord.ChannelMessageSend(message.ChannelID, "too many arguments, check the amount of spaces")
					errCheck("Oepsie woepsie, er was een stukkiewukkie in 't command handler", err)
					return
				}
				res, action := fun.def(args, message)
				switch action {
				case ResponseReply:
					_, err := discord.ChannelMessageSend(message.ChannelID, res)
					errCheck("Oepsie woepsie, er was een stukkiewukkie in 't command handler", err)
					break
				case ResponsePM:
					channel, err := discord.UserChannelCreate(message.Author.ID)
					errCheck("Oepsie woepsie, er was een stukkiewukkie in 't command handler", err)
					_, err = discord.ChannelMessageSend(channel.ID, res)
					errCheck("Oepsie woepsie, er was een stukkiewukkie in 't command handler", err)
					break
				}

			} else {
				discord.ChannelMessageSend(message.ChannelID, "Invalid command")
			}
		}
	})
	discord.AddHandler(func(discord *discordgo.Session, ready *discordgo.Ready) {
		err = discord.UpdateStatus(0, "Gender non-conformity")
		errCheck("Error attempting to set status", err)
		var guilds []string
		for _, s := range ready.Guilds {
			guilds = append(guilds, s.ID)
		}
		startupAddAllGuilds(guilds)
	})
	discord.AddHandler(func(s *discordgo.Session, event *discordgo.GuildCreate) {
		if event.Guild.Unavailable {
			fmt.Printf("\nHi, the bot was here")
			return
		}
	})
	discord.AddHandler(func(s *discordgo.Session, chunk *discordgo.GuildMembersChunk) {
		var users []string
		for i := 0; i < len(chunk.Members); i++ {
			users = append(users, chunk.Members[i].User.ID)
		}
		startupAddAllUsers(users)
	})
	waitForMemberFetch(discord, func(discord *discordgo.Session) {
		ready = true
	})

	err = discord.Open()
	defer discord.Close()

	<-make(chan struct{})
}
