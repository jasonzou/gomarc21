package gomarc21

import "testing"

func TestVersion(t *testing.T) {
	v := version{1, 0, 0}
	if x := v.String(); x != "1.0.0" {
		t.Fatalf("Failed to convert version %v, got: %s", v, x)
	}
}
