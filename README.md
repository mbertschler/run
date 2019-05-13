`run` - build and run Go apps with one command
==============================================

This tool can be used to speed up the iteration time when you are working
on Go applications. It basically combines `go install && your-app`
into a single step. Another feature is to quickly setup a new package.

Installation
------------

```bash
go get -u github.com/mbertschler/run
```

Usage
-----

### Build and run a package

This works if the package mycmd is located in either of those locations:
- `$GORUNDIR/mycmd`

```bash
# basic usage:
run mycmd

# pass any arguments as you usually would:
run mycmd -flag value file.txt

```

### Create a new package

This command creates a new Go command in `$GORUNDIR/mycmd`.
This will create a folder, an empty main.go file, and open this file 
in your editor. This is useful for small tools and scripting applications.

```bash
# create a new package:
run new newcmd

# build and run this package:
run newcmd
```

License
-------

This tool is released under the Apache 2.0 license. See
[LICENSE](https://github.com/mbertschler/run/blob/master/LICENSE).
