package structutil

import (
	"reflect"
	"testing"
)

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

	res, err := GetByJson(target, "hoge")
	if err != nil {
		t.Errorf("Get(%+v, %s) returns err=%+v", target, "hoge", err)
	}
	if res.(string) != "hoge" {
		t.Errorf("Get(%+v, %s) returns wrong value=%+v", target, "hoge", res)
	}

	res, err = GetByJson(target, "foo.bar")
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

func TestGetValueByJson(t *testing.T) {
	type test struct {
		Hoge string `json:"hoge"`
		Foo  *test  `json:"foo"`
	}
	target := &test{}

	v, err := GetValueByJson(reflect.ValueOf(target), "foo.hoge")
	if err != nil {
		t.Errorf("GetValueByJson(%+v, %s) returns err=%+v", target, "foo.bar", err)
	}
	if v.Type().Kind() != reflect.String {
		t.Errorf("GetValueByJson(%+v, %s) returns wrong value: got=%v, want=%v", target, "foo.bar", v.Type().Kind(), reflect.String)
	}

	if err := SetByJson(target, "foo.hoge", "hage"); err != nil {
		t.Errorf("SetByJson(%+v, %s) returns err=%+v", target, "Hoge", err)
	}
	if target.Foo.Hoge != "hage" {
		t.Errorf("SetByJson(%+v, %s) set wrong value: got=%s, want=%s", target, "Hoge", target.Foo.Hoge, "hage")
	}
}

// Running tool: /usr/lib/go-1.13/bin/go test -benchmem -run=^$ n0st.ac/n0stack/n0core/pkg/util/struct -bench ^(BenchmarkGetByJson)$
//
// goos: linux
// goarch: amd64
// pkg: n0st.ac/n0stack/n0core/pkg/util/struct
// BenchmarkGetByJson-8   	 2012806	       543 ns/op	      72 B/op	       5 allocs/op
// PASS
// ok  	n0st.ac/n0stack/n0core/pkg/util/struct	2.800s
// Success: Benchmarks passed.
func BenchmarkGetByJson(b *testing.B) {
	type test struct {
		Hoge string `json:"hoge"`
		Foo  *test  `json:"foo"`
	}
	target := &test{Foo: &test{Hoge: "hoge"}}
	var tmp string

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		t, _ := GetByJson(target, "foo.hoge")
		tmp = t.(string)
	}
	b.StopTimer()

	b.Log(tmp)
}

// Running tool: /usr/lib/go-1.13/bin/go test -benchmem -run=^$ n0st.ac/n0stack/n0core/pkg/util/struct -bench ^(BenchmarkGetByOrigin)$
//
// goos: linux
// goarch: amd64
// pkg: n0st.ac/n0stack/n0core/pkg/util/struct
// BenchmarkGetByOrigin-8   	1000000000	         0.910 ns/op	       0 B/op	       0 allocs/op
// PASS
// ok  	n0st.ac/n0stack/n0core/pkg/util/struct	0.986s
// Success: Benchmarks passed.
func BenchmarkGetByOrigin(b *testing.B) {
	type test struct {
		Hoge string `json:"hoge"`
		Foo  *test  `json:"foo"`
	}
	target := &test{Foo: &test{Hoge: "hoge"}}
	var tmp string

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp = target.Foo.Hoge
	}
	b.StopTimer()

	b.Log(tmp)
}

// Running tool: /usr/lib/go-1.13/bin/go test -benchmem -run=^$ n0st.ac/n0stack/n0core/pkg/util/struct -bench ^(BenchmarkSetByJson)$
//
// goos: linux
// goarch: amd64
// pkg: n0st.ac/n0stack/n0core/pkg/util/struct
// BenchmarkSetByJson-8   	 2040865	       570 ns/op	      56 B/op	       4 allocs/op
// PASS
// ok  	n0st.ac/n0stack/n0core/pkg/util/struct	2.881s
// Success: Benchmarks passed.
func BenchmarkSetByJson(b *testing.B) {
	type test struct {
		Hoge string `json:"hoge"`
		Foo  *test  `json:"foo"`
	}
	target := &test{Foo: &test{Hoge: "hoge"}}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		SetByJson(target, "foo.hoge", "hage")
	}
}

// Running tool: /usr/lib/go-1.13/bin/go test -benchmem -run=^$ n0st.ac/n0stack/n0core/pkg/util/struct -bench ^(BenchmarkSetByOrigin)$
//
// goos: linux
// goarch: amd64
// pkg: n0st.ac/n0stack/n0core/pkg/util/struct
// BenchmarkSetByOrigin-8   	1000000000	         0.783 ns/op	       0 B/op	       0 allocs/op
// PASS
// ok  	n0st.ac/n0stack/n0core/pkg/util/struct	0.884s
// Success: Benchmarks passed.
func BenchmarkSetByOrigin(b *testing.B) {
	type test struct {
		Hoge string `json:"hoge"`
		Foo  *test  `json:"foo"`
	}
	target := &test{Foo: &test{Hoge: "hoge"}}

	for i := 0; i < b.N; i++ {
		target.Foo.Hoge = "hage"
	}
}
