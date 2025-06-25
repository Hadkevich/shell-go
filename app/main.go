package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"strings"
)

const (
	execMod = 0666
)

var (
	STDOUT           = os.Stdout
	STDERR           = os.Stderr
	shellCommands    = []string{"type", "echo", "exit", "pwd", "cd"}
	stdOutCmds       = []string{">", "1>", "2>"}
	stdAppendOutCmds = []string{">>", "1>>"}
)

func main() {
	for {
		fmt.Print("$ ")
		STDOUT = os.Stdout
		line, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		parts := splitByQuotes(strings.TrimRight(line, "\n"))
		command := parts[0]
		args := parts[1:]

		if len(args) >= 3 {
			redirPos := len(args) - 2
			redirFile := args[len(args)-1]

			if slices.Contains(stdOutCmds, args[redirPos]) {
				openFlag := os.O_CREATE | os.O_WRONLY | os.O_TRUNC // learn about all in this line
				switch args[redirPos] {
				case "1>", ">":
					STDOUT, _ = os.OpenFile(redirFile, openFlag, execMod)
				case "2>":
					STDERR, _ = os.OpenFile(redirFile, openFlag, execMod)
				}
				args = args[:redirPos]
			} else if slices.Contains(stdAppendOutCmds, args[redirPos]) {
				openFlag := os.O_CREATE | os.O_WRONLY | os.O_APPEND // Use append flag
				switch args[redirPos] {
				case "1>>", ">>":
					STDOUT, _ = os.OpenFile(redirFile, openFlag, execMod)
				case "2>>":
					STDERR, _ = os.OpenFile(redirFile, openFlag, execMod)
				}
				args = args[:redirPos]
			}
		}

		switch command {
		case "exit":
			exitCommand(args)
		case "echo":
			echoCommand(args, STDOUT)
		case "type":
			typeCommand(args)
		case "pwd":
			pwdCommand()
		case "cd":
			cdCommand(args)
		default:
			filePath, exists := findExecutable(command)

			if exists && filePath != "" {
				cmd := exec.Command(command, args...)
				cmd.Stdout = STDOUT
				cmd.Stderr = STDERR

				cmd.Run()
			} else {
				fmt.Println(command + ": command not found")
			}
		}
	}
}

func exitCommand(args []string) {
	exitCode, err := strconv.Atoi(args[0])

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error arg:", err)
		os.Exit(1)
	}

	os.Exit(exitCode)
}

func echoCommand(args []string, output *os.File) {
	str := strings.Join(args, " ")
	if output != nil {
		output.WriteString(str + "\n")
	} else {
		fmt.Println(str)
	}
}

func typeCommand(args []string) {
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
	return slices.Contains(shellCommands, command)
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
