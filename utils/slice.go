package utils

import (
	"fmt"
	"reflect"
)

// ExtractStringSliceFn v必须是slice
func ExtractStringSliceFn(v interface{}, fn func(it interface{}) string) []string {
	a := make([]string, 0)
	switch reflect.TypeOf(v).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(v)
		for i := 0; i < s.Len(); i++ {
			v := fn(s.Index(i).Interface())
			if v != "" {
				a = append(a, v)
			}
		}
	}
	return a
}

// ExtractInt64SliceFn v必须是slice
func ExtractInt64SliceFn(v interface{}, fn func(it interface{}) int64) []int64 {
	a := make([]int64, 0)
	switch reflect.TypeOf(v).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(v)
		for i := 0; i < s.Len(); i++ {
			v := fn(s.Index(i).Interface())
			if v > 0 {
				a = append(a, v)
			}
		}
	}
	return a
}

// UniqueI64 去重input
func UniqueI64(input []int64) (output []int64) {
	m := make(map[int64]bool)
	for _, v := range input {
		if !m[v] {
			output = append(output, v)
			m[v] = true
		}
	}
	return output
}

// UniqueI32 去重input
func UniqueI32(input []int32) (output []int32) {
	m := make(map[int32]bool)
	for _, v := range input {
		if !m[v] {
			output = append(output, v)
			m[v] = true
		}
	}
	return output
}

// UniqueStr 去重string
func UniqueStr(a []string) (unique []string) {
	m := make(map[string]bool)
	for _, item := range a {
		if _, ok := m[item]; !ok {
			unique = append(unique, item)
		}
		m[item] = true
	}
	return
}

// Contain 判断obj是否在target中，target支持的类型array,slice,map
func Contain(obj interface{}, targets interface{}) bool {
	if targets == nil {
		return false
	}
	targetValue := reflect.ValueOf(targets)
	switch reflect.TypeOf(targets).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				if IsDebug() {
					fmt.Printf("targetValue.Index(i).Interface() %v obj %v\n", targetValue.Index(i).Interface(), obj)
				}
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}

	return false
}
