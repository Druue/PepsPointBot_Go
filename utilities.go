package main

import (
	"database/sql"
	"fmt"
	"log"
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

func parseUserIDFromAt(user string) (string, bool) {
	if len(user) < 2 {
		return "", false
	}
	if !(user[0:2] == "<@" && user[len(user)-1:] == ">") {
		return "", false
	}
	return user[2 : len(user)-1], true
}

func getDescription(funcName string) (string, []string) {
	var slice []string
	switch funcName {
	case "help":
		slice = []string{}
		return "Returns the list of commands and their descriptions", slice
	case "set-prefix":
		slice = []string{"the new prefix"}
		return "Updates the prefix that the bot uses to identify commands", slice
	case "set-nickname":
		slice = []string{"your new nickname"}
		return "Sets your own nickname, which the bot uses when printing how many points people have", slice
	case "get-nickname":
		slice = []string{}
		return "return the nickname this bot uses to refeer to you", slice
	case "give":
		slice = []string{"the user in question", "the amount of points (must be a string)"}
		return "Gives a user an amount of your points", slice
	default:
		panic(fmt.Sprintf("%v", "No such function exists!"))
	}
}
