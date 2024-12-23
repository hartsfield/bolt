package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

func main() {
	var client *http.Client
	var remoteURL string = "http://" + os.Args[1]
	images, err := os.ReadDir("img")
	if err != nil {
		log.Println(err)
	}
	for _, img := range images {
		values := map[string]io.Reader{
			"Email":       strings.NewReader(os.Args[2]),
			"Title":       strings.NewReader(os.Args[3]),
			"Message":     strings.NewReader(os.Args[4]),
			"FileElement": mustOpen("img/" + img.Name()),
			"MediaType":   strings.NewReader(strings.Split(img.Name(), ".")[1]),
		}

		err := Upload(client, remoteURL, values)
		if err != nil {
			panic(err)
		}
	}
}

func Upload(client *http.Client, url string, values map[string]io.Reader) (err error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add an image file
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				return
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return err
		}

	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	res, err := http.Post(url, w.FormDataContentType(), req.Body)
	if err != nil {
		return
	}

	// Check the response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
		return
	}
	log.Println(res.StatusCode)
	return
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
}
