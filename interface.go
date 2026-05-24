package wrap

import (
	"context"
	"net/http"

	"github.com/imroc/req/v3"
)

type Client interface {
	// Do execute an HTTP request with the specified method, URL, headers, and body.
	// - ctx: Context for request cancellation and timeout.
	// - method: HTTP method (e.g., "GET", "POST").
	// - url: The request URL.
	// - headers: Optional headers to include in the request (map of header name to value).
	// - body: Optional request body (can be nil).
	// Returns the HTTP response.
	Do(ctx context.Context, method, url string, headers map[string]string, body any) *req.Response

	// DoWithErrorHandling execute an HTTP request with the specified method, URL, headers, and body.
	// - ctx: Context for request cancellation and timeout.
	// - method: HTTP method (e.g., "GET", "POST").
	// - url: The request URL.
	// - headers: Optional headers to include in the request (map of header name to value).
	// - body: Optional request body (can be nil).
	// Returns the HTTP response.
	// Does additional checks on the response status code and returns an error if it's not successful.
	DoWithErrorHandling(ctx context.Context, method, url string, headers map[string]string, body any) *req.Response

	// DoWithResult execute an HTTP request with the specified method, URL, headers, and body.
	// - ctx: Context for request cancellation and timeout.
	// - method: HTTP method (e.g., "GET", "POST").
	// - url: The request URL.
	// - headers: Optional headers to include in the request (map of header name to value).
	// - body: Optional request body (can be nil).
	// - resultTo: Pointer to the variable where we unmarshal the response body OR EasyJSON Unmarshal Factory(to avoid memory allocations when response is unsuccesful).
	// Returns the HTTP response and any error encountered.
	DoWithResult(ctx context.Context, method, url string, headers map[string]string, body any, resultTo any) *req.Response

	// DoWithRequest execute the given req.Request.
	// - request: The req.Request to be executed.
	// Returns the HTTP response.
	DoWithRequest(request *req.Request) *req.Response

	// DoWithRequest execute the given http.Request.
	// - request: The http.Request to be executed.
	// Returns the HTTP response and any error encountered.
	DoWithHTTPRequest(request *http.Request) (*http.Response, error)

	// SetHeaders Sets the headers for the request
	// headers can be a map[string]string, http.Header or nil if you want to delete all headers
	SetHeaders(headers any) error

	// SetCookies Sets cookies for all requests.
	// Cookies can be []*http.Cookie, []http.Cookie ,map[string]string, or nil to remove all cookies.
	SetCookiesArray(cookies any) error

	// SetProxy Sets proxy for all requests.
	// proxy can be a string, *url.URL, or nil to remove the proxy
	// Also give the checkURL to check if the proxy is working
	SetProxy(proxyURL string, checkURL string) error

	// SetResponseBodyTransformer Sets the response body transformer.
	// fn is a function that takes the raw response body, request, and response as input and returns the transformed response body and an error.
	SetResponseBodyTransformer(fn func(rawBody []byte, req *req.Request, resp *req.Response) (transformedBody []byte, err error))

	// Client returns the underlying req.Client used by the Client.
	Client() *req.Client

	// Get sends an HTTP GET request to the specified URL with optional headers.
	// - ctx: Context for request cancellation and timeout.
	// - url: The request URL.
	// - headers: Optional headers to include in the request (map of header name to value).
	// Returns the HTTP response and any error encountered.
	Get(ctx context.Context, url string, headers map[string]string) *req.Response

	// Post sends an HTTP POST request to the specified URL with optional headers and a request body.
	// - ctx: Context for request cancellation and timeout.
	// - url: The request URL.
	// - headers: Optional headers to include in the request (map of header name to value).
	// - body: Optional request body (can be nil).
	// Returns the HTTP response and any error encountered.
	Post(ctx context.Context, url string, headers map[string]string, body any) *req.Response

	// Put sends an HTTP PUT request to the specified URL with optional headers and a request body.
	// - ctx: Context for request cancellation and timeout.
	// - url: The request URL.
	// - headers: Optional headers to include in the request (map of header name to value).
	// - body: Optional request body (can be nil).
	// Returns the HTTP response and any error encountered.
	Put(ctx context.Context, url string, headers map[string]string, body any) *req.Response

	// Delete sends an HTTP DELETE request to the specified URL with optional headers.
	// - ctx: Context for request cancellation and timeout.
	// - url: The request URL.
	// - headers: Optional headers to include in the request (map of header name to value).
	// Returns the HTTP response and any error encountered.
	Delete(ctx context.Context, url string, headers map[string]string) *req.Response

	// Patch sends an HTTP PATCH request to the specified URL with optional headers and a request body.
	// - ctx: Context for request cancellation and timeout.
	// - url: The request URL.
	// - headers: Optional headers to include in the request (map of header name to value).
	// - body: Optional request body (can be nil).
	// Returns the HTTP response and any error encountered.
	Patch(ctx context.Context, url string, headers map[string]string, body any) *req.Response

	// Head sends an HTTP HEAD request to the specified URL with optional headers.
	// - ctx: Context for request cancellation and timeout.
	// - url: The request URL.
	// - headers: Optional headers to include in the request (map of header name to value).
	// Returns the HTTP response and any error encountered.
	Head(ctx context.Context, url string, headers map[string]string) *req.Response

	// Options sends an HTTP OPTIONS request to the specified URL with optional headers.
	// - ctx: Context for request cancellation and timeout.
	// - url: The request URL.
	// - headers: Optional headers to include in the request (map of header name to value).
	// Returns the HTTP response and any error encountered.
	Options(ctx context.Context, url string, headers map[string]string) *req.Response

	// Trace sends an HTTP TRACE request to the specified URL with optional headers.
	// - ctx: Context for request cancellation and timeout.
	// - url: The request URL.
	// - headers: Optional headers to include in the request (map of header name to value).
	// Returns the HTTP response and any error encountered.
	Trace(ctx context.Context, url string, headers map[string]string) *req.Response

	// Connect sends an HTTP CONNECT request to the specified URL with optional headers.
	// - ctx: Context for request cancellation and timeout.
	// - url: The request URL.
	// - headers: Optional headers to include in the request (map of header name to value).
	// Returns the HTTP response and any error encountered.
	Connect(ctx context.Context, url string, headers map[string]string) *req.Response

	// CloseIdleConnection closes all idle connections, goroutine safe
	CloseIdleConnections()

	// Close closes the HTTP client, it's forbidden to use the client after this method is called.
	Close()
}
