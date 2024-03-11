package tabuas

type RingBuffer struct {
	buffer []interface{}
	size   int
	head   int
	tail   int
	count  int
}

func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		buffer: make([]interface{}, size),
		size:   size,
		head:   0,
		tail:   0,
		count:  0,
	}
}

func (rb *RingBuffer) Add(item interface{}) {
	rb.buffer[rb.head] = item
	rb.head = (rb.head + 1) % rb.size
	if rb.count < rb.size {
		rb.count++
	} else {
		rb.tail = (rb.tail + 1) % rb.size
	}
}

func (rb *RingBuffer) Unroll() []interface{} {
	result := make([]interface{}, 0, rb.count)
	for i := 0; i < rb.count; i++ {
		index := (rb.head + rb.size - i - 1) % rb.size
		result = append(result, rb.buffer[index])
	}
	return result
}
