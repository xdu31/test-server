package util

import (
	"github.com/infobloxopen/atlas-app-toolkit/query"
)

func HasField(f *query.FieldSelection, name string) bool {
	return f == nil || f.Get(name) != nil
}
