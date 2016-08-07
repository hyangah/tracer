// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Serving of pprof-like profiles.

package analysis

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	"github.com/hyangah/tracer/pprof/profile" // copy of cmd/internal/pprof/profile
	"github.com/hyangah/tracer/trace"         // copy of go/src/internal/trace
)

// record represents one entry in pprof-like profiles.
type record struct {
	stk  []*trace.Frame
	n    uint64
	time int64
}

// IOProfile computes IO pprof-like profile (time spent in IO wait).
func IOProfile(w io.Writer) error {
	events := traceEvents
	prof := make(map[uint64]record)
	for _, ev := range events {
		if ev.Type != trace.EvGoBlockNet || ev.Link == nil || ev.StkID == 0 || len(ev.Stk) == 0 {
			continue
		}
		rec := prof[ev.StkID]
		rec.stk = ev.Stk
		rec.n++
		rec.time += ev.Link.Ts - ev.Ts
		prof[ev.StkID] = rec
	}
	return buildProfile(prof).Write(w)
}

// BlockProfile computes blocking pprof-like profile (time spent blocked on synchronization primitives).
func BlockProfile(w io.Writer) error {
	events := traceEvents
	prof := make(map[uint64]record)
	for _, ev := range events {
		switch ev.Type {
		case trace.EvGoBlockSend, trace.EvGoBlockRecv, trace.EvGoBlockSelect,
			trace.EvGoBlockSync, trace.EvGoBlockCond:
		default:
			continue
		}
		if ev.Link == nil || ev.StkID == 0 || len(ev.Stk) == 0 {
			continue
		}
		rec := prof[ev.StkID]
		rec.stk = ev.Stk
		rec.n++
		rec.time += ev.Link.Ts - ev.Ts
		prof[ev.StkID] = rec
	}
	return buildProfile(prof).Write(w)
}

// SyscallProfile computes syscall pprof-like profile (time spent blocked in syscalls).
func SyscallProfile(w io.Writer) error {
	events := traceEvents
	prof := make(map[uint64]record)
	for _, ev := range events {
		if ev.Type != trace.EvGoSysCall || ev.Link == nil || ev.StkID == 0 || len(ev.Stk) == 0 {
			continue
		}
		rec := prof[ev.StkID]
		rec.stk = ev.Stk
		rec.n++
		rec.time += ev.Link.Ts - ev.Ts
		prof[ev.StkID] = rec
	}
	return buildProfile(prof).Write(w)
}

// ScheduleLatencyProfile serves scheduler latency pprof-like profile
// (time between a goroutine become runnable and actually scheduled for execution).
func ScheduleLatencyProfile(w io.Writer) error {
	events := traceEvents
	prof := make(map[uint64]record)
	for _, ev := range events {
		if (ev.Type != trace.EvGoUnblock && ev.Type != trace.EvGoCreate) ||
			ev.Link == nil || ev.StkID == 0 || len(ev.Stk) == 0 {
			continue
		}
		if ev.Type == trace.EvGoCreate && ev.G == 0 { // Fake EvGoCreate event added when starting trace.
			continue
		}
		rec := prof[ev.StkID]
		rec.stk = ev.Stk
		rec.n++
		rec.time += ev.Link.Ts - ev.Ts
		prof[ev.StkID] = rec
	}
	return buildProfile(prof).Write(w)
}

// serveSVGProfile generates pprof-like profile stored in prof and writes in to w.
func serveSVGProfile(prof func(w io.Writer) error) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		blockf, err := ioutil.TempFile("", "block")
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to create temp file: %v", err), http.StatusInternalServerError)
			return
		}
		defer func() {
			blockf.Close()
			os.Remove(blockf.Name())
		}()
		blockb := bufio.NewWriter(blockf)
		if err := prof(blockb); err != nil {
			http.Error(w, fmt.Sprintf("failed to generate profile: %v", err), http.StatusInternalServerError)
			return
		}
		if err := blockb.Flush(); err != nil {
			http.Error(w, fmt.Sprintf("failed to flush temp file: %v", err), http.StatusInternalServerError)
			return
		}
		if err := blockf.Close(); err != nil {
			http.Error(w, fmt.Sprintf("failed to close temp file: %v", err), http.StatusInternalServerError)
			return
		}
		svgFilename := blockf.Name() + ".svg"
		if output, err := exec.Command("go", "tool", "pprof", "-svg", "-output", svgFilename, blockf.Name()).CombinedOutput(); err != nil {
			http.Error(w, fmt.Sprintf("failed to execute go tool pprof: %v\n%s", err, output), http.StatusInternalServerError)
			return
		}
		defer os.Remove(svgFilename)
		w.Header().Set("Content-Type", "image/svg+xml")
		http.ServeFile(w, r, svgFilename)
	}
}

func buildProfile(prof map[uint64]record) *profile.Profile {
	p := &profile.Profile{
		PeriodType: &profile.ValueType{Type: "trace", Unit: "count"},
		Period:     1,
		SampleType: []*profile.ValueType{
			{Type: "contentions", Unit: "count"},
			{Type: "delay", Unit: "nanoseconds"},
		},
	}
	locs := make(map[uint64]*profile.Location)
	funcs := make(map[string]*profile.Function)
	for _, rec := range prof {
		var sloc []*profile.Location
		for _, frame := range rec.stk {
			loc := locs[frame.PC]
			if loc == nil {
				fn := funcs[frame.File+frame.Fn]
				if fn == nil {
					fn = &profile.Function{
						ID:         uint64(len(p.Function) + 1),
						Name:       frame.Fn,
						SystemName: frame.Fn,
						Filename:   frame.File,
					}
					p.Function = append(p.Function, fn)
					funcs[frame.File+frame.Fn] = fn
				}
				loc = &profile.Location{
					ID:      uint64(len(p.Location) + 1),
					Address: frame.PC,
					Line: []profile.Line{
						profile.Line{
							Function: fn,
							Line:     int64(frame.Line),
						},
					},
				}
				p.Location = append(p.Location, loc)
				locs[frame.PC] = loc
			}
			sloc = append(sloc, loc)
		}
		p.Sample = append(p.Sample, &profile.Sample{
			Value:    []int64{int64(rec.n), rec.time},
			Location: sloc,
		})
	}
	return p
}
