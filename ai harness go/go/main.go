package main

import (
    "fmt"
    "html/template"
    "net/http"
)

type Application struct {
    Templates *template.Template
    Results struct {
        Accuracy         float64 `json:"accuracy"`
        Precision        float64 `json:"precision"`
        Recall           float64 `json:"recall"`
        F1Score          float64 `json:"f1_score"`
        HeatmapImagePath string  `json:"heatmap_image_path"`
        AccuracyImagePath string `json:"accuracy_image_path"`
    }
}


func main() {
    app := Application{
        Templates: template.Must(template.ParseFiles("templates/index.html", "templates/results.html")),
    }
    // fs := http.FileServer(http.Dir("../outputs")) 
    // http.Handle("/static/", http.StripPrefix("/static/", fs))
    

    fmt.Println("Server is running on http://localhost:8080")
    http.ListenAndServe(":8080", app.routes())
}
