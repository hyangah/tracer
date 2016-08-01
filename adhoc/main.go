package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/hyangah/tracer/trace"
)

var (
	Events     []*trace.Event
	Goroutines map[uint64]*GDesc
)

type GDesc struct {
	stats  *trace.GDesc
	events []*trace.Event
}

func (gd *GDesc) String() string {
	var buf bytes.Buffer
	if gd == nil {
		return "<nil>"
	}
	if s := gd.stats; s == nil {
		fmt.Fprintln(&buf, "<no stat>")
	} else {
		fmt.Fprintf(&buf, "%+v\n", s)
	}
	fmt.Fprintln(&buf, "Events:")
	for _, ev := range gd.events {
		desc := trace.EventDescriptions[ev.Type]
		fmt.Fprintf(&buf, "\t%v (%v) %v", ev.Ts, ev.G, desc.Name)
		if link := ev.Link; link != nil {
			desc := trace.EventDescriptions[link.Type]
			fmt.Fprintf(&buf, "-> %v (%v) %v", link.Ts, link.G, desc.Name)
		}
		fmt.Fprintf(&buf, "\n")
		for _, f := range ev.Stk {
			fmt.Fprintf(&buf, "\t\t%s:\t%s:%d\n", f.Fn, f.File, f.Line)
		}
	}
	return buf.String()
}

func analyzeGoroutines(events []*trace.Event) map[uint64]*GDesc {
	stats := trace.GoroutineStats(events)
	gs := make(map[uint64]*GDesc, len(stats))
	for _, ev := range events {
		switch ev.Type {
		case trace.EvGoCreate:
			gid := ev.Args[0]
			stat := stats[gid]
			gs[gid] = &GDesc{stats: stat, events: []*trace.Event{ev}}
		case trace.EvGoStart,
			trace.EvGoEnd, trace.EvGoStop,
			trace.EvGoBlockSend, trace.EvGoBlockRecv, trace.EvGoBlockSelect,
			trace.EvGoSched, trace.EvGoPreempt,
			trace.EvGoSleep, trace.EvGoBlock,
			trace.EvGoBlockNet,
			trace.EvGoUnblock,
			trace.EvGoSysBlock,
			trace.EvGoSysExit:
			d := gs[ev.G]
			if d != nil {
				d.events = append(d.events, ev)
				if ev.Link != nil {
					d.events = append(d.events, ev.Link)
				}
			} else if ev.G != 0 {
				log.Printf("event with unknown goroutine id=%d: %+v", ev.G, ev)
			}
		case trace.EvGCSweepStart, trace.EvGCSweepDone,
			trace.EvGCStart, trace.EvGCDone:
			ts := ev.Ts
			d := gs[ev.G]
			if d != nil && d.stats != nil &&
				d.stats.CreationTime <= ts && ts <= d.stats.EndTime {
				d.events = append(d.events, ev)
				if ev.Link != nil {
					d.events = append(d.events, ev.Link)
				}
			}
		}
	}
	for _, g := range gs {
		sort.Sort(byTimestamp(g.events))
		g.events = dedup(g.events)
	}
	return gs
}

type byTimestamp []*trace.Event

func (l byTimestamp) Len() int           { return len(l) }
func (l byTimestamp) Less(i, j int) bool { return l[i].Ts < l[j].Ts }
func (l byTimestamp) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }

func dedup(ev []*trace.Event) []*trace.Event {
	if len(ev) == 0 {
		return nil
	}
	n := 0
	for i := 1; i < len(ev); i++ {
		prev, curr := ev[n], ev[i]
		if prev == curr {
			continue
		}
		n++
		ev[n] = curr
	}
	return ev[:n+1]
}

func main() {

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: ...")
		os.Exit(2)
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
	}
	var err error
	Events, err = parseTrace(flag.Arg(0))
	if err != nil {
		log.Fatalf("failed to parse trace file %q: %v", flag.Arg(0), err)
	}

	Goroutines = analyzeGoroutines(Events)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("> ")
	for scanner.Scan() {
		if err := eval(strings.TrimSpace(scanner.Text())); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("")
		}
		fmt.Printf("> ")
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

func parseTrace(fname string) ([]*trace.Event, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return trace.Parse(bufio.NewReader(f), nil)
}

func eval(cmd string) error {
	f := strings.Fields(cmd)
	if len(f) == 0 {
		return nil
	}

	switch f[0] {
	case "goroutine":
		return goroutine(f[1:])
	case "quit":
		fmt.Println(".... Bye!")
		os.Exit(0)
	default:
		return fmt.Errorf("unknown command: %v", cmd)
	}
	return nil
}

func goroutine(cmd []string) error {
	if len(cmd) == 0 {
		return fmt.Errorf("goroutine command requires extra subcommands")
	}
	switch cmd[0] {
	case "describe":
		if len(cmd) != 2 {
			return fmt.Errorf("goroutine describe <goroutine id>")
		}
		goid, err := strconv.Atoi(cmd[1])
		if err != nil {
			return fmt.Errorf("goroutine describe <goroutine id>")
		}
		gd, ok := Goroutines[uint64(goid)]
		if !ok {
			return fmt.Errorf("no goroutine with id: %d", goid)
		}
		fmt.Printf("%s\n", gd)
	case "list":
		var substr string
		if len(cmd) == 2 {
			substr = cmd[1]
		}
		for goid, desc := range Goroutines {
			if substr == "" {
				fmt.Printf("%d: %v\n", goid, desc)
				continue
			}
			if len(desc.events) == 0 {
				continue
			}
			ev := desc.events[0]
			if ev.Type == trace.EvGoCreate && ev.Link != nil {
				ev = ev.Link
			}
			found := false
			for _, f := range ev.Stk {
				if strings.Contains(f.Fn, substr) || strings.Contains(f.File, substr) {
					found = true
					break
				}
			}
			if found {
				fmt.Printf("%d: %v\n", goid, desc)
			}
		}
	}
	return nil
}
