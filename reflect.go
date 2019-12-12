package reporter

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type reflection struct {
	Type  reflect.Type
	Kind  reflect.Kind
	Value reflect.Value
}

func createReflection(v interface{}) *reflection {
	refVal := reflect.ValueOf(v)
	if refVal.Kind() == reflect.Ptr || refVal.Kind() == reflect.Interface {
		refVal = refVal.Elem() // get the Value contained in an pointer or interface
	}
	var t reflect.Type
	if refVal.IsValid() {
		t = refVal.Type()
	}
	return &reflection{
		Type:  t,
		Kind:  refVal.Kind(),
		Value: refVal,
	}
}

func (ref *reflection) getStructFieldAtIdx(idx int) *reflection {
	structField := ref.Value.Field(idx)
	return createReflection(structField.Interface())
}

func (ref *reflection) tryCommonTypes() string {
	if ref.Type.Name() == reflect.TypeOf(time.Time{}).Name() {
		return reflectTimeToString(&ref.Value)
	}
	return ""
}

// TimeFormatString defaults to "2006-01-02T15:04:05Z07:00", but any valid time.Time format string can be used
var TimeFormatString = "2006-01-02T15:04:05Z07:00"

func reflectTimeToString(structField *reflect.Value) string {
	formatMethod := structField.MethodByName("Format")
	if !formatMethod.IsZero() {
		args := []reflect.Value{reflect.ValueOf(TimeFormatString)}
		timeVal := formatMethod.Call(args)
		if len(timeVal) > 0 {
			return timeVal[0].String()
		}
	}
	return ""
}

func (ref *reflection) getUnderlyingData() string {
	switch ref.Kind {
	case reflect.Struct:
		return strings.Join(MarshalCSV(ref.Value.Interface()), ",")
	case reflect.Slice:
		strVal := ""
		sliceRef := createReflection(ref.Value.Interface())
		if sliceRef.Value.Len() > 0 {
			curRef := createReflection(sliceRef.Value.Index(0).Interface())
			strVal = curRef.getUnderlyingData()
			for i := 1; i < sliceRef.Value.Len(); i++ {
				curRef = createReflection(sliceRef.Value.Index(i).Interface())
				strVal += fmt.Sprintf(",%v", curRef.getUnderlyingData())
			}
		}
		return strVal
	case reflect.Map:
		iter := ref.Value.MapRange()
		strVal := ""
		for iter.Next() {
			keyRef := createReflection(iter.Key().Interface())
			valRef := createReflection(iter.Value().Interface())
			strVal += fmt.Sprintf(",%v:%v", keyRef.getUnderlyingData(), valRef.getUnderlyingData())
		}
		if string(strVal[0]) == "," {
			strVal = strVal[1:]
		}
		return strVal
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprint(ref.Value.Interface())
	default:
		if !ref.Value.IsValid() {
			return ""
		}
		return fmt.Sprintf("%v", ref.Value.Interface())
	}
}
