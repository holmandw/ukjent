package json

import (
	j "encoding/json"

	"github.com/holman_dw/ukjent"
)

func MarshalWord(w ukjent.Word) ([]byte, error) {
	return j.Marshal(w)
}

func UnmarshalWord(b []byte) (ukjent.Word, error) {
	var w ukjent.Word
	err := j.Unmarshal(b, &w)
	return w, err
}
func MarshalWords(ws []ukjent.Word) ([]byte, error) {
	return j.Marshal(ws)
}

func UnmarshalWords(b []byte) ([]ukjent.Word, error) {
	var ws []ukjent.Word
	err := j.Unmarshal(b, &ws)
	return ws, err
}
