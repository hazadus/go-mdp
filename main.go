package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	header = `<!DOCTYPE html
	<html>
		<head>
			<meta http-equiv="content-type" content="text/html; charset=utf-8">
			<title>Markdown Preview CLI Tool</title>
		</head>
		<body>
	`
	footer = `
		</body>
	</html>
	`
)

// parseContent переводит входные данные в формате Markdown в HTML.
func parseContent(input []byte) []byte {
	output := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	var buffer bytes.Buffer

	buffer.WriteString(header)
	buffer.Write(body)
	buffer.WriteString(footer)

	return buffer.Bytes()
}

// saveHTML сохраняет data в файл с именем outputFileName.
func saveHTML(outputFileName string, data []byte) error {
	return os.WriteFile(outputFileName, data, 0644)
}

// run выполняет основной функционал программы.
func run(filename string) error {
	input, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	html := parseContent(input)

	outputFileName := fmt.Sprintf("%s.html", filepath.Base(filename))
	fmt.Println(outputFileName)

	return saveHTML(outputFileName, html)
}

func main() {
	filenameFlag := flag.String("file", "", "Файл в формате Markdown для просмотра")
	flag.Parse()

	if *filenameFlag == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*filenameFlag); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
