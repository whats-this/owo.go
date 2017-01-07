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

type (
	// Response contains json marshalled response from owo
	Response struct {
		Success     bool   `json:"success"`
		Errorcode   int    `json:"errorcode"`
		Description string `json:"description"`
		Files       []File `json:"files"`
	}

	// File represents a single file from json response (if there were no errors)
	File struct {
		Hash        string `json:"hash,omitempty"`
		Name        string `json:"name,omitempty"`
		URL         string `json:"url,omitempty"`
		Size        int    `json:"size,omitempty"`
		Error       bool   `json:"error,omitempty"`
		Errorcode   int    `json:"errorcode,omitempty"`
		Description string `json:"description,omitempty"`
	}
)

// WithCDN returns file url prefixed with the CDN
func (f File) WithCDN(cdn string) string {
	return cdn + f.URL
}
