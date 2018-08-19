package col

import (
	"fmt"
	"reflect"
	"sort"
)

type field struct {
	Index []int
	Name  string
}

type fields []field

func (s fields) Len() int {
	return len(s)
}
func (s fields) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s fields) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

func walkFields(fs fields, v reflect.Value, index []int) fields {
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fieldValue := v.FieldByIndex(f.Index)

		index := append(index, i)

		if f.Anonymous {
			fs = walkFields(fs, fieldValue, index)
		} else {
			name := f.Tag.Get("col")
			if name == "" {
				name = f.Name
			}

			fs = append(fs, field{index, name})
		}
	}

	return fs
}

type Index [][]int

type column struct {
	Position int
	Name     string
}

type columns []column

func (s columns) Len() int {
	return len(s)
}
func (s columns) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s columns) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

func NewIndex(typ interface{}, names ...string) Index {
	if len(names) == 0 {
		return make(Index, 0)
	}

	cols := make(columns, 0, len(names))

	for i, n := range names {
		cols = append(cols, column{i, n})
	}

	index := make(Index, len(names))

	v := reflect.ValueOf(typ)

	fs := walkFields(make(fields, 0, v.NumField()), v, []int{})

	sort.Sort(fs)
	sort.Sort(cols)

	for _, f := range fs {
		if col := cols[0]; f.Name == col.Name {
			index[col.Position] = f.Index

			// try to match next column
			cols = cols[1:]
		}

		if len(cols) == 0 {
			// all done
			return index
		}
	}

	panic(fmt.Sprintf("can't find field %s", cols[0].Name))
}

func (index Index) Values(src interface{}) []interface{} {
	v := reflect.ValueOf(src)

	res := make([]interface{}, 0, len(index))

	for _, fieldIndex := range index {
		fieldValue := v.FieldByIndex(fieldIndex)
		res = append(res, fieldValue.Interface())
	}

	return res
}

func Values(src interface{}, names ...string) []interface{} {
	return NewIndex(src, names...).Values(src)
}
