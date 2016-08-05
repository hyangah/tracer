package shared

import (
	"fmt"
	"io"
	"time"
	"unsafe"
)

var (
	_ = unsafe.Sizeof(0)
	_ = io.ReadFull
	_ = time.Now()
)

type Header struct {
	Events int64
	Frames int64
	GDescs int64
}

func (d *Header) Size() (s uint64) {

	s += 24
	return
}
func (d *Header) Marshal(buf []byte) ([]byte, error) {
	size := d.Size()
	{
		if uint64(cap(buf)) >= size {
			buf = buf[:size]
		} else {
			buf = make([]byte, size)
		}
	}
	i := uint64(0)

	{

		buf[0+0] = byte(d.Events >> 0)

		buf[1+0] = byte(d.Events >> 8)

		buf[2+0] = byte(d.Events >> 16)

		buf[3+0] = byte(d.Events >> 24)

		buf[4+0] = byte(d.Events >> 32)

		buf[5+0] = byte(d.Events >> 40)

		buf[6+0] = byte(d.Events >> 48)

		buf[7+0] = byte(d.Events >> 56)

	}
	{

		buf[0+8] = byte(d.Frames >> 0)

		buf[1+8] = byte(d.Frames >> 8)

		buf[2+8] = byte(d.Frames >> 16)

		buf[3+8] = byte(d.Frames >> 24)

		buf[4+8] = byte(d.Frames >> 32)

		buf[5+8] = byte(d.Frames >> 40)

		buf[6+8] = byte(d.Frames >> 48)

		buf[7+8] = byte(d.Frames >> 56)

	}
	{

		buf[0+16] = byte(d.GDescs >> 0)

		buf[1+16] = byte(d.GDescs >> 8)

		buf[2+16] = byte(d.GDescs >> 16)

		buf[3+16] = byte(d.GDescs >> 24)

		buf[4+16] = byte(d.GDescs >> 32)

		buf[5+16] = byte(d.GDescs >> 40)

		buf[6+16] = byte(d.GDescs >> 48)

		buf[7+16] = byte(d.GDescs >> 56)

	}
	return buf[:i+24], nil
}

func (d *Header) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	{

		d.Events = 0 | (int64(buf[0+0]) << 0) | (int64(buf[1+0]) << 8) | (int64(buf[2+0]) << 16) | (int64(buf[3+0]) << 24) | (int64(buf[4+0]) << 32) | (int64(buf[5+0]) << 40) | (int64(buf[6+0]) << 48) | (int64(buf[7+0]) << 56)

	}
	{

		d.Frames = 0 | (int64(buf[0+8]) << 0) | (int64(buf[1+8]) << 8) | (int64(buf[2+8]) << 16) | (int64(buf[3+8]) << 24) | (int64(buf[4+8]) << 32) | (int64(buf[5+8]) << 40) | (int64(buf[6+8]) << 48) | (int64(buf[7+8]) << 56)

	}
	{

		d.GDescs = 0 | (int64(buf[0+16]) << 0) | (int64(buf[1+16]) << 8) | (int64(buf[2+16]) << 16) | (int64(buf[3+16]) << 24) | (int64(buf[4+16]) << 32) | (int64(buf[5+16]) << 40) | (int64(buf[6+16]) << 48) | (int64(buf[7+16]) << 56)

	}
	return i + 24, nil
}

type Event struct {
	Type  byte
	Ts    int64
	P     int32
	G     uint64
	StkID uint64
	Stk   []int64
	Args  [3]uint64
	Link  int64
	Off   int32
}

func (d *Event) Size() (s uint64) {

	{

		t := uint32(d.P)
		t <<= 1
		if d.P < 0 {
			t = ^t
		}
		for t >= 0x80 {
			t >>= 7
			s++
		}
		s++

	}
	{
		l := uint64(len(d.Stk))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}

		for k := range d.Stk {

			{

				t := uint64(d.Stk[k])
				t <<= 1
				if d.Stk[k] < 0 {
					t = ^t
				}
				for t >= 0x80 {
					t >>= 7
					s++
				}
				s++

			}

		}

	}
	{
		for k := range d.Args {
			_ = k

			s += 8

		}
	}
	{

		t := uint64(d.Link)
		t <<= 1
		if d.Link < 0 {
			t = ^t
		}
		for t >= 0x80 {
			t >>= 7
			s++
		}
		s++

	}
	{

		t := uint32(d.Off)
		t <<= 1
		if d.Off < 0 {
			t = ^t
		}
		for t >= 0x80 {
			t >>= 7
			s++
		}
		s++

	}
	s += 25
	return
}
func (d *Event) Marshal(buf []byte) ([]byte, error) {
	size := d.Size()
	{
		if uint64(cap(buf)) >= size {
			buf = buf[:size]
		} else {
			buf = make([]byte, size)
		}
	}
	i := uint64(0)

	{
		buf[0] = d.Type
	}
	{

		buf[0+1] = byte(d.Ts >> 0)

		buf[1+1] = byte(d.Ts >> 8)

		buf[2+1] = byte(d.Ts >> 16)

		buf[3+1] = byte(d.Ts >> 24)

		buf[4+1] = byte(d.Ts >> 32)

		buf[5+1] = byte(d.Ts >> 40)

		buf[6+1] = byte(d.Ts >> 48)

		buf[7+1] = byte(d.Ts >> 56)

	}
	{

		t := uint32(d.P)

		t <<= 1
		if d.P < 0 {
			t = ^t
		}

		for t >= 0x80 {
			buf[i+9] = byte(t) | 0x80
			t >>= 7
			i++
		}
		buf[i+9] = byte(t)
		i++

	}
	{

		buf[i+0+9] = byte(d.G >> 0)

		buf[i+1+9] = byte(d.G >> 8)

		buf[i+2+9] = byte(d.G >> 16)

		buf[i+3+9] = byte(d.G >> 24)

		buf[i+4+9] = byte(d.G >> 32)

		buf[i+5+9] = byte(d.G >> 40)

		buf[i+6+9] = byte(d.G >> 48)

		buf[i+7+9] = byte(d.G >> 56)

	}
	{

		buf[i+0+17] = byte(d.StkID >> 0)

		buf[i+1+17] = byte(d.StkID >> 8)

		buf[i+2+17] = byte(d.StkID >> 16)

		buf[i+3+17] = byte(d.StkID >> 24)

		buf[i+4+17] = byte(d.StkID >> 32)

		buf[i+5+17] = byte(d.StkID >> 40)

		buf[i+6+17] = byte(d.StkID >> 48)

		buf[i+7+17] = byte(d.StkID >> 56)

	}
	{
		l := uint64(len(d.Stk))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+25] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+25] = byte(t)
			i++

		}
		for k := range d.Stk {

			{

				t := uint64(d.Stk[k])

				t <<= 1
				if d.Stk[k] < 0 {
					t = ^t
				}

				for t >= 0x80 {
					buf[i+25] = byte(t) | 0x80
					t >>= 7
					i++
				}
				buf[i+25] = byte(t)
				i++

			}

		}
	}
	{
		for k := range d.Args {

			{

				buf[i+0+25] = byte(d.Args[k] >> 0)

				buf[i+1+25] = byte(d.Args[k] >> 8)

				buf[i+2+25] = byte(d.Args[k] >> 16)

				buf[i+3+25] = byte(d.Args[k] >> 24)

				buf[i+4+25] = byte(d.Args[k] >> 32)

				buf[i+5+25] = byte(d.Args[k] >> 40)

				buf[i+6+25] = byte(d.Args[k] >> 48)

				buf[i+7+25] = byte(d.Args[k] >> 56)

			}

			i += 8

		}
	}
	{

		t := uint64(d.Link)

		t <<= 1
		if d.Link < 0 {
			t = ^t
		}

		for t >= 0x80 {
			buf[i+25] = byte(t) | 0x80
			t >>= 7
			i++
		}
		buf[i+25] = byte(t)
		i++

	}
	{

		t := uint32(d.Off)

		t <<= 1
		if d.Off < 0 {
			t = ^t
		}

		for t >= 0x80 {
			buf[i+25] = byte(t) | 0x80
			t >>= 7
			i++
		}
		if uint64(len(buf)) <= i+25 {
			fmt.Println("buf len", len(buf), "i+25", i+25)
			fmt.Printf("Event: %+v\n", d)
			fmt.Printf("Size: %d\n", d.Size())
		}

		buf[i+25] = byte(t)
		i++

	}
	return buf[:i+25], nil
}

func (d *Event) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	{
		d.Type = buf[i+0]
	}
	{

		d.Ts = 0 | (int64(buf[i+0+1]) << 0) | (int64(buf[i+1+1]) << 8) | (int64(buf[i+2+1]) << 16) | (int64(buf[i+3+1]) << 24) | (int64(buf[i+4+1]) << 32) | (int64(buf[i+5+1]) << 40) | (int64(buf[i+6+1]) << 48) | (int64(buf[i+7+1]) << 56)

	}
	{

		bs := uint8(7)
		t := uint32(buf[i+9] & 0x7F)
		for buf[i+9]&0x80 == 0x80 {
			i++
			t |= uint32(buf[i+9]&0x7F) << bs
			bs += 7
		}
		i++

		d.P = int32(t >> 1)
		if t&1 != 0 {
			d.P = ^d.P
		}

	}
	{

		d.G = 0 | (uint64(buf[i+0+9]) << 0) | (uint64(buf[i+1+9]) << 8) | (uint64(buf[i+2+9]) << 16) | (uint64(buf[i+3+9]) << 24) | (uint64(buf[i+4+9]) << 32) | (uint64(buf[i+5+9]) << 40) | (uint64(buf[i+6+9]) << 48) | (uint64(buf[i+7+9]) << 56)

	}
	{

		d.StkID = 0 | (uint64(buf[i+0+17]) << 0) | (uint64(buf[i+1+17]) << 8) | (uint64(buf[i+2+17]) << 16) | (uint64(buf[i+3+17]) << 24) | (uint64(buf[i+4+17]) << 32) | (uint64(buf[i+5+17]) << 40) | (uint64(buf[i+6+17]) << 48) | (uint64(buf[i+7+17]) << 56)

	}
	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+25] & 0x7F)
			for buf[i+25]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+25]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		if uint64(cap(d.Stk)) >= l {
			d.Stk = d.Stk[:l]
		} else {
			d.Stk = make([]int64, l)
		}
		for k := range d.Stk {

			{

				bs := uint8(7)
				t := uint64(buf[i+25] & 0x7F)
				for buf[i+25]&0x80 == 0x80 {
					i++
					t |= uint64(buf[i+25]&0x7F) << bs
					bs += 7
				}
				i++

				d.Stk[k] = int64(t >> 1)
				if t&1 != 0 {
					d.Stk[k] = ^d.Stk[k]
				}

			}

		}
	}
	{
		for k := range d.Args {

			{

				d.Args[k] = 0 | (uint64(buf[i+0+25]) << 0) | (uint64(buf[i+1+25]) << 8) | (uint64(buf[i+2+25]) << 16) | (uint64(buf[i+3+25]) << 24) | (uint64(buf[i+4+25]) << 32) | (uint64(buf[i+5+25]) << 40) | (uint64(buf[i+6+25]) << 48) | (uint64(buf[i+7+25]) << 56)

			}

			i += 8

		}
	}
	{

		bs := uint8(7)
		t := uint64(buf[i+25] & 0x7F)
		for buf[i+25]&0x80 == 0x80 {
			i++
			t |= uint64(buf[i+25]&0x7F) << bs
			bs += 7
		}
		i++

		d.Link = int64(t >> 1)
		if t&1 != 0 {
			d.Link = ^d.Link
		}

	}
	{

		bs := uint8(7)
		t := uint32(buf[i+25] & 0x7F)
		for buf[i+25]&0x80 == 0x80 {
			i++
			t |= uint32(buf[i+25]&0x7F) << bs
			bs += 7
		}
		i++

		d.Off = int32(t >> 1)
		if t&1 != 0 {
			d.Off = ^d.Off
		}

	}
	return i + 25, nil
}

type Frame struct {
	PC   uint64
	Fn   string
	File string
	Line int64
}

func (d *Frame) Size() (s uint64) {

	{
		l := uint64(len(d.Fn))

		{

			t := l
			for t >= 0x80 {
				t <<= 7
				s++
			}
			s++

		}
		s += l
	}
	{
		l := uint64(len(d.File))

		{

			t := l
			for t >= 0x80 {
				t <<= 7
				s++
			}
			s++

		}
		s += l
	}
	s += 16
	return
}
func (d *Frame) Marshal(buf []byte) ([]byte, error) {
	size := d.Size()
	{
		if uint64(cap(buf)) >= size {
			buf = buf[:size]
		} else {
			buf = make([]byte, size)
		}
	}
	i := uint64(0)

	{

		buf[0+0] = byte(d.PC >> 0)

		buf[1+0] = byte(d.PC >> 8)

		buf[2+0] = byte(d.PC >> 16)

		buf[3+0] = byte(d.PC >> 24)

		buf[4+0] = byte(d.PC >> 32)

		buf[5+0] = byte(d.PC >> 40)

		buf[6+0] = byte(d.PC >> 48)

		buf[7+0] = byte(d.PC >> 56)

	}
	{
		l := uint64(len(d.Fn))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+8] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+8] = byte(t)
			i++

		}
		copy(buf[i+8:], d.Fn)
		i += l
	}
	{
		l := uint64(len(d.File))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+8] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+8] = byte(t)
			i++

		}
		copy(buf[i+8:], d.File)
		i += l
	}
	{

		buf[i+0+8] = byte(d.Line >> 0)

		buf[i+1+8] = byte(d.Line >> 8)

		buf[i+2+8] = byte(d.Line >> 16)

		buf[i+3+8] = byte(d.Line >> 24)

		buf[i+4+8] = byte(d.Line >> 32)

		buf[i+5+8] = byte(d.Line >> 40)

		buf[i+6+8] = byte(d.Line >> 48)

		buf[i+7+8] = byte(d.Line >> 56)

	}
	return buf[:i+16], nil
}

func (d *Frame) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	{

		d.PC = 0 | (uint64(buf[i+0+0]) << 0) | (uint64(buf[i+1+0]) << 8) | (uint64(buf[i+2+0]) << 16) | (uint64(buf[i+3+0]) << 24) | (uint64(buf[i+4+0]) << 32) | (uint64(buf[i+5+0]) << 40) | (uint64(buf[i+6+0]) << 48) | (uint64(buf[i+7+0]) << 56)

	}
	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+8] & 0x7F)
			for buf[i+8]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+8]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		d.Fn = string(buf[i+8 : i+8+l])
		i += l
	}
	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+8] & 0x7F)
			for buf[i+8]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+8]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		d.File = string(buf[i+8 : i+8+l])
		i += l
	}
	{

		d.Line = 0 | (int64(buf[i+0+8]) << 0) | (int64(buf[i+1+8]) << 8) | (int64(buf[i+2+8]) << 16) | (int64(buf[i+3+8]) << 24) | (int64(buf[i+4+8]) << 32) | (int64(buf[i+5+8]) << 40) | (int64(buf[i+6+8]) << 48) | (int64(buf[i+7+8]) << 56)

	}
	return i + 16, nil
}

type GDesc struct {
	ID            uint64
	Name          string
	PC            uint64
	CreationTime  int64
	StartTime     int64
	EndTime       int64
	ExecTime      int64
	SchedWaitTime int64
	IOTime        int64
	BlockTime     int64
	SyscallTime   int64
	GCTime        int64
	SweepTime     int64
	TotalTime     int64
}

func (d *GDesc) Size() (s uint64) {

	{
		l := uint64(len(d.Name))

		{

			t := l
			for t >= 0x80 {
				t <<= 7
				s++
			}
			s++

		}
		s += l
	}
	s += 104
	return
}
func (d *GDesc) Marshal(buf []byte) ([]byte, error) {
	size := d.Size()
	{
		if uint64(cap(buf)) >= size {
			buf = buf[:size]
		} else {
			buf = make([]byte, size)
		}
	}
	i := uint64(0)

	{

		buf[0+0] = byte(d.ID >> 0)

		buf[1+0] = byte(d.ID >> 8)

		buf[2+0] = byte(d.ID >> 16)

		buf[3+0] = byte(d.ID >> 24)

		buf[4+0] = byte(d.ID >> 32)

		buf[5+0] = byte(d.ID >> 40)

		buf[6+0] = byte(d.ID >> 48)

		buf[7+0] = byte(d.ID >> 56)

	}
	{
		l := uint64(len(d.Name))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+8] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+8] = byte(t)
			i++

		}
		copy(buf[i+8:], d.Name)
		i += l
	}
	{

		buf[i+0+8] = byte(d.PC >> 0)

		buf[i+1+8] = byte(d.PC >> 8)

		buf[i+2+8] = byte(d.PC >> 16)

		buf[i+3+8] = byte(d.PC >> 24)

		buf[i+4+8] = byte(d.PC >> 32)

		buf[i+5+8] = byte(d.PC >> 40)

		buf[i+6+8] = byte(d.PC >> 48)

		buf[i+7+8] = byte(d.PC >> 56)

	}
	{

		buf[i+0+16] = byte(d.CreationTime >> 0)

		buf[i+1+16] = byte(d.CreationTime >> 8)

		buf[i+2+16] = byte(d.CreationTime >> 16)

		buf[i+3+16] = byte(d.CreationTime >> 24)

		buf[i+4+16] = byte(d.CreationTime >> 32)

		buf[i+5+16] = byte(d.CreationTime >> 40)

		buf[i+6+16] = byte(d.CreationTime >> 48)

		buf[i+7+16] = byte(d.CreationTime >> 56)

	}
	{

		buf[i+0+24] = byte(d.StartTime >> 0)

		buf[i+1+24] = byte(d.StartTime >> 8)

		buf[i+2+24] = byte(d.StartTime >> 16)

		buf[i+3+24] = byte(d.StartTime >> 24)

		buf[i+4+24] = byte(d.StartTime >> 32)

		buf[i+5+24] = byte(d.StartTime >> 40)

		buf[i+6+24] = byte(d.StartTime >> 48)

		buf[i+7+24] = byte(d.StartTime >> 56)

	}
	{

		buf[i+0+32] = byte(d.EndTime >> 0)

		buf[i+1+32] = byte(d.EndTime >> 8)

		buf[i+2+32] = byte(d.EndTime >> 16)

		buf[i+3+32] = byte(d.EndTime >> 24)

		buf[i+4+32] = byte(d.EndTime >> 32)

		buf[i+5+32] = byte(d.EndTime >> 40)

		buf[i+6+32] = byte(d.EndTime >> 48)

		buf[i+7+32] = byte(d.EndTime >> 56)

	}
	{

		buf[i+0+40] = byte(d.ExecTime >> 0)

		buf[i+1+40] = byte(d.ExecTime >> 8)

		buf[i+2+40] = byte(d.ExecTime >> 16)

		buf[i+3+40] = byte(d.ExecTime >> 24)

		buf[i+4+40] = byte(d.ExecTime >> 32)

		buf[i+5+40] = byte(d.ExecTime >> 40)

		buf[i+6+40] = byte(d.ExecTime >> 48)

		buf[i+7+40] = byte(d.ExecTime >> 56)

	}
	{

		buf[i+0+48] = byte(d.SchedWaitTime >> 0)

		buf[i+1+48] = byte(d.SchedWaitTime >> 8)

		buf[i+2+48] = byte(d.SchedWaitTime >> 16)

		buf[i+3+48] = byte(d.SchedWaitTime >> 24)

		buf[i+4+48] = byte(d.SchedWaitTime >> 32)

		buf[i+5+48] = byte(d.SchedWaitTime >> 40)

		buf[i+6+48] = byte(d.SchedWaitTime >> 48)

		buf[i+7+48] = byte(d.SchedWaitTime >> 56)

	}
	{

		buf[i+0+56] = byte(d.IOTime >> 0)

		buf[i+1+56] = byte(d.IOTime >> 8)

		buf[i+2+56] = byte(d.IOTime >> 16)

		buf[i+3+56] = byte(d.IOTime >> 24)

		buf[i+4+56] = byte(d.IOTime >> 32)

		buf[i+5+56] = byte(d.IOTime >> 40)

		buf[i+6+56] = byte(d.IOTime >> 48)

		buf[i+7+56] = byte(d.IOTime >> 56)

	}
	{

		buf[i+0+64] = byte(d.BlockTime >> 0)

		buf[i+1+64] = byte(d.BlockTime >> 8)

		buf[i+2+64] = byte(d.BlockTime >> 16)

		buf[i+3+64] = byte(d.BlockTime >> 24)

		buf[i+4+64] = byte(d.BlockTime >> 32)

		buf[i+5+64] = byte(d.BlockTime >> 40)

		buf[i+6+64] = byte(d.BlockTime >> 48)

		buf[i+7+64] = byte(d.BlockTime >> 56)

	}
	{

		buf[i+0+72] = byte(d.SyscallTime >> 0)

		buf[i+1+72] = byte(d.SyscallTime >> 8)

		buf[i+2+72] = byte(d.SyscallTime >> 16)

		buf[i+3+72] = byte(d.SyscallTime >> 24)

		buf[i+4+72] = byte(d.SyscallTime >> 32)

		buf[i+5+72] = byte(d.SyscallTime >> 40)

		buf[i+6+72] = byte(d.SyscallTime >> 48)

		buf[i+7+72] = byte(d.SyscallTime >> 56)

	}
	{

		buf[i+0+80] = byte(d.GCTime >> 0)

		buf[i+1+80] = byte(d.GCTime >> 8)

		buf[i+2+80] = byte(d.GCTime >> 16)

		buf[i+3+80] = byte(d.GCTime >> 24)

		buf[i+4+80] = byte(d.GCTime >> 32)

		buf[i+5+80] = byte(d.GCTime >> 40)

		buf[i+6+80] = byte(d.GCTime >> 48)

		buf[i+7+80] = byte(d.GCTime >> 56)

	}
	{

		buf[i+0+88] = byte(d.SweepTime >> 0)

		buf[i+1+88] = byte(d.SweepTime >> 8)

		buf[i+2+88] = byte(d.SweepTime >> 16)

		buf[i+3+88] = byte(d.SweepTime >> 24)

		buf[i+4+88] = byte(d.SweepTime >> 32)

		buf[i+5+88] = byte(d.SweepTime >> 40)

		buf[i+6+88] = byte(d.SweepTime >> 48)

		buf[i+7+88] = byte(d.SweepTime >> 56)

	}
	{

		buf[i+0+96] = byte(d.TotalTime >> 0)

		buf[i+1+96] = byte(d.TotalTime >> 8)

		buf[i+2+96] = byte(d.TotalTime >> 16)

		buf[i+3+96] = byte(d.TotalTime >> 24)

		buf[i+4+96] = byte(d.TotalTime >> 32)

		buf[i+5+96] = byte(d.TotalTime >> 40)

		buf[i+6+96] = byte(d.TotalTime >> 48)

		buf[i+7+96] = byte(d.TotalTime >> 56)

	}
	return buf[:i+104], nil
}

func (d *GDesc) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	{

		d.ID = 0 | (uint64(buf[i+0+0]) << 0) | (uint64(buf[i+1+0]) << 8) | (uint64(buf[i+2+0]) << 16) | (uint64(buf[i+3+0]) << 24) | (uint64(buf[i+4+0]) << 32) | (uint64(buf[i+5+0]) << 40) | (uint64(buf[i+6+0]) << 48) | (uint64(buf[i+7+0]) << 56)

	}
	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+8] & 0x7F)
			for buf[i+8]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+8]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		d.Name = string(buf[i+8 : i+8+l])
		i += l
	}
	{

		d.PC = 0 | (uint64(buf[i+0+8]) << 0) | (uint64(buf[i+1+8]) << 8) | (uint64(buf[i+2+8]) << 16) | (uint64(buf[i+3+8]) << 24) | (uint64(buf[i+4+8]) << 32) | (uint64(buf[i+5+8]) << 40) | (uint64(buf[i+6+8]) << 48) | (uint64(buf[i+7+8]) << 56)

	}
	{

		d.CreationTime = 0 | (int64(buf[i+0+16]) << 0) | (int64(buf[i+1+16]) << 8) | (int64(buf[i+2+16]) << 16) | (int64(buf[i+3+16]) << 24) | (int64(buf[i+4+16]) << 32) | (int64(buf[i+5+16]) << 40) | (int64(buf[i+6+16]) << 48) | (int64(buf[i+7+16]) << 56)

	}
	{

		d.StartTime = 0 | (int64(buf[i+0+24]) << 0) | (int64(buf[i+1+24]) << 8) | (int64(buf[i+2+24]) << 16) | (int64(buf[i+3+24]) << 24) | (int64(buf[i+4+24]) << 32) | (int64(buf[i+5+24]) << 40) | (int64(buf[i+6+24]) << 48) | (int64(buf[i+7+24]) << 56)

	}
	{

		d.EndTime = 0 | (int64(buf[i+0+32]) << 0) | (int64(buf[i+1+32]) << 8) | (int64(buf[i+2+32]) << 16) | (int64(buf[i+3+32]) << 24) | (int64(buf[i+4+32]) << 32) | (int64(buf[i+5+32]) << 40) | (int64(buf[i+6+32]) << 48) | (int64(buf[i+7+32]) << 56)

	}
	{

		d.ExecTime = 0 | (int64(buf[i+0+40]) << 0) | (int64(buf[i+1+40]) << 8) | (int64(buf[i+2+40]) << 16) | (int64(buf[i+3+40]) << 24) | (int64(buf[i+4+40]) << 32) | (int64(buf[i+5+40]) << 40) | (int64(buf[i+6+40]) << 48) | (int64(buf[i+7+40]) << 56)

	}
	{

		d.SchedWaitTime = 0 | (int64(buf[i+0+48]) << 0) | (int64(buf[i+1+48]) << 8) | (int64(buf[i+2+48]) << 16) | (int64(buf[i+3+48]) << 24) | (int64(buf[i+4+48]) << 32) | (int64(buf[i+5+48]) << 40) | (int64(buf[i+6+48]) << 48) | (int64(buf[i+7+48]) << 56)

	}
	{

		d.IOTime = 0 | (int64(buf[i+0+56]) << 0) | (int64(buf[i+1+56]) << 8) | (int64(buf[i+2+56]) << 16) | (int64(buf[i+3+56]) << 24) | (int64(buf[i+4+56]) << 32) | (int64(buf[i+5+56]) << 40) | (int64(buf[i+6+56]) << 48) | (int64(buf[i+7+56]) << 56)

	}
	{

		d.BlockTime = 0 | (int64(buf[i+0+64]) << 0) | (int64(buf[i+1+64]) << 8) | (int64(buf[i+2+64]) << 16) | (int64(buf[i+3+64]) << 24) | (int64(buf[i+4+64]) << 32) | (int64(buf[i+5+64]) << 40) | (int64(buf[i+6+64]) << 48) | (int64(buf[i+7+64]) << 56)

	}
	{

		d.SyscallTime = 0 | (int64(buf[i+0+72]) << 0) | (int64(buf[i+1+72]) << 8) | (int64(buf[i+2+72]) << 16) | (int64(buf[i+3+72]) << 24) | (int64(buf[i+4+72]) << 32) | (int64(buf[i+5+72]) << 40) | (int64(buf[i+6+72]) << 48) | (int64(buf[i+7+72]) << 56)

	}
	{

		d.GCTime = 0 | (int64(buf[i+0+80]) << 0) | (int64(buf[i+1+80]) << 8) | (int64(buf[i+2+80]) << 16) | (int64(buf[i+3+80]) << 24) | (int64(buf[i+4+80]) << 32) | (int64(buf[i+5+80]) << 40) | (int64(buf[i+6+80]) << 48) | (int64(buf[i+7+80]) << 56)

	}
	{

		d.SweepTime = 0 | (int64(buf[i+0+88]) << 0) | (int64(buf[i+1+88]) << 8) | (int64(buf[i+2+88]) << 16) | (int64(buf[i+3+88]) << 24) | (int64(buf[i+4+88]) << 32) | (int64(buf[i+5+88]) << 40) | (int64(buf[i+6+88]) << 48) | (int64(buf[i+7+88]) << 56)

	}
	{

		d.TotalTime = 0 | (int64(buf[i+0+96]) << 0) | (int64(buf[i+1+96]) << 8) | (int64(buf[i+2+96]) << 16) | (int64(buf[i+3+96]) << 24) | (int64(buf[i+4+96]) << 32) | (int64(buf[i+5+96]) << 40) | (int64(buf[i+6+96]) << 48) | (int64(buf[i+7+96]) << 56)

	}
	return i + 104, nil
}
