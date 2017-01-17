// MIT License

// Copyright (c) 2017 @cking / @Kura Bloodlust#8777 / <im@z0ne.moe>

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

package owo_test

import (
	. "github.com/whats-this/owo.go"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	Describe("The Constructor", func() {
		It("should create an empty client if no parameters are specified", func() {
			client := NewClient("", "", "", "", nil)
			Expect(client.APIFileUploadEndpoint).To(Equal(APIFileUploadEndpoint))
			Expect(DefaultClient().APIRoot).To(Equal(OfficialAPIRoot))
			Expect(client.APIShortenEndpoint).To(Equal(APIShortenEndpoint))
			Expect(client.Key).To(BeEmpty())
		})

		It("should create a new client with the passed parameters", func() {
			client := NewClient("key", "root", "upload", "shorten", nil)
			Expect(client.APIFileUploadEndpoint).To(Equal("upload"))
			Expect(client.APIRoot).To(Equal("root"))
			Expect(client.APIShortenEndpoint).To(Equal("shorten"))
			Expect(client.Key).To(Equal("key"))
		})

		It("should have a default client available, pointing to the official owo endpoint", func() {
			Expect(DefaultClient().APIRoot).To(Equal(OfficialAPIRoot))
		})
	})
})
