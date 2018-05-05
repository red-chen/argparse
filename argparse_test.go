package argparse

import "testing"

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
	parser = ArgumentParser()
	parser.AddOption("mode", "mode help")
	_, err = parser.parseLongOption(params[0], params[1:])

	if err == nil {
		t.Error()
	}

	if err.Error() != "flag needs an argument: --mode" {
		t.Error(err)
	}
}

func TestParseLongOptionBasicWithBoolTrue(t *testing.T) {
	var err error
	var parser *Parser
	var p *Option

	params := []string{"--force"}
	parser = ArgumentParser()
	// 表示在命令行中，如果声明了Flag，但是没有传入参数，那么就采用Bool中设置的值
	parser.AddOption("force", "force do something").Bool(true)
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
	parser = ArgumentParser()
	// 表示在命令行中，如果声明了Flag，但是没有传入参数，那么就采用Bool中设置的值
	parser.AddOption("force", "force do something").Bool(false)
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
	parser = ArgumentParser()
	parser.AddOption("mode", "mode help")
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
	parser = ArgumentParser()
	parser.AddOption("mode", "mode help")
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
	parser = ArgumentParser()
	parser.AddOption("mode", "mode help")

	upload := parser.AddParser("upload", "upload help")
	upload.AddOption("name", "name of file.")

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
	parser = ArgumentParser()
	parser.AddOption("mode", "mode help")

	upload := parser.AddParser("upload", "upload help")
	upload.AddOption("type", "type of file.")
	upload.AddOption("name", "name of file.")

	uploadFile := upload.AddParser("file", "upload file help")
	uploadFile.AddOption("len", "len of file.")

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

/*
测试用例：
*/
