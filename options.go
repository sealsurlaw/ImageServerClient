package gouvre

import "time"

type UploadOpts struct {
	Secret string
}

type DownloadOpts struct {
	Secret string
}

type DownloadByTokenOpts struct {
	Secret string
}

type CreateLinkOpts struct {
	Expires *time.Duration
	Secret  string
}

type CreateUploadLinkOpts struct {
	Expires     *time.Duration
	Secret      string
	Resolutions []int
}

type CreateThumbnailLinkOpts struct {
	Expires *time.Duration
	Square  bool
	Secret  string
}

type CreateBatchThumbnailLinksOpts struct {
	Expires *time.Duration
	Square  bool
	Secret  string
}
