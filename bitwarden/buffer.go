package bitwarden

type Buffer struct {
	buf    []byte
	offset int
}

func (b *Buffer) Write(p []byte) (n int, err error) {
	if cap(b.buf) < len(p)+b.offset {
		// The buffer is too small, create a new one and zero the current one.
		newBuf := make([]byte, b.offset+len(p))
		for i, j := range b.buf {
			b.buf[i] = 0
			newBuf[i] = j
		}

		for i, j := range p {
			newBuf[i+b.offset] = j
			p[i] = 0
		}

		if len(b.buf) != 0 {
			b.buf = b.buf[:0]
		}

		b.buf = newBuf
		b.offset += len(p)
	} else {
		copy(b.buf[b.offset:], p)
		b.offset += len(p)
	}

	return len(p), nil
}

func (b *Buffer) Bytes() []byte {
	return b.buf
}

func (b *Buffer) Clear() {
	for i := range b.buf {
		b.buf[i] = 0
	}
	b.buf = make([]byte, 0, 64)
	b.offset = 0
}

func NewBuffer() *Buffer {
	return &Buffer{
		buf:    make([]byte, 0, 64),
		offset: 0,
	}
}
