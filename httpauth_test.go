package httpauth_test

import (
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/moorereason/httpauth"

	"github.com/goadesign/goa"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Httpauth Middleware", func() {
	var ctx *goa.Context
	var spec *httpauth.Specification
	var req *http.Request
	var rw *TestResponseWriter
	var err error
	var user = "test"
	var pass = "testpass"
	var authString = "Basic " + base64.StdEncoding.EncodeToString([]byte(user+":"+pass))

	validFunc := func(ctx *goa.Context, u, p string) error {
		if u == user && p == pass {
			return nil
		}
		return errors.New("failed")
	}

	handler := func(ctx *goa.Context) error {
		ctx.Respond(200, "ok")
		return nil
	}

	BeforeEach(func() {
		req, err = http.NewRequest("GET", "/goo", nil)
		Ω(err).ShouldNot(HaveOccurred())

		rw = new(TestResponseWriter)
		rw.ParentHeader = make(http.Header)

		s := goa.New("test")
		s.SetEncoder(goa.JSONEncoderFactory(), true, "*/*")
		ctx = goa.NewContext(nil, s, req, rw, nil)

		spec = &httpauth.Specification{
			ValidationProvider: validFunc,
		}
	})

	It("handles valid credentials", func() {
		req.Header.Add("Authorization", authString)

		auth := httpauth.BasicMiddleware(spec)(handler)
		Ω(auth(ctx)).ShouldNot(HaveOccurred())
		Ω(ctx.ResponseStatus()).Should(Equal(http.StatusOK))
		Ω(rw.Body).Should(Equal([]byte("\"ok\"\n")))
	})

	It("handles invalid credentials", func() {
		auth := httpauth.BasicMiddleware(spec)(handler)
		Ω(auth(ctx)).ShouldNot(HaveOccurred())
		Ω(ctx.ResponseStatus()).Should(Equal(http.StatusUnauthorized))
		Ω(ctx.Header()).Should(HaveKey("Www-Authenticate"))
		Ω(ctx.Header().Get("Www-Authenticate")).Should(Equal(`Basic realm="Restricted"`))
		Ω(rw.Body).Should(ContainSubstring("Unauthorized"))
	})

	It("sets a custom realm", func() {
		spec.Realm = "Custom"

		auth := httpauth.BasicMiddleware(spec)(handler)
		Ω(auth(ctx)).ShouldNot(HaveOccurred())
		Ω(ctx.ResponseStatus()).Should(Equal(http.StatusUnauthorized))
		Ω(ctx.Header()).Should(HaveKey("Www-Authenticate"))
		Ω(ctx.Header().Get("Www-Authenticate")).Should(Equal(`Basic realm="Custom"`))
		Ω(rw.Body).Should(ContainSubstring("Unauthorized"))
	})
})

type TestResponseWriter struct {
	ParentHeader http.Header
	Body         []byte
	Status       int
}

func (t *TestResponseWriter) Header() http.Header {
	return t.ParentHeader
}

func (t *TestResponseWriter) Write(b []byte) (int, error) {
	t.Body = append(t.Body, b...)
	return len(b), nil
}

func (t *TestResponseWriter) WriteHeader(s int) {
	t.Status = s
}
