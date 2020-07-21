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
	templates, err = template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatalln("template.ParseFiles:", err)
	}

	// set up tracing
	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: os.Getenv("GOOGLE_CLOUD_PROJECT"),
	})
	if err != nil {
		log.Fatal(err)
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

	http.Handle("/", &ochttp.Handler{
		Propagation:      &propagation.HTTPFormat{},
		Handler:          http.HandlerFunc(handler),
		IsPublicEndpoint: true,
	})
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalln("http.ListenAndServe:", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	if err := templates.Execute(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
