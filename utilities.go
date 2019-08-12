package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
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

func getToken() string {
	file, err := os.Open("TOKEN")
	logErr(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	re := ""

	for scanner.Scan() {
		re += scanner.Text()
	}

	return re
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
