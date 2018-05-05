package main

import (
	"fmt"
	"github.com/red-chen/argparse"
	"os"
)

func rootFunc(c *argparse.Context) {
	fmt.Println("Hello world")
	mode, _ := c.GetString("mode")
	fmt.Println(mode)
	config, _ := c.GetString("config")
	fmt.Println(config)
	force, _ := c.GetBool("force")
	fmt.Println(force)
}

func uploadFunc(c *argparse.Context) {
	fmt.Println("Hello world")
	mode, _ := c.GetString("mode")
	fmt.Println(mode)
	config, _ := c.GetString("config")
	fmt.Println(config)
	force, _ := c.GetBool("force")
	fmt.Println(force)
}

func genParser() *argparse.Parser {
	root := argparse.ArgumentParser("app", "Sample tool for argparse")

	root.AddOption("mode", "testing mode").Short('m').Long("mode").Required().Default("test")
	root.AddOption("config", "config file").Short('c').Long("config").Required().Default("/tmp/config.json")
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
