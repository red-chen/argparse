package goargs

import (
	"errors"
)

// TODO 支持更多的数据类型
// TODO 支持Print重写

// 主要是用于和用户的方法交互
type Context struct {
	options map[string]*Option // 包括了当前Parser以及之上的所有父Pasrser的Option
	parser  *Parser            // 关联的解析器
}

func (self *Context) GetString(dest string) (v string, err error) {
	if v, ok := self.options[dest]; ok {
		return v.getString(), nil
	}
	err = errors.New("NotFound")
	return
}

func (self *Context) GetBool(dest string) (v bool, err error) {
	if v, ok := self.options[dest]; ok {
		return v.getBool(), nil
	}
	err = errors.New("NotFound")
	return
}
