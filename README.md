# Sigdump-Go

Dump gorountine stacktrace whenever signal received

## How to use

```
package main

import (
	"fmt"
	"syscall"
	"os"

	sigdump "github.com/potato2003/sigdump-go"
)

func main() {
	fmt.Printf("Prease try to\n")
	fmt.Printf("\t$ kill -CONT %d\n", os.Getpid())

	sigdump.Setup(syscall.SIGCONT, "./dump.log")
	// sigdump.Setup(syscall.SIGCONT, "") // dump to /tmp/sigdump-${pid}.log
	// sigdump.Setup(syscall.SIGCONT, sigdump.DumpToStdout /* or "-" */) // dump to os.Stdout
	// sigdump.Setup(syscall.SIGCONT, sigdump.DumpToStderr /* or "+" */) // dump to os.Stderr

	select {}
}
```

output
```
Mon, 25 Mar 2019 21:11:51 +0900

=== goroutine dump ===
goroutine 17 [running]:
sigdump-sandbox/vendor/github.com/potato2003/sigdump-go.dumpGoroutine(0x10dea60, 0xc0000a2000)
	/Users/potato2003/src/sigdump-sandbox/vendor/github.com/potato2003/sigdump-go/sigdump.go:49 +0x74
sigdump-sandbox/vendor/github.com/potato2003/sigdump-go.Dump(0x10c9f01, 0xa)
	/Users/potato2003/src/sigdump-sandbox/vendor/github.com/potato2003/sigdump-go/sigdump.go:39 +0xe5
sigdump-sandbox/vendor/github.com/potato2003/sigdump-go.Setup.func1(0xc000058180, 0x10c9f01, 0xa)
	/Users/potato2003/src/sigdump-sandbox/vendor/github.com/potato2003/sigdump-go/sigdump.go:24 +0x55
created by sigdump-sandbox/vendor/github.com/potato2003/sigdump-go.Setup
	/Users/potato2003/src/sigdump-sandbox/vendor/github.com/potato2003/sigdump-go/sigdump.go:22 +0xc7

goroutine 1 [select (no cases)]:
main.main()
	/Users/potato2003/src/sigdump-sandbox/main.go:19 +0xe1

goroutine 5 [syscall]:
os/signal.signal_recv(0x10deb80)
	/usr/local/opt/go/libexec/src/runtime/sigqueue.go:139 +0x9f
os/signal.loop()
	/usr/local/opt/go/libexec/src/os/signal/signal_unix.go:23 +0x22
created by os/signal.init.0
	/usr/local/opt/go/libexec/src/os/signal/signal_unix.go:29 +0x41

======
```
