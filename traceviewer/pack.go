// +build ignore

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
)

func main() {
	b, err := ioutil.ReadFile(filepath.Join(runtime.GOROOT(), "misc/trace/trace_viewer_lean.html"))
	if err != nil {
		panic(err)
	}

	w := new(bytes.Buffer)
	fmt.Fprintf(w, "package traceviewer\n\n")
	fmt.Fprintf(w, "var traceViewerLeanHTML = %q\n", b)
	if err = ioutil.WriteFile("html.go", w.Bytes(), 0666); err != nil {
		panic(err)
	}
}
