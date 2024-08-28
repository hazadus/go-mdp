package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

//go:embed assets/template.html
var defaultTemplate string

// content представляет контент для вставки в шаблон
type content struct {
	Title string
	Body  template.HTML
}

// parseContent переводит входные данные input в формате Markdown в HTML.
// Если указано имя файла templFileName, шаблон будет загружен из него.
// Иначе, используется дефолтный шаблон.
func parseContent(input []byte, templFileName string) ([]byte, error) {
	output := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	// Загружаем дефолтный шаблон
	templ, err := template.New("mdp").Parse(defaultTemplate)
	if err != nil {
		return nil, err
	}

	// Если указан файл с шаблоном, загрузим его
	if templFileName != "" {
		templ, err = template.ParseFiles(templFileName)
		if err != nil {
			return nil, err
		}
	}

	c := content{
		Title: "Markdown Preview CLI Tool",
		Body:  template.HTML(body),
	}

	var buffer bytes.Buffer

	// Наполнить шаблон данными
	if err := templ.Execute(&buffer, c); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
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
	err = exec.Command(cmdPath, cmdParams...).Run()

	// Дать время внешней программе для открытия файла,
	// прежде чем удалить его
	// TODO: подчищать за собой при помощи сигналов, а не задержки.
	time.Sleep(2 * time.Second)
	return err
}

// run выполняет основной функционал программы.
// filename - Markdown-файл, который нужно просмотреть.
// out - Writer, куда будет выведено имя временного файла с HTML,
// skipPreview позволяет пропустить запуск системной программы просмотра
// файла.
// templFileName - "" или имя файла шаблона для выходного HTML-файла.
func run(filename string, out io.Writer, skipPreview bool, templFileName string) error {
	input, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	html, err := parseContent(input, templFileName)
	if err != nil {
		return err
	}

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

	defer os.Remove(outputFileName)
	return preview(outputFileName)
}

func main() {
	filenameFlag := flag.String("file", "", "Файл в формате Markdown для просмотра")
	skipPreviewFlag := flag.Bool("s", false, "Не запускать программу просмотра файла")
	templFileName := flag.String("t", "", "Шаблон для формирования HTML-файла")
	flag.Parse()

	if *filenameFlag == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*filenameFlag, os.Stdout, *skipPreviewFlag, *templFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
