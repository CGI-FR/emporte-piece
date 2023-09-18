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
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type ContextReaderYAML struct {
	read  bool
	input io.Reader
}

func NewContextReaderYAML(input io.Reader) *ContextReaderYAML {
	return &ContextReaderYAML{read: false, input: input}
}

func NewContextReaderYAMLFromFile(filepath string) (*ContextReaderYAML, error) {
	input, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}

	return NewContextReaderYAML(input), nil
}

func (cr *ContextReaderYAML) HasNext() bool {
	return !cr.read
}

func (cr *ContextReaderYAML) Next() (any, error) {
	if cr.read {
		return nil, ErrContextReaderEmpty
	}

	cr.read = true

	bytes, err := io.ReadAll(cr.input)
	if err != nil {
		return nil, fmt.Errorf("error reading JSON input: %w", err)
	}

	context := make(map[string]any)

	if err := yaml.Unmarshal(bytes, &context); err != nil {
		return nil, fmt.Errorf("error parsing JSON input: %w", err)
	}

	return context, nil
}
