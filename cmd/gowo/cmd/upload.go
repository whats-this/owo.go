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
	"os"

	"github.com/spf13/cobra"

	"fmt"

	"github.com/spf13/viper"
	"github.com/whats-this/owo.go"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:     "upload",
	Aliases: []string{"up"},
	Short:   "Upload files to OwO",
	Run: func(cmd *cobra.Command, args []string) {
		cdn := viper.GetString("cdn")
		files := make([]owo.NamedReader, len(args))
		for idx, arg := range args {
			file, err := os.Open(arg)
			if err != nil {
				log.Fatal(err)
			}
			files[idx] = owo.NamedReader{file, arg}
			defer file.Close()
		}
		response, err := owo.UploadFiles(context.Background(), files)
		if err != nil {
			log.Fatal(err)
		}
		if !response.Success {
			log.Fatalf("%d: %s", response.Errorcode, response.Description)
		}
		for _, file := range response.Files {
			if file.Error {
				log.Printf("%d: %s", file.Errorcode, file.Description)
				continue
			}
			fmt.Println(file.WithCDN(cdn))
		}
	},
}

func init() {
	RootCmd.AddCommand(uploadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uploadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uploadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
