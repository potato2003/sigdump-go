package sigdump

import (
	"bytes"
	"strings"
	"testing"
)

func TestDumpGoroutine(t *testing.T) {
	f := new(bytes.Buffer)

	dumpGoroutine(f)

	actual := string(f.Bytes())
	expected := "sigdump-go.TestDumpGoroutine" // self test case name

	if !strings.Contains(actual, expected) {
		t.Fatalf("bad: not contains %s", expected)
		t.Fatalf("actual: %s", actual)
	}
}
