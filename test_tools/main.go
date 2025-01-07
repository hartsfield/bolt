package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	_ "net/http/cookiejar"
	"os"
	"strings"
)

var token string = "TOKEN_GOES_HERE_IF_NEEDED"

func main() {
	log.Println(len(os.Args))
	if len(os.Args) > 1 {
		if len(os.Args) < 4 {
			fmt.Println("Requires 1 - 3 arguments: ")
			fmt.Println("  btst remote_url title message")
			fmt.Println("using generic test data...")
			os.Args = append(os.Args, []string{"a thingy", "a picture of a thing", "fake@email.x"}...)
		}
	} else {
		fmt.Println("Usage:")
		fmt.Println("  btst remote_url title_txt message_txt")
		fmt.Println("Exiting.")
		os.Exit(0)
	}
	var client *http.Client = &http.Client{}
	var remoteURL string = os.Args[1]
	images, err := os.ReadDir("img")
	if err != nil {
		log.Println(err)
	}
	for _, img := range images {
		values := map[string]io.Reader{
			"Title":     strings.NewReader(os.Args[2]),
			"MyText":    strings.NewReader(os.Args[3]),
			"Email":     strings.NewReader(os.Args[4]),
			"Media":     mustOpen("img/" + img.Name()),
			"mediaType": strings.NewReader(strings.Split(img.Name(), ".")[1]),
		}

		b, w := writeMultipart(values)
		res, err := client.Do(mkReq(remoteURL, b, w))
		if err != nil {
			log.Println("Request Error:", err)
		}
		log.Println(res.Status)
	}
}

func mkReq(remoteURL string, b *bytes.Buffer, w *multipart.Writer) *http.Request {
	req, err := http.NewRequest("POST", remoteURL, b)
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.AddCookie(&http.Cookie{Name: "token", Value: token})
	return req
}

func writeMultipart(values map[string]io.Reader) (*bytes.Buffer, *multipart.Writer) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	if func() (err error) {
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
				return
			}
		}
		return nil
	}() != nil {
		log.Panicln("form error")
	}
	w.Close()
	return &b, w
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
}
