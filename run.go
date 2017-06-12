package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	args := os.Args[1:]
	os.Exit(Main(args))
}

func Main(args []string) int {
	if len(args) < 1 {
		fmt.Println("Add the name of the package in the run folder you want to run,\n" +
			"or run 'run new {packagename}'. ")
		return 1
	}
	if len(args) > 1 && args[0] == "new" {
		err := newPackage(args[1])
		if err != nil {
			fmt.Println(err)
			return 1
		}
		return 0
	}
	var pkg = args[0]
	pkg, bin, err := findPackageFolder(pkg)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	err = runCmd("go install " + pkg)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	err = runCmd(bin + " " + strings.Join(args[1:], " "))
	if err != nil {
		fmt.Println(err)
	}
	return 0
}

func findPackageFolder(in string) (pkg string, bin string, err error) {
	gopath := os.Getenv("GOPATH")
	gopaths := strings.Split(gopath, ":")
	for _, gp := range gopaths {
		dir := filepath.Join(gp, "src", "run", in)
		var bin = in
		lastSlash := strings.LastIndex(in, "/")
		if lastSlash != -1 {
			if lastSlash == len(in)-1 {
				return "", "", errors.New(fmt.Sprint("package name can't end with a slash", in))
			}
			bin = in[lastSlash+1:]
		}
		info, err := os.Stat(dir)
		if err == nil && info.IsDir() {
			return "run/" + in, bin, nil
		}
		dir = filepath.Join(gp, "src", in)
		info, err = os.Stat(dir)
		if err == nil && info.IsDir() {
			return in, bin, nil
		}
	}
	return "", "", errors.New(fmt.Sprint("couldn't find the source folder for the package ", in))
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

func newPackage(name string) error {
	gopath := os.Getenv("GOPATH")
	gopath = strings.Split(gopath, ":")[0]
	dir := filepath.Join(gopath, "src", "run", name)
	file := filepath.Join(gopath, "src", "run", name, name+".go")
	info, err := os.Stat(file)
	if err == nil && !info.IsDir() {
		return errors.New(fmt.Sprint("The file ", file, " already exists, doing nothing."))
	}
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(file, []byte(content), 0644)
	if err != nil {
		return err
	}
	if runtime.GOOS == "darwin" {
		runCmd("open " + file)
	} else {
		runCmd("xdg-open " + file)
	}
	return nil
}
