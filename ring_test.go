package ring

import (
	"fmt"
	"testing"
)

type testElement struct {
	v int
}

func (t testElement) String() string {
	return fmt.Sprintf("%d", t.v)
}

func TestAddCap(t *testing.T) {
	cap := 100
	r := New(cap)
	if r.Capacity() != cap {
		t.Fatalf("Capacity() failed: got %v, want %v", r.Capacity(), cap)
	}
	for i := 0; i < cap; i++ {
		tElem := testElement{v: i}
		r.Enqueue(tElem)
		if got := r.PeekN(1)[0].(testElement).v; got != i {
			t.Errorf("Enqueue(%v) failed: got %v, want %v", tElem, got, i)
		}
	}
}

func TestAddCapHundred(t *testing.T) {
	cap := 100
	r := New(cap)
	if r.Capacity() != cap {
		t.Fatalf("Capacity() failed: got %v, want %v", r.Capacity(), cap)
	}
	for i := 0; i < cap*100; i++ {
		tElem := testElement{v: i}
		r.Enqueue(tElem)
		if got := r.PeekN(1)[0].(testElement).v; got != i {
			t.Errorf("Enqueue(%v) failed: got %v, want %v", tElem, got, i)
		}
	}
	v := r.Tail()
	if v == nil {
		t.Error("Tail() returned nil with elements left.")
	}
	val := v.(testElement).v
	for i := 0; i < cap; i++ {
		tElem := r.Dequeue()
		if got := tElem.(testElement).v; got != val {
			t.Errorf("Dequeue() failed: got %v, want %v", got, val)
		}
		val++
	}
	if r.Tail() != nil {
		t.Errorf("Tail() failed: got %v, want nil", r.Tail())
	}
}

func TestOperations(t *testing.T) {
	tests := []struct {
		e int // Enqueue operations.
		c int // Capacity of Ring.
		d int // Dequeue operations.
		p int // PeekN depth.
	}{
		{e: 1, c: 1, d: 1, p: 1},
		{e: 10, c: 1, d: 10, p: 1},
		{e: 10, c: 10, d: 10, p: 10},
		{e: 1, c: 10, d: 10, p: 100},
		{e: 10, c: 1, d: 0, p: 0},
		{e: 10, c: 2, d: 4, p: 4},
		{e: 10, c: 3, d: 4, p: 4},
		{e: 12, c: 7, d: 4, p: 4},
		{e: 13, c: 7, d: 4, p: 4},
		{e: 14, c: 7, d: 4, p: 4},
	}

	for _, tt := range tests {
		t.Logf("Iteration: %+v", tt)
		r := New(tt.c)
		dElem := r.Dequeue()
		if dElem != nil {
			t.Errorf("%+v: empty Dequeue() failed: got %v", tt, dElem)
		}
		dElem = r.Tail()
		if dElem != nil {
			t.Errorf("%+v: empty Tail() failed: got %v", tt, dElem)
		}
		dElem = r.Peek()
		if dElem != nil {
			t.Errorf("%+v: empty Peek() failed: got %v", tt, dElem)
		}
		dList := r.PeekN(tt.p)
		if len(dList) != 0 {
			t.Errorf("%+v: empty PeekN(%d) failed: got %v", tt, tt.c, dList)
		}
		if got, want := r.Capacity(), tt.c; got != want {
			t.Errorf("%+v: Capacity() failed: got %v, want %v", tt, got, want)
		}
		h := []*testElement{}
		for i := 1; i <= tt.e; i++ {
			e := &testElement{v: i}
			h = append(h, e)
			r.Enqueue(*e)
		}
		dElem = r.Peek()
		if got, want := dElem.(testElement).v, h[len(h)-1].v; got != want {
			t.Errorf("%+v: Peek() failed: got %v, want %v", tt, got, want)
		}
		dElem = r.Tail()
		var elem *testElement
		switch {
		case tt.e == 0:
			elem = nil
		case tt.c == len(h):
			elem = h[0]
		case tt.c > len(h):
			elem = h[len(h)-1]
		default:
			elem = h[len(h)-r.Capacity()]
		}
		if got, want := dElem.(testElement).v, elem.v; got != want {
			t.Errorf("%+v: Tail() failed: got %v, want %v", tt, got, want)
		}
		dList = r.PeekN(tt.p)
		eDepth := func() int {
			if tt.c <= tt.e && tt.c <= tt.p {
				return tt.c
			}
			if tt.e <= tt.c && tt.e <= tt.p {
				return tt.e
			}
			return tt.p
		}()
		if got, want := len(dList), eDepth; got != want {
			t.Errorf("%+v: PeekN(%d) failed: got %v, want %v", tt, tt.p, got, want)
		}
		dElem = r.Dequeue()
		if got, want := dElem.(testElement).v, elem.v; got != want {
			t.Errorf("%+v: Dequeue() failed: got %v, want %v", tt, got, want)
		}
	}
}
