package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
)

var templates *template.Template

func main() {
	var err error
	templates, err = template.ParseFiles(
		"templates/index.shtml",
		"templates/header.shtml",
		"templates/menu.shtml",
		"templates/footer.shtml",
		"templates/contact.shtml",
		"templates/dev-info.shtml",
		"templates/documentation.shtml",
		"templates/download.shtml",
		"templates/faq.shtml",
		"templates/features.shtml",
		"templates/related.shtml",
	)
	if err != nil {
		log.Fatalln("template.ParseFiles:", err)
	}

	// set up tracing
	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: os.Getenv("GOOGLE_CLOUD_PROJECT"),
	})
	if err != nil {
		log.Fatalln(err)
	}
	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{
		DefaultSampler: trace.AlwaysSample(),
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	handlers := []struct {
		Prefix, Template, Title string
	}{
		{"/", "index.shtml", "Start page"},
		{"/index.shtml", "index.shtml", "Start page"},
		{"/contact.shtml", "contact.shtml", "Contact"},
		{"/dev-info.shtml", "dev-info.shtml", "Development Information"},
		{"/documentation.shtml", "documentation.shtml", "Documentation"},
		{"/download.shtml", "download.shtml", "Download"},
		{"/faq.shtml", "faq.shtml", "Frequently asked questions"},
		{"/features.shtml", "features.shtml", "Features"},
		{"/related.shtml", "related.shtml", "Related sites"},
	}
	for _, hndl := range handlers {
		http.Handle(hndl.Prefix, &ochttp.Handler{
			Propagation:      &propagation.HTTPFormat{},
			Handler:          templateHandler{hndl.Template, hndl.Title},
			IsPublicEndpoint: true,
		})
	}

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalln("http.ListenAndServe:", err)
	}
}

type templateHandler struct {
	Name  string
	Title string
}

func (h templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := templates.ExecuteTemplate(w, h.Name, h); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
