package main

import (
	"fmt"
	"github.com/red-chen/goargs"
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
	root.AddOption("mode", "testing mode").Short('m').Long("mode").Required()
	root.AddOption("config", "config file").Short('c').Long("config").Default("/tmp/config.json")
	root.AddOption("force", "force do something").Short('f').Bool(true)
	root.SetDefaults(rootFunc)
	return root
}

func main() {
	parser := genParser()
	result := parser.ParseArgs(os.Args[1:])
	result.HandleError()
}
