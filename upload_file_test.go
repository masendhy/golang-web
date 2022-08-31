package golangweb

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// create function for form upload template
func UploadForm(writer http.ResponseWriter, request *http.Request) {
	myTemplates.ExecuteTemplate(writer, "upload.form.gohtml", nil)
}

func Upload(w http.ResponseWriter, r *http.Request) {

	//untuk upload file yang entry
	// request.ParseMultipartForm(32 << 20) untuk menentukan besar ukuran file yang dikirimkan
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}

	//untuk upload file image yang dikirim
	fileDestination, err := os.Create("./resources/" + fileHeader.Filename)
	if err != nil {
		panic(err)
	}

	// me redirect file image yang dikirimkan
	_, err = io.Copy(fileDestination, file)
	if err != nil {
		panic(err)
	}

	//membuat template untuk menampung file yang diupload
	name := r.PostFormValue("name")
	myTemplates.ExecuteTemplate(w, "upload.success.gohtml", map[string]interface{}{
		"Name": name,
		"File": "/static/" + fileHeader.Filename,
	})
}

func TestUploadForm(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", UploadForm)
	mux.HandleFunc("/upload", Upload)
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./resources"))))

	server := http.Server{
		Addr:    "localhost:8888",
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

// go:embed resources/lemari.jpeg
var uploadfiletest []byte

func TestUploadFile(t *testing.T) {
	//i. membuat body untuk dikirim kedalam request
	body := new(bytes.Buffer)

	//ii.cara untuk memasukkan kedalam body nya
	//   a.buat writer dengan multipart file
	writer := multipart.NewWriter(body)

	//   b.untuk menangkap data berupa entry file gunakan WriteField
	writer.WriteField("name", "masendhy")

	//   c.untuk menangkap data berupa form gunakan CreateFormFile
	file, _ := writer.CreateFormFile("file", "hasilupload.png")
	file.Write(uploadfiletest)
	writer.Close()

	request := httptest.NewRequest(http.MethodPost, "http://localhost:8888 ", body)
	request.Header.Set("Content-Type", writer.FormDataContentType())

	recorder :=
		httptest.NewRecorder()

	Upload(recorder, request)
	bodyfile, _ := io.ReadAll(recorder.Result().Body)
	fmt.Println(string(bodyfile))

}
