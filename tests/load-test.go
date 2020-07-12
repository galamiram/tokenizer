package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type result struct {
	response     string
	file         string
	duration     float64
	responseCode int
}

func listTestFiles(folder string) ([]string, error) {
	var files []string

	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})

	return files, err
}

func sendPayload(url *string, payload *[]byte) (*result, error) {
	start := time.Now()
	client := &http.Client{}

	req, err := http.NewRequest("GET", *url, bytes.NewReader(*payload))
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	elapsed := time.Since(start)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	resObj := result{
		duration:     float64(elapsed) / float64(time.Millisecond),
		response:     string(body),
		responseCode: res.StatusCode,
	}

	return &resObj, nil
}

func main() {
	fs := pflag.NewFlagSet("default", pflag.ContinueOnError)
	fs.String("dir", "", "Path to test cases files (e.g. /path/to/folder)")
	fs.String("url", "", "The URL of the tested endpoint, including port if not the default (e.g. http://example.com:8080/test)")
	fs.Int("parallelism", 1, "Number of concurent requests")

	err := fs.Parse(os.Args[1:])
	switch {
	case err == pflag.ErrHelp:
		os.Exit(0)
	case err != nil:
		fmt.Fprintf(os.Stderr, "Error: %s\n\n", err.Error())
		fs.PrintDefaults()
		os.Exit(1)
	}

	viper.BindPFlags(fs)

	files, err := listTestFiles(viper.GetString("dir"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n\n", err.Error())
		os.Exit(1)
	}

	guard := make(chan struct{}, viper.GetInt("parallelism"))
	filesLen := len(files)

	errChan := make(chan error, filesLen)
	resChan := make(chan *result, filesLen)

	for _, f := range files {
		guard <- struct{}{} // would block if guard channel is already filled
		go func(file, url string, resChan chan<- *result, errChan chan<- error) {
			//worker(files[i])
			dat, err := ioutil.ReadFile(file)
			if err != nil {
				errChan <- err
				resChan <- &result{}
			}

			res, err := sendPayload(&url, &dat)
			if err != nil {
				errChan <- err
				resChan <- res
				return
			}
			errChan <- err
			resChan <- res

			res.file = file
			fmt.Println(res)

			<-guard
		}(f, viper.GetString("url"), resChan, errChan)
	}
}
