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

// APIRoot defines base url for the service's API endpoints
const APIRoot = "https://api.awau.moe"

var (
	// APIFileUploadURL to send POSTs with files to
	APIFileUploadURL = APIRoot + "/upload/pomf"
	// APIShortenURL to send GETs with urls to
	APIShortenURL = APIRoot + "/shorten/polr"

	// CDNs represnts a list of CDNs available as of 27/12/2016
	CDNs = []string{
		"https://owo.whats-th.is/",
		"https://i.am-a.ninja/",
		"https://buttsare.sexy/",
		"https://nyanyanya.moe/",
		"https://all.foxgirlsare.sexy/",
		"https://i.stole-a-me.me/",
		"https://can-i-ask-dean-on-a.date/",
	}
)
