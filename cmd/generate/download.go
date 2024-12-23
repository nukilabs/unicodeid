package main

import (
	"archive/zip"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const (
	latestURL  = "https://www.unicode.org/Public/zipped/latest/UCD.zip"
	unicodeZip = "UCD.zip"
	unicodeDir = "UCD"
)

func download() error {
	res, err := http.Get(latestURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	f, err := os.Create(unicodeZip)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, res.Body)
	return err
}

func unzip() error {
	r, err := zip.OpenReader(unicodeZip)
	if err != nil {
		return err
	}
	defer r.Close()

	os.Mkdir(unicodeDir, 0755)

	extract := func(f *zip.File) error {
		path := filepath.Join(unicodeDir, f.Name)
		if f.FileInfo().IsDir() {
			return os.MkdirAll(path, f.Mode())
		} else {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			w, err := os.Create(path)
			if err != nil {
				return err
			}
			defer w.Close()

			_, err = io.Copy(w, rc)
			return err
		}
	}

	for _, f := range r.File {
		if err := extract(f); err != nil {
			return err
		}
	}

	return nil
}
