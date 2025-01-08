package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	binaryName      = "neba"
	buildDirectory  = "./build"
	mainPackagePath = "./cmd/neba"
)

var (
	binaryPath             = filepath.Join(buildDirectory, binaryName)
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

	taskName := "production"
	if len(os.Args) > 1 {
		taskName = os.Args[1]
	}

	task, exists := tasks[taskName]
	if !exists {
		log.Fatalf("unknown task: %s", taskName)
	}

	task()
}

func executeCommand(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), environmentalVariables...)

	if err := cmd.Run(); err != nil {
		log.Fatalf("command: %s %s", name, args)
		log.Fatalf("command %s failed: %v", name, err)
	}
}

func ensureBuildDirectory() {
	os.MkdirAll(buildDirectory, os.ModePerm)
}

func clean() {
	fmt.Println("cleaning build directory...")
	os.RemoveAll(buildDirectory)
}

func tidy() {
	fmt.Println("tidying and formatting code...")
	executeCommand("go", "mod", "tidy", "-v")
	executeCommand("go", "fmt", "./...")
}

func run() {
	fmt.Println("running application...")
	ensureBuildDirectory()
	executeCommand("go", "build", "-v", "-o", binaryPath, mainPackagePath)
	executeCommand(binaryPath)
}

func production() {
	fmt.Println("building application for production...")
	clean()
	tidy()
	ensureBuildDirectory()
	executeCommand("go", "build", "-o", binaryPath, mainPackagePath)
}
