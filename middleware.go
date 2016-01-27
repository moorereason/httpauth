package httpauth

import (
	"fmt"
	"net/http"

	"github.com/goadesign/goa"
)

// ValidationFunc is a function type that takes a goa context along with a
// username and password.  Returning nil denotes success.  The goa context is
// available in case information needs to be inserted into the context.
type ValidationFunc func(ctx *goa.Context, username, password string) error

// Specification describes the HTTP auth properties.
type Specification struct {
	ValidationProvider ValidationFunc

	// The auth realm
	Realm string

	// Log successful authentications.
	LogSuccess bool

	// Log failed authentications.
	LogFailure bool
}

// BasicMiddleware is a middleware that provides Basic HTTP authentication.
func BasicMiddleware(spec *Specification) goa.Middleware {
	if spec.Realm == "" {
		spec.Realm = "Restricted"
	}
	return func(h goa.Handler) goa.Handler {
		return func(ctx *goa.Context) error {
			if spec.ValidationProvider == nil {
				return ctx.Bug("Basic auth validation function undefined")
			}

			username, password, ok := parseBasicAuth(ctx)
			if !ok {
				return unauthorized(ctx, spec)
			}

			var err error
			err = spec.ValidationProvider(ctx, username, password)
			if err != nil {
				if spec.LogFailure {
					ctx.Info("httpauth.basic failure", "user", username, "err", err)
				}
				return unauthorized(ctx, spec)
			}

			if spec.LogSuccess {
				ctx.Info("httpauth.basic success", "user", username)
			}

			return h(ctx)
		}
	}
}

// unauthorized sets the appropriate WWW-Authenticate header prior to sending an
// Unauthorized HTTP response.
func unauthorized(ctx *goa.Context, spec *Specification) error {
	ctx.Header().Set("WWW-Authenticate", fmt.Sprintf("Basic realm=%q", spec.Realm))
	// return ctx.Respond(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
	return ctx.Respond(http.StatusUnauthorized, map[string]interface{}{"ID": -1, "Title": "Unauthorized", "Msg": "Unauthorized Request"})
}
