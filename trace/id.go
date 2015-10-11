package trace

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"unsafe"
)

type ID uint64

// appdash String()
func (id *ID) String() string {
	return fmt.Sprintf("%016x", uint64(*id))
}

// appdash ParseID
func ParseID(s string) (ID, error) {
	i, err := strconv.ParseUint(s, 16, 64)
	if err != nil {
		return 0, err
	}
	return ID(i), nil
}

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

// appdash generateID
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

type SpanID struct {
	Trace  ID
	Span   ID
	Parent ID
}

func (s *SpanID) IsRoot() bool {
	return s.Trace == 0
}

func (s *SpanID) String() string {
	return fmt.Sprintf("<T:%s,S:%s,P:%s>", s.Trace.String(), s.Span.String(), s.Parent.String())
}

func NewRootSpanID() SpanID {
	return SpanID{
		Trace: NewID(),
		Span:  NewID(),
	}
}

func NewSpanID(parent SpanID) SpanID {
	return SpanID{
		Trace:  parent.Trace,
		Span:   NewID(),
		Parent: parent.Span,
	}
}

func NewSpanIDFromRequest(req *http.Request) SpanID {
	parentSpanID := ParseSpanID(req)
	if parentSpanID.IsRoot() {
		return NewRootSpanID()
	}
	return NewSpanID(parentSpanID)
}

func ParseSpanID(req *http.Request) SpanID {
	spanIDStr := req.Header.Get(HeaderSpanID)
	traceIDStr := req.Header.Get(HeaderTraceID)
	spanID, err := ParseID(spanIDStr)
	if err != nil {
		log.Println("Failed to parse spanID string value, using 0.")
	}
	traceID, err := ParseID(traceIDStr)
	if err != nil {
		log.Println("Failed to parse traceID string value, using 0.")
	}
	return SpanID{
		Trace: traceID,
		Span:  spanID,
	}
}

func (i *ID) MarshalJSON() ([]byte, error) {
	return []byte("\"" + i.String() + "\""), nil
}

func (i *ID) UnmarshalJSON(b []byte) error {
	s := string(b)
	newID, err := ParseID(s[1 : len(s)-1])
	i = &newID
	return err
}
