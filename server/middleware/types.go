package middleware

type Error struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}
