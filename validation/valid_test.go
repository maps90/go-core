package validation

import (
	"regexp"
	"testing"
	"time"
)

func TestRequired(t *testing.T) {
	valid := Validation{}

	if valid.Required(nil, "nil").Ok {
		t.Error("nil object should be false")
	}
	if !valid.Required(true, "bool").Ok {
		t.Error("Bool value should always return true")
	}
	if !valid.Required(false, "bool").Ok {
		t.Error("Bool value should always return true")
	}
	if !valid.Required("gocore", "string").Ok {
		t.Error("string should be true")
	}
	if valid.Required(0, "zero").Ok {
		t.Error("Integer should not be equal 0")
	}
	if !valid.Required(1, "int").Ok {
		t.Error("Integer except 0 should be true")
	}
	if !valid.Required(time.Now(), "time").Ok {
		t.Error("time should be true")
	}
	if valid.Required([]string{}, "emptySlice").Ok {
		t.Error("empty slice should be false")
	}
	if !valid.Required([]interface{}{"ok"}, "slice").Ok {
		t.Error("slice should be true")
	}
}

func TestMin(t *testing.T) {
	valid := Validation{}

	if valid.Min(-1, 0, "min0").Ok {
		t.Error("-1 is less than the minimum value of 0 should be false")
	}
	if !valid.Min(1, 0, "min0").Ok {
		t.Error("1 is greater or equal than the minimum value of 0 should be true")
	}
}

func TestMax(t *testing.T) {
	valid := Validation{}

	if valid.Max(1, 0, "max0").Ok {
		t.Error("1 is greater than the minimum value of 0 should be false")
	}
	if !valid.Max(-1, 0, "max0").Ok {
		t.Error("-1 is less or equal than the maximum value of 0 should be true")
	}
}

func TestRange(t *testing.T) {
	valid := Validation{}

	if valid.Range(-1, 0, 1, "range0_1").Ok {
		t.Error("-1 is between 0 and 1 should be false")
	}
	if !valid.Range(1, 0, 1, "range0_1").Ok {
		t.Error("1 is between 0 and 1 should be true")
	}
}

func TestMinSize(t *testing.T) {
	valid := Validation{}

	if !valid.MinSize("ok", 1, "minSize1").Ok {
		t.Error("the length of \"ok\" is greater or equal than the minimum value of 1 should be true")
	}
}

func TestMaxSize(t *testing.T) {
	valid := Validation{}

	if valid.MaxSize("ok", 1, "maxSize1").Ok {
		t.Error("the length of \"ok\" is greater than the maximum value of 1 should be false")
	}
	if !valid.MaxSize("", 1, "maxSize1").Ok {
		t.Error("the length of \"\" is less or equal than the maximum value of 1 should be true")
	}
	if valid.MaxSize([]interface{}{"ok", false}, 1, "maxSize1").Ok {
		t.Error("the length of [\"ok\", false] is greater than the maximum value of 1 should be false")
	}
	if !valid.MaxSize([]string{}, 1, "maxSize1").Ok {
		t.Error("the length of empty slice is less or equal than the maximum value of 1 should be true")
	}
}

func TestLength(t *testing.T) {
	valid := Validation{}

	if !valid.Length("1", 1, "length1").Ok {
		t.Error("the length of \"1\" must equal 1 should be true")
	}
	if !valid.Length([]interface{}{"ok"}, 1, "length1").Ok {
		t.Error("the length of [\"ok\"] must equal 1 should be true")
	}
}

func TestAlpha(t *testing.T) {
	valid := Validation{}

	if valid.Alpha("a,1-@ $", "alpha").Ok {
		t.Error("\"a,1-@ $\" are valid alpha characters should be false")
	}
	if !valid.Alpha("abCD", "alpha").Ok {
		t.Error("\"abCD\" are valid alpha characters should be true")
	}
}

func TestNumeric(t *testing.T) {
	valid := Validation{}

	if valid.Numeric("a,1-@ $", "numeric").Ok {
		t.Error("\"a,1-@ $\" are valid numeric characters should be false")
	}
	if !valid.Numeric("1234", "numeric").Ok {
		t.Error("\"1234\" are valid numeric characters should be true")
	}
}

func TestAlphaNumeric(t *testing.T) {
	valid := Validation{}

	if valid.AlphaNumeric("a,1-@ $", "alphaNumeric").Ok {
		t.Error("\"a,1-@ $\" are valid alpha or numeric characters should be false")
	}
	if !valid.AlphaNumeric("1234aB", "alphaNumeric").Ok {
		t.Error("\"1234aB\" are valid alpha or numeric characters should be true")
	}
}

func TestMatch(t *testing.T) {
	valid := Validation{}

	if valid.Match("dimas@gmail", regexp.MustCompile("^\\w+@\\w+\\.\\w+$"), "match").Ok {
		t.Error("\"dimas@gmail\" match \"^\\w+@\\w+\\.\\w+$\"  should be false")
	}
	if !valid.Match("dimas@gmail.com", regexp.MustCompile("^\\w+@\\w+\\.\\w+$"), "match").Ok {
		t.Error("\"dimas@gmail\" match \"^\\w+@\\w+\\.\\w+$\"  should be true")
	}
}

func TestNoMatch(t *testing.T) {
	valid := Validation{}

	if valid.NoMatch("123@gmail", regexp.MustCompile("[^\\w\\d]"), "nomatch").Ok {
		t.Error("\"123@gmail\" not match \"[^\\w\\d]\"  should be false")
	}
	if !valid.NoMatch("123gmail", regexp.MustCompile("[^\\w\\d]"), "match").Ok {
		t.Error("\"123@gmail\" not match \"[^\\w\\d@]\"  should be true")
	}
}

func TestAlphaDash(t *testing.T) {
	valid := Validation{}

	if valid.AlphaDash("a,1-@ $", "alphaDash").Ok {
		t.Error("\"a,1-@ $\" are valid alpha or numeric or dash(-_) characters should be false")
	}
	if !valid.AlphaDash("1234aB-_", "alphaDash").Ok {
		t.Error("\"1234aB\" are valid alpha or numeric or dash(-_) characters should be true")
	}
}

func TestValid(t *testing.T) {
	type user struct {
		ID   int
		Name string `valid:"Required;Match(/^(test)?\\w*@(/test/);com$/)"`
		Age  int    `valid:"Required;Range(1, 140)"`
	}
	valid := Validation{}

	u := user{Name: "test@/test/;com", Age: 40}
	b, err := valid.Valid(u)
	if err != nil {
		t.Fatal(err)
	}
	if !b {
		t.Error("validation should be passed")
	}

	uptr := &user{Name: "test", Age: 40}
	valid.Clear()
	b, err = valid.Valid(uptr)
	if err != nil {
		t.Fatal(err)
	}
	if b {
		t.Error("validation should not be passed")
	}
	if len(valid.Errors) != 1 {
		t.Fatalf("valid errors len should be 1 but got %d", len(valid.Errors))
	}

	u = user{Name: "test@/test/;com", Age: 180}
	valid.Clear()
	b, err = valid.Valid(u)
	if err != nil {
		t.Fatal(err)
	}
	if b {
		t.Error("validation should not be passed")
	}
	if len(valid.Errors) != 1 {
		t.Fatalf("valid errors len should be 1 but got %d", len(valid.Errors))
	}
}

func TestRecursiveValid(t *testing.T) {
	type User struct {
		ID   int
		Name string `valid:"Required;Match(/^(test)?\\w*@(/test/);com$/)"`
		Age  int    `valid:"Required;Range(1, 140)"`
	}

	type AnonymouseUser struct {
		ID2   int
		Name2 string `valid:"Required;Match(/^(test)?\\w*@(/test/);com$/)"`
		Age2  int    `valid:"Required;Range(1, 140)"`
	}

	type Account struct {
		Password string `valid:"Required"`
		U        User
		AnonymouseUser
	}
	valid := Validation{}

	u := Account{Password: "abc123_", U: User{}}
	b, err := valid.RecursiveValid(u)
	if err != nil {
		t.Fatal(err)
	}
	if b {
		t.Error("validation should not be passed")
	}
}
