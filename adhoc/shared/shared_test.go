package shared

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/hyangah/tracer/trace"
)

func TestChunk(t *testing.T) {
	oldBufSize := bufSize
	defer func() { bufSize = oldBufSize }()

	bufSize = 16

	var v [32]byte
	for i := 0; i < 32; i++ {
		v[i] = byte(i)
	}

	var got []byte
	c := &chunk{r: bytes.NewReader(v[:])}
	for i := 1; i < 64; i *= 2 {
		b, err := c.Bytes(i)
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Fatalf("c.Byte(%d) returned unexpected error %v after reading %d bytes (%v)", i, err, len(got), c)
		}
		n := i
		if len(b) < n {
			n = len(b)
		}

		if err := c.Advance(len(b) + 1); err == nil {
			t.Errorf("c.Advance(invalid) ucceeded: %+v", c)
		}
		got = append(got, b[:n]...)
		if err := c.Advance(n); err != nil {
			t.Errorf("c.Advance(%d) = %v;  %+v", n, err, c)
		}
	}
	if want := v[:]; !bytes.Equal(got, want) {
		t.Errorf("Got %v, want %v", got, want)
	}
}

func TestMarshalUnmarshal(t *testing.T) {
	frames := []trace.Frame{
		{Fn: "func1"},
		{Fn: "func2"},
		{Fn: "func3"},
	}
	ev0 := trace.Event{
		Off:   0,
		Type:  1,
		Ts:    3,
		P:     4,
		G:     5,
		StkID: 6,
		Stk:   []*trace.Frame{&frames[0], &frames[1]},
		Args:  [3]uint64{1, 2, 3},
	}
	ev1 := trace.Event{
		Off:   10,
		Type:  11,
		Ts:    13,
		P:     14,
		G:     15,
		StkID: 16,
		Stk:   []*trace.Frame{&frames[0], &frames[1], &frames[2]},
		Link:  &ev0,
	}
	ev2 := trace.Event{
		Off:  20,
		Type: 21,
		Ts:   23,
	}
	gdesc := trace.GDesc{
		ID: 12345,
	}

	events := []*trace.Event{&ev0, &ev1, &ev2}
	gdescs := map[uint64]*trace.GDesc{12345: &gdesc}
	var buf bytes.Buffer
	if err := Marshal(&buf, events, gdescs); err != nil {
		t.Fatalf("Marshal = %v", err)
	}
	t.Logf("size = %d", buf.Len())
	gotEvents, gotGDesc, err := Unmarshal(&buf)
	if err != nil {
		t.Fatalf("Unmarshal = %v", err)
	}
	if !reflect.DeepEqual(gotEvents, events) {
		t.Errorf("Got events %+v, want %+v", gotEvents, events)
	}
	if !reflect.DeepEqual(gotGDesc, gdescs) {
		t.Errorf("Got gdescs %+v, want %+v", gotGDesc, gdescs)
	}

}
