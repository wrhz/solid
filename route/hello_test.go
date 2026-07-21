package route_test

import (
	"net/http/httptest"
	"testing"

	"github.com/wrhz/solid/test"
)

func TestHello(t *testing.T) {
	err := test.TestRoute(func() {
		req := httptest.NewRequest("GET", "/hello", nil)
		rec := httptest.NewRecorder()

		test.ServeHTTP(rec, req)

		t.Log(rec.Body.String())
	})

	if err != nil {
		t.Fatal(err)
	}
}