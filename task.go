package main

import (
	"log/slog"
	"os"
	"os/exec"
)

const (
	binaryName      = "neba"
	buildDirectory  = "./build"
	binaryPath      = "./build/neba.exe"
	mainPackagePath = "./cmd/neba"
)

var (
	environmentalVariables = []string{
		"CGO_ENABLED=0",
		"GOOS=windows",
		"GOARCH=amd64",
	}
)

func main() {
	tasks := map[string]func(){
		"clean":      clean,
		"tidy":       tidy,
		"run":        run,
		"production": production,
	}

	taskName := "run"
	if len(os.Args) > 1 {
		taskName = os.Args[1]
	}

	task, exists := tasks[taskName]
	if !exists {
		slog.Error("specified task is not defined", "task", taskName)
		os.Exit(1)
	}

	task()
}

func executeCommand(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), environmentalVariables...)

	if err := cmd.Run(); err != nil {
		slog.Error("command exited with error",
			slog.String("name", name),
			slog.Any("args", args),
			slog.Any("error", err))
	}
}

func clean() {
	os.RemoveAll(buildDirectory)
}

func tidy() {
	executeCommand("go", "mod", "tidy", "-v")
	executeCommand("go", "fmt", "./...")
}

func run() {
	executeCommand("go", "build", "-o", binaryPath, mainPackagePath)
	executeCommand(binaryPath)
}

func production() {
	clean()
	tidy()
	// TODO: cross-platform build with ldflags
	// TODO: auto versioning
}
