package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	for {
		fmt.Print("$ ")
		line, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		parts := splitByQuotes(strings.TrimRight(line, "\n"))
		command := parts[0]
		args := parts[1:]
		var output *os.File

		for i, arg := range args {
			if (arg == ">" || arg == "1>") && i < len(args)-1 {
				output, err = os.Create(args[i+1])
				check(err)
				args = args[:i]
				break
			}
		}

		switch command {
		case exit.String():
			ExitCommand(args)
		case echo.String():
			EchoCommand(args, output)
		case type_.String():
			TypeCommand(args)
		case pwd.String():
			pwdCommand()
		case cd.String():
			cdCommand(args)
		default:
			filePath, exists := findExecutable(command)

			if exists && filePath != "" {
				cmd := exec.Command(command, args...)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				if output != nil {
					defer output.Close()
					cmd.Stdout = output
				}

				cmd.Run()

			} else {
				fmt.Println(command + ": command not found")
			}
		}
	}
}

type Command int

const (
	exit Command = iota
	echo
	type_
	pwd
	cd
)

var commandName = map[Command]string{
	exit:  "exit",
	echo:  "echo",
	type_: "type",
	pwd:   "pwd",
	cd:    "cd",
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

func EchoCommand(args []string, output *os.File) {
	str := strings.Join(args, " ")
	if output != nil {
		defer output.Close()
		output.WriteString(str + "\n")
	} else {
		fmt.Println(str)
	}
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
	builtIns := []string{"echo", "exit", "type", "pwd"}
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

func pwdCommand() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error arg:", err)
	}
	fmt.Println(dir)
}

func cdCommand(args []string) {
	path := strings.Join(args, " ")

	if len(path) == 0 {
		pwdCommand()
		return
	}

	if strings.Compare(path, "~") == 0 {
		homeDir, error := os.UserHomeDir()
		if error == nil {
			os.Chdir(homeDir)
		}
		return
	}

	stat, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cd: %s: No such file or directory\n", path)
		return
	}

	if stat.IsDir() {
		os.Chdir(path)
	}
}

func splitByQuotes(s string) []string {
	var result []string
	var current string
	inQuote, inDoubleQuote, escapeNext := false, false, false

	for i := 0; i < len(s); i++ {
		if escapeNext && inDoubleQuote {
			if s[i] == '$' || s[i] == '"' || s[i] == '\\' {
				current += string(s[i])
			} else {
				current += "\\"
				current += string(s[i])
			}
			escapeNext = !escapeNext
		} else if escapeNext {
			current += string(s[i])
			escapeNext = !escapeNext
		} else if s[i] == '\'' && !inDoubleQuote {
			inQuote = !inQuote
		} else if s[i] == '\\' && !inQuote {
			escapeNext = !escapeNext
		} else if s[i] == '"' && !inQuote {
			inDoubleQuote = !inDoubleQuote
		} else if s[i] == ' ' && !inQuote && !inDoubleQuote {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(s[i])
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}
