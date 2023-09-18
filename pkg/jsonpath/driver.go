// Copyright (C) 2023 CGI France
//
// This file is part of emporte-piece.
//
// Emporte-piece is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Emporte-piece is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with emporte-piece.  If not, see <http://www.gnu.org/licenses/>.

package jsonpath

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var patternStack = regexp.MustCompile(`^\$\[(-?\d+)\]$`)

func Get(path string, contexts ...any) ([]Result, error) {
	context := contexts[0]
	paths := strings.Split(path, ".")

	if paths[0][0] == '$' {
		context, contexts = refreshContext(paths, contexts)
		paths = paths[1:]
	}

	return get(context, paths, contexts)
}

func refreshContext(paths []string, contexts []any) (any, []any) {
	if patternStack.MatchString(paths[0]) {
		stackindexstr := patternStack.FindStringSubmatch(paths[0])[1]
		stackindex, _ := strconv.Atoi(stackindexstr)

		if stackindex < 0 {
			return contexts[len(contexts)+stackindex], contexts[:len(contexts)+stackindex+1]
		}

		return contexts[stackindex], contexts[:stackindex]
	}

	return contexts[len(contexts)-1], contexts[:len(contexts)-1]
}

func get(context any, paths []string, stack []any) ([]Result, error) {
	path := paths[0]

	if len(paths) == 1 {
		if obj, ok := context.(map[string]any); ok {
			return []Result{{Selected: obj[path], Stack: append(stack, obj[path])}}, nil
		}

		return nil, nil
	}

	if path == "[]" || path == "*" {
		if array, ok := context.([]map[string]any); ok {
			return getArrayMap(array, paths, stack)
		}

		if array, ok := context.([]any); ok {
			return getArrayAny(array, paths, stack)
		}
	}

	if obj, ok := context.(map[string]any); ok {
		return get(obj[path], paths[1:], append(stack, obj[path]))
	}

	return nil, nil
}

func getArrayMap(array []map[string]any, paths []string, stack []any) ([]Result, error) {
	allResults := []Result{}

	for _, item := range array {
		newstack := make([]any, len(stack)+1)
		copy(newstack, stack)
		newstack[len(stack)] = item

		results, err := get(item, paths[1:], newstack)
		if err != nil {
			return nil, err
		}

		allResults = append(allResults, results...)
	}

	return allResults, nil
}

func getArrayAny(array []any, paths []string, stack []any) ([]Result, error) {
	allResults := []Result{}

	for _, item := range array {
		newstack := make([]any, len(stack)+1)
		copy(newstack, stack)
		newstack[len(stack)] = item

		results, err := get(item, paths[1:], newstack)
		if err != nil {
			return nil, err
		}

		allResults = append(allResults, results...)
	}

	return allResults, nil
}

func Develop(template string, contexts ...any) ([]ResultString, error) {
	resultstrings := []ResultString{}
	path, pathBegin, pathEnd := extractPath(template)

	if len(path) == 0 {
		return []ResultString{
			{
				Selected: template,
				Stack:    contexts,
			},
		}, nil
	}

	results, err := Get(path, contexts...)
	if err != nil {
		return resultstrings, err
	}

	for _, result := range results {
		selectedString := strings.Builder{}
		selectedString.WriteString(template[0:pathBegin])
		selectedString.WriteString(toString(result.Selected))
		selectedString.WriteString(template[pathEnd:])
		resultstring := ResultString{
			Selected: selectedString.String(),
			Stack:    result.Stack,
		}
		resultstrings = append(resultstrings, resultstring)
	}

	return resultstrings, nil
}

//nolint:gomnd
func extractPath(template string) (string, int, int) {
	var step, pathBegin, pathEnd int

	path := strings.Builder{}

	for index, char := range template {
		switch step {
		case 0:
			if char == '{' {
				step = 1
				pathBegin = index
			}
		case 1:
			if char == '{' {
				step = 2
			} else {
				step = 0
			}
		case 2:
			if char == '}' {
				step = 3
			} else {
				path.WriteRune(char)
			}
		case 3:
			if char == '}' {
				pathEnd = index + 1

				return path.String(), pathBegin, pathEnd
			}

			path.WriteRune('}')
			path.WriteRune(char)

			step = 2
		}
	}

	return "", 0, 0
}

//nolint:gocyclop,cyclop
func toString(v any) string {
	switch vTyped := v.(type) {
	case string:
		return vTyped
	case fmt.Stringer:
		return vTyped.String()
	case int:
		return strconv.Itoa(vTyped)
	case int64:
		return strconv.FormatInt(vTyped, 10)
	case int32:
		return strconv.FormatInt(int64(vTyped), 10)
	case int16:
		return strconv.FormatInt(int64(vTyped), 10)
	case int8:
		return strconv.FormatInt(int64(vTyped), 10)
	case uint:
		return strconv.FormatUint(uint64(vTyped), 10)
	case uint64:
		return strconv.FormatUint(vTyped, 10)
	case uint32:
		return strconv.FormatUint(uint64(vTyped), 10)
	case uint16:
		return strconv.FormatUint(uint64(vTyped), 10)
	case uint8:
		return strconv.FormatUint(uint64(vTyped), 10)
	case bool:
		return strconv.FormatBool(vTyped)
	default:
		return fmt.Sprintf("%v", vTyped)
	}
}
