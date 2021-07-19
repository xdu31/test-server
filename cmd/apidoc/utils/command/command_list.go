// command_list module contains method to parse and manipulate multiple commands.
package command

import (
	"regexp"
)

var (
	cmdRe = regexp.MustCompile(`^@(?P<Type>{}|\$)?(?P<Path>[_a-zA-Z0-9.]+)(?:$|\s+<<<\s*(?P<More>\w+)|(?P<Value>.*))`)
)

// CommandList structure ...
type CommandList []*Command

func (c *CommandList) tail() *Command {
	if len(*c) == 0 {
		return nil
	}

	return (*c)[len(*c)-1]
}

// Parse function accepts a string and depending on its contents takes
// following actions:
//   - creates new command and appends it to the list
//   - if command is defined with extended syntax appends value to the
//     end of the last command till end-marker reached
//   - if no active commands found and string does not contain any return
//     false
func (c *CommandList) Parse(s string) bool {

	if t := c.tail(); t != nil && !t.Complete {
		if s == t.More {
			t.Complete = true
			return true
		} else {
			t.Value = append(t.Value, s)
			return true
		}
	}

	match := cmdRe.FindStringSubmatch(s)

	cmd := &Command{}
	for i, n := range cmdRe.SubexpNames() {
		if i != 0 && len(match) > i {
			switch n {
			case "Path":
				cmd.Path = match[i]
			case "Type":
				cmd.Type = match[i]
			case "Value":
				cmd.Value = []string{match[i]}
			case "More":
				cmd.More = match[i]
				if match[i] == "" {
					cmd.Complete = true
				}
			}
		}
	}

	if cmd.Path != "" {
		*c = append(*c, cmd)
		return true
	}

	return false
}
