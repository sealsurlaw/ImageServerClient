package imgreq

type CreateLinkRequest struct {
	Filename string `json:"filename"`
	Secret   string `json:"secret,omitempty"`
}

type CreateThumbnailRequest struct {
	Resolution int    `json:"resolution"`
	Filename   string `json:"filename"`
	Secret     string `json:"secret,omitempty"`
}

type CreateThumbnailsRequest struct {
	Resolution int      `json:"resolution"`
	Filenames  []string `json:"filenames"`
	Secret     string   `json:"secret,omitempty"`
}
