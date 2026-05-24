# wrap-http

A thin, opinionated wrapper around [imroc/req](https://github.com/imroc/req) that adds structured error handling, content-encoding-aware decompression, typed sentinel errors for every HTTP status code, and optional [easyjson](https://github.com/mailru/easyjson) unmarshalling.

> Provided as-is, without any warranty or support.

## Installation

```bash
go get github.com/broaskaGit/wrap-http
```

## Quick start

```go
client := wrap.NewReqClient()
defer client.Close()

ctx := context.Background()

// Simple GET
resp := client.Get(ctx, "https://api.example.com/items", nil)
if resp.Err != nil {
    log.Fatal(resp.Err)
}

// POST with JSON body, unmarshal response
var result MyResponse
resp = client.DoWithResult(ctx, http.MethodPost, "https://api.example.com/items",
    map[string]string{"Content-Type": "application/json"},
    MyRequest{Name: "foo"},
    &result,
)
if resp.Err != nil {
    log.Fatal(resp.Err)
}
```

## Core methods

| Method | Description |
|---|---|
| `Do` | Execute a request; wraps connection errors with `ErrConn`. |
| `DoWithErrorHandling` | Same as `Do`, plus wraps 4xx/5xx responses with `ErrStatus` and the matching HTTP sentinel error. |
| `DoWithResult` | Same as `DoWithErrorHandling`, plus deserialises the response body into `resultTo`. |
| `DoWithRequest` | Execute a pre-built `*req.Request`. |
| `DoWithHTTPRequest` | Execute a standard `*http.Request` and return a `*http.Response`. |

Convenience shorthands — `Get`, `Post`, `Put`, `Patch`, `Delete`, `Head`, `Options`, `Trace`, `Connect` — all delegate to `Do`.

## Error handling

Every error returned through `resp.Err` is composed with `errors.Join`, so callers can use `errors.Is` to check for any layer:

```go
if errors.Is(resp.Err, wrap.ErrStatus) {
    if errors.Is(resp.Err, wrap.ErrNotFound) {
        // 404
    }
}
if errors.Is(resp.Err, wrap.ErrUnmarshal) { /* decode failed */ }
if errors.Is(resp.Err, wrap.ErrConn)      { /* transport error */ }
```

Sentinel errors exist for every standard 4xx and 5xx status code. Use `wrap.CodeToErr(code)` to map an arbitrary integer status code to its sentinel (returns `ErrUnknownCode` for unrecognised codes).

## Unmarshalling (`DoWithResult`)

`resultTo` accepts two forms:

- **Any pointer** (`*MyStruct`) — standard `encoding/json` unmarshalling.
- **`UnmarshallerFactory`** — a `func() easyjson.Unmarshaler` factory; the factory is only called on a successful response, avoiding allocations on error paths.

```go
// Standard JSON
var result MyStruct
client.DoWithResult(ctx, "GET", url, nil, nil, &result)

// EasyJSON (zero allocation on failure)
client.DoWithResult(ctx, "GET", url, nil, nil,
    wrap.UnmarshallerFactory(func() easyjson.Unmarshaler { return &MyStruct{} }),
)
```

## Configuration helpers

### Headers

```go
client.SetHeaders(map[string]string{"Authorization": "Bearer token"})
client.SetHeaders(httpHeader)  // http.Header
client.SetHeaders(nil)         // clears all common headers
```

### Cookies

```go
client.SetCookiesArray([]*http.Cookie{{Name: "session", Value: "abc"}})
client.SetCookiesArray(map[string]string{"session": "abc"})
client.SetCookiesArray(nil) // clears all cookies
```

### Proxy

```go
// checkURL is optional; when non-empty a GET is issued to verify connectivity.
err := client.SetProxy("http://proxy.example.com:8080", "https://api.example.com")
```

### Body transformer (decompression)

Calling `SetResponseBodyTransformer(nil)` installs a default transformer that decodes the response body based on the `Content-Encoding` response header:

| `Content-Encoding` | Algorithm |
|---|---|
| `gzip` | gzip |
| `deflate` | zlib (RFC 1950), falling back to raw DEFLATE (RFC 1951) |
| `zlib` | zlib |
| `br` | Brotli |
| `zstd` | Zstandard |
| `identity` / absent | pass-through |

When no `Content-Encoding` header is present (or it carries an unrecognised value), the transformer probes each format in order of prevalence and returns the raw body unchanged if none match.

A custom function can be provided instead:

```go
client.SetResponseBodyTransformer(func(raw []byte, req *req.Request, resp *req.Response) ([]byte, error) {
    // custom transformation
    return raw, nil
})
```

### Escape hatch

```go
client.Client() // returns the underlying *req.Client for direct configuration
```