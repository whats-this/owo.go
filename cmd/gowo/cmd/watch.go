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
	"log"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/fsnotify/fsnotify"
	"github.com/whats-this/owo.go"
)

var queue = make(map[string]struct{})
var queueLock = sync.Mutex{}

// watchCmd represents the upload command
var watchCmd = &cobra.Command{
	Use:     "watch",
	Aliases: []string{"w", "this"},
	Short:   "Watches directory and uploads new files to OwO",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("Need at least one location.")
		}

		cdn := viper.GetString("cdn")

		w, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}

		for _, arg := range args {
			err = w.Add(arg)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf(`Added "%s" to watcher.`, arg)
		}

		go func() {
			ticker := time.NewTicker(time.Second)
			for range ticker.C {
				var names []string
				queueLock.Lock()
				for name := range queue {
					if len(names) == owo.FileCountLimit {
						break
					}
					names = append(names, name)
					delete(queue, name)
				}
				queueLock.Unlock()

				if len(names) == 0 {
					continue
				}

				go func() {
					err := doUpload(cdn, names)
					if err != nil {
						log.Print(err)
					}
				}()
			}
		}()

		log.Print("Started watching for events in given destinations.")
		for ev := range w.Events {
			if ev.Op&fsnotify.Create == fsnotify.Create {
				name := ev.Name
				go func() {
					time.Sleep(2 * time.Second)
					queueLock.Lock()
					queue[name] = struct{}{}
					queueLock.Unlock()
				}()
			}
		}
	},
}

func init() {
	uploadCmd.AddCommand(watchCmd)
}
