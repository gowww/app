// gowww is the CLI of the gowww/app framework.
package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gowww/cli"
)

func build() {
	if flagBuildDocker {
		buildDocker()
	} else {
		buildGo()
	}
}

func buildName() string {
	if flagBuildDocker {
		return flagBuildName + "_linux_amd64"
	}
	goos, ok := os.LookupEnv("GOOS")
	if !ok {
		goos = runtime.GOOS
	}
	goarch, ok := os.LookupEnv("GOARCH")
	if !ok {
		goarch = runtime.GOARCH
	}
	return flagBuildName + "_" + goos + "_" + goarch
}

func buildExec(info string, cmdName string, cmdArgs ...string) error {
	log.Println(info)
	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err == nil {
		cli.CleanLines(1)
	}
	return err
}

func buildGo() error {
	return buildExec("Building...",
		"go", "build", "-o", buildName())
}

func buildDocker() error {
	return buildExec("Building with Docker...",
		"docker", "run", "--rm", "-v", getwd(true)+":/go/src/"+flagBuildName, "-w", "/go/src/"+flagBuildName, "golang:latest", "sh", "-c", "go get . && go build -o "+buildName())
}

func buildScriptsGopherJS(file string) error {
	return buildExec("Building scripts with GopherJS...",
		"gopherjs", "build", "./"+filepath.Dir(file), "--output", filepath.Join("static", filepath.Dir(file), "main.js"), "--minify")
}

func buildStylesSass(file string) error {
	outFile := filepath.Join("static", strings.TrimSuffix(file, filepath.Ext(file))+".css")
	os.MkdirAll(filepath.Dir(outFile), os.ModePerm)
	return buildExec("Building styles with Sass...",
		"sassc", file, outFile, "--sourcemap")
}

func buildStylesStylus(file string) error {
	outDir := filepath.Join("static", filepath.Dir(file))
	os.MkdirAll(outDir, os.ModePerm)
	return buildExec("Building styles with Stylus...",
		"stylus", file, "--out", outDir, "--compress", "--sourcemap")
}
