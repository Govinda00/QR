package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/skip2/go-qrcode"
)

var generatedQRs []string

func servePage(w http.ResponseWriter, r *http.Request) {
	html := `
	<html>
		<head>
			<title>QR Code Generator</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					text-align: center;
					padding: 20px;
				}
				input[type="text"] {
					padding: 10px;
					width: 300px;
					margin: 10px 0;
				}
				button {
					padding: 10px 20px;
					font-size: 16px;
				}
				img {
					margin-top: 20px;
					max-width: 300px;
					height: auto;
				}
				ul {
					list-style-type: none;
				}
				li {
					margin: 10px;
				}
				a {
					text-decoration: none;
				}
			</style>
		</head>
		<body>
			<h1>QR Code Generator</h1>
			<form method="POST" action="/generate">
				<input type="text" name="url" placeholder="Enter URL" required />
				<button type="submit">Generate QR Code</button>
			</form>
			<h2>Generated QR Codes:</h2>
			<ul>
			{{range .}}
				<li>
					<img src="/qr/{{.}}" alt="QR Code" />
					<br/>
					<a href="/qr/{{.}}" download>Download QR Code ({{.}})</a>
				</li>
			{{end}}
			</ul>
		</body>
	</html>
	`

	t, err := template.New("qr").Parse(html)
	if err != nil {
		log.Fatal(err)
	}

	err = t.Execute(w, generatedQRs)
	if err != nil {
		log.Fatal(err)
	}
}

func generateQRCode(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	url := r.FormValue("url")

	if url == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	fileName := fmt.Sprintf("%d_qr.png", time.Now().Unix())
	err := qrcode.WriteFile(url, qrcode.Medium, 256, fileName)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Error generating QR code", http.StatusInternalServerError)
		return
	}

	generatedQRs = append(generatedQRs, fileName)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	http.HandleFunc("/", servePage)
	http.HandleFunc("/generate", generateQRCode)
	http.Handle("/qr/", http.StripPrefix("/qr/", http.FileServer(http.Dir("."))))
	fmt.Println("Starting server on http://localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
