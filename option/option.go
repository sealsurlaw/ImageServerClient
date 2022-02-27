package option

import "time"

type CreateLinkOpts struct {
	Expires *time.Duration
}

type CreateThumbnailLinkOpts struct {
	Expires *time.Duration
	Cropped bool
}
