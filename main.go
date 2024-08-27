package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"

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

// parseContent переводит входные данные input в формате Markdown в HTML.
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

// preview открывает файл filename стандартными средствами системы.
func preview(filename string) error {
	// В переменных будут установлены команда и параметры для запуска
	// стандартного средства просмотра целевой ОС.
	cmdName := ""
	cmdParams := []string{}

	switch runtime.GOOS {
	case "darwin":
		cmdName = "open"
	case "linux":
		cmdName = "xdg-open"
	case "windows":
		cmdName = "cmd.exe"
		cmdParams = []string{"/C", "start"}
	default:
		return fmt.Errorf("OS not supported.")
	}

	cmdParams = append(cmdParams, filename)

	// Найти исполняемый файл в PATH
	cmdPath, err := exec.LookPath(cmdName)
	if err != nil {
		return err
	}

	// Открыть файл при помощи выбранной программы
	return exec.Command(cmdPath, cmdParams...).Run()
}

// run выполняет основной функционал программы.
// filename - Markdown-файл, который нужно просмотреть.
// out - Writer, куда будет выведено имя временного файла с HTML,
// skipPreview позволяет пропустить запуск системной программы просмотра
// файла.
func run(filename string, out io.Writer, skipPreview bool) error {
	input, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	html := parseContent(input)

	// Создадим временный файл для сохранения HTML
	tempFile, err := os.CreateTemp("", "mdp*.html")
	if err != nil {
		return err
	}
	if err := tempFile.Close(); err != nil {
		return err
	}

	outputFileName := tempFile.Name()
	fmt.Fprintln(out, outputFileName)

	err = saveHTML(outputFileName, html)
	if err != nil {
		return err
	}

	if skipPreview {
		return nil
	}

	return preview(outputFileName)
}

func main() {
	filenameFlag := flag.String("file", "", "Файл в формате Markdown для просмотра")
	skipPreviewFlag := flag.Bool("s", false, "Не запускать программу просмотра файла")
	flag.Parse()

	if *filenameFlag == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*filenameFlag, os.Stdout, *skipPreviewFlag); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
