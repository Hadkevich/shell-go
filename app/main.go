package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")
		line, err := bufio.NewReader(os.Stdin).ReadString('\n')

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		line = strings.TrimSpace(line)
		parts := strings.Split(line, " ")
		command := parts[0]
		args := parts[1:]

		switch command {
		case exit.String():
			ExitCommand(args)
		case echo.String():
			EchoCommand(args)
		case type_.String():
			TypeCommand(args)
		default:
			if filePath, exists := findExecutable(command); exists {
				cmd := exec.Command(command, args...)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err := cmd.Run()
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s: command not found\n", cmd, filePath)
				}
			}
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

	if isShellBuiltin(value) {
		fmt.Println(value + " is a shell builtin")
		return
	}

	if file, exists := findExecutable(value); exists {
		fmt.Fprintf(os.Stdout, "%s is %s\n", value, file)
		return
	}

	fmt.Println(value + ": not found")
}

func isShellBuiltin(command string) bool {
	builtIns := []string{"echo", "exit", "type"}
	for _, b := range builtIns {
		if b == command {
			return true
		}
	}
	return false
}

func findExecutable(bin string) (string, bool) {
	paths := os.Getenv("PATH")
	for _, path := range strings.Split(paths, ":") {
		file := path + "/" + bin
		if _, err := os.Stat(file); err == nil {
			return file, true
		}
	}
	return "", false
}
