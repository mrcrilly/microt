package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func main() {
	in := flag.String("in", "", "provide a file to parse")
	out := flag.String("out", "parsed.file", "provide a file to write to")
	flag.Parse()

	if *in == "" || strings.Contains(*in, "=") {
		panic("bad value for in given")
	}

	if *out == "" || strings.Contains(*out, "=") {
		panic("bad value for out given")
	}

	raw, err := ioutil.ReadFile(*in)
	if err != nil {
		panic(err)
	}

	args := os.Args[1:]
	vars := map[string]interface{}{}

	for _, arg := range args {
		if strings.Contains(arg, "=") {
			kv := strings.Split(arg, "=")

			if strings.Contains(kv[1], "::") {
				subKv := strings.Split(kv[1], "::")

				if len(subKv[1]) == 0 {
					panic("bad type given for value: " + kv[0])
				}

				fmt.Printf("value type: %s\n", subKv[1])

				switch subKv[1] {
				case "int32":
					asInt32, err := strconv.ParseInt(subKv[0], 10, 32)
					if err != nil {
						panic(err)
					}

					vars[kv[0]] = asInt32

				case "int64":
					asInt64, err := strconv.ParseInt(subKv[0], 10, 64)
					if err != nil {
						panic(err)
					}

					vars[kv[0]] = asInt64

				case "float32":
					asFloat32, err := strconv.ParseFloat(subKv[0], 32)
					if err != nil {
						panic(err)
					}

					vars[kv[0]] = asFloat32

				case "float64":
					asFloat64, err := strconv.ParseFloat(subKv[0], 64)
					if err != nil {
						panic(err)
					}

					vars[kv[0]] = asFloat64

				case "bool":
					asBool, err := strconv.ParseBool(subKv[0])
					if err != nil {
						panic(err)
					}

					vars[kv[0]] = asBool

				default:
					panic("unknown value type: " + subKv[1])
				}
			} else {
				vars[kv[0]] = kv[1]
			}
		}
	}

	fmt.Printf("%+v\n\n", vars)

	tmpl, err := template.New("x").Parse(string(raw))
	if err != nil {
		panic(err)
	}

	fd, err := os.OpenFile(*out, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	defer fd.Close()
	err = tmpl.Execute(fd, vars)
	if err != nil {
		panic(err)
	}
}
