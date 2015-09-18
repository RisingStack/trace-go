package trace

import "net/http"

func Get(url string, parent *http.Request) (resp *http.Response, err error) {
	c := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	SetIdIfMissing(req, parent)
	return c.Do(req)
}

func Instrument(fn http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		SetIdIfMissing(r, nil)
		fn(rw, r)
	}
}
