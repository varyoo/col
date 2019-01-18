package col

import (
	"reflect"
)

func (index Index) Pointers(dest interface{}) []interface{} {
	v := reflect.ValueOf(dest)

	if v.Kind() != reflect.Ptr {
		panic("destination struct is not a pointer")
	}

	v = v.Elem()
	ret := make([]interface{}, 0, len(index))

	for _, fieldIndex := range index {
		fieldValue := v.FieldByIndex(fieldIndex)
		ret = append(ret, fieldValue.Addr().Interface())
	}

	return ret
}

type Row interface {
	Scan(dest ...interface{}) error
}

// Usage :
//  type Tweet struct {
//  	Username string `col:"tweets.user_name"`
//  	Text string `col:"tweets.text"`
//  }
//
//  scanner := col.NewIndex(Tweet{}, "tweets.user_name", "tweets.text")
//  rows, _ := squirrel.Select("tweets.user_name", "tweets.text").From("tweets").Query()
//
//  for rows.Next() {
//  	tweet := Tweet{}
//
//  	scanner.Scan(rows, &tweet)
//  }
func (index Index) Scan(row Row, dest interface{}) error {
	return row.Scan(index.Pointers(dest)...)
}
