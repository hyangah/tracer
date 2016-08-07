package adhoc

import (
	"regexp"

	"github.com/hyangah/tracer/trace"
)

/*
type Event struct {
	Off  int  // offset in input file (for debugging and error reporting)
	Type byte // one of Ev*

	Ts    int64     // timestamp in nanoseconds
	P     int       // P on which the event happened (can be one of TimerP, NetpollP, SyscallP)
	G     uint64    // G on which the event happened
	StkID uint64    // unique stack ID
	Stk   []*Frame  // stack trace (can be empty)
	Args  [3]uint64 // event-type-specific arguments
	// linked event (can be nil), depends on event type:
	// for GCStart: the GCStop
	// for GCScanStart: the GCScanDone
	// for GCSweepStart: the GCSweepDone
	// for GoCreate: first GoStart of the created goroutine
	// for GoStart: the associated GoEnd, GoBlock or other blocking event
	// for GoSched/GoPreempt: the next GoStart
	// for GoBlock and other blocking events: the unblock event
	// for GoUnblock: the associated GoStart
	// for blocking GoSysCall: the associated GoSysExit
	// for GoSysExit: the next GoStart
	Link *Event
	// Has unexported fields.
}
*/

type EventFilter func(e *trace.Event) bool

func Events(filter EventFilter) (selected []*trace.Event) {
	for _, e := range AllEvents() {
		if filter(e) {
			selected = append(selected, e)
		}
	}
	return selected
}

func StackContains(function, file string) EventFilter {
	if function == "" && file == "" {
		return func(_ *trace.Event) bool { return true }
	}

	var fnRE, fileRE *regexp.Regexp
	if function != "" {
		fnRE = regexp.MustCompile(function)
	}
	if file != "" {
		fileRE = regexp.MustCompile(file)
	}
	return func(e *trace.Event) bool {
		for _, frame := range e.Stk {
			if fnRE != nil && fnRE.MatchString(frame.Fn) {
				return true
			}
			if fileRE != nil && fileRE.MatchString(frame.File) {
				return true
			}
		}
		return false
	}
}

func forGoroutine(id uint64) EventFilter {
	return func(e *trace.Event) bool {
		switch e.Type {
		case trace.EvGoCreate: // who created
			if e.Link != nil && e.Link.G == id {
				return true
			}
		case trace.EvGoUnblock: // who unblocked
			if e.Args[0] == id {
				return true
			}
		}
		return e.G == id
	}
}

func ForGoroutine(id uint64) EventFilter {
	return forGoroutine(id)
}

func ForGoroutines(ids []uint64) EventFilter {
	gids := make(map[uint64]EventFilter, len(ids))
	for _, id := range ids {
		gids[id] = forGoroutine(id)
	}
	return func(e *trace.Event) bool {
		for _, f := range gids {
			if f(e) {
				return true
			}
		}
		return false
	}
}
