package httpauth

import (
	"encoding/base64"
	"strings"

	"github.com/goadesign/goa"
)

const basicScheme string = "Basic "

func parseBasicAuth(ctx *goa.Context) (username, password string, ok bool) {
	auth := ctx.Request().Header.Get("Authorization")
	if auth == "" {
		return
	}

	if !strings.HasPrefix(auth, basicScheme) {
		return
	}

	c, err := base64.StdEncoding.DecodeString(auth[len(basicScheme):])
	if err != nil {
		return
	}

	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return
	}

	return cs[:s], cs[s+1:], true
}
