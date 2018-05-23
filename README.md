### JFMT

jfmt (json formatting) is a small program without any dependencies that can format or print json input. It uses the standard [go text/template](https://golang.org/pkg/text/template)
for formatting the output and if no template argument was provided then it will pretty print the json input to the stdout.

#### install

make build && mv jfmt /usr/local/bin

#### example

```
echo "[{"name": "foo"},{"name": "bar:}]" | jfmt "{{ range . }}>{{ .name }}\n{{end}}"
>foo
>bar
```