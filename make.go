//go:build exclude

package main

import (
	"log"
	"os"
	"os/exec"
)

const (
	buildDir    = "./build"
	binPath     = "./build/neba.exe"
	mainPkgPath = "."
)

var (
	env = []string{
		"CGO_ENABLED=0",
		"GOOS=windows",
		"GOARCH=amd64",
	}
)

func main() {
	tasks := map[string]func(){
		"clean": clean, "tidy": tidy, "run": run, "production": prod,
	}
	taskName := "run"
	if len(os.Args) > 1 {
		taskName = os.Args[1]
	}
	tasks[taskName]()
}

func do(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), env...)

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func clean() {
	if err := os.RemoveAll(buildDir); err != nil {
		return
	}
}

func tidy() {
	do("go", "mod", "tidy", "-v")
	do("go", "vet", "./...")
	do("go", "fmt", "./...")
}

func run() {
	clean()
	do("go", "build", "-o", binPath, mainPkgPath)
	do(binPath)
}

func prod() {
	clean()
	tidy()
	// TODO: cross-platform build with ldflags
	// TODO: auto versioning
}
