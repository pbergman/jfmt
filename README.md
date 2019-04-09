### JFMT

jfmt (json formatting) is a small program without any dependencies that can format or print json input.

It uses the standard [go text/template](https://golang.org/pkg/text/template) for formatting the output and if no template argument was provided then it will pretty print the json input to the stdout.

#### install

make build && mv jfmt /usr/local/bin

#### example

```
echo "[{"name": "foo"},{"name": "bar:}]" | jfmt "{{ range . }}>{{ .name }}\n{{end}}"
>foo
>bar
```

using the template blocks with multiple arguments:

```
echo "[{"name": "foo"},{"name": "bar:}]" | jfmt '{{define "dumper"}}{{ . | dump }}{{end}}' '{{template "dumper" .}}'
[
    {
        "name": "foo"
    },
    {
        "name": "bar"
    }
]
```

or with predefined templates:

```
# create default template folder
[ ! -d "~/.config/jfmt/" ] && mkdir  ~/.config/jfmt/
# create a template
echo '{{define "dumper"}}{{ . | dump }}{{end}}' >  ~/.config/jfmt/dumper.tmpl
# print output
echo "[{"name": "foo"},{"name": "bar:}]" | jfmt '{{template "dumper" .}}'

```