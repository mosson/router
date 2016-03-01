package router

import (
	"net/http"
	"net/url"
	"testing"
)

type handler struct {
	testFn func(http.ResponseWriter, *http.Request)
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.testFn(w, r)
}

func TestRouter(t *testing.T) {
	cnt := 0
	router := New()
	router.RegisterFn("/api/v1/users/:id", func(w http.ResponseWriter, r *http.Request) {
		cnt = cnt + 1

		if r.Form["id"][0] != "123" {
			t.Errorf("expected 123, actual %v", r.Form["id"][0])
		}
	})

	router.Register("/api/v1/users/:id/:name", &handler{testFn: func(w http.ResponseWriter, r *http.Request) {
		cnt = cnt + 1

		if r.Form["id"][0] != "456" {
			t.Errorf("expected 456, actual %v", r.Form["id"][0])
		}

		if r.Form["name"][0] != "789" {
			t.Errorf("expected 789, actual %v", r.Form["id"][0])
		}
	}})

	r, err := http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		t.Error(err)
	}
	r.Form = make(url.Values)

	router.Handle(nil, r, "/api/v1/users/123")

	if cnt != 1 {
		t.Errorf("expected 1, actual %v", cnt)
	}

	router.Handle(nil, r, "/api/v1/users/456/789")

	if cnt != 2 {
		t.Errorf("expected 2, actual %v", cnt)
	}

	if router.Handle(nil, r, "/api/v1/hoge") == true {
		t.Errorf("expected false")
	}

}
