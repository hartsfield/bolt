package main

import (
        "log"
        "net/http"
        "time"
        "io"
        "bytes"
        "encoding/json"
        "os"
        "strings"
)

type item struct { 
        {{- range $k, $v :=  .Inputs }}
        {{- range $k_, $v_ :=  $v }}
        {{ $v_}}{{"\t"}}string{{"\t"}}`json:"{{$v_}}"`
        {{- end }}
        {{- end }}                             
        ID{{"\t"}}string `json:"ID"`
        TS{{"\t"}}time.Time `json:"TS"`
}
var stream []*item
func uploadHandler(w http.ResponseWriter, r *http.Request) {
        mr, err := r.MultipartReader()
        if err != nil {
                log.Println(err)
        }

        var data *item = &item{ID: genPostID(10)}
        for {
                part, err_part := mr.NextPart()
                if err_part == io.EOF {
                        break
                }
                {{- range $k, $v :=  .Inputs }}
                {{- if ne $k "file" }}
                {{- range $k_, $v_ :=  $v }}
                if part.FormName() == "{{$v_}}" {
                        buf := new(bytes.Buffer)
                        buf.ReadFrom(part)
                        data.{{$v_}} = buf.String()
                }
                {{- end }}
                {{- end }}
                {{- end }}      
                //if part.FormName() == "{{/* {{index .Inputs "file" 0}} */}}" {}
        }
        stream = append(stream, data)

        b, err := json.Marshal(data)
        if err != nil {
                log.Println(err)
        }
        ajaxResponse(w, map[string]string{
                "success":   "true",
                "replyID":   data.ID,
                "item":    string(b),
        })
        saveJSON(b)
}
func init() {
        readDB()
}
func readDB() {
        content, err := os.ReadFile("JSON_DB.json")
        if err != nil {
                log.Println(err)
        }

        if len(content) > 0 {
                lines := strings.Split(string(content), "\n")
                for _, l := range  lines {
                        var item item
                        err := json.Unmarshal([]byte(l), &item)
                        if err != nil {
                                log.Println(err)
                        }

                        stream = append(stream, &item)
                }
        }
}

func saveJSON(json_b []byte) {
        f, err := os.OpenFile("JSON_DB.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
        if err != nil {
                log.Println(err)
        }

        defer f.Close()

        if _, err = f.WriteString(string(json_b) + "\n"); err != nil {
                log.Println(err)
        }
}
