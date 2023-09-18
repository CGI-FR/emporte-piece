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

package jsonpath_test

import (
	"fmt"
	"testing"

	"github.com/cgi-fr/emporte-piece/pkg/jsonpath"
	"github.com/stretchr/testify/assert"
)

func TestDevelopLen(t *testing.T) {
	t.Parallel()

	//nolint:lll
	testdatas := []struct {
		template    string
		context     any
		expectedLen int
	}{
		{"hello {{name}} !", map[string]any{"name": "world"}, 1},
		{"hello {{persons.[].name}} !", map[string]any{"persons": []map[string]any{{"name": "world"}, {"name": "space"}}}, 1},
		{"hello {{$.persons.[].name}} !", map[string]any{"persons": []map[string]any{{"name": "world"}, {"name": "space"}}}, 1},
		{"hello {{$[-1].persons.[].name}} !", map[string]any{"persons": []map[string]any{{"name": "world"}, {"name": "space"}}}, 1},
	}

	for _, td := range testdatas {
		t.Run(fmt.Sprintf("%v", td.template), func(t *testing.T) {
			t.Parallel()

			res, err := jsonpath.Develop("hello {{name}} !", map[string]any{"name": "world"})

			assert.NoError(t, err)
			assert.Len(t, res, 1)
		})
	}
}

//nolint:golint,revive,stylecheck
func TestDevelopContext(t *testing.T) {
	t.Parallel()

	column_1_name := "column_1"
	column_2_name := "column_2"
	column_3_name := "column_3"
	column_4_name := "column_4"

	column_1 := map[string]any{"name": column_1_name}
	column_2 := map[string]any{"name": column_2_name}
	column_3 := map[string]any{"name": column_3_name}
	column_4 := map[string]any{"name": column_4_name}

	table_1_columns := []map[string]any{column_1, column_2}
	table_2_columns := []map[string]any{column_3, column_4}

	table_1_name := "table_1"
	table_2_name := "table_2"

	table_1 := map[string]any{"name": table_1_name, "columns": table_1_columns}
	table_2 := map[string]any{"name": table_2_name, "columns": table_2_columns}

	tables := []map[string]any{table_1, table_2}

	root := map[string]any{"tables": tables}

	stack := []any{root, tables, table_2, table_2_name}

	res, err := jsonpath.Develop("{{$[-2].columns.[].name}}.txt", stack...)

	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, "column_3.txt", res[0].Selected)
	assert.Equal(t, "column_4.txt", res[1].Selected)
}
