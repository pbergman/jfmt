package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"text/template"
)

var (
	version, ref string
)

func exit(msg string, code int) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(code)
}

func getData() interface{} {

	var data interface{}

	fi, err := os.Stdin.Stat()

	if err != nil {
		exit(err.Error(), 1)
	}

	if fi.Mode()&os.ModeCharDevice == 0 {
		dec := json.NewDecoder(os.Stdin)

		if err := dec.Decode(&data); err != nil {
			exit(err.Error(), 2)
		}
	} else {
		exit("missing json input", 3)
	}

	return data
}

func usage() {
	fmt.Printf(`jfmt %[2]s(%[3]s)

usage: %[1]s [-h|help] [template...] <INPUT

This is a small program that format json input. It uses the standard go text/template for formating the
output and if no template argument was provided then it will prity print the json input to the stdout.

example:

echo '[{"name": "foo"},{"name": "bar"}]' | %[1]s '{{ range . }}{{ .names }}\n{{end}}'
`, os.Args[0], version, ref)
}

func dump(v interface{}) string {
	out, err := json.MarshalIndent(v, "", "    ")

	if err != nil {
		exit(err.Error(), 8)
	}

	return string(out)
}

func getTmplFuncs() map[string]interface{} {
	return map[string]interface{}{
		"dump": dump,
	}
}

func fmtResult(tmpl *template.Template, source string, data interface{}) {

	var err error

	if tmpl, err = tmpl.Parse(source); err != nil {
		exit(err.Error(), 4)
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, data)

	if err != nil {
		exit(err.Error(), 5)
	}

	pattern := regexp.MustCompile(`\\[tnvf]`)
	output := buf.String()

	if pattern.MatchString(buf.String()) {
		output = pattern.ReplaceAllStringFunc(output, func(match string) string {
			switch match[1] {
			case 't':
				return "\t"
			case 'n':
				return "\n"
			case 'v':
				return "\v"
			case 'f':
				return "\f"
			}
			return ""
		})
	}

	fmt.Fprintln(os.Stdout, output)
}

func hasHelp(size int, args []string) bool {
	for i := 0; i < size; i++ {
		if args[i] == "-h" || args[i] == "help" {
			return true
		}
	}
	return false
}

func main() {
	data := getData()
	args := os.Args[1:]

	if size := len(args); size > 1 {
		if hasHelp(size, args) {
			usage()
		} else {
			tmpl := template.New("out").Option("missingkey=zero").Funcs(getTmplFuncs())
			for i := 0; i < size; i++ {
				fmtResult(tmpl, args[i], data)
			}
		}
	} else {
		fmt.Fprint(os.Stdout, dump(data))
	}
}
