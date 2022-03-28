package gomarc21

import "fmt"

// Version is current version of this library.
var Version = version{0, 0, 9}

// v holds the version of this library.
type version struct {
	Major, Minor, Patch int
}

func (v version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}
