package structutil

import "testing"

func TestGet(t *testing.T) {
	type test struct {
		Hoge string
		Foo  struct {
			Bar string
		}
	}

	target := &test{
		Hoge: "hoge",
	}
	target.Foo.Bar = "bar"

	res, err := Get(target, "Hoge")
	if err != nil {
		t.Errorf("Get(%+v, %s) returns err=%+v", target, "Hoge", err)
	}
	if res.(string) != "hoge" {
		t.Errorf("Get(%+v, %s) returns wrong value=%+v", target, "hoge", res)
	}

	res, err = Get(target, "Foo.Bar")
	if err != nil {
		t.Errorf("Get(%+v, %s) returns err=%+v", target, "Foo.Bar", err)
	}
	if res.(string) != "bar" {
		t.Errorf("Get(%+v, %s) returns wrong value=%+v", target, "bar", res)
	}
}

func TestGetByJsonTag(t *testing.T) {
	type test struct {
		Hoge string `json:"hoge"`
		Foo  struct {
			Bar string `json:"bar"`
		} `json:"foo"`
	}

	target := &test{
		Hoge: "hoge",
	}
	target.Foo.Bar = "bar"

	res, err := GetByJsonTag(target, "hoge")
	if err != nil {
		t.Errorf("Get(%+v, %s) returns err=%+v", target, "hoge", err)
	}
	if res.(string) != "hoge" {
		t.Errorf("Get(%+v, %s) returns wrong value=%+v", target, "hoge", res)
	}

	res, err = GetByJsonTag(target, "foo.bar")
	if err != nil {
		t.Errorf("Get(%+v, %s) returns err=%+v", target, "foo.bar", err)
	}
	if res.(string) != "bar" {
		t.Errorf("Get(%+v, %s) returns wrong value=%+v", target, "bar", res)
	}
}

func TestSet(t *testing.T) {
	type test struct {
		Hoge string `json:"hoge"`
		Foo  struct {
			Bar string `json:"bar"`
		} `json:"foo"`
	}

	target := &test{
		Hoge: "hoge",
	}
	target.Foo.Bar = "bar"

	err := Set(target, "Hoge", "hage")
	if err != nil {
		t.Errorf("Get(%+v, %s, %s) returns err=%+v", target, "Hoge", "hage", err)
	}
	if target.Hoge != "hage" {
		t.Errorf("Set(%+v, %s, %s) returns wrong value: got=%s, want=%s", target, "Hoge", "hage", target.Hoge, "hage")
	}

	err = Set(target, "Foo.Bar", "baz")
	if err != nil {
		t.Errorf("Get(%+v, %s, %s) returns err=%+v", target, "Foo.Bar", "baz", err)
	}
	if target.Foo.Bar != "baz" {
		t.Errorf("Get(%+v, %s, %s) returns wrong value: got=%s, want=%s", target, "Foo.Bar", "baz", target.Foo.Bar, "baz")
	}
}

func TestSetByJson(t *testing.T) {
	type test struct {
		Hoge string `json:"hoge"`
		Foo  struct {
			Bar string `json:"bar"`
		} `json:"foo"`
	}

	target := &test{
		Hoge: "hoge",
	}
	target.Foo.Bar = "bar"

	err := SetByJson(target, "hoge", "hage")
	if err != nil {
		t.Errorf("Get(%+v, %s, %s) returns err=%+v", target, "hoge", "hage", err)
	}
	if target.Hoge != "hage" {
		t.Errorf("Set(%+v, %s, %s) returns wrong value: got=%s, want=%s", target, "hoge", "hage", target.Hoge, "hage")
	}

	err = SetByJson(target, "foo.bar", "baz")
	if err != nil {
		t.Errorf("Get(%+v, %s, %s) returns err=%+v", target, "foo.bar", "baz", err)
	}
	if target.Foo.Bar != "baz" {
		t.Errorf("Get(%+v, %s, %s) returns wrong value: got=%s, want=%s", target, "foo.bar", "baz", target.Foo.Bar, "baz")
	}
}

func TestUpdateWithMaskUsingJson(t *testing.T) {
	type test struct {
		Hoge string `json:"hoge"`
		Foo  struct {
			Bar string `json:"bar"`
		} `json:"foo"`
	}

	target := &test{
		Hoge: "hoge",
	}
	target.Foo.Bar = "bar"

	source := &test{
		Hoge: "hage",
	}
	source.Foo.Bar = "baz"

	if err := UpdateWithMaskUsingJson(target, source, []string{"hoge"}); err != nil {
		t.Errorf("UpdateWithMaskUsingJson(%+v, %+v, %+v) returns err=%+v", target, source, []string{"hoge"}, err)
	}
	if target.Hoge != "hage" {
		t.Errorf("UpdateWithMaskUsingJson(%+v, %+v, %+v) does not update target: got=%s, want=%s", target, source, []string{"hoge"}, target.Hoge, "hage")
	}
	if target.Foo.Bar != "bar" {
		t.Errorf("UpdateWithMaskUsingJson(%+v, %+v, %+v) update target unexpectedly: got=%s, want=%s", target, source, []string{"hoge"}, target.Foo.Bar, "bar")
	}
}
