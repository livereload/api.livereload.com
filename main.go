package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "strings"
)

func importLicenses(w http.ResponseWriter, r *http.Request) {
    components := strings.Split(r.URL.Path, "/")

    fmt.Fprintf(w, "Import %+v!", components)
}

func index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "3000"
    }

    http.HandleFunc("/licensing/admin/import/", importLicenses)
    http.HandleFunc("/", index)
    log.Printf("Listening on port %s.", port)
    http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
