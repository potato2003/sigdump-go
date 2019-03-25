package sigdump

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

const (
	DumpToStdout string = "-"
	DumpToStderr string = "+"
)

func Setup(signalToTrap syscall.Signal, path string) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, signalToTrap)

	go func() {
		_ = <-sigCh
		Dump(path)
	}()
}

func Dump(path string) {
	f := openDumpPath(path)
	defer func() {
		if os.Stdout == f || os.Stderr == f {
			return
		}

		f.Close()
	}()

	dumpTimestamp(f)
	dumpGoroutine(f)
}

func dumpTimestamp(w io.Writer) {
	t := time.Now().Format(time.RFC1123Z)
	w.Write([]byte(t + "\n\n"))
}

func dumpGoroutine(w io.Writer) {
	buf := make([]byte, 1024*1024)
	stacklen := runtime.Stack(buf, true)
	w.Write([]byte("=== goroutine dump ===\n"))
	w.Write(buf[:stacklen])
	w.Write([]byte("\n======\n\n"))
}

func openDumpPath(path string) io.WriteCloser {
	if path == DumpToStdout {
		return os.Stdout
	}
	if path == DumpToStderr {
		return os.Stderr
	}

	if path == "" {
		path = fmt.Sprintf("/tmp/sigdump-%d.log", os.Getpid())
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		w := openDumpPath(DumpToStderr)
		w.Write([]byte(err.Error()))
		return w
	}

	return f
}
