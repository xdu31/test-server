package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-openapi/spec"

	"github.com/xdu31/test-server/cmd/apidoc/utils/command"
)

var (
	apiKey bool
	proto  string
	scheme string
)

func init() {
	flag.StringVar(&proto, "proto", "", "HTTP or HTTPS protocol used")
	flag.BoolVar(&apiKey, "apiKey", false, "Determines whether API Key authorization is present")
	flag.StringVar(&scheme, "scheme", "", "Path to a scheme file")

	flag.Parse()
}

func main() {
	b, err := ioutil.ReadFile(scheme)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading scheme file: %v", err)
		os.Exit(1)
	}

	var v map[string]interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		fmt.Fprintf(os.Stderr, "Error unmarshalling scheme: %v", err)
		os.Exit(1)
	}

	if !apiKey {
		delete(v, "securityDefinitions")
		delete(v, "security")
	}

	if proto != "" {
		v["schemes"] = []string{proto}
	}

	b, err = json.Marshal(walk(v))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling scheme: %v", err)
		os.Exit(1)
	}

	var sw spec.Swagger

	if err := json.Unmarshal(b, &sw); err != nil {
		fmt.Fprintf(os.Stderr, "Error unmarshaling scheme: %v", err)
		os.Exit(1)
	}

	b, _ = json.MarshalIndent(sw, "", "  ")

	fmt.Printf("%s\n", b)
	os.Exit(0)
}

func parseString(s string) (string, *command.CommandList) {
	cl := &command.CommandList{}
	desc := []string{}

	for _, v := range strings.Split(s, "\n") {
		trimmed := strings.TrimSpace(v)

		if !cl.Parse(trimmed) {
			desc = append(desc, v)
		}
	}

	return strings.Join(desc, "\n"), cl
}

func parseCommands(v map[string]interface{}) (map[string]interface{}, *command.CommandList) {
	var desc string
	var cmd *command.CommandList

	if d, dOk := v["description"]; dOk {
		if dStr, ok := d.(string); ok {
			desc, cmd = parseString(dStr)
		}
	} else if t, tOk := v["title"]; tOk {
		if tStr, ok := t.(string); ok {
			desc, cmd = parseString(tStr)
		}
		delete(v, "title")
	}

	if desc != "" {
		v["description"] = desc
	}

	return v, cmd
}

func walk(v interface{}) interface{} {
	switch v := v.(type) {
	case map[string]interface{}:
		var cmdList *command.CommandList
		v, cmdList = parseCommands(v)

		if cmdList != nil {
			for _, cmd := range *cmdList {
				v = cmd.Apply(v).(map[string]interface{})
			}
		}

		for k, vv := range v {
			v[k] = walk(vv)
		}
	case []interface{}:
		for i, vv := range v {
			v[i] = walk(vv)
		}
	default:
		return v
	}

	return v
}
