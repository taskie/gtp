# gtp

macro processor powered by Go (text|html)/template

## Install

```sh
go install github.com/taskie/gtp/cmd/gtp
```

## Usage

```sh
gtp -d john.json hello.tmpl >john.txt
```

or

```sh
gtp -D json hello.tmpl <john.json >john.txt
```

### Template (hello.tmpl)

```
Hello, {{.user.name}}!
```

### Input (john.json)

```json
{"user": {"name": "John"}}
```

### Output (john.txt)

```
Hello, John!
```

## Dependency

![dependency](images/dependency.png)


## License

Apache License 2.0
