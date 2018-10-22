package goargs

import (
	"fmt"
	"testing"
)

/*
测试用例：parseLongOption
1. 基础测试
  1.1 基础测试，保证服务能正常运行 - TestParseLongOptionBasic
  1.2 基础测试，测试Option采用空格的情况 - TestParseLongOptionBasicWithSpace
  1.3 基础测试，测试Option为父Parser定义 - TestParseLongOptionBasicOptInSuperParser
  1.4 基础测试，测试Option为父父Parser定义 - TestParseLongOptionBasicOptInSuperSuperParser
  1.5 基础测试，测试Bool设置是否正常
      - TestParseLongOptionBasicWithBoolTrue
      - TestParseLongOptionBasicWithBoolFalse
2. 异常测试
  2.1 传入参数不合法
      - TestParseLongOptionErrMissValue
  2.2 传入参数为--help
3. 边界测试
4. 特定场景测试
*/

func TestParseLongOptionErrMissValue(t *testing.T) {
	var err error
	var parser *Parser

	params := []string{"--mode"}
	parser = ArgumentParser("", "")
	parser.AddOption("mode", "mode help").Long("mode")
	_, err = parser.parseLongOption(params[0], params[1:])

	if err == nil {
		t.Error()
	}

	fmt.Println(err)
	if err.Error() != "Flag needs an argument: --mode" {
		t.Error(err)
	}
}

func TestParseLongOptionBasicWithBoolTrue(t *testing.T) {
	var err error
	var parser *Parser
	var p *Option

	params := []string{"--force"}
	parser = ArgumentParser("", "")
	// 表示在命令行中，如果声明了Flag，但是没有传入参数，那么就采用Bool中设置的值
	parser.AddOption("force", "force do something").Bool(true)
	parser.preFilterAllOption()
	_, err = parser.parseLongOption(params[0], params[1:])

	if err != nil {
		t.Error(err)
	}

	p, exists := parser.Opts["force"]

	if !exists {
		t.Error()
	}

	if p.getBool() != true {
		t.Error()
	}
}

func TestParseLongOptionBasicWithBoolFalse(t *testing.T) {
	var err error
	var parser *Parser
	var p *Option

	params := []string{"--force"}
	parser = ArgumentParser("", "")
	// 表示在命令行中，如果声明了Flag，但是没有传入参数，那么就采用Bool中设置的值
	parser.AddOption("force", "force do something").Bool(false)
	parser.preFilterAllOption()
	_, err = parser.parseLongOption(params[0], params[1:])

	if err != nil {
		t.Error(err)
	}

	p, exists := parser.Opts["force"]

	if !exists {
		t.Error()
	}

	if p.getBool() != false {
		t.Error()
	}
}

func TestParseLongOptionBasic(t *testing.T) {
	var remain []string
	var err error
	var parser *Parser

	params := []string{"--mode=test", "-p"}
	parser = ArgumentParser("", "")
	parser.AddOption("mode", "mode help").Long("mode")
	remain, err = parser.parseLongOption(params[0], params[1:])

	if err != nil {
		t.Error()
	}

	if len(params)-1 != len(remain) {
		t.Error()
	}

	p, exists := parser.Opts["mode"]

	if !exists {
		t.Error()
	}

	v := p.getString()

	if v != "test" {
		t.Error(v)
	}
}

func TestParseLongOptionBasicWithSpace(t *testing.T) {
	var remain []string
	var err error
	var parser *Parser

	params := []string{"--mode", "test", "-p"}
	parser = ArgumentParser("", "")
	parser.AddOption("mode", "mode help").Long("mode")
	remain, err = parser.parseLongOption(params[0], params[1:])

	if err != nil {
		t.Error()
	}

	if len(params)-2 != len(remain) {
		t.Error()
	}

	p, exists := parser.Opts["mode"]

	if !exists {
		t.Error()
	}

	v := p.getString()

	if v != "test" {
		t.Error(v)
	}
}

func TestParseLongOptionBasicOptInSuperParser(t *testing.T) {
	var remain []string
	var err error
	var parser *Parser
	var p *Option
	var exists bool
	var v string

	params := []string{"upload", "--mode=test", "--name", "readme.md"}
	parser = ArgumentParser("", "")
	parser.AddOption("mode", "mode help").Long("mode")

	upload := parser.AddParser("upload", "upload help")
	upload.AddOption("name", "name of file.").Long("name")

	// 解析参数
	remain, err = upload.parseLongOption(params[1], params[2:])

	if err != nil {
		t.Error()
	}

	if len(params)-2 != len(remain) {
		t.Error()
	}

	//
	p, exists = parser.Opts["mode"]

	if !exists {
		t.Error()
	}

	v = p.getString()

	if v != "test" {
		t.Error(v)
	}

	//
	remain, err = upload.parseLongOption(remain[0], remain[1:])
	p, exists = upload.Opts["name"]

	if !exists {
		t.Error()
	}

	v = p.getString()

	if v != "readme.md" {
		t.Error(v)
	}
}

func TestParseLongOptionBasicOptInSuperSuperParser(t *testing.T) {
	var remain []string
	var err error
	var parser *Parser
	var p *Option
	var exists bool
	var v string

	params := []string{"upload", "file", "--mode=test", "--name", "readme.md", "--len=178"}
	parser = ArgumentParser("", "")
	parser.AddOption("mode", "mode help").Long("mode")

	upload := parser.AddParser("upload", "upload help")
	upload.AddOption("type", "type of file.").Long("type")
	upload.AddOption("name", "name of file.").Long("name")

	uploadFile := upload.AddParser("file", "upload file help")
	uploadFile.AddOption("len", "len of file.").Long("len")

	// 解析参数
	remain, err = uploadFile.parseLongOption(params[2], params[3:])

	if err != nil {
		t.Error()
	}

	if len(params)-3 != len(remain) {
		t.Error()
	}

	//
	p, exists = parser.Opts["mode"]

	if !exists {
		t.Error()
	}

	v = p.getString()

	if v != "test" {
		t.Error(v)
	}

	//
	remain, err = uploadFile.parseLongOption(remain[0], remain[1:])
	p, exists = upload.Opts["name"]

	if !exists {
		t.Error()
	}

	v = p.getString()

	if v != "readme.md" {
		t.Error(v)
	}

	//
	remain, err = uploadFile.parseLongOption(remain[0], remain[1:])
	p, exists = uploadFile.Opts["len"]

	if !exists {
		t.Error()
	}

	v = p.getString()

	if v != "178" {
		t.Error(v)
	}
}

func TestPreFilterAllOption(t *testing.T) {
	//var err error
	var parser *Parser

	parser = ArgumentParser("", "")
	parser.AddOption("mode", "mode type")
	parser.preFilterAllOption()

	if len(parser.Opts) != 1 {
		t.Errorf("Expect the Opts length 1 but %d", len(parser.Opts))
	}

	var opt *Option = nil
	if v, ok := parser.Opts["mode"]; !ok {
		t.Errorf("Can not found the key: 'mode'")
	} else {
		opt = v
	}
	if opt.longV != "mode" {
		t.Errorf("Expect longV is 'mode' but '%s'", opt.longV)
	}

	if opt.shortV != 0 {
		t.Errorf("Expect shortV is None but '%c'", opt.shortV)
	}
}

func TestPreFilterAllOptionWithShort(t *testing.T) {
	//var err error
	var parser *Parser

	parser = ArgumentParser("", "")
	parser.AddOption("mode", "mode type").Short('m')
	parser.preFilterAllOption()

	if len(parser.Opts) != 1 {
		t.Errorf("Expect the Opts length 1 but %d", len(parser.Opts))
	}

	var opt *Option = nil
	if v, ok := parser.Opts["mode"]; !ok {
		t.Errorf("Can not found the key: 'mode'")
	} else {
		opt = v
	}
	if opt.longV != "" {
		t.Errorf("Expect longV is 'mode' but '%s'", opt.longV)
	}
	if opt.shortV != 'm' {
		t.Errorf("Expect shortV is 'm' but '%c'", opt.shortV)
	}
}

func TestGetCmdsAndParams(t *testing.T) {
	//var err error
	var parser *Parser

	input := []string{"test", "--mode", "debug"}

	parser = ArgumentParser("", "")
	cmds, params := parser.getCmdsAndParams(input)

	if len(cmds) != 1 {
		t.Errorf("Expect the cmds length 1 but %d", len(cmds))
	}

	if cmds[0] != "test" {
		t.Errorf("Expect the cmds is 'test' but '%s'", cmds[0])
	}

	if len(params) != 2 {
		t.Errorf("Expect the params length 2 but %d", len(params))
	}

	if params[0] != "--mode" {
		t.Errorf("Expect the params[0] is '--mode' but '%s'", params[0])
	}

	if params[1] != "debug" {
		t.Errorf("Expect the params[1] is 'debug' but '%s'", params[1])
	}
}

func TestGetCmdsAndParamsWithMultiCmds(t *testing.T) {
	//var err error
	var parser *Parser

	input := []string{"test", "add", "--mode", "debug"}

	parser = ArgumentParser("", "")
	cmds, params := parser.getCmdsAndParams(input)

	if len(cmds) != 2 {
		t.Errorf("Expect the cmds length 2 but %d", len(cmds))
	}

	if cmds[0] != "test" {
		t.Errorf("Expect the cmds is 'test' but '%s'", cmds[0])
	}

	if cmds[1] != "add" {
		t.Errorf("Expect the cmds is 'add' but '%s'", cmds[1])
	}

	if len(params) != 2 {
		t.Errorf("Expect the params length 2 but %d", len(params))
	}

	if params[0] != "--mode" {
		t.Errorf("Expect the params[0] is '--mode' but '%s'", params[0])
	}

	if params[1] != "debug" {
		t.Errorf("Expect the params[1] is 'debug' but '%s'", params[1])
	}
}

func TestLookupParser(t *testing.T) {
	var err error
	var parser *Parser

	input := []string{}

	parser = ArgumentParser("root", "help")
	parser.AddOption("mode", "mode type")
	parser.preFilterAllOption()
	options := map[string]*Option{}

	var tmp *Parser
	if tmp, err = parser.lookupParser(parser.Subs, input, options); err != nil {
		t.Error(err)
	}

	if tmp.Name != "root" {
		t.Errorf("Expect the Name is 'root' but '%s'", tmp.Name)
	}
}

func TestLookupParserSub(t *testing.T) {
	var err error
	var parser *Parser

	input := []string{"upload"}

	parser = ArgumentParser("root", "help")
	parser.AddOption("mode", "mode type")

	sub := parser.AddParser("upload", "upload help")
	sub.AddOption("path", "path help")
	options := map[string]*Option{}

	var tmp *Parser
	if tmp, err = parser.lookupParser(parser.Subs, input, options); err != nil {
		t.Error(err)
	}

	if tmp.Title != "upload" {
		t.Errorf("Expect the Title is 'upload' but '%s'", tmp.Title)
	}
}

func TestLookupParserSubSub(t *testing.T) {
	var err error
	var parser *Parser

	input := []string{"upload", "check"}

	parser = ArgumentParser("root", "help")
	parser.AddOption("mode", "mode type")

	sub := parser.AddParser("upload", "upload help")
	sub.AddOption("path", "path help")

	subsub := sub.AddParser("check", "upload check help")
	subsub.AddOption("ischeck", "ischeck help")

	options := map[string]*Option{}

	var tmp *Parser
	if tmp, err = parser.lookupParser(parser.Subs, input, options); err != nil {
		t.Error(err)
	}

	if tmp.Title != "check" {
		t.Errorf("Expect the Title is 'check' but '%s'", tmp.Title)
	}
}

func TestLookupParserErrorSub(t *testing.T) {
	var err error
	var parser *Parser

	input := []string{"upload1"}

	parser = ArgumentParser("root", "help")
	parser.AddOption("mode", "mode type")

	sub := parser.AddParser("upload", "upload help")
	sub.AddOption("path", "path help")

	options := map[string]*Option{}

	if _, err = parser.lookupParser(parser.Subs, input, options); err == nil {
		t.Error("Excpect can not found the parser")
	}
}

func TestLookupParserWithErrSub(t *testing.T) {
	var err error
	var parser *Parser

	input := []string{"upload1", "check"}

	parser = ArgumentParser("root", "help")
	parser.AddOption("mode", "mode type")

	sub := parser.AddParser("upload", "upload help")
	sub.AddOption("path", "path help")

	subsub := sub.AddParser("check", "upload check help")
	subsub.AddOption("ischeck", "ischeck help")

	options := map[string]*Option{}

	if _, err = parser.lookupParser(parser.Subs, input, options); err == nil {
		t.Error("Excpect can not found the parser")
	}
}

func TestLookupParserWithErrSubSub(t *testing.T) {
	var err error
	var parser *Parser

	input := []string{"upload1", "check1"}

	parser = ArgumentParser("root", "help")
	parser.AddOption("mode", "mode type")

	sub := parser.AddParser("upload", "upload help")
	sub.AddOption("path", "path help")

	subsub := sub.AddParser("check", "upload check help")
	subsub.AddOption("ischeck", "ischeck help")

	options := map[string]*Option{}

	if _, err = parser.lookupParser(parser.Subs, input, options); err == nil {
		t.Error("Excpect can not found the parser")
	}
}

func TestLookupParserMultiSub(t *testing.T) {
	var err error
	var parser *Parser

	parser = ArgumentParser("root", "help")
	parser.AddOption("mode", "mode type")

	upload := parser.AddParser("upload", "upload help")
	upload.AddOption("path", "upload path help")

	download := parser.AddParser("download", "download help")
	download.AddOption("path", "download path help")

	{
		options := map[string]*Option{}
		input := []string{"upload"}
		var tmp *Parser
		if tmp, err = parser.lookupParser(parser.Subs, input, options); err != nil {
			t.Error(err)
		}

		if tmp.Title != "upload" {
			t.Errorf("Expect the Title is 'upload' but '%s'", tmp.Title)
		}

		if len(options) != 2 {
			t.Errorf("Expect the options length is 2 but %d", len(options))
		}

		{
			if v, ok := options["mode"]; !ok {
				t.Error("Can not found 'mode'")
			} else {
				assertEqual(t, "mode", v.dest)
				assertEqual(t, "mode type", v.help)
			}
		}
		{
			if v, ok := options["path"]; !ok {
				t.Error("Can not found 'path'")
			} else {
				assertEqual(t, "path", v.dest)
				assertEqual(t, "upload path help", v.help)
			}
		}
	}
	{
		options := map[string]*Option{}
		input := []string{"download"}
		var tmp *Parser
		if tmp, err = parser.lookupParser(parser.Subs, input, options); err != nil {
			t.Error(err)
		}

		assertEqual(t, "download", tmp.Title)
		assertEqualInt(t, 2, len(options))

		{
			if v, ok := options["mode"]; !ok {
				t.Error("Can not found 'mode'")
			} else {
				assertEqual(t, "mode", v.dest)
				assertEqual(t, "mode type", v.help)
			}
		}
		{
			if v, ok := options["path"]; !ok {
				t.Error("Can not found 'path'")
			} else {
				assertEqual(t, "path", v.dest)
				assertEqual(t, "download path help", v.help)
			}
		}
	}
}

func TestBindParams(t *testing.T) {
	var err error
	var parser *Parser

	parser = ArgumentParser("root", "help")
	parser.AddOption("mode", "mode type")

	upload := parser.AddParser("upload", "upload help")
	upload.AddOption("path", "upload path help")

	download := parser.AddParser("download", "download help")
	download.AddOption("path", "download path help")

	options := map[string]*Option{}

	tmpParser, _ := parser.lookupParser(parser.Subs, []string{"upload"}, options)

	err = tmpParser.bindParams([]string{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestBindParamsWithValidArgs(t *testing.T) {
	var err error
	var parser *Parser

	parser = ArgumentParser("root", "help")
	parser.AddOption("mode", "mode type")

	upload := parser.AddParser("upload", "upload help")
	upload.AddOption("path", "upload path help").Short('p')

	download := parser.AddParser("download", "download help")
	download.AddOption("path", "download path help")

	parser.preFilterAllOption()

	options := map[string]*Option{}

	tmpParser, _ := parser.lookupParser(parser.Subs, []string{"upload"}, options)

	err = tmpParser.bindParams([]string{"--mode=test", "-p", "/home/admin"})
	if err != nil {
		t.Fatal(err)
	}

	{
		if o, ok := options["mode"]; !ok {
			t.Fatal()
		} else {
			assertEqual(t, "test", o.getString())
		}
	}
	{
		if o, ok := options["path"]; !ok {
			t.Fatal()
		} else {
			assertEqual(t, "/home/admin", o.getString())
		}
	}
}

func TestBindParamsWithInvalidArgs(t *testing.T) {
	var err error
	var parser *Parser

	parser = ArgumentParser("root", "help")
	parser.AddOption("mode", "mode type")

	upload := parser.AddParser("upload", "upload help")
	upload.AddOption("path", "upload path help")

	download := parser.AddParser("download", "download help")
	download.AddOption("path", "download path help")

	options := map[string]*Option{}

	tmpParser, _ := parser.lookupParser(parser.Subs, []string{"upload"}, options)

	err = tmpParser.bindParams([]string{"--modex=test"})
	if err == nil {
		t.Fatal()
	}
	assertEqual(t, "Unrecognized arguments: --modex", err.Error())
}

func TestBindParamsWithValidShortArgs(t *testing.T) {
	var err error
	var parser *Parser

	parser = ArgumentParser("root", "help")
	parser.AddOption("mode", "mode type").Short('m')

	upload := parser.AddParser("upload", "upload help")
	upload.AddOption("path", "upload path help").Short('p')

	download := parser.AddParser("download", "download help")
	download.AddOption("path", "download path help").Short('p')

	options := map[string]*Option{}

	tmpParser, _ := parser.lookupParser(parser.Subs, []string{"upload"}, options)

	err = tmpParser.bindParams([]string{"-m test"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestBindParamsWithInvalidShortArgs(t *testing.T) {
	var err error
	var parser *Parser

	parser = ArgumentParser("root", "help")
	parser.AddOption("mode", "mode type").Short('m')

	upload := parser.AddParser("upload", "upload help")
	upload.AddOption("path", "upload path help").Short('p')

	download := parser.AddParser("download", "download help")
	download.AddOption("path", "download path help").Short('p')

	options := map[string]*Option{}

	tmpParser, _ := parser.lookupParser(parser.Subs, []string{"upload"}, options)

	err = tmpParser.bindParams([]string{"-x test"})
	if err == nil {
		t.Fatal()
	}
	assertEqual(t, "Unknown short flag: 'x' in -x test", err.Error())
}

func TestPostFilterAllOption(t *testing.T) {
	var err error
	var parser *Parser

	parser = ArgumentParser("root", "help")
	parser.AddOption("mode", "mode type").Short('m').Required()

	upload := parser.AddParser("upload", "upload help")
	upload.AddOption("path", "upload path help").Short('p').Required()

	download := parser.AddParser("download", "download help")
	download.AddOption("path", "download path help").Short('p').Required()

	parser.preFilterAllOption()

	options := map[string]*Option{}

	tmpParser, _ := parser.lookupParser(parser.Subs, []string{"upload"}, options)

	tmpParser.bindParams([]string{"-m test", "-p /home/admin"})

	err = tmpParser.postFilterAllOption()
	if err != nil {
		t.Fatal(err)
	}
}

func TestPostFilterAllOptionNoValue(t *testing.T) {
	var err error
	var parser *Parser

	parser = ArgumentParser("root", "help")
	parser.AddOption("mode", "mode type").Short('m').Required()

	upload := parser.AddParser("upload", "upload help")
	upload.AddOption("path", "upload path help").Short('p').Required()

	download := parser.AddParser("download", "download help")
	download.AddOption("path", "download path help").Short('p').Required()

	parser.preFilterAllOption()

	options := map[string]*Option{}

	tmpParser, _ := parser.lookupParser(parser.Subs, []string{"upload"}, options)

	tmpParser.bindParams([]string{"-m test"})

	err = tmpParser.postFilterAllOption()
	if err == nil {
		t.Fatal()
	}
	assertEqual(t, "Missing required option: '-p'", err.Error())
}

func TestPostFilterAllOptionInvalidValue(t *testing.T) {
	var err error
	var parser *Parser

	parser = ArgumentParser("root", "help")
	parser.AddOption("mode", "mode type").Short('m').Required()

	upload := parser.AddParser("upload", "upload help")
	upload.AddOption("path", "upload path help").Short('p').Required()

	download := parser.AddParser("download", "download help")
	download.AddOption("path", "download path help").Short('p').Required()

	parser.preFilterAllOption()

	options := map[string]*Option{}

	tmpParser, _ := parser.lookupParser(parser.Subs, []string{}, options)

	err = tmpParser.bindParams([]string{"-m test", "-p /home/admin"})
	if err == nil {
		t.Fatal()
	}

	assertEqual(t, "Unknown short flag: 'p' in -p /home/admin", err.Error())
}

// function test

func rootFunc(c *Context) {
	mode, _ := c.GetString("mode")

	if mode != "debug" {
		c.Error(fmt.Errorf("Expect the mode is 'debug' but '%s'", mode))
		return
	}
}

func Test_FT_normal(t *testing.T) {
	var err error
	var parser *Parser

	{
		parser = ArgumentParser("root", "help")
		parser.AddOption("mode", "mode type").Long("modelong").Short('m').Required()
		parser.SetDefaults(rootFunc)

		result := parser.ParseArgs([]string{"-m", "debug"})
		err = result.Handle()

		if err != nil {
			t.Fatal(err)
		}
	}
	{
		parser = ArgumentParser("root", "help")
		parser.AddOption("mode", "mode type").Long("modelong").Short('m').Required()
		parser.SetDefaults(rootFunc)
		result := parser.ParseArgs([]string{})
		err = result.Handle()

		if err == nil {
			t.Fatal()
		}

		assertEqual(t, "Missing required option: '-m/--modelong'", err.Error())
	}
}

func Test_FT_sub(t *testing.T) {

}

func Test_FT_subsub(t *testing.T) {

}

func Test_FT_defaultValue(t *testing.T) {

}

func Test_FT_short(t *testing.T) {

}

func Test_FT_long(t *testing.T) {

}

func Test_FT_short_long(t *testing.T) {

}

func Test_FT_required(t *testing.T) {

}

func Test_FT_bool(t *testing.T) {

}

func Test_FT_flag(t *testing.T) {

}
