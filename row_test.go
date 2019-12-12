package reporter

import (
	"reflect"
	"testing"
	"time"
)

type exampleRow struct {
	StrEx      string     `csv:"string"`
	TimeEx     time.Time  `csv:"time"`
	FloatEx    float64    `csv:"float"`
	BoolEx     bool       `csv:"bool"`
	StructEx   forNesting `csv:"struct"`
	StrSliceEx []string   `csv:"str_slice"`
	IntSliceEx []int      `csv:"int_slice"`
	FlSliceEx  []float64  `csv:"float_slice"`
	IgnoreMe   string     `csv:"-"`
	PtrEx      *string    `csv:"pointer"`
}

type forNesting struct {
	AnotherStr   string         `csv:"street"`
	IntEx        int            `csv:"city"`
	NestedStruct []*alsoNesting `csv:"nested"`
}

type alsoNesting struct {
	Stuff []int `csv:"stuff"`
}

func TestCreateHeader(t *testing.T) {
	expected := []string{"string", "time", "float", "bool", "struct", "str_slice", "int_slice", "float_slice", "pointer"}
	actual := CreateHeader(&exampleRow{})

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Want %#v, got %#v", expected, actual)
	}
}

func TestMarshalCSV(t *testing.T) {
	str := "ptr to string"
	var structList []*alsoNesting
	structList = append(structList, &alsoNesting{Stuff: []int{0, 2}})
	structList = append(structList, &alsoNesting{Stuff: []int{3, 4}})

	r := &exampleRow{
		StrEx:   "I'm a string",
		TimeEx:  time.Date(1978, 04, 30, 0, 0, 0, 0, time.UTC),
		FloatEx: 39.390293,
		BoolEx:  false,
		StructEx: forNesting{
			AnotherStr:   "I'm another string",
			IntEx:        3,
			NestedStruct: structList,
		},
		StrSliceEx: []string{"first thing", "second thing"},
		IntSliceEx: []int{1, 2, 3},
		FlSliceEx:  []float64{0.32903, 300.93},
		IgnoreMe:   "I'd better not end up in the result",
		PtrEx:      &str,
	}
	expected := []string{"I'm a string", "1978-04-30T00:00:00Z", "39.390293", "false", "I'm another string,3,0,2,3,4", "first thing,second thing", "1,2,3", "0.32903,300.93", "ptr to string"}

	actual := MarshalCSV(r)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("\nWant %#v\ngot  %#v", expected, actual)
	}
}

func TestMarshalCSV_nonStructKinds(t *testing.T) {
	s1 := "hi"
	s2 := "there"
	slicePtr := []*string{&s1, &s2}

	var v interface{}

	m := make(map[interface{}]interface{})
	m["key"] = true
	m[-9032.3] = []string{"hi", "there"}

	var interfaceSlice []interface{}
	interfaceSlice = append(interfaceSlice, 90)
	interfaceSlice = append(interfaceSlice, "abc")
	interfaceSlice = append(interfaceSlice, []byte{'a', '#'})

	var u uint16 = 1<<16 - 1

	tests := []struct {
		desc string
		in   interface{}
		want []string
	}{
		{"empty interface", v, []string{""}},
		{"[]*string", slicePtr, []string{"hi,there"}},
		{"[]bool", []bool{true, true, false}, []string{"true,true,false"}},
		{"[]float64", []float64{0.3902, 59039.32}, []string{"0.3902,59039.32"}},
		{"complex number", complex(4, 3), []string{"(4+3i)"}},
		{"map[interface{}]interface{}", m, []string{"key:true,-9032.3:hi,there"}},
		{"[]interface{}", interfaceSlice, []string{"90,abc,97,35"}},
		{"[]interface{}", u, []string{"65535"}},
		{"int", 0, []string{"0"}},
	}

	for _, test := range tests {
		actual := MarshalCSV(test.in)
		if !reflect.DeepEqual(test.want, actual) {
			t.Errorf("%s - Want %q, got %q", test.desc, test.want, actual)
		}
	}
}
