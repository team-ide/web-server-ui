package servers

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestType(t *testing.T) {
	var a *int
	aType := reflect.TypeOf(a)

	fmt.Println(aType.Kind())
	fmt.Println(aType.Elem().Kind())
	aV := reflect.New(aType.Elem())
	fmt.Println(aV.Kind())
	fmt.Println(aV.Elem().Kind())
	var aVV = new(int)
	*aVV = 11
	aV.Elem().Set(reflect.ValueOf(11))
	fmt.Println(aType.Kind())
	fmt.Println(aV.Interface())
	fmt.Println(aV.Elem().Interface())

	_ = json.Unmarshal([]byte("22"), &a)
	fmt.Println(a)
	fmt.Println(*a)

}
