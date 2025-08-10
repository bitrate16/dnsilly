// dnsilly - dns automation utility
// Copyright (C) 2025  bitrate16 (bitrate16@gmail.com)
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package util

import (
	"reflect"
	"strconv"
)

// GPT-based code (I was lazy)
func SetDefaults(v interface{}) {
	val := reflect.ValueOf(v).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		tag := typ.Field(i).Tag.Get("default")

		if tag != "" && field.IsZero() {
			switch field.Kind() {
			case reflect.String:
				field.SetString(tag)
			case reflect.Int:
				if intValue, err := strconv.Atoi(tag); err == nil {
					field.SetInt(int64(intValue))
				}
			case reflect.Bool:
				field.SetBool(tag == "true")
			case reflect.Struct:
				SetDefaults(field.Addr().Interface())
			}
		}
	}
}
