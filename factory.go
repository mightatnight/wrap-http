package wrap

import "github.com/mailru/easyjson"

// UnmarshallerFactory is a function that creates a new instance of an easyjson.Unmarshaler
// Used to avoid unnecessary allocations when response is unsuccessful
type UnmarshallerFactory func() easyjson.Unmarshaler
