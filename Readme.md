# Markdown Preview Project

## Overview

This project is a simple web application built in Go that allows users to upload multiple Markdown files and preview their content rendered as HTML. The application uses the [Blackfriday](https://github.com/russross/blackfriday) library for Markdown parsing.

## Features

- Upload multiple Markdown files at once.
- Preview the rendered HTML content of each file.
- AJAX-based file upload for a seamless user experience.
- Responsive design using Bootstrap for an enhanced UI.

## Technologies Used

- Go (Golang)
- HTML/CSS
- JavaScript (jQuery)
- Blackfriday library for Markdown parsing

## Getting Started

### Prerequisites

- Go (version 1.15 or higher)
- Internet connection for Bootstrap and jQuery CDN links

### Installation

1. Clone the repository:

   ```bash
   git clone <repository-url>
   ```

2. Navigate to the project directory:

   ```bash
   cd <project-directory>
   ```

3. Run the application:

   ```bash
   go run project.go
   ```

4. Open your web browser and go to `http://localhost:8080`.

### Usage

1. Click on "Choose file(s)" to select one or more Markdown files from your computer.
2. Click "Upload" to submit the files.
3. The rendered previews of your Markdown files will be displayed below the upload form.
