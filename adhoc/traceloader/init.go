package traceloader

import (
	"bufio"
	"os"

	"golang.org/x/exp/mmap"

	"github.com/hyangah/tracer/adhoc/shared"
	"github.com/hyangah/tracer/trace"
)

var (
	Events     []*trace.Event
	Goroutines map[uint64]*trace.GDesc
)

type reader struct {
	mr     *mmap.ReaderAt
	offset int64
}

func (r *reader) Read(b []byte) (n int, err error) {
	n, err = r.mr.ReadAt(b, r.offset)
	r.offset += int64(n)
	return n, err
}

func init() {
	fname := os.Getenv("TRACER_ADHOC_TRACE")
	if fname == "" {
		panic("TRACER_ADHOC_TRACE must be set")
	}
	//f, err := os.Open(fname)
	f, err := mmap.Open(fname)
	if err != nil {
		panic("invalid TRACER_ADHOC_TRACE: " + err.Error())
	}
	defer f.Close()

	Events, Goroutines, err = shared.Unmarshal(bufio.NewReader(&reader{mr:f}))
	if err != nil {
		panic("invalid TRACER_ADHOC_TRACE: " + err.Error())
	}
}
