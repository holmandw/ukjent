package ukjent

// ErrorResponse is the data sent back in event of an error
type Error struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}
