package imgopt

import "time"

type UploadOpts struct {
	Secret string
}

type DownloadOpts struct {
	Secret string
}

type CreateLinkOpts struct {
	Expires *time.Duration
	Secret  string
}

type CreateThumbnailLinkOpts struct {
	Expires *time.Duration
	Square  bool
	Secret  string
}

type CreateThumbnailLinksOpts struct {
	Expires *time.Duration
	Square  bool
	Secret  string
}
