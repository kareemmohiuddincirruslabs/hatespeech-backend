package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// func (app *Application) uploadPageHandler(w http.ResponseWriter, r *http.Request) {
// 	if err := app.Templates.ExecuteTemplate(w, "index.html", nil); err != nil {
// 		app.handleError(w, "Failed to render the upload page", err, http.StatusInternalServerError)
// 	}
// }

func (app *Application) handleUpload(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseMultipartForm(10 << 20); err != nil {
        app.handleError(w, "Failed to parse multipart form", err, http.StatusInternalServerError)
        return
    }

    modelname := r.FormValue("modelname")

    file, header, err := r.FormFile("datafile")
    if err != nil {
        app.handleError(w, "Invalid file", err, http.StatusBadRequest)
        return
    }
    defer file.Close()

    fmt.Printf("Received model name: %s\n", modelname)
    fmt.Printf("Received file: %s\n", header.Filename)

    tempFile, err := os.CreateTemp("temp", "upload-*.json")
    if err != nil {
        app.handleError(w, "Error creating temporary file", err, http.StatusInternalServerError)
        return
    }
    defer tempFile.Close()

    _, err = io.Copy(tempFile, file)
    if err != nil {
        app.handleError(w, "Error copying file", err, http.StatusInternalServerError)
        return
    }
    tempFileName := tempFile.Name()

	mainPyPath := filepath.Join("..", "python", "main.py")
	analysisPyPath := filepath.Join("..", "python", "analysis.py")

	cmd := exec.Command("python", mainPyPath, modelname, tempFileName)
	cmdOutput, err := cmd.CombinedOutput()
	outputLines := strings.Split(strings.TrimSpace(string(cmdOutput)), "\n")
	fmt.Println(outputLines)
	if err != nil {
		app.handleError(w, "Error processing file with main.py", err, http.StatusInternalServerError)
		return
	}
	outputFilePath := outputLines[len(outputLines)-1]
	fmt.Println("Main.py ran successfully")

	cmd = exec.Command("python", analysisPyPath, outputFilePath)
	cmdOutput, err = cmd.CombinedOutput()
	fmt.Printf("Python script output: %s\n", string(cmdOutput))
	if err != nil {
		app.handleError(w, "Error processing file with analysis.py", err, http.StatusInternalServerError)
		return
	}
	fmt.Println("analysis.py ran successfully")

	outputFilePath = strings.TrimSuffix(outputFilePath, ".json") + "_results.json"
	heatmapImagePath := strings.Replace(outputFilePath, "_results.json", "_heatmap.png", 1)
	err = app.parseAnalysisResults(outputFilePath, heatmapImagePath)
	if err != nil {
		app.handleError(w, "Error parsing analysis results", err, http.StatusInternalServerError)
		return
	}
	

    w.Header().Set("Content-Type", "application/json")
    if err = json.NewEncoder(w).Encode(app.Results); err != nil {
        app.handleError(w, "Failed to encode results as JSON", err, http.StatusInternalServerError)
    }
	fmt.Println("json file created successfully")
	
	// if err := app.Templates.ExecuteTemplate(w, "results.html", app); err != nil {
	// 	app.handleError(w, "Failed to render results page", err, http.StatusInternalServerError)
	// }
	fmt.Println(app.Results.AccuracyImagePath)
	fmt.Println(app.Results.HeatmapImagePath)
	
}


func (app *Application) parseAnalysisResults(filePath string, heatmapImagePath string) error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, &app.Results)
	if err != nil {
		return err
	}

	app.Results.HeatmapImagePath = "http://localhost:8080/static/" + filepath.Base(heatmapImagePath)
	app.Results.AccuracyImagePath = "http://localhost:8080/static/accuracies_graph.png"

	return nil
}

