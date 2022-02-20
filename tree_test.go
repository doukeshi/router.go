package router

import (
	"log"
	"net/http"
	"reflect"
	"testing"
)

type tt struct {
	*testing.T
}

func (t tt) deepEqual(actual, expected interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("actual: %v, expeted: %v", actual, expected)
	}
}

type fakerResp struct {
	http.ResponseWriter

	headers    http.Header
	pattern    string
	statusCode int
}

func (fa *fakerResp) Write(b []byte) (int, error) {
	fa.pattern = string(b)
	return 0, nil
}

func (fa *fakerResp) Header() http.Header {
	if fa.headers == nil {
		fa.headers = http.Header{}
	}
	return fa.headers
}

func (fa *fakerResp) WriteHeader(statusCode int) {
	fa.statusCode = statusCode
}

func TestNewNode(t *testing.T) {
	actual := newNode("/:p")
	expected := &node{
		part:     "/:p",
		handlers: make(map[string]http.Handler),
		children: make(map[string]*node),
		child:    nil,
	}
	tt{t}.deepEqual(actual, expected)
}

func TestNewTree(t *testing.T) {
	actual := NewTree()
	expected := &tree{
		root: &node{
			part:     "/",
			handlers: make(map[string]http.Handler),
			children: make(map[string]*node),
			child:    nil,
		},
	}
	tt{t}.deepEqual(actual, expected)
}

func TestSplit(t *testing.T) {
	cases := []struct {
		actual   []string
		expected []string
	}{
		{actual: split(""), expected: nil},
		{actual: split("/"), expected: nil},
		{actual: split("//"), expected: nil},
		{actual: split("///"), expected: nil},
		{actual: split("foo"), expected: []string{"foo"}},
		{actual: split("/foo"), expected: []string{"foo"}},
		{actual: split("//foo"), expected: []string{"foo"}},
		{actual: split("//foo//"), expected: []string{"foo"}},
		{actual: split("/foo/bar"), expected: []string{"foo", "bar"}},
		{actual: split("/foo//bar"), expected: []string{"foo", "bar"}},
		{actual: split("/foo/bar/x"), expected: []string{"foo", "bar", "x"}},
		{actual: split("/foo/bar/:x"), expected: []string{"foo", "bar", ":x"}},
		{actual: split("/foo/bar/:x/"), expected: []string{"foo", "bar", ":x"}},
	}

	tt := tt{t}
	for _, c := range cases {
		tt.deepEqual(c.actual, c.expected)
	}
}

func TestTree(t *testing.T) {
	tree := NewTree()
	routes := []struct {
		method  string
		pattern string
	}{
		{http.MethodGet, "/"},
		{http.MethodGet, "/a"},
		{http.MethodGet, "/b"},
		{http.MethodGet, "/posts"},
		{http.MethodGet, "/posts/foo"},
		{http.MethodGet, "/posts/bar"},
		{http.MethodGet, "/posts/:id"},
		{http.MethodPost, "/posts/:id"},
		{http.MethodPut, "/posts/:id"},
		{http.MethodDelete, "/posts/:id"},
		{http.MethodGet, "/posts/:id/author"},
		{http.MethodGet, "/posts/:id/comments"},
		{http.MethodGet, "/posts/:id/comments/foo"},
		{http.MethodGet, "/posts/:id/comments/bar"},
		{http.MethodGet, "/posts/:id/comments/:id"},
		{http.MethodGet, "/posts/:id/comments/:id/foobar"},
	}

	for _, rou := range routes {
		pattern := rou.pattern
		method := rou.method
		tree.i(
			method,
			pattern,
			http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
				_, _ = rw.Write([]byte(pattern))
			}),
		)
	}

	cases := []struct {
		method     string
		path       string
		statusCode int
		pattern    string
	}{
		{http.MethodGet, "/", 0, "/"},
		{http.MethodGet, "/post", 404, ""},
		{http.MethodGet, "/posts", 0, "/posts"},
		{http.MethodGet, "/posts/foo", 0, "/posts/foo"},
		{http.MethodGet, "/posts/bar", 0, "/posts/bar"},
		{http.MethodGet, "/posts/100", 0, "/posts/:id"},
		{http.MethodPost, "/posts/100", 0, "/posts/:id"},
		{http.MethodPut, "/posts/100", 0, "/posts/:id"},
		{http.MethodPatch, "/posts/100", 405, ""},
		{http.MethodDelete, "/posts/100", 0, "/posts/:id"},
		{http.MethodGet, "/posts/100/foobar", 404, ""},
		{http.MethodGet, "/posts/100/author", 0, "/posts/:id/author"},
		{http.MethodGet, "/posts/100/comments", 0, "/posts/:id/comments"},
		{http.MethodGet, "/posts/100/comments/foo", 0, "/posts/:id/comments/foo"},
		{http.MethodGet, "/posts/100/comments/bar", 0, "/posts/:id/comments/bar"},
		{http.MethodGet, "/posts/100/comments/101", 0, "/posts/:id/comments/:id"},
		{http.MethodGet, "/posts/100/comments/101/foobar", 0, "/posts/:id/comments/:id/foobar"},
	}

	tt := tt{t}
	for _, c := range cases {
		h := tree.lookup(c.method, c.path)

		res := &fakerResp{}
		h.ServeHTTP(res, nil)
		log.Printf("req: [%s, %s], resp: {%d, %s}", c.method, c.path, res.statusCode, res.pattern)

		tt.deepEqual(res.statusCode, c.statusCode)
		if res.statusCode == 0 {
			tt.deepEqual(res.pattern, c.pattern)
		}
	}
}
