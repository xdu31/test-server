// sorting module contains _sorting corner-case processing.
package command

import (
	"fmt"
)

const SortingDescription = `A collection of response resources can be sorted by a space separated field-ordering pair.

Ordering can be either 'asc' (for ascended order: 'A' to 'Z' for strings, and 0 to 9 for numbers) or 'desc' (for descended order).

You can sort by following fields:
%s
`

func applySorting(v map[string]interface{}, cmdv interface{}) map[string]interface{} {
	desc := fmt.Sprintf(SortingDescription, cmdv)
	params, ok := v["parameters"].([]interface{})
	if !ok {
		return v
	}
	fixedParams := make([]interface{}, len(params))
	for i, v := range params {
		vMap := v.(map[string]interface{})
		if vMap["name"].(string) == "_order_by" {
			vMap["description"] = desc
		}
		fixedParams[i] = vMap
	}
	v["parameters"] = fixedParams
	return v
}
