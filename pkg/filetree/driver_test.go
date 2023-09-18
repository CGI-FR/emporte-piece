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

package filetree_test

import (
	"io"
	"os"
	"testing"

	"github.com/cgi-fr/emporte-piece/pkg/filetree"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestMain(m *testing.M) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: true}) //nolint:exhaustruct

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	code := m.Run()

	os.Exit(code)
}

func TestDevelopContext(t *testing.T) {
	t.Parallel()

	fsys := filetree.NewInMemoryFileSystem()

	content := `hello {{$person := Stack -2}}{{$person.name}}`

	assert.NoError(t, fsys.Mkdir("template", os.ModePerm))
	assert.NoError(t, fsys.WriteFile("template/{{persons.[].name}}.txt", []byte(content), os.ModePerm))
	assert.NoError(t, fsys.Mkdir("result", os.ModePerm))

	context := map[string]any{
		"persons": []map[string]any{
			{
				"name": "John",
			},
		},
	}

	driver := filetree.NewDriver(fsys)

	assert.NoError(t, driver.Develop("template", "result", context))

	f, err := fsys.Open("result/John.txt")
	assert.NoError(t, err)

	b, err := io.ReadAll(f)
	assert.NoError(t, err)

	assert.Equal(t, "hello John", string(b))
}

//nolint:lll,funlen
func TestStack(t *testing.T) {
	t.Parallel()

	context := map[string]any{}

	yamlb, _ := os.ReadFile("testdata/context_1.yaml")
	_ = yaml.Unmarshal(yamlb, context)

	content0 := `{{Stack 0}}`
	content1 := `{{Stack 1}}`
	content2 := `{{Stack 2}}`
	content3 := `{{Stack 3}}`
	content4 := `{{Stack 4}}`
	content5 := `{{Stack 5}}`
	content6 := `{{Stack 6}}`
	content7 := `{{Stack 7}}`
	contentM1 := `{{Stack -1}}`
	contentM2 := `{{Stack -2}}`
	contentM3 := `{{Stack -3}}`

	fsys := filetree.NewInMemoryFileSystem()

	assert.NoError(t, fsys.Mkdir("template/{{departements.[].nom}}/{{$[-2].villes.[].nom}}/{{$[-2].restaurants.[].nom}}", os.ModePerm))
	assert.NoError(t, fsys.WriteFile("template/{{departements.[].nom}}/{{$[-2].villes.[].nom}}/{{$[-2].restaurants.[].nom}}/stack0.txt", []byte(content0), os.ModePerm))
	assert.NoError(t, fsys.WriteFile("template/{{departements.[].nom}}/{{$[-2].villes.[].nom}}/{{$[-2].restaurants.[].nom}}/stack1.txt", []byte(content1), os.ModePerm))
	assert.NoError(t, fsys.WriteFile("template/{{departements.[].nom}}/{{$[-2].villes.[].nom}}/{{$[-2].restaurants.[].nom}}/stack2.txt", []byte(content2), os.ModePerm))
	assert.NoError(t, fsys.WriteFile("template/{{departements.[].nom}}/{{$[-2].villes.[].nom}}/{{$[-2].restaurants.[].nom}}/stack3.txt", []byte(content3), os.ModePerm))
	assert.NoError(t, fsys.WriteFile("template/{{departements.[].nom}}/{{$[-2].villes.[].nom}}/{{$[-2].restaurants.[].nom}}/stack4.txt", []byte(content4), os.ModePerm))
	assert.NoError(t, fsys.WriteFile("template/{{departements.[].nom}}/{{$[-2].villes.[].nom}}/{{$[-2].restaurants.[].nom}}/stack5.txt", []byte(content5), os.ModePerm))
	assert.NoError(t, fsys.WriteFile("template/{{departements.[].nom}}/{{$[-2].villes.[].nom}}/{{$[-2].restaurants.[].nom}}/stack6.txt", []byte(content6), os.ModePerm))
	assert.NoError(t, fsys.WriteFile("template/{{departements.[].nom}}/{{$[-2].villes.[].nom}}/{{$[-2].restaurants.[].nom}}/stack7.txt", []byte(content7), os.ModePerm))
	assert.NoError(t, fsys.WriteFile("template/{{departements.[].nom}}/{{$[-2].villes.[].nom}}/{{$[-2].restaurants.[].nom}}/stackM1.txt", []byte(contentM1), os.ModePerm))
	assert.NoError(t, fsys.WriteFile("template/{{departements.[].nom}}/{{$[-2].villes.[].nom}}/{{$[-2].restaurants.[].nom}}/stackM2.txt", []byte(contentM2), os.ModePerm))
	assert.NoError(t, fsys.WriteFile("template/{{departements.[].nom}}/{{$[-2].villes.[].nom}}/{{$[-2].restaurants.[].nom}}/stackM3.txt", []byte(contentM3), os.ModePerm))

	driver := filetree.NewDriver(fsys)
	assert.NoError(t, driver.Develop("template/", "result/", context))

	// full context on stack 0
	file0, _ := fsys.Open("result/Ain/Bourg-en-Bresse/Auberge Bressane/stack0.txt")
	file0b, _ := io.ReadAll(file0)
	assert.Equal(t, "map[departements:[map[nom:Ain numero:1 villes:[map[nom:Bourg-en-Bresse population:41681 restaurants:[map[nom:Auberge Bressane] map[nom:Mets et Vins] map[nom:Place Bernard]]] map[nom:Oyonnax population:22271 restaurants:[map[nom:La Renaissance] map[nom:Les saveurs d'Italie] map[nom:Le Chalet Gourmand]]]]] map[nom:Aisne numero:2 villes:[map[nom:Saint-Quentin population:53100 restaurants:[map[nom:Al Taglio Pizza] map[nom:Chez Jean] map[nom:Au Made In France]]] map[nom:Soissons population:43580 restaurants:[map[nom:Le Bouche à Oreilles] map[nom:Léon] map[nom:La Cathédrale]]]]]]]", string(file0b))

	// departements list on stack 1
	file1, _ := fsys.Open("result/Ain/Bourg-en-Bresse/Auberge Bressane/stack1.txt")
	file1b, _ := io.ReadAll(file1)
	assert.Equal(t, "[map[nom:Ain numero:1 villes:[map[nom:Bourg-en-Bresse population:41681 restaurants:[map[nom:Auberge Bressane] map[nom:Mets et Vins] map[nom:Place Bernard]]] map[nom:Oyonnax population:22271 restaurants:[map[nom:La Renaissance] map[nom:Les saveurs d'Italie] map[nom:Le Chalet Gourmand]]]]] map[nom:Aisne numero:2 villes:[map[nom:Saint-Quentin population:53100 restaurants:[map[nom:Al Taglio Pizza] map[nom:Chez Jean] map[nom:Au Made In France]]] map[nom:Soissons population:43580 restaurants:[map[nom:Le Bouche à Oreilles] map[nom:Léon] map[nom:La Cathédrale]]]]]]", string(file1b))

	// first departement on stack 2
	file2, _ := fsys.Open("result/Ain/Bourg-en-Bresse/Auberge Bressane/stack2.txt")
	file2b, _ := io.ReadAll(file2)
	assert.Equal(t, "map[nom:Ain numero:1 villes:[map[nom:Bourg-en-Bresse population:41681 restaurants:[map[nom:Auberge Bressane] map[nom:Mets et Vins] map[nom:Place Bernard]]] map[nom:Oyonnax population:22271 restaurants:[map[nom:La Renaissance] map[nom:Les saveurs d'Italie] map[nom:Le Chalet Gourmand]]]]]", string(file2b))

	// cities of first departement on stack 3
	file3, _ := fsys.Open("result/Ain/Bourg-en-Bresse/Auberge Bressane/stack3.txt")
	file3b, _ := io.ReadAll(file3)
	assert.Equal(t, "[map[nom:Bourg-en-Bresse population:41681 restaurants:[map[nom:Auberge Bressane] map[nom:Mets et Vins] map[nom:Place Bernard]]] map[nom:Oyonnax population:22271 restaurants:[map[nom:La Renaissance] map[nom:Les saveurs d'Italie] map[nom:Le Chalet Gourmand]]]]", string(file3b))

	// first city of first departement on stack 4
	file4, _ := fsys.Open("result/Ain/Bourg-en-Bresse/Auberge Bressane/stack4.txt")
	file4b, _ := io.ReadAll(file4)
	assert.Equal(t, "map[nom:Bourg-en-Bresse population:41681 restaurants:[map[nom:Auberge Bressane] map[nom:Mets et Vins] map[nom:Place Bernard]]]", string(file4b))

	// restaurants of first city of first departement on stack 5
	file5, _ := fsys.Open("result/Ain/Bourg-en-Bresse/Auberge Bressane/stack5.txt")
	file5b, _ := io.ReadAll(file5)
	assert.Equal(t, "[map[nom:Auberge Bressane] map[nom:Mets et Vins] map[nom:Place Bernard]]", string(file5b))

	// first restaurant of first city of first departement on stack 6
	file6, _ := fsys.Open("result/Ain/Bourg-en-Bresse/Auberge Bressane/stack6.txt")
	file6b, _ := io.ReadAll(file6)
	assert.Equal(t, "map[nom:Auberge Bressane]", string(file6b))

	// name of first restaurant of first city of first departement on stack 7
	file7, _ := fsys.Open("result/Ain/Bourg-en-Bresse/Auberge Bressane/stack7.txt")
	file7b, _ := io.ReadAll(file7)
	assert.Equal(t, "Auberge Bressane", string(file7b))

	// name of first restaurant of first city of first departement on stack -1
	fileM1, _ := fsys.Open("result/Ain/Bourg-en-Bresse/Auberge Bressane/stackM1.txt")
	fileM1b, _ := io.ReadAll(fileM1)
	assert.Equal(t, "Auberge Bressane", string(fileM1b))

	// first restaurant of first city of first departement on stack -2
	fileM2, _ := fsys.Open("result/Ain/Bourg-en-Bresse/Auberge Bressane/stackM2.txt")
	fileM2b, _ := io.ReadAll(fileM2)
	assert.Equal(t, "map[nom:Auberge Bressane]", string(fileM2b))

	// restaurants of first city of first departement on stack -3
	fileM3, _ := fsys.Open("result/Ain/Bourg-en-Bresse/Auberge Bressane/stackM3.txt")
	fileM3b, _ := io.ReadAll(fileM3)
	assert.Equal(t, "[map[nom:Auberge Bressane] map[nom:Mets et Vins] map[nom:Place Bernard]]", string(fileM3b))
}
