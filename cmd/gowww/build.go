// gowww is the CLI of the gowww/app framework.
package main

import (
	"log"
	"os"
	"os/exec"
	"runtime"
)

func buildName() string {
	if *flagBuildDocker {
		return *flagName + "_linux_amd64"
	}
	goos, ok := os.LookupEnv("GOOS")
	if !ok {
		goos = runtime.GOOS
	}
	goarch, ok := os.LookupEnv("GOARCH")
	if !ok {
		goarch = runtime.GOARCH
	}
	return *flagName + "_" + goos + "_" + goarch
}

func build(info string, cmdName string, cmdArgs ...string) error {
	log.Println(info)
	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err == nil {
		cleanLines(1)
	}
	return err
}

func buildGo() error {
	return build("Building...",
		"go", "build", "-o", buildName())
}

func buildDocker() error {
	return build("Building with Docker...",
		"docker", "run", "--rm", "-v", getwd(true)+":/go/src/"+*flagName, "-w", "/go/src/"+*flagName, "golang:latest", "sh", "-c", "go get . && go build -o "+buildName())
}

func buildScriptsGopherJS() error {
	return build("Building scripts with GopherJS...",
		"gopherjs", "build", "./scripts", "--output", "static/scripts/main.js", "--minify")
}

func buildStylesStylus(file string) error {
	return build("Building styles with Stylus...",
		"stylus", file, "--out", "static/styles", "--compress", "--sourcemap")
}
