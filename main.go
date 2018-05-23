package main

import (
    "os"
    "fmt"
    "encoding/json"
    "text/template"
    "bytes"
    "regexp"
)

var (
    version, ref string
)

func getData() interface{} {

    var data interface{}

    fi, err := os.Stdin.Stat()

    if err != nil {
        panic(err)
    }

    if fi.Mode() & os.ModeCharDevice == 0 {
        dec := json.NewDecoder(os.Stdin)

        if err := dec.Decode(&data); err != nil {
            panic(err)
        }
    }

    return data
}

func usage() {
    fmt.Printf(`jfmt %[2]s(%[3]s)

usage: %[1]s [-h|help] [template] <INPUT

This is a small program that format json input. It uses the standard go text/template for formating the
output and if no template argument was provided then it will prity print the json input to the stdout.

example:

echo '[{"name": "foo"},{"name": "bar"}]' | %[1]s '{{ range . }}{{ .names }}\n{{end}}'
`, os.Args[0], version, ref)
}

func main() {

    data := getData()

    if len(os.Args) > 1 {

        if os.Args[1] == "-h" || os.Args[1] == "help" {
            usage()
            return
        }

        tmpl, err := template.New("out").Parse(os.Args[1])

        if err != nil {
            panic(err)
        }

        buf := new(bytes.Buffer)

        err = tmpl.Execute(buf, data)

        if err != nil {
            panic(err)
        }

        pattern := regexp.MustCompile(`\\[tnvf]`)
        output := buf.String()

        if pattern.MatchString(buf.String()) {
            output = pattern.ReplaceAllStringFunc(output, func(match string) string {
                switch match[1] {
                case 't':
                    return "\t";
                case 'n':
                    return "\n";
                case 'v':
                    return "\v";
                case 'f':
                    return "\f";
                }
                return "";
            })
        }

        fmt.Print(output)

    } else {
        out, err := json.MarshalIndent(data, "", "    ")

        if err != nil {
            panic(err)
        }

        fmt.Println(string(out))
    }
}
