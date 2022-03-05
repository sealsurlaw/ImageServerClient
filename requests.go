package gouvre

type CreateLinkRequest struct {
	Filename string `json:"filename"`
	Secret   string `json:"secret,omitempty"`
}

type CreateUploadLinkRequest struct {
	Filename string `json:"filename"`
	Secret   string `json:"secret,omitempty"`
}

type CreateThumbnailLinkRequest struct {
	Resolution int    `json:"resolution"`
	Filename   string `json:"filename"`
	Secret     string `json:"secret,omitempty"`
}

type CreateBatchThumbnailLinksRequest struct {
	Resolution int      `json:"resolution"`
	Filenames  []string `json:"filenames"`
	Secret     string   `json:"secret,omitempty"`
}
