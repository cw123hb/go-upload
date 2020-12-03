package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, r.URL)
	if r.Method == "GET" {
		current := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(current, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("upload.tpl")
		t.Execute(w, token)
	} else {
		r.ParseMultipartForm(32 << 10)
		file, handler, err := r.FormFile("uploadfile") // get file handle
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "Uploaded", handler.Header) // response
		f, err := os.OpenFile("./upload/"+handler.Filename,
			os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}

func main() {
	http.HandleFunc("/", upload)
	log.Fatal(http.ListenAndServe(":8899", nil))
}
