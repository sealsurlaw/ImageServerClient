package isopt

import "time"

type CreateLinkOpts struct {
	Expires *time.Duration
}

type CreateThumbnailLinkOpts struct {
	Expires *time.Duration
	Cropped bool
}
