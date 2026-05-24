package wrap

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"context"
	"encoding/json"
	"errors"
	"io"
	"maps"
	"net/http"
	"strings"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/imroc/req/v3"
	"github.com/klauspost/compress/zstd"
	"github.com/mailru/easyjson"
)

// Implementation check
var _ Client = (*httpClient)(nil)

type httpClient struct {
	client *req.Client
}

func NewReqClient() *httpClient {
	return &httpClient{
		client: req.NewClient(),
	}
}

func (c *httpClient) Do(ctx context.Context, method, url string, headers map[string]string, body any) *req.Response {
	req := c.client.R()

	// Request configuration
	req.SetContext(ctx)
	req.Method = method
	req.SetURL(url)
	req.SetHeaders(headers)
	req.SetBody(body)

	resp := req.Do(ctx)
	if resp.Err != nil {
		resp.Err = errors.Join(resp.Err, ErrConn)
	}

	return resp
}

func (c *httpClient) DoWithErrorHandling(ctx context.Context, method, url string, headers map[string]string, body any) *req.Response {
	resp := c.Do(ctx, method, url, headers, body)
	if resp.Err != nil {
		return resp
	}
	if resp.IsErrorState() {
		resp.Err = errors.Join(resp.Err, ErrStatus, CodeToErr(resp.StatusCode))
	}

	return resp
}

func (c *httpClient) DoWithResult(ctx context.Context, method, url string, headers map[string]string, body any, resultTo any) *req.Response {
	resp := c.DoWithErrorHandling(ctx, method, url, headers, body)
	if resp.Err != nil {
		return resp
	}

	bodyBytes := resp.Bytes()
	if bodyBytes == nil {
		respBytes, err := resp.ToBytes() // Should handle all decoding types automatically
		if err != nil {
			resp.Err = errors.Join(resp.Err, err, ErrConn)
		}
		bodyBytes = respBytes
	}

	switch v := resultTo.(type) {
	case UnmarshallerFactory:
		res := v()
		if err := easyjson.Unmarshal(bodyBytes, res); err != nil {
			resp.Err = errors.Join(resp.Err, err, ErrUnmarshal)
		}
		resp.Request.Result = res // Also setting the result here just in case
	default:
		if err := json.Unmarshal(bodyBytes, resultTo); err != nil {
			resp.Err = errors.Join(resp.Err, err, ErrUnmarshal)
		}
		resp.Request.Result = resultTo // Also setting the result here just in case
	}

	return resp
}

func (c httpClient) DoWithRequest(req *req.Request) *req.Response {
	resp := req.Do()
	if resp.Err != nil {
		resp.Err = errors.Join(resp.Err, ErrConn)
	}

	return resp
}

func (c *httpClient) DoWithHTTPRequest(req *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errors.Join(ErrConn, err)
	}

	return resp, nil
}

func (c *httpClient) SetHeaders(headers any) error {
	hd := make(map[string]string)
	switch h := headers.(type) {
	case map[string]string:
		maps.Copy(hd, h)
	case http.Header:
		for k, v := range h {
			hd[k] = v[0]
		}
	case nil:
		for k := range c.client.Headers {
			c.client.Headers.Del(k)
		}
		return nil
	default:
		return ErrInvalidHeadersType
	}
	c.client.SetCommonHeaders(hd)
	return nil
}

func (c *httpClient) SetCookiesArray(cookies any) error {
	cks := make([]*http.Cookie, 0)

	switch cookie := cookies.(type) {
	case []*http.Cookie:
		cks = cookie
	case []http.Cookie:
		for _, val := range cookie {
			cks = append(cks, &val)
		}
	case map[string]string:
		for k, v := range cookie {
			cks = append(cks, &http.Cookie{
				Name:  k,
				Value: v,
			})
		}
	case nil:
		c.client.Cookies = nil
		return nil
	default:
		return ErrInvalidCookiesType
	}
	for _, ck := range cks {
		c.client.Cookies = append(c.client.Cookies, ck) // Use direct append instead of proxy func
	}

	return nil
}

func (c *httpClient) SetProxy(proxyURL string, checkURL string) error {
	c.client.SetProxyURL(proxyURL)

	if checkURL != "" {
		parentCtx := context.TODO()
		ctx, cancel := context.WithTimeout(parentCtx, 30*time.Second)
		defer cancel()
		res := c.DoWithErrorHandling(ctx, http.MethodGet, checkURL, nil, nil)
		return res.Err
	}

	return nil
}

func (c *httpClient) SetResponseBodyTransformer(fn func(rawBody []byte, req *req.Request, resp *req.Response) (transformedBody []byte, err error)) {
	if fn == nil {
		fn = func(rawBody []byte, req *req.Request, resp *req.Response) (transformedBody []byte, err error) {
			encoding := strings.ToLower(strings.TrimSpace(resp.Header.Get("Content-Encoding")))

			switch encoding {
			case "gzip":
				gr, err := gzip.NewReader(bytes.NewReader(rawBody))
				if err != nil {
					return rawBody, err
				}
				defer gr.Close()
				return io.ReadAll(gr)

			case "deflate":
				// Content-Encoding: deflate commonly carries zlib-framed data (RFC 1950);
				// try zlib first, fall back to raw DEFLATE (RFC 1951).
				if zr, err := zlib.NewReader(bytes.NewReader(rawBody)); err == nil {
					defer zr.Close()
					if out, err := io.ReadAll(zr); err == nil {
						return out, nil
					}
				}
				fr := flate.NewReader(bytes.NewReader(rawBody))
				defer fr.Close()
				return io.ReadAll(fr)

			case "zlib":
				zr, err := zlib.NewReader(bytes.NewReader(rawBody))
				if err != nil {
					return rawBody, err
				}
				defer zr.Close()
				return io.ReadAll(zr)

			case "br":
				return io.ReadAll(brotli.NewReader(bytes.NewReader(rawBody)))

			case "zstd":
				zr, err := zstd.NewReader(bytes.NewReader(rawBody))
				if err != nil {
					return rawBody, err
				}
				defer zr.Close()
				return io.ReadAll(zr)

			case "identity", "":
				return rawBody, nil
			}

			// No recognised Content-Encoding — probe each format in order of prevalence.
			if gr, err := gzip.NewReader(bytes.NewReader(rawBody)); err == nil {
				if out, err := io.ReadAll(gr); err == nil {
					_ = gr.Close()
					return out, nil
				}
				_ = gr.Close()
			}

			if zr, err := zlib.NewReader(bytes.NewReader(rawBody)); err == nil {
				if out, err := io.ReadAll(zr); err == nil {
					_ = zr.Close()
					return out, nil
				}
				_ = zr.Close()
			}

			fr := flate.NewReader(bytes.NewReader(rawBody))
			if out, err := io.ReadAll(fr); err == nil {
				_ = fr.Close()
				return out, nil
			}
			_ = fr.Close()

			if out, err := io.ReadAll(brotli.NewReader(bytes.NewReader(rawBody))); err == nil {
				return out, nil
			}

			if zr, err := zstd.NewReader(bytes.NewReader(rawBody)); err == nil {
				defer zr.Close()
				if out, err := io.ReadAll(zr); err == nil {
					return out, nil
				}
			}

			return rawBody, nil
		}
	}

	c.client.SetResponseBodyTransformer(fn)
}

func (c *httpClient) Client() *req.Client {
	return c.client
}

// Get sends an HTTP GET request to the specified URL with optional headers.
// - ctx: Context for request cancellation and timeout.
// - url: The request URL.
// - headers: Optional headers to include in the request (map of header name to value).
// Returns the HTTP response and any error encountered.
func (c *httpClient) Get(ctx context.Context, url string, headers map[string]string) *req.Response {
	return c.Do(ctx, http.MethodGet, url, headers, nil)
}

// Post sends an HTTP POST request to the specified URL with optional headers and a request body.
// - ctx: Context for request cancellation and timeout.
// - url: The request URL.
// - headers: Optional headers to include in the request (map of header name to value).
// - body: Optional request body (can be nil).
// Returns the HTTP response and any error encountered.
func (c *httpClient) Post(ctx context.Context, url string, headers map[string]string, body any) *req.Response {
	return c.Do(ctx, http.MethodPost, url, headers, body)
}

// Put sends an HTTP PUT request to the specified URL with optional headers and a request body.
// - ctx: Context for request cancellation and timeout.
// - url: The request URL.
// - headers: Optional headers to include in the request (map of header name to value).
// - body: Optional request body (can be nil).
// Returns the HTTP response and any error encountered.
func (c *httpClient) Put(ctx context.Context, url string, headers map[string]string, body any) *req.Response {
	return c.Do(ctx, http.MethodPut, url, headers, body)
}

// Delete sends an HTTP DELETE request to the specified URL with optional headers.
// - ctx: Context for request cancellation and timeout.
// - url: The request URL.
// - headers: Optional headers to include in the request (map of header name to value).
// Returns the HTTP response and any error encountered.
func (c *httpClient) Delete(ctx context.Context, url string, headers map[string]string) *req.Response {
	return c.Do(ctx, http.MethodDelete, url, headers, nil)
}

// Patch sends an HTTP PATCH request to the specified URL with optional headers and a request body.
// - ctx: Context for request cancellation and timeout.
// - url: The request URL.
// - headers: Optional headers to include in the request (map of header name to value).
// - body: Optional request body (can be nil).
// Returns the HTTP response and any error encountered.
func (c *httpClient) Patch(ctx context.Context, url string, headers map[string]string, body any) *req.Response {
	return c.Do(ctx, http.MethodPatch, url, headers, body)
}

// Head sends an HTTP HEAD request to the specified URL with optional headers.
// - ctx: Context for request cancellation and timeout.
// - url: The request URL.
// - headers: Optional headers to include in the request (map of header name to value).
// Returns the HTTP response and any error encountered.
func (c *httpClient) Head(ctx context.Context, url string, headers map[string]string) *req.Response {
	return c.Do(ctx, http.MethodHead, url, headers, nil)
}

// Options sends an HTTP OPTIONS request to the specified URL with optional headers.
// - ctx: Context for request cancellation and timeout.
// - url: The request URL.
// - headers: Optional headers to include in the request (map of header name to value).
// Returns the HTTP response and any error encountered.
func (c *httpClient) Options(ctx context.Context, url string, headers map[string]string) *req.Response {
	return c.Do(ctx, http.MethodOptions, url, headers, nil)
}

// Trace sends an HTTP TRACE request to the specified URL with optional headers.
// - ctx: Context for request cancellation and timeout.
// - url: The request URL.
// - headers: Optional headers to include in the request (map of header name to value).
// Returns the HTTP response and any error encountered.
func (c *httpClient) Trace(ctx context.Context, url string, headers map[string]string) *req.Response {
	return c.Do(ctx, http.MethodTrace, url, headers, nil)
}

// Connect sends an HTTP CONNECT request to the specified URL with optional headers.
// - ctx: Context for request cancellation and timeout.
// - url: The request URL.
// - headers: Optional headers to include in the request (map of header name to value).
// Returns the HTTP response and any error encountered.
func (c *httpClient) Connect(ctx context.Context, url string, headers map[string]string) *req.Response {
	return c.Do(ctx, http.MethodConnect, url, headers, nil)
}

func (c *httpClient) CloseIdleConnections() {
	c.client.Transport.CloseIdleConnections()
}

func (c *httpClient) Close() {
	if c.client != nil {
		c.client.CloseIdleConnections()
	}
	c.client = nil
}
