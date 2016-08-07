// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Trace is a tool for viewing trace files.

Trace files can be generated with:
	- runtime/trace.Start
	- net/http/pprof package
	- go test -trace

Example usage:
Generate a trace file with 'go test':
	go test -trace trace.out pkg
View the trace in a web browser:
	go tool trace trace.out
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/exp/mmap"

	"github.com/hyangah/gore/session"
	"github.com/hyangah/tracer/adhoc/shared"
	"github.com/hyangah/tracer/analysis"
	"github.com/hyangah/tracer/trace" // copy of go/src/internal/trace
	"github.com/hyangah/tracer/traceviewer"
)

const usageMessage = "" +
	`Usage of 'go tool trace':
Given a trace file produced by 'go test':
	go test -trace=trace.out pkg

Open a web browser displaying trace:
	tracer [flags] [pkg.test] trace.out
[pkg.test] argument is required for traces produced by Go 1.6 and below.
Go 1.7 does not require the binary argument.

Flags:
	-http=addr: HTTP service address (e.g., ':6060')
`

var (
	httpFlag = flag.String("http", "localhost:0", "HTTP service address (e.g., ':6060')")

	// The binary file name, left here for serveSVGProfile.
	programBinary string
	traceFile     string
	ranges        []traceviewer.Range
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, usageMessage)
		os.Exit(2)
	}
	flag.Parse()

	// Go 1.7 traces embed symbol info and does not require the binary.
	// But we optionally accept binary as first arg for Go 1.5 traces.
	switch flag.NArg() {
	case 1:
		traceFile = flag.Arg(0)
	case 2:
		programBinary = flag.Arg(0)
		traceFile = flag.Arg(1)
	default:
		flag.Usage()
	}

	log.Printf("Parsing trace...")
	events, err := parseEvents()
	if err != nil {
		dief("%v\n", err)
	}
	goroutines := trace.GoroutineStats(events)

	ranges = traceviewer.Init(events, goroutines)

	analysis.RegisterHTTPHandlers(events, goroutines)

	ln, err := net.Listen("tcp", *httpFlag)
	if err != nil {
		dief("failed to create server socket: %v\n", err)
	}

	log.Printf("Opening browser: http://%v", ln.Addr().String())
	if !startBrowser("http://" + ln.Addr().String()) {
		fmt.Fprintf(os.Stderr, "Trace viewer is listening on http://%s\n", ln.Addr().String())
	}

	// Start http server.
	http.HandleFunc("/", httpMain)

	go func() {
		err = http.Serve(ln, nil)
		dief("failed to start http server: %v\n", err)
	}()

	adhoc(events, goroutines)
}

var loader struct {
	once   sync.Once
	events []*trace.Event
	err    error
}

func parseEvents() ([]*trace.Event, error) {
	loader.once.Do(func() {
		tracef, err := os.Open(traceFile)
		if err != nil {
			loader.err = fmt.Errorf("failed to open trace file: %v", err)
			return
		}
		defer tracef.Close()

		// Parse and symbolize.
		events, err := trace.Parse(bufio.NewReader(tracef), trace.Addr2LineSymbolizer(programBinary))
		if err != nil {
			loader.err = fmt.Errorf("failed to parse trace: %v", err)
			return
		}
		loader.events = events
	})
	return loader.events, loader.err
}

// httpMain serves the starting page.
func httpMain(w http.ResponseWriter, r *http.Request) {
	if err := templMain.Execute(w, ranges); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

var templMain = template.Must(template.New("").Parse(`
<html>
<body>
{{if $}}
	{{range $e := $}}
		<a href="/trace?start={{$e.Start}}&end={{$e.End}}">View trace ({{$e.Name}})</a><br>
	{{end}}
	<br>
{{else}}
	<a href="/trace">View trace</a><br>
{{end}}
<a href="/goroutines">Goroutine analysis</a><br>
<a href="/io">Network blocking profile</a><br>
<a href="/block">Synchronization blocking profile</a><br>
<a href="/syscall">Syscall blocking profile</a><br>
<a href="/sched">Scheduler latency profile</a><br>
</body>
</html>
`))

// startBrowser tries to open the URL in a browser
// and reports whether it succeeds.
// Note: copied from x/tools/cmd/cover/html.go
func startBrowser(url string) bool {
	// try to start the browser
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start() == nil
}

func dief(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg, args...)
	os.Exit(1)
}

func adhoc(events []*trace.Event, goroutines map[uint64]*trace.GDesc) {
	f, _ := ioutil.TempFile("", "mmap")
	name := f.Name()
	defer os.RemoveAll(name)

	w := bufio.NewWriter(f)
	if err := shared.Marshal(w, events, goroutines); err != nil {
		log.Fatal(err)
	}
	w.Flush()
	f.Close()

	r, err := mmap.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	os.Setenv("TRACER_ADHOC_TRACE", name)

	s := &session.Session{
		AutoImports: true,
		Pkg:         "github.com/hyangah/tracer/adhoc",
	}
	if err := s.Init(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(":help for help")

	rl := session.NewContLiner()
	defer func() { rl.Close() }()

	rl.SetWordCompleter(s.CompleteWord)

	for {
		in, err := rl.Prompt()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintf(os.Stderr, "fatal: %s", err)
			os.Exit(1)
		}

		if in == "" {
			continue
		}

		if handled, err := customCommand(in, events, goroutines); handled {
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %s\n", err)
			}
			rl.Clear()
			continue
		}

		if err := rl.Reindent(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			rl.Clear()
			continue
		}

		c := make(chan os.Signal, 1)
		go func() { <-c }()
		signal.Notify(c, os.Interrupt)

		err = s.Eval(in)
		if err != nil {
			if err == session.ErrContinue {
				signal.Stop(c)
				close(c)
				continue
			} else if err == session.ErrQuit {
				signal.Stop(c)
				close(c)
				break
			}
			fmt.Println(err)
		}
		signal.Stop(c)
		close(c)
		rl.Accepted()
	}
}

func customCommand(in string, events []*trace.Event, goroutines map[uint64]*trace.GDesc) (handled bool, err error) {
	cmd := strings.Fields(strings.TrimSpace(in))
	if len(cmd) == 0 {
		return false, nil
	}
	switch cmd[0] {
	case ":pprof":
		return pprofCmd(cmd[1:], events, goroutines)
	case ":goroutine":
		return goroutineCmd(cmd[1:], events, goroutines)
	}
	return false, nil
}

func pprofCmd(args []string, events []*trace.Event, goroutines map[uint64]*trace.GDesc) (handled bool, err error) {
	if len(args) != 2 {
		return true, fmt.Errorf("usage: :pprof [io|block|sched|syscall] output_fname")
	}
	var pprof func(w io.Writer) error
	switch args[0] {
	default:
		return true, fmt.Errorf("usage: :pprof [io|block|sched|syscall] output_fname")
	case "io":
		pprof = analysis.IOProfile
	case "block":
		pprof = analysis.BlockProfile
	case "sched":
		pprof = analysis.ScheduleLatencyProfile
	case "syscall":
		pprof = analysis.SyscallProfile
	}
	f, err := os.OpenFile(args[1], os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return true, fmt.Errorf("failed to open output file: %v", err)
	}
	if err := pprof(f); err != nil {
		f.Close()
		return true, err
	}
	return true, f.Close()
}

func goroutineCmd(args []string, events []*trace.Event, goroutines map[uint64]*trace.GDesc) (handled bool, err error) {
	for _, id := range args {
		goid, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return true, fmt.Errorf("usage: :goroutine [id1] [id2] ...")
		}
		if g, ok := goroutines[goid]; !ok {
			fmt.Printf("%d\t<goroutine not found>\n", goid)
		} else {
			fmt.Printf("%v\n", g)
		}
	}
	return true, nil
}
