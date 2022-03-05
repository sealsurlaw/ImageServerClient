package gouvre

type ErrorResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Msg    string `json:"msg"`
	Error  string `json:"error,omitempty"`
}

type CreateLinkResponse struct {
	ExpiresAt string `json:"expiresAt"`
	Url       string `json:"url"`
}

type CreateUploadLinkResponse struct {
	ExpiresAt string `json:"expiresAt"`
	Url       string `json:"url"`
}

type CreateThumbnailLinkResponse struct {
	ExpiresAt string `json:"expiresAt"`
	Url       string `json:"url"`
}

type CreateBatchThumbnailLinksResponse struct {
	ExpiresAt     string            `json:"expiresAt"`
	FilenameToUrl map[string]string `json:"filenameToUrl"`
}
