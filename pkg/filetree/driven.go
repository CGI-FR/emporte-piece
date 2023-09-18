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

import "io/fs"

type FileSystem interface {
	ReadDir(name string) ([]fs.DirEntry, error)
	WriteFile(name string, data []byte, perm fs.FileMode) error
	Mkdir(name string, perm fs.FileMode) error
	Open(name string) (fs.File, error)
}
