/* command module contains facilities to execlute particular replacements/manipulation on JSON definition.

   To define JSON manipulation command you must add it to the description

   Special commands can be embeded inside description strings that may modify existing schema.
   Basic command format is following:
   @<Value Type><Modification Path> <Value>

   Modification Path is a relative path starting from the map where command source 'description' is defined.
   {
	   "description": "
	   some text
	   @example 1
	   @another.example "33"
	   "
   }

   Will turn into following: {"description": "some text", "example": 1, "another": {"example": "33"}}

   Possible value types:

   {}    -- Value is processed as map[string]interface{}
   $     -- Value is processed as a literal string
   <NOP> -- Value is processed as bytearray

   <NOP> is used when we define example '@example "some_string"' will process "some_string" as bytearray
   (including opening/closing quotes) and creates a key "example" in a current position (where description
   defined)

   {}: '@{}someProperty {"a": "c"}' command will create someProperty key in a current position
   (where description defined) and save map[string]interface{}{"a": "c"} inside making possible further
   commands like '@{}someProperty.a.propB {"z": {"x": 10}}' update it to become as follows:
   	map[string]interface{}{
		"a": "c",
		"z": map[string]interface{}{
			"x": 10,
		},
	}

   $: forces value to be a string value.

   In addition extended value syntax available:
   @$someProp <<< EOF
   multi-line
   value
   we want to embed
   Note that it will be casted to string due to '$' presence
   EOF


   Special keys (or parts of path) include:

   _filtering  -- accepts value and casts it to a string
   _error      -- accepts JSON with example values for error message
   description -- allways casted to string
*/
package command

import (
	"encoding/json"
	"strings"
)

// Command structure ...
type Command struct {
	Path     string
	Type     string
	More     string
	Value    []string
	Complete bool
}

// GetValue function returns command argument value casted to a type depending
// on its type modificator.
func (c *Command) GetValue() interface{} {
	switch c.Type {
	case "$":
		return c.GetStringValue()
	case "{}":
		return c.GetJSONValue()
	default:
		return c.GetRawValue()
	}
}

// GetRawValue returns bytearray value.
func (c *Command) GetRawValue() interface{} {
	r := json.RawMessage(strings.Join(c.Value, "\n"))
	return &r
}

// GetJSONValue returns marshalled JSON.
func (c *Command) GetJSONValue() interface{} {
	var v map[string]interface{}
	if err := json.Unmarshal([]byte(strings.Join(c.Value, "\n")), &v); err != nil {
		panic(err)
	}

	return v
}

// GetStringValue returns string.
func (c *Command) GetStringValue() interface{} {
	return strings.Join(c.Value, "\n")
}

// Apply function applies modifications specified by command path, type.
func (c *Command) Apply(v interface{}) interface{} {
	switch c.Path {
	case "_filtering":
		return applyFiltering(v.(map[string]interface{}), c.GetStringValue())
	case "_sorting":
		return applySorting(v.(map[string]interface{}), c.GetStringValue())
	case "_error":
		return applyError(v.(map[string]interface{}), c.GetJSONValue())
	}
	return c.applyAux(v, c.Path)
}

func (c *Command) applyAux(v interface{}, path string) interface{} {
	p := strings.SplitN(path, ".", 2)

	if v == nil {
		v = map[string]interface{}{}
	}

	if path == "_error" {
		return applyError(v.(map[string]interface{}), c.GetJSONValue())
	}

	switch v := v.(type) {
	case map[string]interface{}:
		if len(p) == 2 {
			v[p[0]] = c.applyAux(v[p[0]], p[1])
		} else {
			if path == "description" {
				v[path] = c.GetStringValue()
			} else {
				v[path] = c.GetValue()
			}
		}
	}

	return v
}
