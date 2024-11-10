package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Font struct {
	Family string            `json:"family"`
	Files  map[string]string `json:"files"`
}

type FontResp struct {
	Items  []*Font `json:"items"`
	Errors *bool   `json:"errors"`
}

type FontRequest struct {
	Name    string `json:"name"`
	Weight  string `json:"weight"`
	Variant string `json:"variant"`
}

type FontServiceImpl struct {
	fonts []string
}

func (f *FontServiceImpl) GetFonts() []string {
	return f.fonts
}

func (f *FontServiceImpl) SetFonts(fonts []string) {

}

func (f *FontServiceImpl) Download(fontName string) error {
	fmt.Println("Downloading font from Google", fontName)
	authtoken := "IzaSyDowpwb2Kg_riAn97y7Rcg9LqmxFSEr1SI"

	if authtoken == "" {
		return fmt.Errorf("missing auth token")
	}

	formattedName := strings.Replace(fontName, " ", "+", -1)

	url := fmt.Sprintf("https://www.googleapis.com/webfonts/v1/webfonts?family=%s&key=%s", formattedName, authtoken)
	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	var payload FontResp

	json.Unmarshal(body, &payload)

	if payload.Errors != nil {
		return fmt.Errorf("error downloading font %s", fontName)
	}

	filesToDownload := payload.Items
	os.Mkdir(fmt.Sprintf("./fonts/%s", fontName), 0755)

	for _, v := range filesToDownload {
		for k, kv := range v.Files {
			url := kv
			resp, err := http.Get(url)

			if err != nil {
				return err
			}

			defer resp.Body.Close()

			body, e := io.ReadAll(resp.Body)

			if e != nil {
				return err
			}

			os.WriteFile(fmt.Sprintf("./fonts/%s/%s.ttf", v.Family, k), body, 0644)
		}
	}

	fmt.Println("Download complete")

	return nil
}

func (f *FontServiceImpl) GetFont(name, weight, variant string) ([]byte, error) {

	pattern := fmt.Sprintf("./fonts/%s", name)

	matches, err := filepath.Glob(pattern)

	if err != nil {
		return nil, err
	}

	if len(matches) > 1 {
		return nil, fmt.Errorf("found more than one file for name %s", name)
	}

	if len(matches) == 0 {
		err := f.Download(name)

		if err != nil {
			return nil, err
		}

	}

	file, err := os.ReadFile(fmt.Sprintf("./fonts/%s/%s.ttf", name, weight))

	if err != nil {
		return nil, err
	}

	return file, nil
}
