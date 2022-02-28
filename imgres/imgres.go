package imgres

type ErrorResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Msg    string `json:"msg"`
	Error  string `json:"error,omitempty"`
}

type LinkResponse struct {
	ExpiresAt string `json:"expiresAt"`
	Url       string `json:"url"`
}

type ThumbnailResponse struct {
	ExpiresAt string `json:"expiresAt"`
	Url       string `json:"url"`
}

type ThumbnailsResponse struct {
	ExpiresAt     string            `json:"expiresAt"`
	FilenameToUrl map[string]string `json:"filenameToUrl"`
}
