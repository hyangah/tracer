package adhoc

import (
	"bufio"
	"os"
	"sync"

	"golang.org/x/exp/mmap"

	"github.com/hyangah/tracer/adhoc/shared"
	"github.com/hyangah/tracer/trace"
)

var (
	once          sync.Once
	allEvents     []*trace.Event
	allGoroutines map[uint64]*trace.GDesc
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

func AllEvents() []*trace.Event {
	once.Do(prep)
	return allEvents
}

func AllGoroutines() map[uint64]*trace.GDesc {
	once.Do(prep)
	return allGoroutines
}

func prep() {
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

	allEvents, allGoroutines, err = shared.Unmarshal(bufio.NewReader(&reader{mr: f}))
	if err != nil {
		panic("invalid TRACER_ADHOC_TRACE: " + err.Error())
	}
}
