# GoArgs

[![Build Status](https://travis-ci.org/red-chen/goargs.svg?branch=master)](https://travis-ci.org/red-chen/goargs)
[![Report Status](https://goreportcard.com/badge/github.com/red-chen/goargs)](https://goreportcard.com/report/github.com/red-chen/goargs)

[中文](README_cn.md)

The goargs module makes it easy to write user-friendly command-line interfaces.
The program defines what arguments it requires, and goargs will figure out how
to parse those out of os.Args. The goargs module also automatically generates
help and usage messages and issues errors when users give the program invalid
arguments.

# Example

The following code is sample code

### Example1 : Base sample

```
package main

import (
	"github.com/red-chen/goargs"
	"fmt"
	"os"
)

func rootFunc(c *goargs.Context) {
	fmt.Println("Hello world")
	mode, _ := c.GetString("mode")
	config, _ := c.GetString("config")
	force, _ := c.GetBool("force")

	fmt.Println(mode)
	fmt.Println(config)
	fmt.Println(force)
}

func genParser() *goargs.Parser {
	root := goargs.ArgumentParser("app", "Sample tool for goargs")
	root.AddOption("mode", "testing mode").Short('m').Long("mode").Default("test")
	root.AddOption("config", "config file").Short('c').Long("config").Default("/tmp/config.json")
	root.AddOption("force", "force do something").Short('f').Bool(true)
	root.SetDefaults(rootFunc)
	return root
}

func main() {
	parser := genParser()
	result := parser.ParseArgs(os.Args[1:])
	result.Handle()
}
```

Got useful help messages by '-h'':
```
# ./example1 -h
Sample tool for goargs

Usage:
    app -c/--config CONFIG -m/--mode MODE
        [-f]

Options:
  -c, --config force do something
  -f           testing mode
  -m, --mode   config file
```

Call binary:
```
# ./example1
Hello world
test
/tmp/config.json
false
```

### Example2 : Sub-Command

```
package main

import (
	"fmt"
	"github.com/red-chen/goargs"
	"os"
)

func rootFunc(c *goargs.Context) {
	mode, _ := c.GetString("mode")
	config, _ := c.GetString("config")
	force, _ := c.GetBool("force")

	fmt.Println(mode)
	fmt.Println(config)
	fmt.Println(force)
}

func uploadFunc(c *goargs.Context) {
	mode, _ := c.GetString("mode")
	config, _ := c.GetString("config")
	force, _ := c.GetBool("force")
	file, _ := c.GetBool("file")


	fmt.Println(mode)
	fmt.Println(config)
	fmt.Println(force)
	fmt.Println(file)
}

func genParser() *goargs.Parser {
	root := goargs.ArgumentParser("app", "Sample tool for goargs")

	root.AddOption("mode", "testing mode").Short('m').Long("mode").Default("test")
	root.AddOption("config", "config file").Short('c').Long("config").Default("/tmp/config.json")
	root.AddOption("force", "force do something").Short('f').Bool(true)
	root.SetDefaults(rootFunc)

	upload := root.AddParser("upload", "upload file to cloud")
	upload.AddOption("file", "file name").Required().Default("test")
	upload.SetDefaults(uploadFunc)

	return root

}

func main() {
	parser := genParser()
	result := parser.ParseArgs(os.Args[1:])
	result.Handle()
}
```



