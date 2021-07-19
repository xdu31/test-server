// filtering module contains _filtering corner-case processing.
package command

import (
	"fmt"
)

const FilteringDescription = `A collection of response resources can be filtered by a logical expression string that includes JSON tag references to values in each resource, literal values, and logical operators. If a resource does not have the specified tag, its value is assumed to be null.

Literal values include numbers (integer and floating-point), and quoted (both single- or double-quoted) literal strings, and 'null'.

You can filter by following fields:
%s
`

func applyFiltering(v map[string]interface{}, cmdv interface{}) map[string]interface{} {
	desc := fmt.Sprintf(FilteringDescription, cmdv)
	params := v["parameters"].([]interface{})
	fixedParams := make([]interface{}, len(params))
	for i, v := range params {
		vMap := v.(map[string]interface{})
		if vMap["name"].(string) == "_filter" {
			vMap["description"] = desc
		}
		fixedParams[i] = vMap
	}
	v["parameters"] = fixedParams
	return v
}
