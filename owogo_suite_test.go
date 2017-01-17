package owo_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	. "github.com/whats-this/owo.go"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

const APITestKey = "TEST-KEY"
const APITestFileUpload = "/pomf"
const APITestShorten = "/polr"

var server *httptest.Server

func TestOwoGo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Owo.Go Suite")
}

var _ = BeforeSuite(func() {
	mux := http.NewServeMux()

	mux.HandleFunc(APITestFileUpload, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != APITestKey {
			json.NewEncoder(w).Encode(&Response{
				Success:     false,
				Description: "APIKey",
				Errorcode:   3,
			})
			return
		}

		err := r.ParseMultipartForm(32 << 20) // 32 MB memory
		if err != nil {
			json.NewEncoder(w).Encode(&Response{
				Success:     false,
				Description: "ParseMultipartForm",
				Errorcode:   2,
			})
			return
		}

		files := r.MultipartForm.File["files[]"]
		jsonFiles := make([]File, 0)
		for _, file := range files {
			jsonFiles = append(jsonFiles, File{Name: file.Filename})
		}

		json.NewEncoder(w).Encode(&Response{
			Success: true,
			Files:   jsonFiles,
		})
	})

	mux.HandleFunc(APITestFileUpload+"/fail", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(&Response{
			Success:     false,
			Description: "fail",
			Errorcode:   1,
		})
	})

	mux.HandleFunc(APITestFileUpload+"/json", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("fail"))
	})

	server = httptest.NewServer(mux)
})

var _ = AfterSuite(func() {
	server.CloseClientConnections()
	server.Close()
})
