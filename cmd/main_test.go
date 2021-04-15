package cmd

import (
	"testing"
)

var safeNameTests = []struct {
	unsafe string
	safe   string
}{
	{"this is bad", "this_is_bad"},
	{"unsafe+/;,?><chars", "unsafe_______chars"},
}

func TestMakeSafeFilName(t *testing.T) {
	for _, tt := range safeNameTests {
		actual := makeSafeName(tt.unsafe)
		if actual != tt.safe {
			t.Errorf("Unsafe filename '%s': expected %s got %s", tt.unsafe, tt.safe, actual)
		}
	}
}
