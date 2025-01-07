package main

import (
	"bytes"
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
	b, err := os.ReadFile("testApp/bolt.conf.json")
	if err != nil {
		log.Println(err)
	}
	p := strings.Split(strings.Split(string(b), "port\":\"")[1], "\"")[0]
	log.Println(len(os.Args))
	if len(os.Args) < 3 {
		os.Args = append(os.Args, []string{"a thingy", "a picture of a thing", "fake@email.x"}...)
	}
	var client *http.Client = &http.Client{}
	var remoteURL string = "http://localhost:" + p + "/uploadItem"
	images, err := os.ReadDir("img")
	if err != nil {
		log.Println(err)
	}
	for _, img := range images {
		values := map[string]io.Reader{
			"Title":     strings.NewReader(os.Args[1]),
			"MyText":    strings.NewReader(os.Args[2]),
			"Email":     strings.NewReader(os.Args[3]),
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
