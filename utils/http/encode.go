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
