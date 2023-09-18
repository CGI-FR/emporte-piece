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

package infra

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type ContextReaderJSONL struct {
	input *bufio.Scanner
	err   error
	value []byte
}

func NewContextReaderJSONL(input io.Reader) *ContextReaderJSONL {
	return &ContextReaderJSONL{
		input: bufio.NewScanner(input),
		err:   nil,
		value: nil,
	}
}

func NewContextReaderJSONLFromFile(filepath string) (*ContextReaderJSON, error) {
	input, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}

	return NewContextReaderJSON(input), nil
}

func (cr *ContextReaderJSONL) HasNext() bool {
	if !cr.input.Scan() {
		cr.err = cr.input.Err()

		return false
	}

	cr.value = cr.input.Bytes()

	return cr.err == nil && cr.value != nil
}

func (cr *ContextReaderJSONL) Next() (any, error) {
	if cr.err != nil {
		return nil, cr.err
	}

	context := make(map[string]any)

	if err := json.Unmarshal(cr.value, &context); err != nil {
		return nil, fmt.Errorf("error parsing JSON input: %w", err)
	}

	return context, nil
}
