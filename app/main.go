package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")
		userInput, err := bufio.NewReader(os.Stdin).ReadString('\n')

		commands := strings.Split(userInput[:len(userInput)-1], " ")
		command := commands[0]
		args := commands[1:]

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		switch {
		case command == "exit":
			handleExit(args)
		case command == "echo":
			fmt.Println(strings.Join(args, " "))
		default:
			fmt.Println(command + ": command not found")
		}
	}
}

func handleExit(args []string) {
	exitCode, err := strconv.Atoi(args[0])

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error arg:", err)
		os.Exit(1)
	}

	os.Exit(exitCode)
}
