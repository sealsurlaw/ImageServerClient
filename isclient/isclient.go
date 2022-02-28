package isclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/sealsurlaw/ImageServerClient/isopt"
	"github.com/sealsurlaw/ImageServerClient/isres"
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

func (c *Client) UploadImage(filename string, file *os.File) error {
	res, err := c.postFormFile(filename, file)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyData, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf(string(bodyData))
	}

	return nil
}

func (c *Client) DownloadImage(filename string) ([]byte, error) {
	downloadUrl := fmt.Sprintf("%s/download/%s", c.baseUrl, filename)
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
	downloadUrl := fmt.Sprintf("%s/link/%d", c.baseUrl, token)
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

func (c *Client) CreateLink(filename string, opts ...isopt.CreateLinkOpts) (*isres.LinkResponse, error) {
	var opt isopt.CreateLinkOpts
	if len(opts) == 0 {
		opt = isopt.CreateLinkOpts{}
	} else {
		opt = opts[0]
	}

	data := url.Values{
		"filename": {filename},
	}

	queryParams := url.Values{}
	if opt.Expires != nil {
		queryParams.Add("expires", opt.Expires.String())
	}

	linkUrl := fmt.Sprintf("%s/link?%s", c.baseUrl, queryParams.Encode())

	res, err := c.postForm(linkUrl, data)
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

	linkResponse := &isres.LinkResponse{}
	err = json.Unmarshal(bodyData, linkResponse)
	if err != nil {
		return nil, err
	}

	return linkResponse, nil
}

func (c *Client) CreateThumbnailLink(filename string, resolution int, opts ...isopt.CreateThumbnailLinkOpts) (*isres.LinkResponse, error) {
	var opt isopt.CreateThumbnailLinkOpts
	if len(opts) == 0 {
		opt = isopt.CreateThumbnailLinkOpts{}
	} else {
		opt = opts[0]
	}

	data := url.Values{
		"filename":   {filename},
		"resolution": {strconv.Itoa(resolution)},
	}

	queryParams := url.Values{
		"cropped": {strconv.FormatBool(opt.Cropped)},
	}
	if opt.Expires != nil {
		queryParams.Add("expires", opt.Expires.String())
	}

	thumbnailUrl := fmt.Sprintf("%s/thumbnail?%s", c.baseUrl, queryParams.Encode())

	res, err := c.postForm(thumbnailUrl, data)
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

	linkResponse := &isres.LinkResponse{}
	err = json.Unmarshal(bodyData, linkResponse)
	if err != nil {
		return nil, err
	}

	return linkResponse, nil
}

func (c *Client) postFormFile(filename string, file *os.File) (*http.Response, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(fw, file)
	if err != nil {
		return nil, err
	}

	err = writer.WriteField("filename", filename)
	if err != nil {
		return nil, err
	}

	writer.Close()

	uploadUrl := fmt.Sprintf("%s/upload", c.baseUrl)
	req, err := http.NewRequest(http.MethodPost, uploadUrl, bytes.NewReader(body.Bytes()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	return http.DefaultClient.Do(req)
}

func (c *Client) postForm(url string, data url.Values) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	return http.DefaultClient.Do(req)
}

func (c *Client) get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	return http.DefaultClient.Do(req)
}
