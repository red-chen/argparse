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

func downloadFunc(c *goargs.Context) {
	mode, _ := c.GetString("mode")
	config, _ := c.GetString("config")
	force, _ := c.GetBool("force")
	out, _ := c.GetBool("out")

	fmt.Println(mode)
	fmt.Println(config)
	fmt.Println(force)
	fmt.Println(out)
}

func genParser() *goargs.Parser {
	root := goargs.ArgumentParser("app", "Sample tool for goargs")

	root.AddOption("mode", "testing mode").Short('m').Long("mode").Default("test")
	root.AddOption("config", "config file").Short('c').Long("config").Default("/tmp/config.json")
	root.AddOption("force", "force do something").Short('f').Bool(true)
	root.SetDefaults(rootFunc)

	upload := root.AddParser("upload", "Upload file to cloud")
	upload.AddOption("file", "file name").Required().Default("test")
	upload.SetDefaults(uploadFunc)

	download := root.AddParser("download", "Download file from cloud")
	download.AddOption("out", "file name").Required()
	download.SetDefaults(downloadFunc)

	return root

}

func main() {
	parser := genParser()
	result := parser.ParseArgs(os.Args[1:])
	result.Handle()
}
