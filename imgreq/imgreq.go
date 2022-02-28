package imgreq

type CreateLinkRequest struct {
	Filename string `json:"filename"`
}

type CreateThumbnailRequest struct {
	Resolution int    `json:"resolution"`
	Filename   string `json:"filename"`
}

type CreateThumbnailsRequest struct {
	Resolution int      `json:"resolution"`
	Filenames  []string `json:"filenames"`
}
