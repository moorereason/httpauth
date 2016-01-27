httpauth
========

Package httpauth provides a goa middleware that implements Basic HTTP authentication.

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/moorereason/httpauth/blob/master/LICENSE)
[![Godoc](https://godoc.org/github.com/moorereason/httpauth?status.svg)](http://godoc.org/github.com/moorereason/httpauth)

Middleware
----------

The package provides a middleware that can be mounted to services or controllers that require authentication.
The basic authentication middleware is instantiated using the BasicMiddleware function.
This function accepts a specification that describes how the middleware should operate.

```go
	spec := &httpauth.Specification{
		LogFailure:	true,
		LogSuccess:	true,
		Realm:		"Restricted",
		ValidationFunc:	authHandler,
	}
	service.Use(httpauth.BasicMiddleware(spec))
	// or
	protectedController.Use(httpauth.BasicMiddleware(spec))
```
