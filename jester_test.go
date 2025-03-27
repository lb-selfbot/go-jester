package jester_test

import (
	"reflect"
	"testing"

	"github.com/goccy/go-json"
	"github.com/lb-selfbot/go-jester"
)

func TestJester(t *testing.T) {
	var err error

	js, err := jester.NewJson([]byte(`{
		"test": {
			"string_array": ["asdf", "ghjk", "zxcv"],
			"string_array_null": ["abc", null, "efg"],
			"array": [1, "2", 3],
			"arraywithsubs": [{"subkeyone": 1},
			{"subkeytwo": 2, "subkeythree": 3}],
			"int": 10,
			"float": 5.150,
			"string": "simplejson",
			"bool": true,
			"sub_obj": {"a": 1}
		}
	}`))
	if js == nil {
		t.Fatal("got nil")
	}
	if err != nil {
		t.Fatalf("got err %#v", err)
	}

	aws := js.Get("test").Get("arraywithsubs")
	if aws == nil {
		t.Fatal("got nil")
	}

	if got, _ := aws.Get(0, "subkeyone").Int(); got != 1 {
		t.Errorf("got %#v", got)
	}
	if got, _ := aws.Get(1, "subkeytwo").Int(); got != 2 {
		t.Errorf("got %#v", got)
	}
	if got, _ := aws.Get(1, "subkeythree").Int(); got != 3 {
		t.Errorf("got %#v", got)
	}

	if i, _ := js.Get("test").Get("int").Int(); i != 10 {
		t.Errorf("got %#v", i)
	}

	if f, _ := js.Get("test").Get("float").Float64(); f != 5.150 {
		t.Errorf("got %#v", f)
	}

	if s, _ := js.Get("test").Get("string").String(); s != "simplejson" {
		t.Errorf("got %#v", s)
	}

	if b, _ := js.Get("test").Get("bool").Bool(); b != true {
		t.Errorf("got %#v", b)
	}

	if mi := js.Get("test").Get("int").MustInt(); mi != 10 {
		t.Errorf("got %#v", mi)
	}

	if mi := js.Get("test").Get("missing_int").MustInt(5150); mi != 5150 {
		t.Errorf("got %#v", mi)
	}

	if s := js.Get("test").Get("string").MustString(); s != "simplejson" {
		t.Errorf("got %#v", s)
	}

	if s := js.Get("test").Get("missing_string").MustString("fyea"); s != "fyea" {
		t.Errorf("got %#v", s)
	}

	a := js.Get("test").Get("missing_array").MustSlice([]any{"1", 2, "3"})
	if !reflect.DeepEqual(a, []any{"1", 2, "3"}) {
		t.Errorf("got %#v", a)
	}

	msa := js.Get("test").Get("string_array").MustStringSlice()
	if !reflect.DeepEqual(msa, []string{"asdf", "ghjk", "zxcv"}) {
		t.Errorf("got %#v", msa)
	}

	msa = js.Get("test").Get("string_array").MustStringSlice([]string{"1", "2", "3"})
	if !reflect.DeepEqual(msa, []string{"asdf", "ghjk", "zxcv"}) {
		t.Errorf("got %#v", msa)
	}

	msa = js.Get("test").Get("missing_array").MustStringSlice([]string{"1", "2", "3"})
	if !reflect.DeepEqual(msa, []string{"1", "2", "3"}) {
		t.Errorf("got %#v", msa)
	}

	mm := js.Get("test").Get("missing_map").MustMap(map[string]any{"found": false})
	if !reflect.DeepEqual(mm, map[string]any{"found": false}) {
		t.Errorf("got %#v", mm)
	}

	sa, err := js.Get("test").Get("string_array").StringSlice()
	if err != nil {
		t.Fatalf("got err %#v", err)
	}
	if !reflect.DeepEqual(sa, []string{"asdf", "ghjk", "zxcv"}) {
		t.Errorf("got %#v", sa)
	}

	sa, err = js.Get("test").Get("string_array_null").StringSlice()
	if err != nil {
		t.Fatalf("got err %#v", err)
	}
	if !reflect.DeepEqual(sa, []string{"abc", "", "efg"}) {
		t.Errorf("got %#v", sa)
	}

	if s, _ := js.Get("test", "string").String(); s != "simplejson" {
		t.Errorf("got %#v", s)
	}

	if i, _ := js.Get("test", "int").Int(); i != 10 {
		t.Errorf("got %#v", i)
	}

	if b := js.Get("test").Get("bool").MustBool(); b != true {
		t.Errorf("got %#v", b)
	}

	js.Set("float2", 300.0)
	if f := js.Get("float2").MustFloat64(); f != 300.0 {
		t.Errorf("got %#v", f)
	}

	js.Set("test2", "setTest")
	if s := js.Get("test2").MustString(); s != "setTest" {
		t.Errorf("got %#v", s)
	}

	js.Delete("test2")
	if s := js.Get("test2").MustString(); s == "setTest" {
		t.Errorf("got %#v", s)
	}

	js.Get("test").Get("sub_obj").Set("a", 2)
	if i := js.Get("test").Get("sub_obj").Get("a").MustInt(); i != 2 {
		t.Errorf("got %#v", i)
	}

	js.Get("test", "sub_obj").Set("a", 3)
	if i := js.Get("test").Get("sub_obj").Get("a").MustInt(); i != 3 {
		t.Errorf("got %#v", i)
	}
}

func TestStdlibInterfaces(t *testing.T) {
	val := new(struct {
		Name   string       `json:"name"`
		Params *jester.Data `json:"params"`
	})
	val2 := new(struct {
		Name   string       `json:"name"`
		Params *jester.Data `json:"params"`
	})

	raw := `{"name":"myobject","params":{"string":"simplejson"}}`

	if err := json.Unmarshal([]byte(raw), val); err != nil {
		t.Fatalf("err %#v", err)
	}
	if val.Name != "myobject" {
		t.Errorf("got %#v", val.Name)
	}
	if val.Params.Interface() == nil {
		t.Errorf("got %#v", val.Params.Interface())
	}
	if s, _ := val.Params.Get("string").String(); s != "simplejson" {
		t.Errorf("got %#v", s)
	}

	p, err := json.Marshal(val)
	if err != nil {
		t.Fatalf("err %#v", err)
	}
	if err = json.Unmarshal(p, val2); err != nil {
		t.Fatalf("err %#v", err)
	}
	if !reflect.DeepEqual(val, val2) { // stable
		t.Errorf("got %#v expected %#v", val2, val)
	}
}

func TestSet(t *testing.T) {
	js, err := jester.NewJson([]byte(`{}`))
	if err != nil {
		t.Fatalf("err %#v", err)
	}

	js.Set("baz", "bing")

	s, err := js.Get("baz").String()
	if err != nil {
		t.Fatalf("err %#v", err)
	}
	if s != "bing" {
		t.Errorf("got %#v", s)
	}
}

func TestSetPath(t *testing.T) {
	js, err := jester.NewJson([]byte(`{}`))
	if err != nil {
		t.Fatalf("err %#v", err)
	}

	js.SetPath([]any{"foo", "bar"}, "baz")

	s, err := js.Get("foo", "bar").String()
	if err != nil {
		t.Fatalf("err %#v", err)
	}
	if s != "baz" {
		t.Errorf("got %#v", s)
	}
}

func TestSetPathNoPath(t *testing.T) {
	js, err := jester.NewJson([]byte(`{"some":"data","some_number":1.0,"some_bool":false}`))
	if err != nil {
		t.Fatalf("err %#v", err)
	}

	if f := js.Get("some_number").MustFloat64(99.0); f != 1.0 {
		t.Errorf("got %#v", f)
	}

	js.SetPath([]any{}, map[string]any{"foo": "bar"})

	s, err := js.Get("foo").String()
	if err != nil {
		t.Fatalf("err %#v", err)
	}
	if s != "bar" {
		t.Errorf("got %#v", s)
	}

	if f := js.Get("some_number").MustFloat64(99.0); f != 99.0 {
		t.Errorf("got %#v", f)
	}
}

func TestPathWillAugmentExisting(t *testing.T) {
	js, err := jester.NewJson([]byte(`{"this":{"a":"aa","b":"bb","c":"cc"}}`))
	if err != nil {
		t.Fatalf("err %#v", err)
	}

	js.SetPath([]any{"this", "d"}, "dd")

	cases := []struct {
		path    []any
		outcome string
	}{
		{
			path:    []any{"this", "a"},
			outcome: "aa",
		},
		{
			path:    []any{"this", "b"},
			outcome: "bb",
		},
		{
			path:    []any{"this", "c"},
			outcome: "cc",
		},
		{
			path:    []any{"this", "d"},
			outcome: "dd",
		},
	}

	for _, tc := range cases {
		s, err := js.Get(tc.path...).String()
		if err != nil {
			t.Fatalf("err %#v", err)
		}
		if s != tc.outcome {
			t.Errorf("got %#v expected %#v", s, tc.outcome)
		}
	}
}

func TestPathWillOverwriteExisting(t *testing.T) {
	// notice how "a" is 0.1 - but then we'll try to set at path a, foo
	js, err := jester.NewJson([]byte(`{"this":{"a":0.1,"b":"bb","c":"cc"}}`))
	if err != nil {
		t.Fatalf("err %#v", err)
	}

	js.SetPath([]any{"this", "a", "foo"}, "bar")

	s, err := js.Get("this", "a", "foo").String()
	if err != nil {
		t.Fatalf("err %#v", err)
	}
	if s != "bar" {
		t.Errorf("got %#v", s)
	}
}

func TestLen(t *testing.T) {
	// Test Len with map
	js, err := jester.NewJson([]byte(`{"a":1,"b":2,"c":3}`))
	if err != nil {
		t.Fatalf("err %#v", err)
	}
	if l := js.Len(); l != 3 {
		t.Errorf("map length: got %d, expected 3", l)
	}

	// Test Len with array/slice
	js, err = jester.NewJson([]byte(`[1,2,3,4,5]`))
	if err != nil {
		t.Fatalf("err %#v", err)
	}
	if l := js.Len(); l != 5 {
		t.Errorf("slice length: got %d, expected 5", l)
	}

	// Test Len with string
	js = jester.New("test string")
	if l := js.Len(); l != 11 {
		t.Errorf("string length: got %d, expected 11", l)
	}

	// Test Len with byte slice
	js = jester.New([]byte("test bytes"))
	if l := js.Len(); l != 10 {
		t.Errorf("bytes length: got %d, expected 10", l)
	}

	// Test Len with empty map
	js = jester.NewEmpty()
	if l := js.Len(); l != 0 {
		t.Errorf("empty map length: got %d, expected 0", l)
	}

	// Test Len with nil
	js = jester.New(nil)
	if l := js.Len(); l != 0 {
		t.Errorf("nil length: got %d, expected 0", l)
	}

	// Test Len with nested structure
	js, err = jester.NewJson([]byte(`{"nested":{"a":1,"b":2},"array":[1,2,3]}`))
	if err != nil {
		t.Fatalf("err %#v", err)
	}
	if l := js.Len(); l != 2 {
		t.Errorf("nested structure length: got %d, expected 2", l)
	}
	if l := js.Get("nested").Len(); l != 2 {
		t.Errorf("nested map length: got %d, expected 2", l)
	}
	if l := js.Get("array").Len(); l != 3 {
		t.Errorf("nested array length: got %d, expected 3", l)
	}
}

func TestIterator(t *testing.T) {
	// Test slice iteration
	js, err := jester.NewJson([]byte(`[1, 2, 3, 4, 5]`))
	if err != nil {
		t.Fatalf("err %#v", err)
	}

	// Collect all elements using iterator
	var collected []int
	for v := range js.Iterator() {
		val, err := v.Int()
		if err != nil {
			t.Fatalf("err %#v", err)
		}
		collected = append(collected, val)
	}

	// Verify collected values
	expected := []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(collected, expected) {
		t.Errorf("slice iteration: got %#v, expected %#v", collected, expected)
	}

	// Test mixed type slice
	js, err = jester.NewJson([]byte(`[1, "test", true, 4.5]`))
	if err != nil {
		t.Fatalf("err %#v", err)
	}

	var types []string
	for v := range js.Iterator() {
		switch {
		case v.MustInt(0) != 0:
			types = append(types, "int")
		case v.MustString("") != "":
			types = append(types, "string")
		case v.MustBool(false):
			types = append(types, "bool")
		case v.MustFloat64(0) != 0:
			types = append(types, "float")
		default:
			types = append(types, "unknown")
		}
	}

	expectedTypes := []string{"int", "string", "bool", "float"}
	if !reflect.DeepEqual(types, expectedTypes) {
		t.Errorf("mixed type iteration: got %#v, expected %#v", types, expectedTypes)
	}

	// Test empty data
	js = jester.New(nil)
	count := 0
	for range js.Iterator() {
		count++
	}
	if count != 0 {
		t.Errorf("nil data iteration should yield 0 values, got %d", count)
	}

	// Test nested structure access via iteration
	js, err = jester.NewJson([]byte(`[{"name": "item1", "value": 10}, {"name": "item2", "value": 20}]`))
	if err != nil {
		t.Fatalf("err %#v", err)
	}

	var items []struct {
		Name  string
		Value int
	}

	for item := range js.Iterator() {
		name, err := item.Get("name").String()
		if err != nil {
			t.Fatalf("err %#v", err)
		}

		value, err := item.Get("value").Int()
		if err != nil {
			t.Fatalf("err %#v", err)
		}

		items = append(items, struct {
			Name  string
			Value int
		}{Name: name, Value: value})
	}

	expectedItems := []struct {
		Name  string
		Value int
	}{
		{Name: "item1", Value: 10},
		{Name: "item2", Value: 20},
	}

	if !reflect.DeepEqual(items, expectedItems) {
		t.Errorf("nested object iteration: got %#v, expected %#v", items, expectedItems)
	}
}
