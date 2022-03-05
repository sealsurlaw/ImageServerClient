package gouvre

type CreateLinkRequest struct {
	Filename string `json:"filename"`
	Secret   string `json:"secret"`
}

type CreateUploadLinkRequest struct {
	Filename    string `json:"filename"`
	Secret      string `json:"secret"`
	Resolutions []int  `json:"resolutions"`
}

type CreateThumbnailLinkRequest struct {
	Resolution int    `json:"resolution"`
	Filename   string `json:"filename"`
	Secret     string `json:"secret"`
}

type CreateBatchThumbnailLinksRequest struct {
	Resolution int      `json:"resolution"`
	Filenames  []string `json:"filenames"`
	Secret     string   `json:"secret"`
}
