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
	"fmt"
	"io"
	"os"
	"path"

	"github.com/cgi-fr/emporte-piece/pkg/jsonpath"
	"github.com/cgi-fr/emporte-piece/pkg/template"
	"github.com/rs/zerolog/log"
)

type Driver struct {
	fs FileSystem
}

func NewDriver(fsys FileSystem) Driver {
	return Driver{
		fs: fsys,
	}
}

func (d Driver) Develop(templatePath string, targetPath string, contexts ...any) error {
	files, _ := d.fs.ReadDir(templatePath)
	for _, file := range files {
		rs, err := jsonpath.Develop(file.Name(), contexts...)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		for _, devpath := range rs {
			subTemplatePath := path.Join(templatePath, file.Name())
			subTargetPath := path.Join(targetPath, devpath.Selected)

			log.Info().Str("from", subTemplatePath).Msg("generating " + subTargetPath)

			if file.IsDir() {
				if err := d.developDir(subTargetPath, subTemplatePath, devpath); err != nil {
					return err
				}
			} else {
				if err := d.developFile(subTargetPath, subTemplatePath, devpath); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (d Driver) developFile(subTargetPath string, subTemplatePath string, devpath jsonpath.ResultString) error {
	tmplFile, err := d.fs.Open(subTemplatePath)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	tmplContent, err := io.ReadAll(tmplFile)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	content, err := template.Generate(string(tmplContent), devpath.Stack)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if err := d.fs.WriteFile(subTargetPath, content, os.ModePerm); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (d Driver) developDir(subTargetPath string, subTemplatePath string, devpath jsonpath.ResultString) error {
	if err := d.fs.Mkdir(subTargetPath, os.ModePerm); err != nil && !os.IsExist(err) {
		return fmt.Errorf("%w", err)
	} else if err := d.Develop(subTemplatePath, subTargetPath, devpath.Stack...); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
