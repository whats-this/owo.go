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
	"context"
	"log"

	"github.com/spf13/cobra"

	"fmt"

	"bytes"

	"github.com/atotto/clipboard"
	"github.com/spf13/viper"
	"github.com/whats-this/owo.go"
)

func DoUpload(cdn string, names []string) {
	files, err := owo.FilesToNamedReaders(names)
	if err != nil {
		log.Println("[upload]", err)
		return
	}
	response, err := owo.UploadFiles(context.Background(), files)
	if err != nil {
		log.Println("[upload]", err)
		return
	}
	if !response.Success {
		log.Printf("[upload] %d: %s", response.Errorcode, response.Description)
		return
	}
	buf := bytes.Buffer{}
	for _, file := range response.Files {
		if file.Error {
			log.Printf("%d: %s", file.Errorcode, file.Description)
			continue
		}
		fmt.Fprintf(&buf, "%s\n", file.WithCDN(cdn))
	}
	if shouldClipboard {
		err = clipboard.WriteAll(buf.String())
		if err != nil {
			log.Println("[upload]", err)
			return
		}
		log.Printf("Wrote %d URLs to clipboard", len(response.Files))
	} else {
		fmt.Print(buf.String())
	}
	response = nil
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
		DoUpload(cdn, args)
	},
}

func init() {
	RootCmd.AddCommand(uploadCmd)
}
