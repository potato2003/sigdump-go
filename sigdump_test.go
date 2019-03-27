package sigdump

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestSetup(t *testing.T) {
	path := fmt.Sprintf("/tmp/sigdump-test-testsetup-%d.log", os.Getpid())

	Setup(syscall.SIGHUP, path)

	// Send signal to this process
	syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
	time.Sleep(100 * time.Millisecond)

	_, err := os.Stat(path)
	if err != nil {
		t.Fatalf("bad: failed to TestStdup - not exists %s", path)
	}
	os.Remove(path)
}

func TestDump(t *testing.T) {
	// Suppress Stdout/Stderr
	stdoutBk := os.Stdout
	stderrBk := os.Stderr
	os.Stdout, _ = os.Open(os.DevNull)
	os.Stderr, _ = os.Open(os.DevNull)

	Dump("-")
	Dump("+")
	Dump("")

	path := fmt.Sprintf("/tmp/sigdump-test-testsetup-%d.log", os.Getpid())
	Dump(path)
	_, err := os.Stat(path)
	if err != nil {
		t.Fatalf("bad: failed to TestStdup - not exists %s", path)
	}
	os.Remove(path)

	os.Stdout = stdoutBk
	os.Stderr = stderrBk
}

func TestDumpTimestamp(t *testing.T) {
	f := new(bytes.Buffer)

	dumpTimestamp(f)

	actual := string(f.Bytes())
	_, err := time.Parse(time.RFC1123Z, strings.TrimRight(actual, "\n"))

	if err != nil {
		t.Fatalf("bad: %s", err)
	}
}

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

func TestOpenDumpPath(t *testing.T) {
	// open Stdout
	actual := openDumpPath("-")
	expected := os.Stdout

	if actual != expected {
		t.Fatalf("bad: %+v", actual)
	}

	if "-" != DumpToStdout {
		t.Fatalf("bad: DumpToStdout not equals to \"-\"")
	}

	// open Stderr
	actual = openDumpPath("+")
	expected = os.Stderr

	if actual != expected {
		t.Fatalf("bad: %+v", actual)
	}

	if "+" != DumpToStderr {
		t.Fatalf("bad: DumpToStderr not equals to \"+\"")
	}

	// open file in /tmp dir
	f := openDumpPath("")
	f.Close()

	path := fmt.Sprintf("/tmp/sigdump-%d.log", os.Getpid())
	_, err := os.Stat(path)
	if err != nil {
		t.Fatalf("bad: failed to open dump file - %s", path)
	}

	// open specified file
	path = fmt.Sprintf("/tmp/sigdump-test-%d.log", os.Getpid())
	f = openDumpPath(path)
	f.Close()

	_, err = os.Stat(path)
	if err != nil {
		t.Fatalf("bad: failed to open dump file - %s", path)
	}
	os.Remove(path)

	// open Stderr if file not exists
	actual = openDumpPath("/tmp/not_found_dir/not_found_file")
	expected = os.Stderr

	if actual != expected {
		t.Fatalf("bad: %+v", actual)
	}
}
