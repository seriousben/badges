package internal

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

const indexTmplString = `<html>
<head>
<title>Badges</title>
</head>

<body>
<table>
<tr><td>Badge<td><img src="/badge?label=label&description=description" />
<tr><td>Badge<td><img src="/badge?description=description" />
<tr><td>Badge<td><img src="/badge?label=label" />
<tr><td>Badge<td><img src="/badge?label=The quick brown fox jumps over the lazy dog&description=The quick brown fox jumps over the lazy dog" />
<tr><td>Progress<td><img src="/progress?success_rate=0.75" />
</table>
</body>
</html>`

var indexTmpl = template.Must(template.New("index").Parse(indexTmplString))

const progressSvgTmplString = `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="{{.width}}" height="20" role="img"
aria-label="{{.title}}">
<title>{{.title}}</title>
<clipPath id="r">
	<rect width="{{.width}}" height="20" rx="3" fill="#fff" />
</clipPath>
<g clip-path="url(#r)">
	<rect width="{{ .successWidth }}" height="20" fill="#4c1" />
	<rect x="{{ .successWidth }}" width="{{ .failWidth }}" height="20" fill="#9f9f9f" />
	<rect width="{{.width}}" height="20" fill="url(#s)" />
</g>
</svg>`

var progressSvgTmpl = template.Must(template.New("progress").Parse(progressSvgTmplString))

func Serve(port int) error {
	server := &badgesServer{}
	return http.ListenAndServe(fmt.Sprintf(":%d", port), server)
}

type badgesServer struct{}

func (s *badgesServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	mux := http.NewServeMux()
	mux.HandleFunc("/badge", s.serveBadge)
	mux.HandleFunc("/progress", s.serveProgress)
	mux.HandleFunc("/", s.serveRoot)
	mux.ServeHTTP(w, req)
}

func (s *badgesServer) serveRoot(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	// TODO: improve cachig with etag
	w.Header().Add("cache-control", "no-cache")
	w.Header().Add("content-type", "text/html;charset=utf-8")
	err := indexTmpl.Execute(w, map[string]string{
		"title": q.Get("title"),
	})
	if err != nil {
		log.Println("[ERROR] ", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *badgesServer) serveBadge(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	// 30 days
	w.Header().Add("cache-control", "max-age=44640, immutable")
	w.Header().Add("content-security-policy", "default-src 'none'; img-src data:")
	w.Header().Add("content-type", "image/svg+xml;charset=utf-8")
	err := renderBadge(w, q.Get("label"), q.Get("description"))
	if err != nil {
		log.Println("[ERROR] ", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *badgesServer) serveProgress(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	successRate, err := strconv.ParseFloat(q.Get("success_rate"), 64)
	if err != nil || successRate < 0 || successRate > 1 {
		successRate = 0
	}
	w.Header().Add("cache-control", "max-age=44640, immutable")
	w.Header().Add("content-security-policy", "default-src 'none'; img-src data:")
	w.Header().Add("content-type", "image/svg+xml;charset=utf-8")
	err = progressSvgTmpl.Execute(w, map[string]interface{}{
		"title":        q.Get("title"),
		"width":        "100",
		"successWidth": successRate * 100,
		"failWidth":    100 - (successRate * 100),
		"successRate":  successRate,
	})
	if err != nil {
		log.Println("[ERROR] ", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
