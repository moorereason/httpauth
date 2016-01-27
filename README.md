httpauth
========

Package httpauth provides a goa middleware that implements Basic HTTP authentication.

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
