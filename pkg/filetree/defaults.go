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
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/alediaferia/prefixmap"
	"github.com/rs/zerolog/log"
)

type InMemoryFileSystem struct {
	disk *prefixmap.PrefixMap
}

func NewInMemoryFileSystem() *InMemoryFileSystem {
	return &InMemoryFileSystem{
		disk: prefixmap.New(),
	}
}

func (fsys *InMemoryFileSystem) ReadDir(name string) ([]fs.DirEntry, error) {
	result := []fs.DirEntry{}

	if !strings.HasSuffix(name, string(os.PathSeparator)) {
		name += string(os.PathSeparator)
	}

	for _, item := range fsys.disk.GetByPrefix(name) {
		file := item.(*File) //nolint:forcetypeassert
		if strings.Count(strings.TrimPrefix(file.path, name), string(os.PathSeparator)) == 0 {
			result = append(result, file)
		}
	}

	log.Trace().Str("param", name).Interface("result", result).Msg("InMemoryFileSystem - ReadDir")

	return result, nil
}

func (fsys *InMemoryFileSystem) WriteFile(name string, data []byte, perm fs.FileMode) error {
	var file *File

	files := fsys.disk.Get(name)
	if len(files) == 0 {
		file = NewFile(name, false, perm)

		for dir, _ := path.Split(name); len(dir) > 0; dir, _ = path.Split(dir) {
			dir = path.Clean(dir)
			if err := fsys.Mkdir(dir, fs.ModePerm); err != nil {
				return err
			}
		}
	} else {
		file = files[0].(*File) //nolint:forcetypeassert
	}

	if _, err := file.Write(data); err != nil {
		return err
	}

	fsys.disk.Replace(name, file)

	return nil
}

func (fsys *InMemoryFileSystem) Mkdir(name string, perm fs.FileMode) error {
	var file *File

	name = strings.TrimSuffix(name, string(os.PathSeparator))

	files := fsys.disk.Get(name)
	if len(files) == 0 {
		file = NewFile(name, true, perm)
	} else {
		file = files[0].(*File) //nolint:forcetypeassert
	}

	if !file.isDir {
		return fs.ErrExist
	}

	log.Trace().Str("param", name).Msg("InMemoryFileSystem - MkDir")

	fsys.disk.Replace(name, file)

	return nil
}

func (fsys *InMemoryFileSystem) Open(name string) (fs.File, error) {
	var file *File

	files := fsys.disk.Get(name)
	if len(files) == 0 {
		return nil, os.ErrNotExist
	}

	file = files[0].(*File) //nolint:forcetypeassert

	return file, nil
}
