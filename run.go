package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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
	err = runCmd(".", append([]string{bin}, args[1:]...)...)
	if err != nil {
		fmt.Println(err)
	}
	return 0
}

func findPackageFolder(in string) (pkg string, bin string, err error) {
	// convert to absolute path (if necessary)
	dir, err := filepath.Abs(in)
	if err != nil {
		return "", "", err
	}

	// get package name from absolute path (also works if in is ".")
	bin = filepath.Base(dir)

	// check if path exists and contains go files
	info, err := os.Stat(dir)
	if err == nil && info.IsDir() {
		if !containsGoFiles(dir) {
			return "", "", fmt.Errorf("%s does not contain any go files", dir)
		}
		return dir, bin, nil
	}

	// try package name in gorundir
	dir = filepath.Join(runDir(), in)
	info, err = os.Stat(dir)
	if err == nil && info.IsDir() {
		if !containsGoFiles(dir) {
			return "", "", fmt.Errorf("%s does not contain any go files", dir)
		}
		return dir, bin, nil
	}

	// no package found
	return "", "", fmt.Errorf(`couldn't find the source folder for package %s (create it with "run %s new")`,
		in, in)
}

func containsGoFiles(path string) bool {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if filepath.Ext(f.Name()) == ".go" {
			return true
		}
	}

	return false
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
