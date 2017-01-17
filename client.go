// MIT License

// Copyright (c) 2016 @zet4 / @Zeta#2229 / <my-name-is-zeta@and.my.foxgirlsare.sexy>

// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package owo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	nhttp "net/http"
	"net/textproto"
	"net/url"
	"path"
	"strings"
	"time"
)

type (
	// Client stores http client and key
	Client struct {
		Key                   string
		APIRoot               string
		APIFileUploadEndpoint string
		APIShortenEndpoint    string
		http                  *nhttp.Client
	}

	// NamedReader wrapper for a single file to upload
	NamedReader struct {
		Reader   io.Reader
		Filename string
	}
)

var (
	global = NewClient("", OfficialAPIRoot, "", "", &nhttp.Client{Timeout: time.Minute})
)

// NewClient returns a fully configured client
func NewClient(key, root, upload, shorten string, http *nhttp.Client) *Client {
	if http == nil {
		http = &nhttp.Client{}
	}

	c := &Client{http: http}
	if key != "" {
		c.Key = key
	}

	if root != "" {
		c.APIRoot = root
	} else {
		c.APIRoot = OfficialAPIRoot
	}

	if upload != "" {
		c.APIFileUploadEndpoint = upload
	} else {
		c.APIFileUploadEndpoint = APIFileUploadEndpoint
	}

	if shorten != "" {
		c.APIShortenEndpoint = shorten
	} else {
		c.APIShortenEndpoint = APIShortenEndpoint
	}

	return c
}

// DefaultClient returns a client with the default, official parameters
func DefaultClient() *Client {
	return global
}

// UploadFile uploads a file
func (o *Client) UploadFile(ctx context.Context, r NamedReader) (response *Response, err error) {
	return o.UploadFiles(ctx, []NamedReader{r})
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

// UploadFiles uploads multiple files using the client
func (o *Client) UploadFiles(ctx context.Context, rs []NamedReader) (response *Response, err error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for _, r := range rs {
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="files[]"; filename="%s"`, escapeQuotes(r.Filename)))

		var buf bytes.Buffer
		var b []byte
		tee := io.TeeReader(r.Reader, &buf)
		b, err = ioutil.ReadAll(tee)
		if err != nil {
			return
		}
		contenttype := nhttp.DetectContentType(b)
		if contenttype == "application/octet-stream" {
			contenttype = mime.TypeByExtension(path.Ext(r.Filename))
			if contenttype == "" {
				contenttype = "application/octet-stream"
			}
		}
		h.Set("Content-Type", contenttype)

		// no error checking necessary
		// - `writer` uses `body` (writes to bytes.Buffer). It only throws if it runs out of memory.
		part, _ := writer.CreatePart(h)

		// no error checking necessary
		// - `part` is a valid destination (indirect write to bytes.Buffer). It only throws if it runs out of memory.
		// - `buf` is always a valid source (bytes.Buffer). It always returns data or EOF (which is not an error for Copy).
		io.Copy(part, &buf)
	}

	// no error checking necessary
	// - `writer` uses `body` (writes to bytes.Buffer). It only throws if it runs out of memory.
	err = writer.Close()

	req, err := nhttp.NewRequest("POST", o.APIRoot+o.APIFileUploadEndpoint, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", o.Key)

	resp, err := o.http.Do(req)
	if err != nil {
		return
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return
	}
	if !response.Success {
		return nil, &ErrUploadFailed{response.Description, response.Errorcode}
	}

	return
}

/*
   Polr compatible url shortening wrapper
*/

// ShortenURLs shortens multiple urls
func (o *Client) ShortenURLs(ctx context.Context, urls []string) (shortened []string, err error) {
	shortened = make([]string, len(urls))
	for idx, u := range urls {
		result, err := o.ShortenURL(ctx, u)
		shortened[idx] = result
		if err != nil {
			return shortened, err
		}
	}
	return
}

// ShortenURL shortens a single url
func (o *Client) ShortenURL(ctx context.Context, u string) (shortened string, err error) {
	v := url.Values{}
	v.Set("key", o.Key)
	v.Set("action", "shorten")
	v.Add("url", u)
	au, err := url.Parse(o.APIRoot + APIShortenEndpoint)
	if err != nil {
		return
	}
	au.RawQuery = v.Encode()

	req, err := nhttp.NewRequest("GET", au.String(), nil)
	if err != nil {
		return
	}
	req = req.WithContext(ctx)

	resp, err := o.http.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		resp.Body.Close()
		return "", fmt.Errorf("Unexpected status code '%d', expected 200.", resp.StatusCode)
	}
	respstr, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return
	}
	shortened = string(respstr)
	return
}
