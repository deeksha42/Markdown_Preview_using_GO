package main

import (
	"fmt"
	"io"
	"mime/multipart" // Add this import
	"net/http"
	"sync"

	"github.com/russross/blackfriday/v2"
)

var previewsMutex sync.Mutex
var previews map[string]string

func renderMarkdown(filename, content string) string {
	html := blackfriday.Run([]byte(content))
	return fmt.Sprintf("<h1>%s</h1>%s", filename, html)
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Display the HTML form to upload Markdown files
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `
			<!DOCTYPE html>
			<html>
			<head>
				<title>Markdown Preview</title>
				<link href="https://stackpath.bootstrapcdn.com/bootstrap/4.5.2/css/bootstrap.min.css" rel="stylesheet">
			</head>
			<body>
				<div class="container">
					<h1 class="mt-5 mb-4">Upload Markdown Files</h1>
					<form id="uploadForm" action="/" method="post" enctype="multipart/form-data">
						<div class="custom-file">
							<input type="file" class="custom-file-input" id="files" name="files" multiple>
							<label class="custom-file-label" for="files">Choose file(s)</label>
						</div>
						<button type="submit" class="btn btn-primary mt-3">Upload</button>
					</form>
					<div id="previews" class="mt-5"></div>
				</div>
				<script src="https://code.jquery.com/jquery-3.5.1.min.js"></script>
				<script>
					// Update file input label with selected file names
					$('#files').on('change', function() {
						var filenames = Array.from(this.files).map(file => file.name);
						$('.custom-file-label').text(filenames.join(', '));
					});

					// Submit form via AJAX and display previews
					$('#uploadForm').submit(function(event) {
						event.preventDefault();
						var formData = new FormData(this);
						$.ajax({
							url: '/',
							type: 'POST',
							data: formData,
							cache: false,
							contentType: false,
							processData: false,
							success: function(data) {
								$('#previews').html(data);
							},
							error: function(xhr, status, error) {
								console.error(xhr.responseText);
							}
						});
					});
				</script>
			</body>
			</html>
		`)
		return
	}

	if r.Method == http.MethodPost {
		// Handle file upload
		r.ParseMultipartForm(10 << 20) // 10 MB limit
		files := r.MultipartForm.File["files"]

		var wg sync.WaitGroup
		wg.Add(len(files))

		errChan := make(chan error, len(files))

		for _, fileHeader := range files {
			go func(fileHeader *multipart.FileHeader) { // Corrected type
				defer wg.Done()

				file, err := fileHeader.Open()
				if err != nil {
					errChan <- fmt.Errorf("Error reading file: %s", err)
					return
				}
				defer file.Close()

				// Read file content
				content, err := io.ReadAll(file)
				if err != nil {
					errChan <- fmt.Errorf("Error reading file content: %s", err)
					return
				}

				// Render Markdown and add preview to map
				previewsMutex.Lock()
				previews[fileHeader.Filename] = renderMarkdown(fileHeader.Filename, string(content))
				previewsMutex.Unlock()

			}(fileHeader)
		}

		wg.Wait()
		close(errChan)

		// Check if any error occurred
		if len(errChan) > 0 {
			for err := range errChan {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// Display uploaded file previews
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "<h1 class='mt-5'>Uploaded Markdown Previews</h1>")
		for _, preview := range previews {
			fmt.Fprintf(w, "<div class='card mt-3'><div class='card-body'>%s</div></div>", preview)
		}
		return
	}

	// Method not allowed
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func main() {
	previews = make(map[string]string)

	http.HandleFunc("/", handleUpload)
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
