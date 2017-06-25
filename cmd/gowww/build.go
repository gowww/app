// gowww is the CLI of the gowww/app framework.
package main

import (
	"log"
	"os"
	"os/exec"
	"runtime"
)

func buildName() string {
	if *cmdBuildFlagDocker {
		return *cmdFlagName + "_linux_amd64"
	}
	return *cmdFlagName + "_" + runtime.GOOS + "_" + runtime.GOARCH
}

func build() error {
	log.Println("Building...")
	cmd := exec.Command("go", "build", "-o", buildName())
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
	cmd := exec.Command("docker", "run", "--rm", "-v", getwd(true)+":/go/src/"+*cmdFlagName, "-w", "/go/src/"+*cmdFlagName, "golang:latest", "sh", "-c", "go get . && go build -o "+buildName())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
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
