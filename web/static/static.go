package static

import (
	"bytes"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	hContentLength = "Content-Length"
)

type Asset interface {
	Name() string
	ModTime() time.Time
	Data() []byte
}

type AssetStore interface {
	AssetFile(name string) Asset
}

func Handler(assets AssetStore) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !(r.Method == http.MethodGet || r.Method == http.MethodHead) {
			sendMessage(w, http.StatusMethodNotAllowed)
			return
		}

		pageURL := r.URL.Path

		if len(pageURL) == 0 {
			http.Redirect(w, r, "/", http.StatusMovedPermanently) // 301 not 308
			return
		}

		if strings.HasSuffix(pageURL, "/index.html") {
			pageURL = strings.TrimSuffix(pageURL, "index.html")
			http.Redirect(w, r, pageURL, http.StatusMovedPermanently) // 301 not 308
			return
		}

		if strings.HasSuffix(pageURL, "/") {
			pageURL = pageURL + "index.html"
		}

		f := assets.AssetFile(pageURL)
		if f == nil {
			// return index html5 route
			f = assets.AssetFile("/index.html")
		}

		// w.Header().Set(hEtag, `"` + f.Hash() + `"`)
		http.ServeContent(w, r, f.Name(), f.ModTime(), bytes.NewReader(f.Data()))

	})

}

func sendMessage(w http.ResponseWriter, statusCode int) {
	msg := []byte(http.StatusText(statusCode) + "\n")
	w.Header().Set(hContentLength, strconv.Itoa(len(msg)))
	w.WriteHeader(statusCode)
	w.Write(msg)
}
