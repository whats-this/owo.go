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
	"context"
	"errors"
	"fmt"
	"time"

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

	Describe("The file uploader", func() {
		var client *Client
		JustBeforeEach(func() {
			client = NewClient(APITestKey, server.URL, APITestFileUpload, APITestShorten, nil)
		})

		It("should succefully upload a valid, single file", func() {
			files, err := FilesToNamedReaders([]string{"./test-files/valid1.txt"})
			Expect(err).ToNot(HaveOccurred())
			res, err := client.UploadFile(nil, files[0])
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Success).To(BeTrue())
			Expect(res.Files[0].Name).To(Equal("valid1.txt"))
		})

		It("should not upload an invalid, single file", func() {
			res, err := client.UploadFile(nil, NamedReader{
				Reader:   &FailReader{errorMessage: "single file fail", allowedTries: 0},
				Filename: "",
			})
			Expect(err).To(MatchError("single file fail"))
			Expect(res).To(BeNil())
		})

		It("should fail if the file can't be read", func() {
			res, err := client.UploadFiles(nil, []NamedReader{
				NamedReader{
					Reader:   &FailReader{errorMessage: "failed to readAll", allowedTries: 0},
					Filename: "",
				},
			})
			Expect(err).To(MatchError("failed to readAll"))
			Expect(res).To(BeNil())
		})

		It("should fail on context timeout", func() {
			files, err := FilesToNamedReaders([]string{"./test-files/valid1.txt", "./test-files/valid2.txt"})
			Expect(err).ToNot(HaveOccurred())
			client.APIRoot = "http://example"
			ctx, cancelContext := context.WithTimeout(context.Background(), 200*time.Millisecond)
			res, err := client.UploadFiles(ctx, files)
			cancelContext()
			Expect(res).To(BeNil())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("net/http: request canceled while waiting for connection"))
		})

		It("should fail on invalid json", func() {
			files, err := FilesToNamedReaders([]string{"./test-files/valid1.txt", "./test-files/valid2.txt"})
			Expect(err).ToNot(HaveOccurred())
			client.APIFileUploadEndpoint = client.APIFileUploadEndpoint + "/json"
			res, err := client.UploadFiles(nil, files)
			Expect(res).To(BeNil())
			Expect(err).To(MatchError("invalid character 'i' in literal false (expecting 'l')"))
		})

		It("should fail if the server returns unsuccessful", func() {
			files, err := FilesToNamedReaders([]string{"./test-files/valid1.txt", "./test-files/valid2.txt"})
			Expect(err).ToNot(HaveOccurred())
			client.APIFileUploadEndpoint = client.APIFileUploadEndpoint + "/fail"
			res, err := client.UploadFiles(nil, files)
			Expect(res).To(BeNil())
			Expect(err).To(MatchError(fmt.Sprintf("Upload failed with code %d and message '%s'", 1, "fail")))
		})

		It("should succeed", func() {
			files, err := FilesToNamedReaders([]string{"./test-files/valid1.txt", "./test-files/valid.uha", "./test-files/valid.png"})
			Expect(err).ToNot(HaveOccurred())
			res, err := client.UploadFiles(nil, files)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Success).To(BeTrue())
			Expect(res.Files[0].Name).To(Equal("valid1.txt"))
			Expect(res.Files[1].Name).To(Equal("valid.uha"))
			Expect(res.Files[2].Name).To(Equal("valid.png"))
		})
	})
})

type FailReader struct {
	errorMessage string
	count        int
	allowedTries int
}

func (f *FailReader) Read(p []byte) (n int, err error) {
	f.count = f.count + 1
	if f.count > f.allowedTries {
		return 0, errors.New(f.errorMessage)
	} else {
		p = []byte{0}
		return 1, nil
	}
}
