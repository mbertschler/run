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
	os.Exit(command(args))
}

func runDir() string {
	path := os.Getenv("GORUNDIR")
	if len(path) == 0 {
		path = os.Getenv("GOPATH")
		if len(path) == 0 {
			fmt.Println("GORUNDIR or GOPATH environment variable needs to be set to scripts directory")
			os.Exit(1)
		}
		path = filepath.Join(path, "src")
	}
	return path
}

func command(args []string) int {
	if len(args) < 1 {
		fmt.Println("Add the name of the package in the run folder you want to run,\n" +
			"or run 'run {packagename} new'. ")
		return 1
	}

	var pkg = args[0]
	dir, bin, err := findPackageFolder(pkg)
	if err != nil {
		if len(args) > 1 && args[1] == "new" {
			err := newPackage(args[0])
			if err != nil {
				fmt.Println(err)
				return 1
			}
			return 0
		}

		fmt.Println(err)
		return 1
	}
	err = runCmd(dir, "go", "install")
	if err != nil {
		fmt.Println(err)
		return 1
	}
	err = runCmd(dir, append([]string{bin}, args[1:]...)...)
	if err != nil {
		fmt.Println(err)
	}
	return 0
}

func findPackageFolder(in string) (pkg string, bin string, err error) {
	path := runDir()
	dir := filepath.Join(path, in)
	bin = in
	lastSlash := strings.LastIndex(in, "/")
	if lastSlash != -1 {
		if lastSlash == len(in)-1 {
			return "", "", errors.New(fmt.Sprint("package name can't end with a slash", in))
		}
		bin = in[lastSlash+1:]
	}
	info, err := os.Stat(dir)
	if err == nil && info.IsDir() {
		return dir, bin, nil
	}

	return "", "", errors.New(fmt.Sprint("couldn't find the source folder for the package ",
		in, " (create it with \"run ", in, " new\")"))
}

func runCmd(dir string, cmd ...string) error {
	c := exec.Command(cmd[0], cmd[1:]...)
	c.Dir = dir
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	return c.Run()
}

var content = `package main

func main() {

}`

func newPackage(name string) error {
	path := runDir()
	dir := filepath.Join(path, name)
	file := filepath.Join(path, name, name+".go")
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
		runCmd("/", "open", file)
	} else {
		runCmd("/", "xdg-open", file)
	}
	return nil
}
