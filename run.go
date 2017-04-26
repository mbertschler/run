package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Add the name of the package in the run folder you want to run,\n" +
			"or run 'run new {packagename}'. ")
		return
	}
	if len(os.Args) > 2 && os.Args[1] == "new" {
		newPackage(os.Args[2])
		return
	}
	var pkg = os.Args[1]
	err := runCmd("go install run/" + pkg)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = runCmd(pkg + " " + strings.Join(os.Args[2:], " "))
	if err != nil {
		fmt.Println(err)
	}
}

func runCmd(cmd string) error {
	// log.Println(cmd)
	parts := strings.Fields(cmd)
	c := exec.Command(parts[0], parts[1:]...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	return c.Run()
}

var content = `package main

func main() {

}`

func newPackage(name string) {
	gopath := os.Getenv("GOPATH")
	gopath = strings.Split(gopath, ":")[0]
	dir := filepath.Join(gopath, "src", "run", name)
	file := filepath.Join(gopath, "src", "run", name, name+".go")
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ioutil.WriteFile(file, []byte(content), 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	if runtime.GOOS == "darwin" {
		runCmd("open " + file)
	} else {
		runCmd("xdg-open " + file)
	}
}
