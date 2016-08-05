package shared

//go:generate gencode go -schema data.schema -package shared

import (
	"fmt"
	"io"

	"github.com/hyangah/tracer/trace"
)

func toFrameIDs(framesID map[*trace.Frame]int, stk []*trace.Frame) []int64 {
	if len(stk) == 0 {
		return nil
	}
	r := make([]int64, len(stk))
	for i, f := range stk {
		r[i] = int64(framesID[f])
	}
	return r
}

func toFrames(frames []trace.Frame, stk []int64) []*trace.Frame {
	if len(stk) == 0 {
		return nil
	}
	r := make([]*trace.Frame, len(stk))
	for i, f := range stk {
		if f >= int64(len(frames)) || f < 0 {
			panic(fmt.Sprintf("f=%v len(frames)=%v", f, len(frames)))
		}
		r[i] = &frames[f]
	}
	return r
}

const invalidEventID = int64(-1)

func toEventID(events map[*trace.Event]int, e *trace.Event) int64 {
	if e == nil {
		return invalidEventID
	}
	if v, ok := events[e]; !ok {
		return invalidEventID
	} else {
		return int64(v)
	}
}

func toEvent(events []trace.Event, id int64) *trace.Event {
	if id == invalidEventID {
		return nil
	}
	return &events[id]
}

type chunk struct {
	r io.Reader

	buf    []byte
	offset int
	last   int
	eof    bool
}

var bufSize = 1 << 20

func (c *chunk) Bytes(minsz int) ([]byte, error) {
	if n := cap(c.buf); n < minsz || n < bufSize {
		b := c.buf
		if minsz < bufSize {
			c.buf = make([]byte, bufSize)
		} else {
			c.buf = make([]byte, minsz)
		}
		if c.last-c.offset > 0 {
			c.last = copy(c.buf, b[c.offset:c.last])
			c.offset = 0
		}
	}
	if !c.eof && c.last-c.offset < minsz {
		c.last = copy(c.buf[0:], c.buf[c.offset:c.last])
		c.offset = 0

		n, err := io.ReadFull(c.r, c.buf[c.last:cap(c.buf)])
		if err == io.EOF {
			c.eof = true
		}
		c.last += n
	}
	if c.last-c.offset == 0 {
		return nil, io.EOF
	}
	return c.buf[c.offset:c.last], nil
}

func (c *chunk) Advance(sz int) error {
	if c.last-c.offset < sz {
		return fmt.Errorf("invalid")
	}
	c.offset += sz
	return nil
}

func Unmarshal(r io.Reader) ([]*trace.Event, map[uint64]*trace.GDesc, error) {
	c := &chunk{r: r}
	sz := 1 << 10

	// Header
	b, err := c.Bytes(sz)
	if err != nil {
		return nil, nil, err
	}
	var h Header
	if n, err := h.Unmarshal(b); err != nil {
		return nil, nil, err
	} else {
		c.Advance(int(n))
	}

	events := make([]trace.Event, h.Events)
	frames := make([]trace.Frame, h.Frames)
	gdescs := make(map[uint64]*trace.GDesc, h.GDescs)

	// Events
	for i := int64(0); i < h.Events; i++ {
		b, err := c.Bytes(sz)
		if err != nil {
			return nil, nil, err
		}
		var ev Event
		if n, err := ev.Unmarshal(b); err != nil {
			return nil, nil, err
		} else {
			c.Advance(int(n))
		}
		events[i] = trace.Event{
			Off:   int(ev.Off),
			Type:  ev.Type,
			Ts:    ev.Ts,
			P:     int(ev.P),
			G:     ev.G,
			StkID: ev.StkID,
			Stk:   toFrames(frames, ev.Stk),
			Args:  ev.Args,
			Link:  toEvent(events, ev.Link),
		}
	}

	// Frames
	for i := int64(0); i < h.Frames; i++ {
		b, err := c.Bytes(sz)
		if err != nil {
			return nil, nil, err
		}
		var frame Frame
		if n, err := frame.Unmarshal(b); err != nil {
			return nil, nil, err
		} else {
			c.Advance(int(n))
		}
		frames[i] = trace.Frame{
			PC:   frame.PC,
			Fn:   frame.Fn,
			File: frame.File,
			Line: int(frame.Line),
		}
	}

	// GDescs
	for i := int64(0); i < h.GDescs; i++ {
		b, err := c.Bytes(sz)
		if err != nil {
			return nil, nil, err
		}
		var gd GDesc
		if n, err := gd.Unmarshal(b); err != nil {
			return nil, nil, err
		} else {
			c.Advance(int(n))
		}
		gdescs[gd.ID] = &trace.GDesc{
			ID:            gd.ID,
			Name:          gd.Name,
			PC:            gd.PC,
			CreationTime:  gd.CreationTime,
			StartTime:     gd.StartTime,
			EndTime:       gd.EndTime,
			ExecTime:      gd.ExecTime,
			SchedWaitTime: gd.SchedWaitTime,
			IOTime:        gd.IOTime,
			BlockTime:     gd.BlockTime,
			SyscallTime:   gd.SyscallTime,
			GCTime:        gd.GCTime,
			SweepTime:     gd.SweepTime,
			TotalTime:     gd.TotalTime,
		}
	}

	eventPtrs := make([]*trace.Event, len(events))
	for i := range events {
		eventPtrs[i] = &events[i]
	}

	return eventPtrs, gdescs, nil
}

func Marshal(w io.Writer, events []*trace.Event, gdesc map[uint64]*trace.GDesc) error {
	var frames []*trace.Frame
	framesID := make(map[*trace.Frame]int)
	eventsID := make(map[*trace.Event]int, len(events))

	for i, ev := range events {
		eventsID[ev] = i
		for _, f := range ev.Stk {
			if _, ok := framesID[f]; !ok {
				frames = append(frames, f)
				framesID[f] = len(frames) - 1
			}
		}
	}

	// Header
	h := Header{
		Events: int64(len(events)),
		Frames: int64(len(frames)),
		GDescs: int64(len(gdesc)),
	}
	buf, err := h.Marshal(nil)
	if err != nil {
		return err
	}
	if _, err := w.Write(buf); err != nil {
		return err
	}

	// Events
	for _, ev := range events {
		e := Event{
			Type:  ev.Type,
			Ts:    ev.Ts,
			P:     int32(ev.P),
			G:     ev.G,
			StkID: ev.StkID,
			Stk:   toFrameIDs(framesID, ev.Stk),
			Args:  ev.Args,
			Link:  toEventID(eventsID, ev.Link),
			Off:   int32(ev.Off),
		}
		buf, err = e.Marshal(buf)
		if err != nil {
			return err
		}
		if _, err := w.Write(buf); err != nil {
			return err
		}
	}

	// Frames
	for _, f := range frames {
		f := Frame{
			PC:   f.PC,
			Fn:   f.Fn,
			File: f.File,
			Line: int64(f.Line),
		}
		buf, err = f.Marshal(buf)
		if err != nil {
			return err
		}
		if _, err := w.Write(buf); err != nil {
			return err
		}
	}

	// GDesc
	for _, gd := range gdesc {
		g := GDesc{
			ID:            gd.ID,
			Name:          gd.Name,
			PC:            gd.PC,
			CreationTime:  gd.CreationTime,
			StartTime:     gd.StartTime,
			EndTime:       gd.EndTime,
			ExecTime:      gd.ExecTime,
			SchedWaitTime: gd.SchedWaitTime,
			IOTime:        gd.IOTime,
			BlockTime:     gd.BlockTime,
			SyscallTime:   gd.SyscallTime,
			GCTime:        gd.GCTime,
			SweepTime:     gd.SweepTime,
			TotalTime:     gd.TotalTime,
		}
		buf, err = g.Marshal(buf)
		if err != nil {
			return err
		}
		if _, err := w.Write(buf); err != nil {
			return err
		}
	}
	return nil
}
