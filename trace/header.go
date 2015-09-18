package trace

import (
	"net/http"

	"github.com/twinj/uuid"
)

const (
	HeaderId = "x-request-id"

	HeaderParentId = "x-parent-id"
)

func GetId(r *http.Request) string {
	return r.Header.Get(HeaderId)
}

func SetIdIfMissing(r *http.Request, parent *http.Request) {
	id := GetId(r)
	if id == "" {
		if parent == nil {
			id = NewEvent(RequestReceived).RequestId
		} else {
			parentId := GetId(parent)
			if parentId == "" {
				id = uuid.NewV4().String()
			} else {
				id = parentId
			}
		}
		r.Header.Add(HeaderId, id)
	}
}
