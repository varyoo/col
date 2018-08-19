package col

import (
	"reflect"
	"testing"
)

type (
	Anon struct {
		NestedString string `col:"nested"`
	}

	Wrapper struct {
		Anon
		OuterString string `col:"outer"`
		OtherString string
	}
)

func TestWalk(t *testing.T) {
	s := Wrapper{Anon{"nested"}, "outer", "other"}
	var fs []field
	v := reflect.ValueOf(s)

	res := []string{"nested", "outer", "other"}

	fs = walkFields(fs, v, []int{})
	if len(fs) != len(res) {
		t.Fatal("Bad number of fields")
	}

	for i, f := range fs {
		v := v.FieldByIndex(f.Index)

		if have, _ := v.Interface().(string); have != res[i] {
			t.Errorf("Want %s but have %s", res[i], have)
		}
	}
}

type valueTest struct {
	Src    interface{}
	Values interface{}
	Names  []string
}

func TestValues(t *testing.T) {
	s := Wrapper{Anon{"nested"}, "outer", "other"}

	tests := []valueTest{
		{
			Src:    s,
			Values: []interface{}{"other", "nested"},
			Names:  []string{"OtherString", "nested"},
		},
		{
			Src:    s,
			Values: []interface{}{"other", "outer"},
			Names:  []string{"OtherString", "outer"},
		},
		{
			Src:    s,
			Values: []interface{}{"outer", "other"},
			Names:  []string{"outer", "OtherString"},
		},
		{
			Src:    s,
			Values: []interface{}{"nested", "other", "outer"},
			Names:  []string{"nested", "OtherString", "outer"},
		},
	}

	for i, tc := range tests {
		vs := Values(tc.Src, tc.Names...)

		if !reflect.DeepEqual(tc.Values, vs) {
			t.Errorf("Test %d failed", i)
		}
	}
}
