package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// declare a custom runtime type, which has the underlying type int32( the same as our movie struct field)
type Runtime int32

var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

//implement a MarshalJSON() method on the runtime type as that it satisfies the
// json.Marshaler interface. this should return the JSON-encoded value for the movie
// runtime(in our case, it will return in hte format "<runtime> mins").

func (r Runtime) MarshalJSON() ([]byte, error) {
	//generate a string containing the movie runtime in the required format
	jsonValue := fmt.Sprintf("%d mins", r)

	//use the strconv.Quote() function on the string to wrap it in double quotes. it
	// needs to be surrounded by double quotes in the order to be a valid *JSON string*.
	quotedJSONValue := strconv.Quote(jsonValue)

	// convert the quoted string value to a byte slice and return it.
	return []byte(quotedJSONValue), nil
}

func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}
	parts := strings.Split(unquotedJSONValue, " ")

	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}
	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}
	*r = Runtime(i)

	return nil
}
