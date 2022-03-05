package gouvre

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
)

type Client struct {
	baseUrl string
	token   string
}

func NewClient(baseUrl, token string) *Client {
	return &Client{
		baseUrl: baseUrl,
		token:   token,
	}
}

func (c *Client) UploadImage(filename string, fileData []byte, opts ...UploadOpts) error {
	opt := UploadOpts{}
	if len(opts) > 0 {
		opt = opts[0]
	}

	uploadUrl := fmt.Sprintf("%s/uploads", c.baseUrl)
	res, err := c.postFormFile(uploadUrl, filename, fileData, opt.Secret)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if _, err = tryParseData(res); err != nil {
		return err
	}
	return nil
}

func (c *Client) UploadImageFromToken(token string, fileData []byte, opts ...UploadOpts) error {
	opt := UploadOpts{}
	if len(opts) > 0 {
		opt = opts[0]
	}

	uploadUrl := fmt.Sprintf("%s/uploads/%s", c.baseUrl, token)
	res, err := c.postFormFile(uploadUrl, "", fileData, opt.Secret)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if _, err = tryParseData(res); err != nil {
		return err
	}
	return nil
}

func (c *Client) DownloadImage(filename string, opts ...DownloadOpts) ([]byte, error) {
	opt := DownloadOpts{}
	if len(opts) > 0 {
		opt = opts[0]
	}

	queryParams := url.Values{}
	if opt.Secret != "" {
		queryParams.Add("secret", opt.Secret)
	}

	downloadUrl := fmt.Sprintf("%s/images/%s?%s", c.baseUrl, filename, queryParams.Encode())
	res, err := c.get(downloadUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	fileData, err := tryParseData(res)
	if err != nil {
		return nil, err
	}

	return fileData, nil
}

func (c *Client) DownloadImageByToken(token string, opts ...DownloadByTokenOpts) ([]byte, error) {
	opt := DownloadByTokenOpts{}
	if len(opts) > 0 {
		opt = opts[0]
	}

	queryParams := url.Values{}
	if opt.Secret != "" {
		queryParams.Add("secret", opt.Secret)
	}

	downloadUrl := fmt.Sprintf("%s/uploads/%s?%s", c.baseUrl, token, queryParams.Encode())
	res, err := c.get(downloadUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	fileData, err := tryParseData(res)
	if err != nil {
		return nil, err
	}

	return fileData, nil
}

func (c *Client) CreateLink(
	filename string,
	opts ...CreateLinkOpts,
) (*CreateLinkResponse, error) {
	opt := CreateLinkOpts{}
	if len(opts) > 0 {
		opt = opts[0]
	}

	req := &CreateLinkRequest{
		Filename: filename,
		Secret:   opt.Secret,
	}

	queryParams := url.Values{}
	if opt.Expires != nil {
		queryParams.Add("expires", opt.Expires.String())
	}

	linkUrl := fmt.Sprintf("%s/links?%s", c.baseUrl, queryParams.Encode())

	res, err := c.post(linkUrl, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	bodyData, err := tryParseData(res)
	if err != nil {
		return nil, err
	}

	linkResponse := &CreateLinkResponse{}
	err = json.Unmarshal(bodyData, linkResponse)
	if err != nil {
		return nil, err
	}

	return linkResponse, nil
}

func (c *Client) CreateUploadLink(
	filename string,
	opts ...CreateUploadLinkOpts,
) (*CreateUploadLinkResponse, error) {
	opt := CreateUploadLinkOpts{}
	if len(opts) > 0 {
		opt = opts[0]
	}

	req := &CreateUploadLinkRequest{
		Filename:    filename,
		Secret:      opt.Secret,
		Resolutions: opt.Resolutions,
	}

	queryParams := url.Values{}
	if opt.Expires != nil {
		queryParams.Add("expires", opt.Expires.String())
	}

	linkUrl := fmt.Sprintf("%s/links/upload?%s", c.baseUrl, queryParams.Encode())

	res, err := c.post(linkUrl, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	bodyData, err := tryParseData(res)
	if err != nil {
		return nil, err
	}

	linkResponse := &CreateUploadLinkResponse{}
	err = json.Unmarshal(bodyData, linkResponse)
	if err != nil {
		return nil, err
	}

	return linkResponse, nil
}

func (c *Client) CreateThumbnailLink(
	resolution int,
	filename string,
	opts ...CreateThumbnailLinkOpts,
) (*CreateThumbnailLinkResponse, error) {
	opt := CreateThumbnailLinkOpts{}
	if len(opts) > 0 {
		opt = opts[0]
	}

	req := &CreateThumbnailLinkRequest{
		Resolution: resolution,
		Filename:   filename,
		Secret:     opt.Secret,
	}

	queryParams := url.Values{
		"square": {strconv.FormatBool(opt.Square)},
	}
	if opt.Expires != nil {
		queryParams.Add("expires", opt.Expires.String())
	}

	thumbnailUrl := fmt.Sprintf("%s/links/thumbnails?%s", c.baseUrl, queryParams.Encode())

	res, err := c.post(thumbnailUrl, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	bodyData, err := tryParseData(res)
	if err != nil {
		return nil, err
	}

	linkResponse := &CreateThumbnailLinkResponse{}
	err = json.Unmarshal(bodyData, linkResponse)
	if err != nil {
		return nil, err
	}

	return linkResponse, nil
}

func (c *Client) CreateBatchThumbnailLinks(
	resolution int,
	filenames []string,
	opts ...CreateBatchThumbnailLinksOpts,
) (*CreateBatchThumbnailLinksResponse, error) {
	opt := CreateBatchThumbnailLinksOpts{}
	if len(opts) > 0 {
		opt = opts[0]
	}

	req := &CreateBatchThumbnailLinksRequest{
		Resolution: resolution,
		Filenames:  filenames,
		Secret:     opt.Secret,
	}

	queryParams := url.Values{
		"square": {strconv.FormatBool(opt.Square)},
	}
	if opt.Expires != nil {
		queryParams.Add("expires", opt.Expires.String())
	}

	thumbnailUrl := fmt.Sprintf("%s/links/thumbnails/batch?%s", c.baseUrl, queryParams.Encode())

	res, err := c.post(thumbnailUrl, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	bodyData, err := tryParseData(res)
	if err != nil {
		return nil, err
	}

	linkResponse := &CreateBatchThumbnailLinksResponse{}
	err = json.Unmarshal(bodyData, linkResponse)
	if err != nil {
		return nil, err
	}

	return linkResponse, nil
}

func (c *Client) get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	return http.DefaultClient.Do(req)
}

func (c *Client) post(url string, jsonRequest interface{}) (*http.Response, error) {
	data, err := json.Marshal(jsonRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	return http.DefaultClient.Do(req)
}

func (c *Client) postFormFile(url string, filename string, fileData []byte, secret string) (*http.Response, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormFile("file", "file")
	if err != nil {
		return nil, err
	}

	r := bytes.NewBuffer(fileData)
	_, err = io.Copy(fw, r)
	if err != nil {
		return nil, err
	}

	if filename != "" {
		err = writer.WriteField("filename", filename)
		if err != nil {
			return nil, err
		}
	}

	if secret != "" {
		err = writer.WriteField("secret", secret)
		if err != nil {
			return nil, err
		}
	}

	writer.Close()

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body.Bytes()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	return http.DefaultClient.Do(req)
}

func tryParseData(res *http.Response) ([]byte, error) {
	bodyData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf(string(bodyData))
	}

	return bodyData, nil
}
