// This utility is used to filter OpenAPI operations by tag.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-openapi/spec"
)

var (
	tags   map[string]bool
	opIds  map[string]bool
	scheme string
)

func init() {
	tagsPtr := flag.String("tags", "", "tags to be filtered")
	operationIds := flag.String("operationIds", "", "operationIds to be filtered")
	schemePtr := flag.String("scheme", "", "path to a scheme file")

	flag.Parse()

	tags = map[string]bool{}
	for _, v := range strings.Split(*tagsPtr, ",") {
		tags[v] = true
	}

	opIds = map[string]bool{}
	for _, v := range strings.Split(*operationIds, ",") {
		opIds[v] = true
	}

	scheme = *schemePtr
}

func main() {
	v := filterScheme()

	// Marshal scheme to a bytearray
	b, err := json.Marshal(v)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling scheme: %v", err)
		os.Exit(1)
	}

	// Marshal and Umarshal using go-openapi for
	// nicier formatting.
	var sw spec.Swagger
	json.Unmarshal(b, &sw)
	b, _ = json.MarshalIndent(sw, "", "  ")

	fmt.Printf("%s\n", b)
	os.Exit(0)
}

func filterScheme() interface{} {
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

	filter(v["paths"])

	defs := v["definitions"].(map[string]interface{})
	newDefs := map[string]interface{}{}

	// Update definitions to include only referenced schemes.
	for rk, _ := range gatherRefs(v["paths"]) {
		rName := refToName(rk)
		newDefs[rName] = defs[rName]
		for rrName, _ := range gatherDefinitionRefs(rk, v["definitions"].(map[string]interface{})) {
			newDefs[rrName] = defs[rrName]
		}
	}

	v["definitions"] = newDefs
	return v
}

func gatherDefinitionRefs(ref string, defs map[string]interface{}) map[string]struct{} {
	var seen = map[string]bool{}
	var refs = map[string]struct{}{}

	gatherDefinitionRefsAux(ref, defs, refs, seen)
	return refs
}

func gatherDefinitionRefsAux(ref string, defs map[string]interface{}, refs map[string]struct{}, seen map[string]bool) {
	if seen[ref] {
		return
	}

	seen[ref] = true

	for r, _ := range gatherRefs(defs[refToName(ref)]) {
		refs[refToName(r)] = struct{}{}
		gatherDefinitionRefsAux(r, defs, refs, seen)
	}

	return
}

func refToName(ref string) string {
	return strings.TrimPrefix(ref, "#/definitions/")
}

func gatherRefs(v interface{}) map[string]struct{} {
	refs := map[string]struct{}{}
	switch v := v.(type) {
	case map[string]interface{}:
		for k, vv := range v {
			if k == "$ref" {
				refs[vv.(string)] = struct{}{}
			}

			for rk, _ := range gatherRefs(vv) {
				refs[rk] = struct{}{}
			}
		}
	case []interface{}:
		for _, vv := range v {
			for rk, _ := range gatherRefs(vv) {
				refs[rk] = struct{}{}
			}
		}
	}

	return refs
}

func matches(v map[string]interface{}) bool {
	if vv, ok := v["tags"]; ok {
		if tags[vv.([]interface{})[0].(string)] {
			return true
		}
	}

	if vv, ok := v["operationId"]; ok {
		if opIds[vv.(string)] {
			return true
		}
	}

	return false
}

func filter(v interface{}) bool {
	switch v := v.(type) {
	case map[string]interface{}:
		if matches(v) {
			return true
		}

		for k, vv := range v {
			if filter(vv) {
				delete(v, k)
			}
		}

		if len(v) == 0 {
			return true
		}
	}

	return false
}
