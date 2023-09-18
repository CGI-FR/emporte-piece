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
	"io/fs"
	"os"
)

type FileSystem struct{}

func (fsys FileSystem) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(name) //nolint:wrapcheck
}

func (fsys FileSystem) WriteFile(name string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(name, data, perm) //nolint:wrapcheck
}

func (fsys FileSystem) Mkdir(name string, perm fs.FileMode) error {
	return os.MkdirAll(name, perm) //nolint:wrapcheck
}

func (fsys FileSystem) Open(name string) (fs.File, error) {
	return os.Open(name) //nolint:wrapcheck
}
