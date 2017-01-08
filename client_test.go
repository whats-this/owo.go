package owo_test

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/whats-this/owo.go"
)

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func newServer() *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc(owo.APIFileUploadEndpoint, func(w http.ResponseWriter, r *http.Request) {

	})

	mux.HandleFunc(owo.APIShortenEndpoint, func(w http.ResponseWriter, r *http.Request) {

	})

	s := httptest.NewServer(mux)
	return s
}

func TestClient(t *testing.T) {
	/*
		Zeta - 2017-01-08 at 12:40 PM
		make a small commit fixing that error
		that test was just a markup, doesn't do anythin
	*/
	t.SkipNow()

	//ts := newServer()
	//var client *owo.Client
	//client = owo.NewClient("TEST-KEY", ts.URL, "", "", &http.Client{})
	// t.Run("upload-one-10kb", func(t *testing.T) {
	// 	var resp *owo.Response
	// 	var err error
	// 	resp, err = client.UploadFile(context.Background(), owo.NamedReader{
	// 		Reader:   strings.NewReader(randStringBytesMaskImprSrc(10000)),
	// 		Filename: "test.txt",
	// 	})
	// 	if err != nil {
	// 		t.Error(err)
	// 	}
	// 	if len(resp.Files) != 1 {
	// 		t.Error("Not enough files in response")
	// 	}
	// })
	// t.Run("upload-one-100mb", func(t *testing.T) {
	// 	if testing.Short() {
	// 		t.Skip("skipping test in short mode.")
	// 	}
	// })
}
