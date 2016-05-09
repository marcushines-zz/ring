// Package ring is a simple slice based implementation of a ring buffer.  Ring
// reads from head and writes to head like a stack but dequeues are done from
// tail. This allows for deletion of the oldest values once capacity has been
// reached.
package ring

// Ring implements a circular buffer. It has one different implementation from
// a standard ring buffer in that most reads are expected to be done from head
// rather than tail.
type Ring struct {
	head int // most recent value position
	tail int // oldest value position
	buff []interface{}
}

// New returns a newly initialized Ring of capacity c.
func New(c int) *Ring {
	return &Ring{
		head: -1,
		tail: 0,
		buff: make([]interface{}, c),
	}
}

// Capacity returns the current capacity of r.
func (r Ring) Capacity() int {
	return len(r.buff)
}

// Enqueue adds value e to the head of r.
func (r *Ring) Enqueue(e interface{}) {
	r.set(r.head+1, e)
	old := r.head
	r.head = r.mod(r.head + 1)
	if old != -1 && r.head == r.tail {
		r.tail = r.mod(r.tail + 1)
	}
}

// Dequeue removes and returns the tail value from r. Returns nil, if r is empty.
func (r *Ring) Dequeue() interface{} {
	if r.head == -1 {
		return nil
	}
	v := r.get(r.tail)
	if r.tail == r.head {
		r.head = -1
		r.tail = 0
	} else {
		r.tail = r.mod(r.tail + 1)
	}
	return v
}

// Peek returns the value that Dequeue would have returned without acutally
// removing it. If r is empty return nil.
func (r *Ring) Peek() interface{} {
	if r.head == -1 {
		return nil
	}
	return r.get(r.head)
}

// PeekN peeks n interface{}s deep into r.  If r's capacity is less than n or r
// r contains less than n interface{}s only min(r capacity, r current) is returned.
// The slice returned is copy of r's buffer however the contents of the buffer
// are not copied. If n <= 0 or r is empty, nil is returned.
func (r *Ring) PeekN(n int) []interface{} {
	if n <= 0 {
		return nil
	}
	if r.head == -1 {
		return nil
	}
	b := []interface{}{}
	i := r.head
	for {
		b = append(b, r.buff[i])
		if i == r.tail || len(b) == n {
			break
		}
		i--
		if i < 0 {
			i = r.Capacity() - 1
		}
	}
	return b
}

// Tail returns the value that Dequeue would have returned without acutally
// removing it. If r is empty return nil.
func (r *Ring) Tail() interface{} {
	if r.head == -1 {
		return nil
	}
	return r.get(r.tail)
}

// sets a value in r at the given unmodified index.
func (r *Ring) set(p int, v interface{}) {
	r.buff[r.mod(p)] = v
}

// gets a value in r based at a given unmodified index.
func (r *Ring) get(p int) interface{} {
	return r.buff[r.mod(p)]
}

// returns the modified index of an unmodified index
func (r *Ring) mod(p int) int {
	return p % r.Capacity()
}
