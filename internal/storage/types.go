package storage

import (
	"strconv"
)

// SchemaVersion is a simple int representation of the schema version.
type SchemaVersion int

// ToString converts the schema version into a string and returns that converted value.
func (s SchemaVersion) ToString() string {
	return strconv.Itoa(int(s))
}
