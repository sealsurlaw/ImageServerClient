package imgclient

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

	"github.com/sealsurlaw/gouvre-go-client/imgopt"
	"github.com/sealsurlaw/gouvre-go-client/imgreq"
	"github.com/sealsurlaw/gouvre-go-client/imgres"
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

func (c *Client) UploadImage(filename string, fileData []byte, opts ...imgopt.UploadOpts) error {
	var opt imgopt.UploadOpts
	if len(opts) == 0 {
		opt = imgopt.UploadOpts{}
	} else {
		opt = opts[0]
	}

	res, err := c.postFormFile(filename, fileData, opt.Secret)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		bodyData, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf(string(bodyData))
	}
	return nil
}

func (c *Client) DownloadImage(filename string, opts ...imgopt.DownloadOpts) ([]byte, error) {
	var opt imgopt.DownloadOpts
	if len(opts) == 0 {
		opt = imgopt.DownloadOpts{}
	} else {
		opt = opts[0]
	}

	downloadUrl := fmt.Sprintf("%s/images/%s", c.baseUrl, filename)
	if opt.Secret != "" {
		downloadUrl = fmt.Sprintf("%s?secret=%s", downloadUrl, opt.Secret)
	}
	res, err := c.get(downloadUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	fileData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return fileData, nil
}

func (c *Client) DownloadImageByToken(token int64) ([]byte, error) {
	downloadUrl := fmt.Sprintf("%s/links/%d", c.baseUrl, token)
	res, err := c.get(downloadUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	fileData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return fileData, nil
}

func (c *Client) CreateLink(
	filename string,
	opts ...imgopt.CreateLinkOpts,
) (*imgres.LinkResponse, error) {
	var opt imgopt.CreateLinkOpts
	if len(opts) == 0 {
		opt = imgopt.CreateLinkOpts{}
	} else {
		opt = opts[0]
	}

	req := &imgreq.CreateLinkRequest{
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

	bodyData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf(string(bodyData))
	}

	linkResponse := &imgres.LinkResponse{}
	err = json.Unmarshal(bodyData, linkResponse)
	if err != nil {
		return nil, err
	}

	return linkResponse, nil
}

func (c *Client) CreateThumbnailLink(
	resolution int,
	filename string,
	opts ...imgopt.CreateThumbnailLinkOpts,
) (*imgres.ThumbnailResponse, error) {
	var opt imgopt.CreateThumbnailLinkOpts
	if len(opts) == 0 {
		opt = imgopt.CreateThumbnailLinkOpts{}
	} else {
		opt = opts[0]
	}

	req := &imgreq.CreateThumbnailRequest{
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

	thumbnailUrl := fmt.Sprintf("%s/thumbnails?%s", c.baseUrl, queryParams.Encode())

	res, err := c.post(thumbnailUrl, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	bodyData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(string(bodyData))
	}

	linkResponse := &imgres.ThumbnailResponse{}
	err = json.Unmarshal(bodyData, linkResponse)
	if err != nil {
		return nil, err
	}

	return linkResponse, nil
}

func (c *Client) CreateBatchThumbnailLinks(
	resolution int,
	filenames []string,
	opts ...imgopt.CreateThumbnailLinksOpts,
) (*imgres.ThumbnailsResponse, error) {
	var opt imgopt.CreateThumbnailLinksOpts
	if len(opts) == 0 {
		opt = imgopt.CreateThumbnailLinksOpts{}
	} else {
		opt = opts[0]
	}

	req := &imgreq.CreateThumbnailsRequest{
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

	thumbnailUrl := fmt.Sprintf("%s/thumbnails/batch?%s", c.baseUrl, queryParams.Encode())

	res, err := c.post(thumbnailUrl, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	bodyData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(string(bodyData))
	}

	linkResponse := &imgres.ThumbnailsResponse{}
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

func (c *Client) postFormFile(filename string, fileData []byte, secret string) (*http.Response, error) {
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

	err = writer.WriteField("filename", filename)
	if err != nil {
		return nil, err
	}

	if secret != "" {
		err = writer.WriteField("secret", secret)
		if err != nil {
			return nil, err
		}
	}

	writer.Close()

	uploadUrl := fmt.Sprintf("%s/images", c.baseUrl)
	req, err := http.NewRequest(http.MethodPost, uploadUrl, bytes.NewReader(body.Bytes()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	return http.DefaultClient.Do(req)
}
