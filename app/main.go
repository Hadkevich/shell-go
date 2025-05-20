package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
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
			ExitCommand(args)
		case echo.String():
			EchoCommand(args)
		case type_.String():
			TypeCommand(args)
		default:
			fmt.Println(command + ": command not found")
		}
	}
}

type Command int

const (
	exit Command = iota
	echo
	type_
)

var commandName = map[Command]string{
	exit:  "exit",
	echo:  "echo",
	type_: "type",
}

var builtIns = []string{"echo", "exit", "type"}

func (ss Command) String() string {
	return commandName[ss]
}

func ExitCommand(args []string) {
	exitCode, err := strconv.Atoi(args[0])

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error arg:", err)
		os.Exit(1)
	}

	os.Exit(exitCode)
}

func EchoCommand(args []string) {
	fmt.Println(strings.Join(args, " "))
}

func TypeCommand(args []string) {
	value := args[0]

	if slices.Contains(builtIns, value) {
		fmt.Println(value + " is a shell builtin")
		return
	}

	if file, exists := findBinInPath(value); exists {
		fmt.Fprintf(os.Stdout, "%s is %s\n", value, file)
		return
	}

	fmt.Println(value + ": not found")
}

func findBinInPath(bin string) (string, bool) {
	paths := os.Getenv("PATH")
	for _, path := range strings.Split(paths, ":") {
		file := path + "/" + bin
		if _, err := os.Stat(file); err == nil {
			return file, true
		}
	}
	return "", false
}
