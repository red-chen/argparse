# GoArgs

[![Build Status](https://travis-ci.org/red-chen/goargs.svg?branch=master)](https://travis-ci.org/red-chen/goargs)
[![Report Status](https://goreportcard.com/badge/github.com/red-chen/goargs)](https://goreportcard.com/report/github.com/red-chen/goargs)
[![Coverage Status](https://coveralls.io/repos/github/red-chen/goargs/badge.svg?branch=master)](https://coveralls.io/github/red-chen/goargs?branch=master)


GoArgs非常便于创建用户友好的命令行命令。比如设置可选和必选参数，添加不同的操作方法，GoArgs会自动
解析os.Args的输入参数。同时，GoArgs会自动根据用户编程的设定，显示友好的帮助信息。


GoArgs的思想来自于Python的argparser。


# 实例


### 实例一：最基础的用法

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

### 示例二 : 多个子方法

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



