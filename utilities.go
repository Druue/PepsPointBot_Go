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
