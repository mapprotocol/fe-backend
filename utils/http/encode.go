package http

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

func URLEncode(s interface{}) (string, error) {
	if s == nil {
		return "", errors.New("provided value is nil")
	}

	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return "", errors.New("provided value is not a struct")
	}

	typ := val.Type()
	urls := url.Values{}
	for i := 0; i < val.NumField(); i++ {
		if !val.Field(i).CanInterface() {
			continue
		}
		name := typ.Field(i).Name
		//urls.Add(typ.Field(i).Tag.Get("json"), fmt.Sprintf("%v", val.Field(i).Interface()))
		tag := typ.Field(i).Tag.Get("json")
		if tag != "" {
			index := strings.Index(tag, ",")
			if index == -1 {
				name = tag
			} else {
				name = tag[:index]
			}
		}
		urls.Set(name, fmt.Sprintf("%v", val.Field(i).Interface()))
	}
	return urls.Encode(), nil
}

//func URLEncode(params map[string]interface{}) (string, error) {
//	if params == nil {
//		return "", errors.New("provided value is nil")
//	}
//
//	urls := url.Values{}
//	for k, v := range params {
//		switch v := v.(type) {
//		case string:
//			urls.Set(k, v)
//		case int, int8, int16, int32, int64:
//			urls.Set(k, strconv.FormatInt(reflect.ValueOf(v).Int(), 10))
//		case uint, uint8, uint16, uint32, uint64:
//			urls.Set(k, strconv.FormatUint(reflect.ValueOf(v).Uint(), 10))
//		case float32, float64:
//			urls.Set(k, strconv.FormatFloat(reflect.ValueOf(v).Float(), 'f', -1, 64))
//		case bool:
//			urls.Set(k, strconv.FormatBool(v))
//		default:
//			return "", errors.New("unsupported type")
//		}
//	}
//	return urls.Encode(), nil
//}
