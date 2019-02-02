package json

import (
	j "encoding/json"
	"fmt"

	"github.com/holman_dw/ukjent"
)

func MarshalError(e ukjent.Error) ([]byte, error) {
	return j.Marshal(e)
}

func UnmarshalError(b []byte) (ukjent.Error, error) {
	return ukjent.Error{}, fmt.Errorf("unmarshalling errors not implemented")
}
