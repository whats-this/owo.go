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

package cmd

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
	"github.com/whats-this/owo.go"
)

func doUpload(cdn string, names []string) (err error) {
	var files []owo.NamedReader
	files, err = owo.FilesToNamedReaders(names)
	if err != nil {
		return
	}
	var response *owo.Response
	response, err = owo.DefaultClient().UploadFiles(context.Background(), files)
	if err != nil {
		return
	}
	buf := bytes.Buffer{}
	var url string
	for _, file := range response.Files {
		url, err = file.WithCDN(cdn)
		if err != nil {
			log.Println("[upload]", err)
		}
		fmt.Fprintf(&buf, "%s\n", url)
	}
	err = output(buf.String(), len(response.Files))
	if err != nil {
		return
	}
	response = nil
	return
}

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:     "upload",
	Aliases: []string{"up", "whats"},
	Short:   "Upload files to OwO",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("Need at least one file.")
		}
		cdn := viper.GetString("cdn")
		if err := doUpload(cdn, args); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(uploadCmd)
}
