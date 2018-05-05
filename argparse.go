package argparse

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"os"
	"sort"
	"strings"
)

var (
	ERR_Usage    = errors.New("Usage")
	ERR_NotFound = errors.New("NotFound")
)

type Handler func(c *Context)

type Parser struct {
	Name        string
	Title       string
	Help        string
	Super       *Parser
	Root        *Parser
	Subs        map[string]*Parser // method - parser
	Opts        map[string]*Option // dest - option
	ShortOpts   map[rune]*Option   // short - option
	LongOpts    map[string]*Option // long - option
	HandlerFunc Handler
}

type Result struct {
	cx          *Context
	HandlerFunc Handler
}

func ArgumentParser(n string, h string) *Parser {
	self := &Parser{
		Name:        n,
		Title:       "root",
		Help:        h,
		Super:       nil,
		Subs:        map[string]*Parser{},
		Opts:        map[string]*Option{},
		ShortOpts:   map[rune]*Option{},
		LongOpts:    map[string]*Option{},
		HandlerFunc: nil,
	}

	self.Root = self
	return self
}

func (self *Parser) AddParser(cmd string, help string) *Parser {
	p := &Parser{
		Title: cmd,
		Help:  help,
		Super: self,
		Root:  self.Root,
		Subs:  map[string]*Parser{},
		Opts:  map[string]*Option{},
	}
	self.Subs[cmd] = p
	return p
}

func (self *Parser) AddOption(dest string, help string) *Option {
	arg := newOption(dest, help, self)
	self.Opts[dest] = arg
	return arg
}

func (self *Parser) SetDefaults(handler Handler) {
	self.HandlerFunc = handler
}

// ./app upload -m test -c /tmp/conifg.json
func (self *Parser) lookupParser(sources map[string]*Parser, cmds []string) (*Parser, error) {
	if sources == nil {
		return nil, ERR_NotFound
	}

	if len(cmds) == 0 {
		return self, nil
	} else if len(cmds) == 1 {
		m := cmds[0]
		if parser, ok := sources[m]; ok {
			return parser, nil
		}
		return nil, ERR_NotFound
	} else {
		m := cmds[0]
		if parser, ok := sources[m]; ok {
			return self.lookupParser(parser.Subs, cmds[1:])
		}
		return nil, ERR_NotFound
	}
}

// 逐层回溯父Parser的Opts
func (self *Parser) getLongOption(paser *Parser, opt string) (option *Option, exists bool) {
	option, exists = paser.LongOpts[opt]
	if !exists {
		if paser.Super == nil {
			return
		}
		return self.getLongOption(paser.Super, opt)
	}
	return
}

func (self *Parser) parseLongOption(opt string, params []string) (remain []string, err error) {
	var value interface{}
	opt = opt[2:]
	remain = params

	if len(opt) == 0 || opt[0] == '-' || opt[0] == '=' {
		err = errors.New("bad flag syntax")
		return
	}

	split := strings.SplitN(opt, "=", 2)
	opt = split[0]

	option, exists := self.getLongOption(self, opt)
	if !exists {
		if opt == "help" {
			fmt.Println(self.getUsage())
			err = ERR_Usage
			return
		}
		// skip unknown option
		return
	}

	if len(split) == 2 {
		// '--flag=arg'
		value = split[1]
	} else if option.setBool {
		// '--flag' (arg was optional)
		value = option.boolV
	} else if len(remain) > 0 {
		value = remain[0]
		remain = remain[1:]
	} else {
		// '--flag' (arg was required)
		err = errors.New(fmt.Sprintf("flag needs an argument: --%s", opt))
		return
	}

	err = option.parse(value)
	return
}

// 逐层回溯父Parser的Opts
func (self *Parser) getShortOption(paser *Parser, opt rune) (option *Option, exists bool) {
	option, exists = paser.ShortOpts[opt]
	if !exists {
		if paser.Super == nil {
			return
		}
		return self.getShortOption(paser.Super, opt)
	}
	return
}

func (self *Parser) parserSingleShortOption(opt string, params []string) (outOpt string, remain []string, err error) {
	remain = params
	outOpt = opt[1:]

	c := opt[0]
	option, exists := self.getShortOption(self, rune(c))

	if !exists {
		switch {
		case c == 'h':
			fmt.Println(self.getUsage())
			err = ERR_Usage
			return
		default:
			err = errors.New(fmt.Sprintf("unknown short flag: %q in -%s", c, opt))
			return
		}
	}

	var value interface{}
	if len(opt) > 2 && opt[1] == '=' {
		value = opt[2:]
		outOpt = ""
	} else if option.setBool {
		// '-f' (arg was optional)
		value = option.boolV
		outOpt = ""
	} else if len(opt) > 1 {
		value = opt[1:]
		outOpt = ""
	} else if len(params) > 0 {
		value = params[0]
		remain = params[1:]
	} else {
		err = errors.New(fmt.Sprintf("flag needs an argument: %q in -%s", c, opt))
		return
	}

	err = option.parse(value)
	return

}

func (self *Parser) parseShortOption(opt string, params []string) (remain []string, err error) {
	opt = opt[1:]
	remain = params

	// http://www.gnu.org/software/libc/manual/html_node/Argument-Syntax.html
	// Multiple options may follow a hyphen delimiter in a single token if the options do not take arguments.
	// Thus, ‘-abc’ is equivalent to ‘-a -b -c’.
	// short opt can bo a series of opt letters of flags (e.g "-abc")
	for len(opt) > 0 {
		opt, remain, err = self.parserSingleShortOption(opt, params)
		if err != nil {
			return
		}
	}
	return
}

func (self *Parser) bindParams(params []string) (err error) {
	for len(params) > 0 {
		s := params[0]
		params = params[1:]

		// 样式：--name
		if s[1] == '-' {
			if len(s) == 2 { // --
				return errors.New("NotSupport")
			}
			params, err = self.parseLongOption(s, params)
		} else {
			params, err = self.parseShortOption(s, params)
		}
		if err != nil {
			return
		}
	}
	return
}

func (self *Parser) getCmdsAndParams(input []string) (cmds []string, params []string) {
	for index, v := range input {
		//
		if len(v) == 0 || v[0] != '-' || len(v) == 1 {
			cmds = append(cmds, v)
		} else {
			params = input[index:]
			return
		}
	}
	return
}

func (self *Parser) preFilterAllOption(p *Parser) (options map[string]*Option) {
	options = map[string]*Option{}
	for _, v := range p.Opts {
		v.pre()
		options[v.dest] = v
	}

	for _, v := range p.Subs {
		tmp := self.preFilterAllOption(v)
		for kk, vv := range tmp {
			options[kk] = vv
		}
	}
	return
}

func (self *Parser) postFilterAllOption(p *Parser) (err error) {
	for _, v := range p.Opts {
		err = v.post()
	}

	for _, v := range p.Subs {
		err = self.postFilterAllOption(v)
	}

	return
}

func (self *Parser) ParseArgs(input []string) (result *Result) {
	var err error
	var parser *Parser

	options := self.preFilterAllOption(self)

	// 解析为对应的参数
	// 先找到method，然后根据method检查对应的参数
	cmds, params := self.getCmdsAndParams(input)

	// 查找对应的Parser
	if parser, err = self.lookupParser(self.Subs, cmds); err != nil {
		panic(err)
	}
	// 执行参数检查和绑定
	if err = parser.bindParams(params); err != nil {
		switch err {
		case ERR_Usage:
			os.Exit(0)
		default:
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if err = self.postFilterAllOption(self); err != nil {
		panic(err)
	}

	c := &Context{
		options: options,
		parser:  self,
	}
	result = &Result{
		cx:          c,
		HandlerFunc: self.HandlerFunc,
	}

	return
}

func (self *Result) Handle() {
	self.HandlerFunc(self.cx)
}

// for usage
func (self *Parser) OptionText() string {
	tmp := []string{}
	// 开始位置
	var startPoint int
	startPoint = 3 + len(self.Root.Name)

	for k, v := range self.Opts {
		if v.setBool {
			tmp = append(tmp, "["+v.getOptString()+"]")
		} else {
			if !v.requiredV {
				tmp = append(tmp, "["+v.getOptString()+" "+strings.ToUpper(k)+"]")
			} else {
				tmp = append(tmp, v.getOptString()+" "+strings.ToUpper(k))
			}
		}

	}

	sort.Strings(tmp)

	count := 0
	out := []string{}
	for _, s := range tmp {
		if count == 2 {
			out = append(out, "\n")
			out = append(out, strings.Repeat(" ", startPoint))
			count = 0
		}
		out = append(out, s)
		count++
	}

	return strings.Join(out, " ")
}

// TODO 格式化
func (self *Parser) OptionDetailText() string {
	buf := new(bytes.Buffer)
	lines := []string{}
	maxlen := 0
	for _, v := range self.Opts {
		line := ""

		if v.shortV != 0 {
			if v.setBool {
				line = fmt.Sprintf("  -%c      ", v.shortV)
			} else {
				line = fmt.Sprintf("  -%c, --%s", v.shortV, v.longV)
			}

		} else {
			line = fmt.Sprintf("      --%s", v.longV)
		}

		line += " " + v.help

		//line += "\x00"
		if len(line) > maxlen {
			maxlen = len(line)
		}
		lines = append(lines, line)

	}

	sort.Strings(lines)

	for _, line := range lines {
		//sidx := strings.Index(line, "\x00")
		//spacing := strings.Repeat(" ", maxlen-sidx)
		fmt.Fprintln(buf, line)
	}

	return buf.String()
}

func (self *Parser) SubCommandText() string {
	buf := new(bytes.Buffer)
	lines := []string{}
	for _, v := range self.Subs {
		line := fmt.Sprintf("  %s	%s", v.Title, v.Help)
		lines = append(lines, line)
	}

	sort.Strings(lines)

	for _, line := range lines {
		//sidx := strings.Index(line, "\x00")
		//spacing := strings.Repeat(" ", maxlen-sidx)
		fmt.Fprintln(buf, line)
	}

	return buf.String()
}

func (self *Parser) getUsage() string {
	var b bytes.Buffer
	var tmpl string

	tmpl = `{{ .Help }}

Usage:
    {{ .Root.Name }} {{ .OptionText }} {{ with .OptionDetailText}}

Options:
{{.}}{{end}} {{ with .SubCommandText}}
SubCommands:
{{.}}{{end}}
`

	t := template.New("top")
	template.Must(t.Parse(tmpl))
	t.Execute(&b, &self)

	return b.String()
}
