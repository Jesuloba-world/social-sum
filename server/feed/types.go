package feed

type Error struct {
	Message string `json:"message"`
	Errors  string `json:"error"`
}
