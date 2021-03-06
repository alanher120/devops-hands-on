/**
 * Copyright 2017 Google Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// [START all]
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"os/signal"
	"syscall"
	"sync"
)

func main() {

	started := time.Now()

	// use PORT environment variable, or default to 8080
	port := "8080"
	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = fromEnv
	}

	// register hello function to handle all requests
	server := http.NewServeMux()
	server.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Serving request: %s", r.URL.Path)
		duration := time.Since(started)
		//30s前，此接口自動無法運作
		if duration.Seconds() < 15 {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("error: %v", duration.Seconds())))
		} else {
			w.WriteHeader(200)
			host, _ := os.Hostname()
			fmt.Fprintf(w, "Hello, world!  Hostname: %s  time: %v \n", host,duration.Seconds())
		}
	})
	server.HandleFunc("/long", func(w http.ResponseWriter, r *http.Request) {
		for i:=0; i < 20; i++ {
			log.Printf("Sleep %d", (i+1))
			time.Sleep(time.Second)
			
		}
		fmt.Fprintf(w, "Process finished \n")
	})

    //gracefully shutdown 1 Beg
	srv := &http.Server{
        Addr:           ":8080",
        Handler:        server,
	}

	//使用WaitGroup同步Goroutine
	var wg sync.WaitGroup
    // Handle SIGINT and SIGTERM.
    ch := make(chan os.Signal)
    signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-ch
		wg.Add(1)
		//使用context控制srv.Shutdown的超時時間
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		err := srv.Shutdown(ctx)
        if err != nil {
            fmt.Println(err)
        }
        wg.Done()
	}()
	//gracefully shutdown 1 End

	log.Printf("Server listening on port %s", port)
	
	//gracefully shutdown 2 Beg
    err := srv.ListenAndServe()

	fmt.Println("waiting for the remaining connections to finish...")
    wg.Wait()
    if err != nil && err != http.ErrServerClosed {
        panic(err)
    }
    fmt.Println("gracefully shutdown the http server...")
	//gracefully shutdown 2 End

}


// [END all]
