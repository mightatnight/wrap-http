// Package wrap provides sentinel errors for HTTP client operations,
// covering connection failures, response unmarshaling, proxy configuration,
// and the full range of 4xx/5xx HTTP status codes.
package wrap

import (
	"errors"
	"net/http"
)

var (
	// ErrConn is returned when a connection cannot be established or maintained.
	ErrConn = errors.New("connection")

	// ErrStatus is returned when the server responds with a non-successful HTTP status code.
	ErrStatus = errors.New("status")

	// ErrUnknownCode is returned by [CodeToErr] when the given HTTP status code has no
	// corresponding sentinel error defined in this package.
	ErrUnknownCode = errors.New("unknown status code")

	// ErrUnmarshal is returned when the response body cannot be decoded into the target type.
	ErrUnmarshal = errors.New("unmarshal")

	// ErrInvalidProxy is returned when the provided proxy value is not usable.
	ErrInvalidProxy = errors.New("invalid proxy")

	// ErrInvalidProxyURL is returned when the proxy URL is malformed or cannot be parsed.
	ErrInvalidProxyURL = errors.New("invalid proxy url")

	// ErrInvalidProxyCheckURL is returned when the URL used to verify proxy connectivity is
	// malformed or cannot be parsed.
	ErrInvalidProxyCheckURL = errors.New("invalid proxy check url")
)

var (
	// ErrInvalidHeadersType is returned when the value supplied for request headers is not
	// a supported type (e.g. map[string]string or http.Header).
	ErrInvalidHeadersType = errors.New("invalid headers type")

	// ErrInvalidCookiesType is returned when the value supplied for request cookies is not
	// a supported type.
	ErrInvalidCookiesType = errors.New("invalid cookies type")
)

var (
	// 4xx Client Errors

	// ErrBadRequest corresponds to HTTP 400 Bad Request.
	ErrBadRequest = errors.New("bad request")
	// ErrUnauthorized corresponds to HTTP 401 Unauthorized.
	ErrUnauthorized = errors.New("unauthorized")
	// ErrPaymentRequired corresponds to HTTP 402 Payment Required.
	ErrPaymentRequired = errors.New("payment required")
	// ErrForbidden corresponds to HTTP 403 Forbidden.
	ErrForbidden = errors.New("forbidden")
	// ErrNotFound corresponds to HTTP 404 Not Found.
	ErrNotFound = errors.New("not found")
	// ErrMethodNotAllowed corresponds to HTTP 405 Method Not Allowed.
	ErrMethodNotAllowed = errors.New("method not allowed")
	// ErrNotAcceptable corresponds to HTTP 406 Not Acceptable.
	ErrNotAcceptable = errors.New("not acceptable")
	// ErrProxyAuthRequired corresponds to HTTP 407 Proxy Authentication Required.
	ErrProxyAuthRequired = errors.New("proxy authentication required")
	// ErrRequestTimeout corresponds to HTTP 408 Request Timeout.
	ErrRequestTimeout = errors.New("request timeout")
	// ErrConflict corresponds to HTTP 409 Conflict.
	ErrConflict = errors.New("conflict")
	// ErrGone corresponds to HTTP 410 Gone.
	ErrGone = errors.New("gone")
	// ErrLengthRequired corresponds to HTTP 411 Length Required.
	ErrLengthRequired = errors.New("length required")
	// ErrPreconditionFailed corresponds to HTTP 412 Precondition Failed.
	ErrPreconditionFailed = errors.New("precondition failed")
	// ErrRequestEntityTooLarge corresponds to HTTP 413 Request Entity Too Large.
	ErrRequestEntityTooLarge = errors.New("request entity too large")
	// ErrRequestURITooLong corresponds to HTTP 414 Request-URI Too Long.
	ErrRequestURITooLong = errors.New("request uri too long")
	// ErrUnsupportedMediaType corresponds to HTTP 415 Unsupported Media Type.
	ErrUnsupportedMediaType = errors.New("unsupported media type")
	// ErrRequestedRangeNotSatisfiable corresponds to HTTP 416 Requested Range Not Satisfiable.
	ErrRequestedRangeNotSatisfiable = errors.New("requested range not satisfiable")
	// ErrExpectationFailed corresponds to HTTP 417 Expectation Failed.
	ErrExpectationFailed = errors.New("expectation failed")
	// ErrTeapot corresponds to HTTP 418 I'm a Teapot.
	ErrTeapot = errors.New("i'm a teapot")
	// ErrMisdirectedRequest corresponds to HTTP 421 Misdirected Request.
	ErrMisdirectedRequest = errors.New("misdirected request")
	// ErrUnprocessableEntity corresponds to HTTP 422 Unprocessable Entity.
	ErrUnprocessableEntity = errors.New("unprocessable entity")
	// ErrLocked corresponds to HTTP 423 Locked.
	ErrLocked = errors.New("locked")
	// ErrFailedDependency corresponds to HTTP 424 Failed Dependency.
	ErrFailedDependency = errors.New("failed dependency")
	// ErrTooEarly corresponds to HTTP 425 Too Early.
	ErrTooEarly = errors.New("too early")
	// ErrUpgradeRequired corresponds to HTTP 426 Upgrade Required.
	ErrUpgradeRequired = errors.New("upgrade required")
	// ErrPreconditionRequired corresponds to HTTP 428 Precondition Required.
	ErrPreconditionRequired = errors.New("precondition required")
	// ErrTooManyRequests corresponds to HTTP 429 Too Many Requests.
	ErrTooManyRequests = errors.New("too many requests")
	// ErrRequestHeaderFieldsTooLarge corresponds to HTTP 431 Request Header Fields Too Large.
	ErrRequestHeaderFieldsTooLarge = errors.New("request header fields too large")
	// ErrUnavailableForLegalReasons corresponds to HTTP 451 Unavailable For Legal Reasons.
	ErrUnavailableForLegalReasons = errors.New("unavailable for legal reasons")

	// 5xx Server Errors

	// ErrInternalServerError corresponds to HTTP 500 Internal Server Error.
	ErrInternalServerError = errors.New("internal server error")
	// ErrNotImplemented corresponds to HTTP 501 Not Implemented.
	ErrNotImplemented = errors.New("not implemented")
	// ErrBadGateway corresponds to HTTP 502 Bad Gateway.
	ErrBadGateway = errors.New("bad gateway")
	// ErrServiceUnavailable corresponds to HTTP 503 Service Unavailable.
	ErrServiceUnavailable = errors.New("service unavailable")
	// ErrGatewayTimeout corresponds to HTTP 504 Gateway Timeout.
	ErrGatewayTimeout = errors.New("gateway timeout")
	// ErrHTTPVersionNotSupported corresponds to HTTP 505 HTTP Version Not Supported.
	ErrHTTPVersionNotSupported = errors.New("http version not supported")
	// ErrVariantAlsoNegotiates corresponds to HTTP 506 Variant Also Negotiates.
	ErrVariantAlsoNegotiates = errors.New("variant also negotiates")
	// ErrInsufficientStorage corresponds to HTTP 507 Insufficient Storage.
	ErrInsufficientStorage = errors.New("insufficient storage")
	// ErrLoopDetected corresponds to HTTP 508 Loop Detected.
	ErrLoopDetected = errors.New("loop detected")
	// ErrNotExtended corresponds to HTTP 510 Not Extended.
	ErrNotExtended = errors.New("not extended")
	// ErrNetworkAuthenticationRequired corresponds to HTTP 511 Network Authentication Required.
	ErrNetworkAuthenticationRequired = errors.New("network authentication required")
)

// httpstatusErrors maps HTTP status codes to their corresponding sentinel errors.
var httpstatusErrors = map[int]error{
	// 4xx Client Errors
	http.StatusBadRequest:                   ErrBadRequest,
	http.StatusUnauthorized:                 ErrUnauthorized,
	http.StatusPaymentRequired:              ErrPaymentRequired,
	http.StatusForbidden:                    ErrForbidden,
	http.StatusNotFound:                     ErrNotFound,
	http.StatusMethodNotAllowed:             ErrMethodNotAllowed,
	http.StatusNotAcceptable:                ErrNotAcceptable,
	http.StatusProxyAuthRequired:            ErrProxyAuthRequired,
	http.StatusRequestTimeout:               ErrRequestTimeout,
	http.StatusConflict:                     ErrConflict,
	http.StatusGone:                         ErrGone,
	http.StatusLengthRequired:               ErrLengthRequired,
	http.StatusPreconditionFailed:           ErrPreconditionFailed,
	http.StatusRequestEntityTooLarge:        ErrRequestEntityTooLarge,
	http.StatusRequestURITooLong:            ErrRequestURITooLong,
	http.StatusUnsupportedMediaType:         ErrUnsupportedMediaType,
	http.StatusRequestedRangeNotSatisfiable: ErrRequestedRangeNotSatisfiable,
	http.StatusExpectationFailed:            ErrExpectationFailed,
	http.StatusTeapot:                       ErrTeapot,
	http.StatusMisdirectedRequest:           ErrMisdirectedRequest,
	http.StatusUnprocessableEntity:          ErrUnprocessableEntity,
	http.StatusLocked:                       ErrLocked,
	http.StatusFailedDependency:             ErrFailedDependency,
	http.StatusTooEarly:                     ErrTooEarly,
	http.StatusUpgradeRequired:              ErrUpgradeRequired,
	http.StatusPreconditionRequired:         ErrPreconditionRequired,
	http.StatusTooManyRequests:              ErrTooManyRequests,
	http.StatusRequestHeaderFieldsTooLarge:  ErrRequestHeaderFieldsTooLarge,
	http.StatusUnavailableForLegalReasons:   ErrUnavailableForLegalReasons,

	// 5xx Server Errors
	http.StatusInternalServerError:           ErrInternalServerError,
	http.StatusNotImplemented:                ErrNotImplemented,
	http.StatusBadGateway:                    ErrBadGateway,
	http.StatusServiceUnavailable:            ErrServiceUnavailable,
	http.StatusGatewayTimeout:                ErrGatewayTimeout,
	http.StatusHTTPVersionNotSupported:       ErrHTTPVersionNotSupported,
	http.StatusVariantAlsoNegotiates:         ErrVariantAlsoNegotiates,
	http.StatusInsufficientStorage:           ErrInsufficientStorage,
	http.StatusLoopDetected:                  ErrLoopDetected,
	http.StatusNotExtended:                   ErrNotExtended,
	http.StatusNetworkAuthenticationRequired: ErrNetworkAuthenticationRequired,
}

// CodeToErr maps an HTTP status code to its corresponding sentinel error.
// It returns [ErrUnknownCode] if the status code is not covered by this package.
func CodeToErr(code int) error {
	if err, ok := httpstatusErrors[code]; ok {
		return err
	}
	return ErrUnknownCode
}
