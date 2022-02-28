package imgopt

import "time"

type CreateLinkOpts struct {
	Expires *time.Duration
}

type CreateThumbnailLinkOpts struct {
	Expires *time.Duration
	Square  bool
}

type CreateThumbnailLinksOpts struct {
	Expires *time.Duration
	Square  bool
}
