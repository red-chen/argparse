package goargs

import (
	"errors"
	"fmt"
	"strings"
)

type Option struct {
	shortV    rune        // 简写选项，比如 -m
	longV     string      // 完整选项，比如 --mode
	dest      string      // 关键字，用户获取Option的数据
	requiredV bool        // 标记当前参数是否是必选
	help      string      // 帮助信息
	defValue  interface{} // 参数的默认值
	value     interface{} // 参数的值
	stored    bool        // 标记是否已经处理过option
	setBool   bool        // 标记是否设置了BoolV
	boolV     bool        // Bool的默认值
	father    *Parser     // Option的关联解析器， 目的是回写Short和Long Option
}

func newOption(dest string, help string, father *Parser) *Option {
	self := &Option{
		//longV: dest,  // 不能设定，只有当short和long都没有设定时，longV才被设置为dest
		dest:      dest,
		help:      help,
		requiredV: false,
		stored:    false,
		setBool:   false,
		father:    father,
	}
	return self
}

func (self *Option) Short(s rune) *Option {
	self.shortV = s
	// 将自身挂载到所属Parser的ShortOpts下面，用于short查找
	self.father.ShortOpts[s] = self
	return self
}

func (self *Option) Long(l string) *Option {
	self.longV = l
	// 将自身挂载到所属Parser的LongOpts下面，用于long查找
	self.father.LongOpts[l] = self
	return self
}

func (self *Option) Bool(b bool) *Option {
	self.boolV = b
	self.setBool = true
	return self
}

func (self *Option) Required() *Option {
	self.requiredV = true
	return self
}

func (self *Option) Default(v interface{}) *Option {
	self.defValue = v
	return self
}

func (self *Option) getString() string {
	if self.value == nil {
		if self.defValue != nil {
			return self.defValue.(string)
		}
		return ""
	} else {
		return self.value.(string)
	}
}

func (self *Option) getBool() bool {
	if self.setBool {
		return self.boolV
	}
	return self.value.(bool)
}

func (self *Option) getOptString() string {
	out := []string{}
	if self.shortV == 0 && self.longV == "" {
		out = append(out, fmt.Sprintf("--%s", self.dest))
	}
	if self.shortV != 0 {
		out = append(out, fmt.Sprintf("-%c", self.shortV))
	}
	if self.longV != "" {
		out = append(out, fmt.Sprintf("--%s", self.longV))
	}
	return strings.Join(out, "/")
}

func (self *Option) valid() error {
	if self.stored && self.value == nil && self.defValue == nil {
		return errors.New(fmt.Sprintf("Missing required option: '%s'", self.getOptString()))
	}

	if !self.stored && self.defValue == nil {
		return errors.New(fmt.Sprintf("Missing required option: '%s'", self.getOptString()))
	}
	return nil
}

// 检查Value是否合法，并赋值给Option
func (self *Option) parse(v interface{}) (err error) {
	self.stored = true
	self.value = v

	err = self.valid()
	return
}

// 预处理Option
func (self *Option) pre() {
	// 当时用户未显示设置Short和Long时，Long默认和Dest一样
	if self.longV == "" && self.shortV == 0 {
		self.Long(self.dest)
	}
}

// 后处理Option
func (self *Option) post() (err error) {
	// 检查所有必选参数是否已经设置
	if self.requiredV && self.defValue == nil && self.value == nil {
		err = errors.New(fmt.Sprintf("Missing required option: '%s'", self.getOptString()))
		return
	}
	// 检查Bool值是否已经设置
	if !self.stored {
		self.boolV = !self.boolV
	}
	return
}
