App CLI parameters
```
The flibooks inpx app for processing ebooks inpx archives

Usage:
  flibooks [flags]
  flibooks [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  parse       Parse the inpx file
  serve       Serves the REST API

Flags:
      --config-path string   Path to config file
  -h, --help                 help for flibooks

Use "flibooks [command] --help" for more information about a command.
```

App properties

```
Database.Type=mysql
Database.Connection=user:pass@tcp(127.0.0.1:3306)/db?charset=utf8&parseTime=True&loc=Local
Database.LogLevel=Info

Server.Listen=:8000

Data.Dir=
```