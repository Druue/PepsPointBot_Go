package main

import (
	"bufio"
	"log"
	"os"
)

func getToken() string {
	file, err := os.Open("TOKEN")
	if err != nil {
		log.Fatal(err)
	}
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
