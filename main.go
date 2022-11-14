package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Command struct {
	Name    string
	Handler func(flags *flag.FlagSet)
	flags   *flag.FlagSet
}

type CommandRegistry struct {
	commands       map[string]Command
	defaultCommand Command
}

func (reg *CommandRegistry) addCommand(cmd Command) {
	reg.commands[cmd.Name] = cmd
}

func AbsolutePath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	return exPath
}

func LitPath() string {
	return AbsolutePath() + "\\lit\\"
}

func WorkingDirectory() string {
	return LitPath() + "working\\"
}

func AddHandler(flags *flag.FlagSet) {
	fmt.Println("Adding:")
	filePaths := flags.Args()
	files := make([]string, 0)

	for _, file := range filePaths {
		newFiles, _ := filepath.Glob(file)
		files = append(files, newFiles...)

	}
	os.MkdirAll(WorkingDirectory(), os.ModePerm)
	for _, file := range files {
		src, _ := os.Open(file)
		dst, _ := os.Create(WorkingDirectory() + file)
		_, err := io.Copy(dst, src)
		if err == nil {
			fmt.Println("Added: " + file)
		} else {
			fmt.Println("Error Adding: " + file)
			fmt.Println(err)
		}
	}

}

func CreateHelpCommand() Command {
	return Command{
		"help",
		func(flags *flag.FlagSet) {},
		flag.NewFlagSet("help", flag.ExitOnError),
	}
}

func CreateAddCommand() Command {
	flags := flag.NewFlagSet("add", flag.ExitOnError)
	return Command{
		Name:    "add",
		Handler: AddHandler,
		flags:   flags,
	}
}

func NewCommandRegistry() CommandRegistry {
	reg := CommandRegistry{defaultCommand: CreateHelpCommand(), commands: make(map[string]Command)}
	reg.addCommand(CreateAddCommand())
	return reg
}

func (reg *CommandRegistry) invoke(command string, args []string) {
	cmd, found := reg.commands[command]
	if found {
		cmd.flags.Parse(args)
		cmd.Handler(cmd.flags)
	} else {
		reg.defaultCommand.flags.Parse(args)
		reg.defaultCommand.Handler(cmd.flags)
	}
}

func main() {
	flag.Parse()
	commands := flag.Args()

	reg := NewCommandRegistry()
	if len(commands) > 0 {
		reg.invoke(commands[0], commands[1:])
	} else {
		reg.invoke("", []string{})
	}
}
