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
	"net/http"
	"net/textproto"
	"net/url"
	"path"
	"strings"
	"time"
)

type (
	// Client stores http client and key
	Client struct {
		Key    string
		Domain string
		http   *http.Client
	}

	// NamedReader wrapper for a single file to upload
	NamedReader struct {
		Reader   io.Reader
		Filename string
	}
)

var (
	global = Client{
		http: &http.Client{
			Timeout: time.Minute,
		},
	}
)

// SetKey changes global client's API key
func SetKey(k string) {
	global.Key = k
}

// UploadFile uploads a file using the global client
func UploadFile(ctx context.Context, r NamedReader) (response *Response, err error) {
	return global.UploadFile(ctx, r)
}

// UploadFile uploads a file
func (o *Client) UploadFile(ctx context.Context, r NamedReader) (response *Response, err error) {
	return o.UploadFiles(ctx, []NamedReader{r})
}

// UploadFiles uploads multiple files using the global client
func UploadFiles(ctx context.Context, rs []NamedReader) (response *Response, err error) {
	return global.UploadFiles(ctx, rs)
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

// UploadFiles uploads multiple files using the global client
func (o *Client) UploadFiles(ctx context.Context, rs []NamedReader) (response *Response, err error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for _, r := range rs {
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="files[]" filename="%s"`, escapeQuotes(r.Filename)))

		var buf bytes.Buffer
		var b []byte
		tee := io.TeeReader(r.Reader, &buf)
		b, err = ioutil.ReadAll(&buf)
		if err != nil {
			return
		}
		contenttype := http.DetectContentType(b)
		if contenttype == "application/octet-stream" {
			contenttype = mime.TypeByExtension(path.Ext(r.Filename))
			if contenttype == "" {
				contenttype = "application/octet-stream"
			}
		}
		h.Set("Content-Type", contenttype)

		part, err := writer.CreatePart(h)
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(part, tee)
		if err != nil {
			return nil, err
		}
	}
	err = writer.Close()
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", APIFileUploadURL, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", o.Key)
	req = req.WithContext(ctx)

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

// ShortenURL shortens a single url using the global client
func ShortenURL(ctx context.Context, url string) (shortened string, err error) {
	return global.ShortenURL(ctx, url)
}

// ShortenURLs shortens multiple urls using the global client
func ShortenURLs(ctx context.Context, urls []string) (shortened []string, err error) {
	return global.ShortenURLs(ctx, urls)
}

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
	au, err := url.Parse(APIShortenURL)
	if err != nil {
		return
	}
	au.RawQuery = v.Encode()

	req, err := http.NewRequest("GET", au.String(), nil)
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
