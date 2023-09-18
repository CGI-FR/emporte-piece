![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/CGI-FR/emporte-piece/ci.yml?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/cgi-fr/emporte-piece)](https://goreportcard.com/report/github.com/cgi-fr/emporte-piece)
![GitHub all releases](https://img.shields.io/github/downloads/CGI-FR/emporte-piece/total)
![GitHub](https://img.shields.io/github/license/CGI-FR/emporte-piece)
![GitHub Repo stars](https://img.shields.io/github/stars/CGI-FR/emporte-piece)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/CGI-FR/emporte-piece)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/CGI-FR/emporte-piece)

# Emporte-Pièce

Create project file/directory structure from data + template definition.

## Usage

Example

```console
$ ep template < context.yml
8:41AM INF start color=auto debug=false log-json=false verbosity=info
8:41AM INF generating MyProject from=template/{{projectName}}
8:41AM INF generating MyProject/mask from=template/{{projectName}}/mask
8:41AM INF generating MyProject/mask/masking_TABLE1.yml from=template/{{projectName}}/mask/masking_{{$[-2].tables.[].name}}.yml
8:41AM INF generating MyProject/mask/masking_TABLE2.yml from=template/{{projectName}}/mask/masking_{{$[-2].tables.[].name}}.yml
8:41AM INF generating MyProject/pseudo_MyProject.sh from=template/{{projectName}}/pseudo_{{$[-2].projectName}}.sh
8:41AM INF end return=0
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

[GPLv3](https://choosealicense.com/licenses/gpl-3.0/)

Copyright (C) 2023 CGI France

emporte-piece is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

emporte-piece is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
 along with emporte-piece.  If not, see <http://www.gnu.org/licenses/>.
