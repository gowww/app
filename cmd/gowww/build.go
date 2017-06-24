// gowww is the CLI of the gowww/app framework.
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func getBuildName() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return filepath.Base(wd)
}

func build() error {
	log.Println("Building...")
	cmd := exec.Command("go", "build", "-o", fmt.Sprintf("%s_%s_%s", *buildName, runtime.GOOS, runtime.GOARCH))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err == nil {
		cleanLines(1)
	}
	return err
}

func buildDocker() error {
	log.Println("Building with Docker...")
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	buildName := *buildName + "_linux_amd64"
	cmd := exec.Command("docker", "run", "--rm", "-v", wd+":/go/src/"+buildName, "-w", "/go/src/"+buildName, "golang:latest", "sh", "-c", "go get . && go build -o "+buildName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err == nil {
		cleanLines(1)
	}
	return err
}

func buildScriptsGopherJS() error {
	log.Println("Building scripts with GopherJS...")
	cmd := exec.Command("gopherjs", "build", "./scripts", "-o", "static/main.js", "-m")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err == nil {
		cleanLines(1)
	}
	return err
}

// TODO: Build scripts (Babel, TypeScript...) and styles (LESS, SCSS, Stylus...) before and during watching.
