package zerver_rest

import . "github.com/cosiner/golib/errors"

type (
	// HandlerFunc is the common request handler function type
	HandlerFunc func(Request, Response)

	// Handler is an common interface of request handler
	// it will be inited on server started, destroyed on server stopped
	// and use it's method to process income request
	// if want to custom the relation of method and handler, just embed
	// EmptyHandler to your Handler, then implement interface MethodIndicator
	Handler interface {
		Init(*Server) error
		Destroy()
		Get(Request, Response)    // query
		Post(Request, Response)   // add
		Delete(Request, Response) // delete
		Put(Request, Response)    // create or replace
		Patch(Request, Response)  // update
	}

	// MethodIndicator is an interface for user handler to
	// custom method handle functions
	MethodIndicator interface {
		// Handler return an method handle function by method name
		// if nill returned, means access forbidden
		Handler(method string) HandlerFunc
	}

	// EmptyHandler is an empty handler for user to embed
	EmptyHandler struct{}

	// funcHandler is a handler that use user customed handler function
	// it a funcHandler is defined after a normal handler with same pattern,
	// regardless of only one method is used in funcHandler like Get, it's impossible
	// to use the Post of normal handler, whole normal handler is hidded by this
	// funcHandler, and if other method handler like Post, Put is not set,
	// user access of these method is forbiddened
	funcHandler struct {
		EmptyHandler
		get    HandlerFunc
		post   HandlerFunc
		delete HandlerFunc
		put    HandlerFunc
		patch  HandlerFunc
	}
)

// IndicateHandler indicate handler function from a handler and method
func IndicateHandler(method string, handler Handler) HandlerFunc {
	switch handler := handler.(type) {
	case MethodIndicator:
		return handler.Handler(method)
	default:
		return standardIndicate(method, handler)
	}
}

// standardIndicate normal indicate method handle function
// each method indicate the function with same name, such as GET->Get...
func standardIndicate(method string, handler Handler) (handlerFunc HandlerFunc) {
	switch method {
	case GET:
		handlerFunc = handler.Get
	case POST:
		handlerFunc = handler.Post
	case DELETE:
		handlerFunc = handler.Delete
	case PUT:
		handlerFunc = handler.Put
	case PATCH:
		handlerFunc = handler.Patch
	}
	return
}

// funcHandler implements MethodIndicator interface for custom method handler
func (fh *funcHandler) Handler(method string) (handlerFunc HandlerFunc) {
	switch method {
	case GET:
		handlerFunc = fh.get
	case POST:
		handlerFunc = fh.post
	case DELETE:
		handlerFunc = fh.delete
	case PUT:
		handlerFunc = fh.put
	case PATCH:
		handlerFunc = fh.patch
	}
	return
}

// setMethod setup method handler for funcHandler
func (fh *funcHandler) setMethod(method string, handlerFunc HandlerFunc) error {
	switch method {
	case GET:
		fh.get = handlerFunc
	case POST:
		fh.post = handlerFunc
	case PUT:
		fh.put = handlerFunc
	case DELETE:
		fh.delete = handlerFunc
	case PATCH:
		fh.patch = handlerFunc
	default:
		return Err("Not supported request method")
	}
	return nil
}

// EmptyHandler methods
func (EmptyHandler) Init(*Server) error       { return nil }
func (EmptyHandler) Destroy()                 {}
func (EmptyHandler) Get(Request, Response)    {}
func (EmptyHandler) Post(Request, Response)   {}
func (EmptyHandler) Delete(Request, Response) {}
func (EmptyHandler) Put(Request, Response)    {}
func (EmptyHandler) Patch(Request, Response)  {}
