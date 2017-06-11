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
	pkg, bin := findPackageFolder(pkg)
	err := runCmd("go install " + pkg)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = runCmd(bin + " " + strings.Join(os.Args[2:], " "))
	if err != nil {
		fmt.Println(err)
	}
}

func findPackageFolder(in string) (pkg string, bin string) {
	gopath := os.Getenv("GOPATH")
	gopaths := strings.Split(gopath, ":")
	for _, gp := range gopaths {
		dir := filepath.Join(gp, "src", "run", in)
		info, err := os.Stat(dir)
		if err == nil && info.IsDir() {
			return "run/" + in, in
		}
		dir = filepath.Join(gp, "src", in)
		info, err = os.Stat(dir)
		if err == nil && info.IsDir() {
			return in, in
		}
	}
	fmt.Println("couldn't find the source folder for the package", pkg)
	os.Exit(1)
	return "", ""
}

func runCmd(cmd string) error {
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
	info, err := os.Stat(file)
	if err == nil && !info.IsDir() {
		fmt.Println("The file", file, "already exists, doing nothing.")
		return
	}
	err = os.MkdirAll(dir, 0755)
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
