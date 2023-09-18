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

package template

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"unicode"

	"github.com/Masterminds/sprig/v3"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func Generate(tmplstr string, stack []any) ([]byte, error) {
	funcmap := generateFuncMap()

	funcmap["Stack"] = generateStackFunc(stack)

	tmpl, err := template.New("template").Funcs(sprig.TxtFuncMap()).Funcs(funcmap).Parse(tmplstr)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	result := &bytes.Buffer{}
	err = tmpl.Execute(result, stack[0])

	return result.Bytes(), err
}

func generateFuncMap() template.FuncMap {
	funcMap := template.FuncMap{}

	funcMap["ToUpper"] = strings.ToUpper
	funcMap["ToLower"] = strings.ToLower
	funcMap["NoAccent"] = rmAcc

	return funcMap
}

func generateStackFunc(stack []any) func(index int) any {
	return func(index int) any {
		if index < 0 {
			return stack[len(stack)+index]
		}

		return stack[index]
	}
}

// rmAcc removes accents from string
// Function derived from: http://blog.golang.org/normalization
func rmAcc(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, s)

	return result
}
