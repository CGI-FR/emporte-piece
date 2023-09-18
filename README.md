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

[MIT](https://choosealicense.com/licenses/mit/)
