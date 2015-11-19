package trace

import "testing"

const (
	TestIDStr = "094f13e642b9b07f"
	TestIDVal = 670776749184364671
)

func TestNewID(t *testing.T) {
	i1 := NewID()
	i2 := NewID()
	if i1 == i2 {
		t.Errorf("Two generated IDs are the same: %d, %d", i1, i2)
	}
	if i1 == 0 || i2 == 0 {
		t.Errorf("One of the ids are 0: %d, %d", i1, i2)
	}
}

func TestParseID(t *testing.T) {
	i, err := ParseID(TestIDStr)
	if err != nil {
		t.Error("Error occured parsing the ID: ", err.Error())
	}
	if i != TestIDVal {
		t.Errorf("The parsed value is not matching: %d", i)
	}
}

func TestWrongParseID(t *testing.T) {
	_, err := ParseID("x")
	if err == nil {
		t.Error("Error should have happened for wrong input parsing")
	}
}

func TestString(t *testing.T) {
	i := ID(TestIDVal)
	if i.String() != TestIDStr {
		t.Errorf("String ID value does not match: %s", i.String())
	}
}

func TestEquality(t *testing.T) {
	i1 := ID(TestIDVal)
	i2 := ID(TestIDVal)
	if i1 != i2 {
		t.Errorf("IDs should be equal but they are not: %016x != %016x", uint64(i1), uint64(i2))
	}
	if i1.String() != i2.String() {
		t.Errorf("IDs should be equal but they are not: %s != %s", i1.String(), i2.String())
	}
	i1v, err := ParseID(i1.String())
	if err != nil {
		t.Error("Error occured parsing the ID: ", err.Error())
	}
	i2v, err := ParseID(i2.String())
	if err != nil {
		t.Error("Error occured parsing the ID: ", err.Error())
	}
	if i1v != i2v {
		t.Errorf("IDs should be equal but they are not: %s != %s", i1.String(), i2.String())
	}
}

func TestNewRootSpanID(t *testing.T) {
	i := NewRootSpanID()
	if i.Trace == i.Span {
		t.Errorf("Trace and Span should not be equal: %016x !=  %016x", uint64(i.Trace), uint64(i.Span))
	}
	if i.Parent != 0 {
		t.Errorf("Root SpanID shouldn't have Parent ID: %016x != %016x", uint64(i.Parent), uint64(0))
	}
}

func TestRootIsRoot(t *testing.T) {
	i := NewRootSpanID()
	if i.Empty() {
		t.Error("Newly created spanID should not be Empty", i.String())
	}
}

func TestNewSpanID(t *testing.T) {
	p := NewRootSpanID()
	i := NewSpanID(p)
	if i.Trace == i.Span {
		t.Errorf("Trace and Span should not be equal: %016x !=  %016x", uint64(i.Trace), uint64(i.Span))
	}
	if i.Trace != p.Trace {
		t.Errorf("SpanID created from parent should inherit Trace: %s", i.String())
	}
	if i.Parent != p.Span {
		t.Errorf("SpanID created from parent should have Parent: %s", i.String())
	}
}

func TestNotRootIsNotRoot(t *testing.T) {
	p := NewRootSpanID()
	i := NewSpanID(p)
	if i.Empty() {
		t.Error("Newly created spanID should not be Empty", i.String())
	}
}

func TestEmpty(t *testing.T) {
	i := SpanID{}
	if !i.Empty() {
		t.Error("Newly created spanID without anything should be Empty", i.String())
	}
}

func TestSpanIDString(t *testing.T) {
	i := SpanID{
		Trace:  ID(TestIDVal),
		Span:   ID(TestIDVal),
		Parent: ID(TestIDVal),
	}
	e := "<T:094f13e642b9b07f,S:094f13e642b9b07f,P:094f13e642b9b07f>"
	a := i.String()
	if e != a {
		t.Errorf("String method of SpanID did not return the expected result. \nE: %s\nA: %s", e, a)
	}
}

func TestMarshallID(t *testing.T) {
	i := ID(TestIDVal)
	b, err := i.MarshalJSON()
	if err != nil {
		t.Error("Failed to marshal ID. ", err)
	}
	e := "\"094f13e642b9b07f\""
	a := string(b)
	if a != e {
		t.Errorf("Marshalled JSON is not matching what is expected: %s != %s", e, a)
	}
}

func TestUnMarshallID(t *testing.T) {
	var i ID
	b := []byte("\"" + TestIDStr + "\"")
	err := i.UnmarshalJSON(b)
	if err != nil {
		t.Error("Failed to unmarshal ID. ", err)
	}
	e := ID(TestIDVal)
	if i != e {
		t.Errorf("Unmarshalled ID doesn't match with expected: %016x != %016x", e, i)
	}
}
