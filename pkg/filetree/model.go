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

package filetree

import (
	"bytes"
	"io/fs"
	"path"
)

type File struct {
	path    string
	isDir   bool
	mode    fs.FileMode
	content *bytes.Buffer
}

func NewFile(path string, isDir bool, mode fs.FileMode) *File {
	return &File{
		path:    path,
		isDir:   isDir,
		mode:    mode,
		content: new(bytes.Buffer),
	}
}

func (f *File) Name() string {
	return path.Base(f.path)
}

func (f *File) IsDir() bool {
	return f.isDir
}

func (f *File) Type() fs.FileMode {
	return f.mode
}

func (f *File) Info() (fs.FileInfo, error) {
	return nil, nil
}

func (f *File) Read(p []byte) (int, error) {
	return f.content.Read(p) //nolint:wrapcheck
}

func (f *File) Write(p []byte) (int, error) {
	return f.content.Write(p) //nolint:wrapcheck
}

func (f *File) Stat() (fs.FileInfo, error) {
	return nil, nil
}

func (f *File) Close() error {
	return nil
}
