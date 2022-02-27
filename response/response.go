package response

type ErrorResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Msg    string `json:"msg"`
	Error  string `json:"error,omitempty"`
}

type LinkResponse struct {
	Url       string `json:"url"`
	ExpiresAt string `json:"expiresAt"`
}
