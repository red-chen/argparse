package goargs

import "testing"

func Test_Option_Basic(t *testing.T) {
	var err error
	// 最简单初始化测试
	{
		parser := &Parser{
			ShortOpts: map[rune]*Option{},
			LongOpts:  map[string]*Option{},
		}
		arg := newOption("mode", "help mode", parser)
		if err = arg.parse("test"); err != nil {
			t.Error(err)
		}

		// check
		if "mode" != arg.dest {
			t.Error()
		}

		if "help mode" != arg.help {
			t.Error()
		}

		if 0 != arg.shortV {
			t.Error()
		}

		if "" != arg.longV {
			t.Error()
		}

		if false != arg.requiredV {
			t.Error()
		}

		if nil != arg.defValue {
			t.Error()
		}

		if !arg.stored {
			t.Error()
		}

		if v := arg.getString(); "test" != v {
			t.Error()
		}
	}
	// 单独设置Long
	{
		parser := &Parser{
			ShortOpts: map[rune]*Option{},
			LongOpts:  map[string]*Option{},
		}
		arg := newOption("mode", "help mode", parser).Long("modeopt")
		if err = arg.parse("test"); err != nil {
			t.Error(err)
		}

		// check
		if "mode" != arg.dest {
			t.Error()
		}

		if "help mode" != arg.help {
			t.Error()
		}

		if 0 != arg.shortV {
			t.Error()
		}

		if "modeopt" != arg.longV {
			t.Error()
		}

		if false != arg.requiredV {
			t.Error()
		}

		if nil != arg.defValue {
			t.Error()
		}

		if !arg.stored {
			t.Error()
		}

		if v := arg.getString(); "test" != v {
			t.Error()
		}
	}
	// 单独设置Short
	{
		parser := &Parser{
			ShortOpts: map[rune]*Option{},
			LongOpts:  map[string]*Option{},
		}
		arg := newOption("mode", "help mode", parser).Short('x')
		if err = arg.parse("test"); err != nil {
			t.Error(err)
		}

		// check
		if "mode" != arg.dest {
			t.Error()
		}

		if "help mode" != arg.help {
			t.Error()
		}

		if 'x' != arg.shortV {
			t.Error()
		}

		if "" != arg.longV {
			t.Error()
		}

		if false != arg.requiredV {
			t.Error()
		}

		if nil != arg.defValue {
			t.Error()
		}

		if !arg.stored {
			t.Error()
		}

		if v := arg.getString(); "test" != v {
			t.Error()
		}
	}
	// 开启Required
	{
		parser := &Parser{
			ShortOpts: map[rune]*Option{},
			LongOpts:  map[string]*Option{},
		}
		arg := newOption("mode", "help mode", parser).Required()
		if err = arg.parse("test"); err != nil {
			t.Error(err)
		}

		// check
		if "mode" != arg.dest {
			t.Error()
		}

		if "help mode" != arg.help {
			t.Error()
		}

		if 0 != arg.shortV {
			t.Error()
		}

		if "" != arg.longV {
			t.Error()
		}

		if true != arg.requiredV {
			t.Error()
		}

		if nil != arg.defValue {
			t.Error()
		}

		if !arg.stored {
			t.Error()
		}

		if v := arg.getString(); "test" != v {
			t.Error()
		}
	}
	// 设置默认值
	{
		parser := &Parser{
			ShortOpts: map[rune]*Option{},
			LongOpts:  map[string]*Option{},
		}
		arg := newOption("mode", "help mode", parser).Default("admin")

		if d := arg.getString(); "admin" != d {
			t.Error()
		}

		if err = arg.parse("test"); err != nil {
			t.Error(err)
		}

		// check
		if "mode" != arg.dest {
			t.Error()
		}

		if "help mode" != arg.help {
			t.Error()
		}

		if 0 != arg.shortV {
			t.Error()
		}

		if "" != arg.longV {
			t.Error()
		}

		if false != arg.requiredV {
			t.Error()
		}

		if "admin" != arg.defValue.(string) {
			t.Error()
		}

		if !arg.stored {
			t.Error()
		}

		if v := arg.getString(); "test" != v {
			t.Error()
		}
	}
	// Bool值
	{
		parser := &Parser{
			ShortOpts: map[rune]*Option{},
			LongOpts:  map[string]*Option{},
		}
		arg := newOption("mode", "help mode", parser).Bool(true)

		if d := arg.getBool(); true != d {
			t.Error()
		}

		if err = arg.parse("test"); err != nil {
			t.Error(err)
		}

		// check
		if "mode" != arg.dest {
			t.Error()
		}

		if "help mode" != arg.help {
			t.Error()
		}

		if 0 != arg.shortV {
			t.Error()
		}

		if "" != arg.longV {
			t.Error()
		}

		if false != arg.requiredV {
			t.Error()
		}

		if nil != arg.defValue {
			t.Error()
		}

		if !arg.stored {
			t.Error()
		}

		if v := arg.getString(); "test" != v {
			t.Error()
		}
	}

	// 混合
	{
		parser := &Parser{
			ShortOpts: map[rune]*Option{},
			LongOpts:  map[string]*Option{},
		}
		arg := newOption("mode", "help mode", parser).Short('m').Long("modex").Required().Default("admin")

		if d := arg.getString(); "admin" != d {
			t.Error()
		}

		if err = arg.parse("test"); err != nil {
			t.Error(err)
		}

		// check
		if "mode" != arg.dest {
			t.Error()
		}

		if "help mode" != arg.help {
			t.Error()
		}

		if 'm' != arg.shortV {
			t.Error()
		}

		if "modex" != arg.longV {
			t.Error()
		}

		if true != arg.requiredV {
			t.Error()
		}

		if "admin" != arg.defValue.(string) {
			t.Error()
		}

		if !arg.stored {
			t.Error()
		}

		if v := arg.getString(); "test" != v {
			t.Error()
		}
	}
}

func Test_Option_Required(t *testing.T) {
	var err error
	var arg *Option

	parser := &Parser{
		ShortOpts: map[rune]*Option{},
		LongOpts:  map[string]*Option{},
	}

	arg = newOption("mode", "help mode", parser).Required()

	// check
	if "mode" != arg.dest {
		t.Error()
	}

	if "help mode" != arg.help {
		t.Error()
	}

	if 0 != arg.shortV {
		t.Error()
	}

	if false == arg.requiredV {
		t.Error()
	}

	if nil != arg.defValue {
		t.Error()
	}

	if true == arg.stored {
		t.Error()
	}

	err = arg.valid()

	if err == nil {
		t.Error()
	}

	if err.Error() != "Missing required option: '--mode'" {
		t.Error(err)
	}

	arg.Default("test")

	err = arg.valid()

	if err != nil {
		t.Error()
	}

	arg = newOption("mode", "help mode", parser).Required().Short('m')
	if 'm' != arg.shortV {
		t.Error()
	}
	err = arg.valid()
	if err == nil {
		t.Error()
	}
	if err.Error() != "Missing required option: '-m'" {
		t.Error(err)
	}
	arg = newOption("mode", "help mode", parser).Required().Short('m').Long("mymode")
	if 'm' != arg.shortV {
		t.Error()
	}
	err = arg.valid()
	if err == nil {
		t.Error()
	}
	if err.Error() != "Missing required option: '-m/--mymode'" {
		t.Error(err)
	}
}
