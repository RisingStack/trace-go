package trace

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"strconv"
	"sync"
	"unsafe"
)

const (
	idSize  = aes.BlockSize / 2 // 64 bits
	keySize = aes.BlockSize     // 128 bits
)

var (
	ctr []byte
	n   int
	b   []byte
	c   cipher.Block
	m   sync.Mutex
)

func init() {
	buf := make([]byte, keySize+aes.BlockSize)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		panic(err) // /dev/urandom had better work
	}
	c, err = aes.NewCipher(buf[:keySize])
	if err != nil {
		panic(err) // AES had better work
	}
	n = aes.BlockSize
	ctr = buf[keySize:]
	b = make([]byte, aes.BlockSize)
}

// ID represents a unique ID
type ID uint64

// String satisfies the Stringer interface.
// The implementaion comes from the appdash String implementation.
func (i *ID) String() string {
	return fmt.Sprintf("%016x", uint64(*i))
}

// ParseID parses a string into an ID. Return error if the parsing failed.
// The implementaion comes from the appdash ParseID implementation.
func ParseID(s string) (ID, error) {
	i, err := strconv.ParseUint(s, 16, 64)
	if err != nil {
		return 0, err
	}
	return ID(i), nil
}

// MarshalJSON marshal the ID into []byte, in the format returned
// by String. It returns an error if it failed to marshal.
func (i *ID) MarshalJSON() ([]byte, error) {
	return []byte("\"" + i.String() + "\""), nil
}

// UnmarshalJSON reads a JSON and sets the value according to it.
func (i *ID) UnmarshalJSON(b []byte) error {
	s := string(b)
	newID, err := ParseID(s[1 : len(s)-1])
	*i = newID
	return err
}

// NewID returns a new generated unique ID.
// The implementaion comes from the appdash ParseID implementation.
func NewID() ID {
	m.Lock()
	if n == aes.BlockSize {
		c.Encrypt(b, ctr)
		for i := aes.BlockSize - 1; i >= 0; i-- { // increment ctr
			ctr[i]++
			if ctr[i] != 0 {
				break
			}
		}
		n = 0
	}
	id := *(*ID)(unsafe.Pointer(&b[n])) // zero-copy b/c we're arch-neutral
	n += idSize
	m.Unlock()
	return id
}

// SpanID represents a Span with Trace and Parent.
type SpanID struct {
	// Trace contains the ID of the current trace itself.
	Trace ID
	// Span contains the ID of the current span.
	Span ID
	// Parent contains the Span ID of the parent.
	Parent ID
}

// Empty return true if this SpanID is empty, elements are 0s.
func (s *SpanID) Empty() bool {
	return s.Trace == 0
}

// String return the string representation of the SpanID.
// For human consumption.
func (s *SpanID) String() string {
	return fmt.Sprintf("<T:%s,S:%s,P:%s>", s.Trace.String(), s.Span.String(), s.Parent.String())
}

// NewRootSpanID generates a new root SpanID with new Trace ID and Span ID.
// Parent ID is 0.
func NewRootSpanID() SpanID {
	return SpanID{
		Trace: NewID(),
		Span:  NewID(),
	}
}

// NewSpanID generates a new SpanID from parent SpanID.
// The newly generated SpanID inherits the Trace ID from
// parent and Parent ID will be set to the parent's Span ID.
func NewSpanID(parent SpanID) SpanID {
	return SpanID{
		Trace:  parent.Trace,
		Span:   NewID(),
		Parent: parent.Span,
	}
}
