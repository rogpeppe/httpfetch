package httpfetch_test

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	
	gc "gopkg.in/check.v1"
	jujutesting "github.com/juju/testing"

	"github.com/rogpeppe/httpfetch"
)

type suite struct {
	jujutesting.CleanupSuite
}

var _ = gc.Suite(&suite{})

func (s *suite) TestGetURLAsStringSuccess(c *gc.C) {
	text := "hello, world\n"
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(text))
	}
	srv := httptest.NewServer(http.HandlerFunc(handler))
	defer srv.Close()

	got, err := httpfetch.GetURLAsString(srv.URL)
	c.Assert(err, gc.Equals, nil)
	c.Assert(got, gc.Equals, text)
}

func (s *suite) TestGetURLAsStringNotFound(c *gc.C) {
	srv := httptest.NewServer(http.HandlerFunc(http.NotFound))
	defer srv.Close()

	got, err := httpfetch.GetURLAsString(srv.URL)
	c.Assert(err, gc.ErrorMatches, `GET returned unexpected status "404 Not Found"`)
	c.Assert(got, gc.Equals, "")
}

func (s *suite) TestGetURLAsStringHTTPGetError(c *gc.C) {
	s.PatchValue(httpfetch.HTTPGet, func(u string) (*http.Response, error) {
		return nil, errors.New("crash and burn")
	})
	got, err := httpfetch.GetURLAsString("http://0.1.2.3/")
	c.Assert(err, gc.ErrorMatches, "GET failed: crash and burn")
	c.Assert(got, gc.Equals, "")
}

func (s *suite) TestGetURLAsStringClosesBody(c *gc.C) {
	body := &closeChecker{
		Reader: strings.NewReader("hello"),
	}
	s.PatchValue(httpfetch.HTTPGet, func(u string) (*http.Response, error) {
		return &http.Response{
			Body: body,
			StatusCode: http.StatusOK,
		}, nil
	})
	got, err := httpfetch.GetURLAsString("http://0.1.2.3/")
	c.Assert(err, gc.IsNil)
	c.Assert(got, gc.Equals, "hello")
	c.Assert(body.closed, gc.Equals, true)
}

func (s *suite) TestGetURLAsStringReadBodyError(c *gc.C) {
	body := ioutil.NopCloser(errorReader("an error"))
	s.PatchValue(httpfetch.HTTPGet, func(u string) (*http.Response, error) {
		return &http.Response{
			Body: body,
			StatusCode: http.StatusOK,
		}, nil
	})
	got, err := httpfetch.GetURLAsString("http://0.1.2.3/")
	c.Assert(err, gc.ErrorMatches, "cannot read body: an error")
	c.Assert(got, gc.Equals, "")
}

type errorReader string

func (r errorReader) Read([]byte) (int, error) {
	return 0, errors.New(string(r))
}

type closeChecker struct {
	io.Reader
	closed bool
}

func (r *closeChecker) Close() error {
	r.closed = true
	return nil
}
