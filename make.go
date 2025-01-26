//go:build exclude

package main

import (
	"log"
	"os"
	"os/exec"
)

const (
	buildDirectory  = "./build"
	binaryPath      = "./build/neba.exe"
	mainPackagePath = "."
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
		log.Fatalf(taskName)
	}

	task()
}

func executeCommand(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), environmentalVariables...)

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func clean() {
	if err := os.RemoveAll(buildDirectory); err != nil {
		return
	}
}

func tidy() {
	executeCommand("go", "mod", "tidy", "-v")
	executeCommand("go", "fmt", "./...")
}

func run() {
	clean()
	tidy()
	executeCommand("go", "build", "-o", binaryPath, mainPackagePath)
	executeCommand(binaryPath)
}

func production() {
	clean()
	tidy()
	// TODO: cross-platform build with ldflags
	// TODO: auto versioning
}
