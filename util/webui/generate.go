package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/benbjohnson/genesis"
)

func main() {
	if err := generate(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func generate() error {

	cwd := "public"
	out := "assets/public.go"
	pkg := "assets"
	tags := ""
	args := []string{"."}

	if fileExists(out) && missing(cwd) {
		return nil // skip because we are running in production
	}

	// Determine output writer.
	if err := os.Mkdir(pkg, 0755); err != nil {
		// ignore
	}

	f, errCreate := os.OpenFile(out, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if errCreate != nil {
		return errCreate
	}
	defer f.Close()

	// Change working directory, if specified.
	if err := os.Chdir(cwd); err != nil {
		return err
	}

	// Find all matching files.
	var paths []string
	for _, arg := range args {
		a, err := expand(arg)
		if err != nil {
			return err
		}
		paths = append(paths, a...)
	}

	enc := genesis.NewEncoder(f)
	enc.Package = pkg
	enc.Tags = strings.Split(tags, ",")

	// Encode all assets.
	for _, path := range paths {
		// Fetch mod time from stats.
		fi, errStat := os.Stat(path)
		if errStat != nil {
			return errStat
		}

		// Read entire file into memory.
		data, errRead := ioutil.ReadFile(path)
		if errRead != nil {
			return errRead
		}

		// Encode asset to writer.
		if err := enc.Encode(&genesis.Asset{
			Name:    prependSlash(filepath.ToSlash(path)),
			Data:    data,
			ModTime: fi.ModTime(),
		}); err != nil {
			return err
		}
	}

	// Close out encoder.
	if err := enc.Close(); err != nil {
		return err
	}

	return nil

}

// expand converts path into a list of all files within path.
// If path is a file then path is returned.
func expand(path string) ([]string, error) {
	if fi, err := os.Stat(path); err != nil {
		return nil, err
	} else if !fi.IsDir() {
		return []string{path}, nil
	}

	// Read files in directory.
	fis, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	// Iterate over files and expand.
	expanded := make([]string, 0, len(fis))
	for _, fi := range fis {
		a, err := expand(filepath.Join(path, fi.Name()))
		if err != nil {
			return nil, err
		}
		expanded = append(expanded, a...)
	}
	return expanded, nil
}

func prependSlash(s string) string {
	if strings.HasPrefix(s, "/") {
		return s
	}
	return "/" + s
}

func fileExists(filename string) bool {
	if s, err := os.Stat(filename); err == nil && s != nil && !s.IsDir() {
		return true
	}
	return false
}

func missing(dir string) bool {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return true
	}
	return false
}
