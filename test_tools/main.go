package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	_ "net/http/cookiejar"
	"os"
	"strings"
	"time"
)

type resp struct {
	Success    string `json:"success"`
	ReplyID    string `json:"replyID"`
	ItemString string `json:"itemString"`
	Item       *item  `json:"item"`
}

type item struct {
	Media        string    `json:"Media"`
	Title        string    `json:"Title"`
	Email        string    `json:"Email"`
	MyText       string    `json:"MyText"`
	ID           string    `json:"ID"`
	TS           time.Time `json:"TS"`
	MediaType    string    `json:"mediaType"`
	TempFileName string    `json:"tempFileName"`
}

var token string = "TOKEN_GOES_HERE_IF_NEEDED"
var port string

func main() {
	b, err := os.ReadFile("testApp/bolt.conf.json")
	if err != nil {
		log.Println(err)
	}
	port = strings.Split(strings.Split(string(b), "port\":\"")[1], "\"")[0]
	sendPost()
	fmt.Println("\n > > > >    http://localhost:" + port)
	fmt.Println()
}

func sendPost() {
	if len(os.Args) < 3 {
		os.Args = append(os.Args, []string{"a thingy", "a picture of a thing", "fake@email.x"}...)
	}
	var client *http.Client = &http.Client{}
	var remoteURL string = "http://localhost:" + port + "/uploadItem"
	images, err := os.ReadDir("img")
	if err != nil {
		log.Println(err)
	}
	fmt.Print("\n\n\n\n\n")
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
		var r resp = resp{}
		by, err := io.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
		}
		err = json.Unmarshal(by, &r)
		if err != nil {
			log.Println("test", err)
		}
		err = json.Unmarshal([]byte(r.ItemString), &r.Item)
		if err != nil {
			log.Println("test", err)
		}
		fmt.Println("Success:", r.Success, r.ReplyID)
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
