# kocli

cli for knockout app

## Install

```bash
go install github.com/woocoos/kocli
```

### knockout

generate a knockout project in current directory:

```bash
kocli init -p github.com/woocoos/helloworld
```

or go get kocli and go run the script

```bash
go run github.com/woocoos/kocli/script/initproject/main.go
```

### application resource and action

generate app resource info from ent schema, before run this you need `go mod tidy` first.

```bash
# -a app code,need match your AppCode stored in Knockout application
# -f knockout.yaml is the portal config file which contains the portal db config and snowflake config.
# -e resource is ent schema path
kocli res ent -a resource -e ./ent/schema -f knockout.yaml
```

generate app action info from graphql file.

```bash
# -a and -f same as above 
# -g 99design/gqlgen config file path
kocli res gql-action -a resource -g gqlgen.yaml -f knockout.yaml
```