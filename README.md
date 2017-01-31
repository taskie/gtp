# gtp

Macro Processor using Go text/template with JSON / MessagePack

## Usage

```
gtp [OPTION] [FILE]...

OPTION:
    -j [JSON]: specify JSON path
    -m [Messagepack]: specify MessagePack path
```

## Example (JSON)


### data.json

```json
{"where":"world"}
```

### template.txt

```
Hello, {{.where}}!
```

### Execute

```bash
gtp -j data.json template.txt
```

### Output

```
Hello, world!
```

## License

Public Domain
