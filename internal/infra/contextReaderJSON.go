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
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type ContextReaderJSON struct {
	read  bool
	input io.Reader
}

func NewContextReaderJSON(input io.Reader) *ContextReaderJSON {
	return &ContextReaderJSON{read: false, input: input}
}

func NewContextReaderJSONFromFile(filepath string) (*ContextReaderJSON, error) {
	input, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}

	return NewContextReaderJSON(input), nil
}

func (cr *ContextReaderJSON) HasNext() bool {
	return !cr.read
}

func (cr *ContextReaderJSON) Next() (any, error) {
	if cr.read {
		return nil, ErrContextReaderEmpty
	}

	bytes, err := io.ReadAll(cr.input)
	if err != nil {
		return nil, fmt.Errorf("error reading JSON input: %w", err)
	}

	context := make(map[string]any)

	if err := json.Unmarshal(bytes, &context); err != nil {
		return nil, fmt.Errorf("error parsing JSON input: %w", err)
	}

	return context, nil
}
