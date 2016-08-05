package analysis

import (
	"net/http"
	"sync"

	"github.com/hyangah/tracer/trace" // copy of go/src/internal/trace
)

var (
	initOnce    sync.Once
	traceEvents []*trace.Event
	gs          map[uint64]*trace.GDesc
)

func RegisterHTTPHandlers(events []*trace.Event, goroutines map[uint64]*trace.GDesc) {
	initOnce.Do(func() {
		traceEvents = events
		gs = goroutines

		// register http handlers
		http.HandleFunc("/io", serveSVGProfile(IOProfile))
		http.HandleFunc("/block", serveSVGProfile(BlockProfile))
		http.HandleFunc("/syscall", serveSVGProfile(SyscallProfile))
		http.HandleFunc("/sched", serveSVGProfile(ScheduleLatencyProfile))

		http.HandleFunc("/goroutines", httpGoroutines)
		http.HandleFunc("/goroutine", httpGoroutine)
	})
}
