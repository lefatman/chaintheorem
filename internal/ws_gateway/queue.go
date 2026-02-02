package ws_gateway

import "sync"

type outboundQueue struct {
	mu               sync.Mutex
	buf              [][]byte
	head             int
	tail             int
	size             int
	pendingDroppable []byte
	closed           bool
}

func newOutboundQueue(capacity int) *outboundQueue {
	return &outboundQueue{buf: make([][]byte, capacity)}
}

func (q *outboundQueue) Enqueue(frame []byte) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return false
	}
	if q.size == len(q.buf) {
		return false
	}
	q.buf[q.tail] = frame
	q.tail = (q.tail + 1) % len(q.buf)
	q.size++
	return true
}

func (q *outboundQueue) SetDroppable(frame []byte) (replaced []byte, ok bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return nil, false
	}
	replaced = q.pendingDroppable
	q.pendingDroppable = frame
	return replaced, true
}

func (q *outboundQueue) Next() []byte {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.size > 0 {
		frame := q.buf[q.head]
		q.buf[q.head] = nil
		q.head = (q.head + 1) % len(q.buf)
		q.size--
		return frame
	}
	if q.pendingDroppable != nil {
		frame := q.pendingDroppable
		q.pendingDroppable = nil
		return frame
	}
	return nil
}

func (q *outboundQueue) Close() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.closed = true
}
