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
go run github.com/woocoos/kocli/script/initproject/main.go --package github.com/woocoos/helloworld target .
```

in codegen directory, scripts use `//go:build ignore` to keep ent generate run correctly, 
it can ignore build failure when your code is not ready and has some errors.

if the package tail name contains `-`, will be renamed `_`.
### application data

#### resource

generate app resource info from ent schema, before run this you need `go mod tidy` first.

> note: a knockout project usually uses a relative path based on project root,so you need to run this command in your project root directory. 

```bash
# -a app code,need match your AppCode stored in Knockout application
# -f knockout.yaml is the portal config file which contains the portal db config and snowflake config.
# -e resource is ent schema path
kocli res ent -a resource -e ./ent/schema -f knockout.yaml
```

#### action
generate app action info from graphql file.

```bash
# -a and -f same as above 
# -g 99design/gqlgen config file path
kocli res gql-action -a resource -g gqlgen.yaml -f knockout.yaml
```

#### menu
generate app menu info from web menu data.

```bash
# -a and -f same as above
# -d web menu data path
kocli res web-menu -a resource -d ./web/menu.json -f knockout.yaml
```

menu data example:

```json
[
  {
    "name": "home",
    "icon": "home",
    "children": [
      {
        "name": "dashboard",
        "icon": "dashboard",
        "path": "/dashboard"
      }
    ]
  }
]
```

1. `path` will be auto generator an action.
2. `children` will be generated as a menu dir.
