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

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		commands := strings.Split(userInput[:len(userInput)-1], " ")
		command := commands[0]
		args := commands[1:]

		switch command {
		case exit.String():
			handleExit(args)
		case echo.String():
			handleEcho(args)
		case type_.String():
			handleType(args)
		default:
			fmt.Println(command + ": command not found")
		}
	}
}

type builtin int

const (
	exit builtin = iota
	echo
	type_
)

var commandName = map[builtin]string{
	exit:  "exit",
	echo:  "echo",
	type_: "type",
}

func (ss builtin) String() string {
	return commandName[ss]
}

func handleExit(args []string) {
	exitCode, err := strconv.Atoi(args[0])

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error arg:", err)
		os.Exit(1)
	}

	os.Exit(exitCode)
}

func handleEcho(args []string) {
	fmt.Println(strings.Join(args, " "))
}

func handleType(args []string) {
	switch args[0] {
	case exit.String():
		fmt.Println(args[0] + " is a shell builtin")
	case echo.String():
		fmt.Println(args[0] + " is a shell builtin")
	case type_.String():
		fmt.Println(args[0] + " is a shell builtin")
	default:
		fmt.Println(strings.Join(args, " ") + ": not found")
	}
}
