package goargs

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
	Title       string
	err         error
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
		Name:        cmd,
		Title:       cmd,
		Help:        help,
		Super:       self,
		Root:        self.Root,
		Subs:        map[string]*Parser{},
		Opts:        map[string]*Option{},
		ShortOpts:   map[rune]*Option{},
		LongOpts:    map[string]*Option{},
		HandlerFunc: nil,
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
func (self *Parser) lookupParser(sources map[string]*Parser, cmds []string, options map[string]*Option) (*Parser, error) {
	if sources == nil {
		return nil, ERR_NotFound
	}

	for _, v := range self.Opts {
		options[v.dest] = v
	}

	if len(cmds) == 0 {
		return self, nil
	} else {
		m := cmds[0]
		if parser, ok := sources[m]; ok {
			return parser.lookupParser(parser.Subs, cmds[1:], options)
		}
		return nil, ERR_NotFound
	}
}

// 逐层回溯父Parser的Opts
func (self *Parser) getLongOption(parser *Parser, opt string) (option *Option, exists bool) {
	option, exists = parser.LongOpts[opt]
	if !exists {
		if parser.Super == nil {
			return
		}
		return self.getLongOption(parser.Super, opt)
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
		// skip unknown option ?
		err = fmt.Errorf("Unrecognized arguments: --%s", opt)
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
		err = errors.New(fmt.Sprintf("Flag needs an argument: --%s", opt))
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
			err = errors.New(fmt.Sprintf("Unknown short flag: %q in -%s", c, opt))
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
		err = errors.New(fmt.Sprintf("Flag needs an argument: %q in -%s", c, opt))
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

// 前置动作，填写默认参数，默认Dest等
func (self *Parser) preFilterAllOption() {
	for _, v := range self.Opts {
		v.pre()
	}

	for _, v := range self.Subs {
		v.preFilterAllOption()
	}
	return
}

func (self *Parser) postFilterAllOption() (err error) {
	for _, v := range self.Opts {
		err = v.post()
	}

	for _, v := range self.Subs {
		err = v.postFilterAllOption()
	}

	return
}

func (self *Parser) ParseArgs(input []string) (result *Result) {
	var err error
	var parser *Parser

	cx := &Context{
		parser: self,
	}
	result = &Result{
		err:         err,
		cx:          cx,
		HandlerFunc: self.HandlerFunc,
	}

	// Pre操作, 设置默认的longV等
	self.preFilterAllOption()

	// 解析为对应的参数
	// 先找到method，然后根据method检查对应的参数
	cmds, params := self.getCmdsAndParams(input)

	// 查找对应的Parser
	options := map[string]*Option{}
	if parser, err = self.lookupParser(self.Subs, cmds, options); err != nil {
		result.err = err
		return
	}
	result.Title = parser.Title
	cx.options = options

	// 执行参数检查和绑定
	if err = parser.bindParams(params); err != nil {
		result.err = err
		return
	}

	// Post 操作，检查必选等
	if err = self.postFilterAllOption(); err != nil {
		result.err = err
		return
	}
	return
}

func (self *Result) Handle() (err error) {
	if self.err != nil {
		return self.err
	}
	if self.HandlerFunc == nil {
		return fmt.Errorf("missing handler in %s", self.Title)
	}
	self.HandlerFunc(self.cx)
	err = self.cx.err
	return
}

func (self *Result) HandleError() {
	self.Handle()
	if self.err != nil {
		switch self.err {
		case ERR_Usage:
			os.Exit(0)
		default:
			fmt.Printf("err: %s\n", self.err)
			os.Exit(1)
		}
	}
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
	prefixs := []string{}
	maxlen := 0

	for _, v := range self.Opts {
		line := ""
		if v.shortV != 0 {
			if v.setBool {
				sOpt := fmt.Sprintf("-%c", v.shortV)
				line = fmt.Sprintf("%4s", sOpt)
			} else {
				sOpt := fmt.Sprintf("-%c", v.shortV)
				line = fmt.Sprintf("%4s, --%s", sOpt, v.longV)
			}

		} else {
			line = fmt.Sprintf("--%s", v.longV)
		}

		if len(line) > maxlen {
			maxlen = len(line)
		}
		prefixs = append(prefixs, line)
	}

	f := fmt.Sprintf("%%-%ds", maxlen)

	lines := []string{}
	index := 0
	for _, v := range self.Opts {
		line := prefixs[index]
		out := fmt.Sprintf(f+" %s", line, v.help)
		lines = append(lines, out)
		index++
	}

	sort.Strings(lines)

	for _, line := range lines {
		fmt.Fprintln(buf, line)
	}

	return buf.String()
}

func (self *Parser) SubCommandText() string {
	buf := new(bytes.Buffer)
	lines := []string{}
	maxlen := 0
	for _, v := range self.Subs {
		if len(v.Title) > maxlen {
			maxlen = len(v.Title)
		}
	}
	f := fmt.Sprintf("%%-%ds", maxlen)
	for _, v := range self.Subs {
		line := fmt.Sprintf(f+" %s", v.Title, v.Help)
		lines = append(lines, line)
	}

	sort.Strings(lines)

	for _, line := range lines {
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
