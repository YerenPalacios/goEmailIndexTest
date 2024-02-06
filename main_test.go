package main

import (
	"bytes"
	"fmt"
	"handlers"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-chi/chi/v5"
)

func getFilePayload(t *testing.T, name string) (*bytes.Buffer, *multipart.Writer) {
	file, err := os.Open(name)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	// Crear un nuevo buffer para almacenar el contenido del archivo
	buf := new(bytes.Buffer)
	// _, err = io.Copy(buf, file)
	if err != nil {
		t.Fatal(err)
	}

	// Crear un nuevo objeto de multipart para la solicitud POST
	mpWriter := multipart.NewWriter(buf)
	fileWriter, err := mpWriter.CreateFormFile("file", filepath.Base(name))
	if err != nil {
		t.Fatal(err)
	}
	file.Seek(0, 0)
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		t.Fatal(err)
	}
	mpWriter.Close()
	return buf, mpWriter
}

func TestImportFileHandler(t *testing.T) {
	r := chi.NewRouter()
	r.Post("/upload", handlers.HandleImportFile)
	ts := httptest.NewServer(r)
	defer ts.Close()

	buf, mpWriter := getFilePayload(t, "./t.tar.gz")

	req, err := http.NewRequest("POST", ts.URL+"/upload", buf)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", mpWriter.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var respBody bytes.Buffer
	_, err = io.Copy(&respBody, resp.Body)

	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("Expected status OK: got %v", resp.Status)
	}

	if respBody.String() != "{\"status\": \"Ok\"}" {
		t.Errorf("Expected message {\"status\": \"Ok\"}: got %v", respBody.String())
	}
}

func TestImportFileHandlerBadFile(t *testing.T) {
	r := chi.NewRouter()
	r.Post("/upload", handlers.HandleImportFile)
	ts := httptest.NewServer(r)
	defer ts.Close()

	buf, mpWriter := getFilePayload(t, "./blank.pdf")

	req, err := http.NewRequest("POST", ts.URL+"/upload", buf)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", mpWriter.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var respBody bytes.Buffer
	_, err = io.Copy(&respBody, resp.Body)

	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(respBody.String())
	if resp.StatusCode != 400 {
		t.Errorf("Expected status Bad Request: got %v", resp.Status)
	}
	if respBody.String() != "error opening file\n" {
		t.Errorf("Expected message 'error opening file': got %v", respBody.String())
	}
}
