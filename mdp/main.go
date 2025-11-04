package main

import (
	"bytes"
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

const (
	defaultTemplate = `<!DOCTYPE html>
	<html>
		<head>
			<meta http-equiv="content-type" content="text/html; charset=utf-8">
			<title>{{.Title}}</title>
		</head>
		<body>
		{{.Body}}
	   </body>
	</html>
	`
)

type content struct {
	Title string
	Body  template.HTML
}

func main() {
	fileName := flag.String("file", "", "Markdown file to preview")
	skipPreview := flag.Bool("s", false, "Skip auto-preview")
	templateFileName := flag.String("t", "", "Alternate template name")
	flag.Parse()

	if *fileName == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*fileName, *templateFileName, os.Stdout, *skipPreview); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(fileName, templateFileName string, out io.Writer, skipPreview bool) error {
	input, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	htmlData, err := parseContent(input, templateFileName)
	if err != nil {
		return err
	}

	temp, err := os.CreateTemp("", "mdp*.html")
	if err != nil {
		return err
	}

	if err := temp.Close(); err != nil {
		return err
	}

	outName := temp.Name()
	fmt.Fprintln(out, outName)

	if err := saveHTML(outName, htmlData); err != nil {
		return err
	}

	if skipPreview {
		return nil
	}

	defer os.Remove(outName)

	return preview(outName)
}

func parseContent(input []byte, templateFileName string) ([]byte, error) {
	output := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	t, err := template.New("mdp").Parse(defaultTemplate)
	if err != nil {
		return nil, err
	}

	if templateFileName != "" {
		t, err = template.ParseFiles(templateFileName)
		if err != nil {
			return nil, err
		}
	}

	c := content{
		Title: "Markdown Preview Tool",
		Body:  template.HTML(body),
	}

	var buffer bytes.Buffer

	if err := t.Execute(&buffer, c); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func saveHTML(outputFileName string, data []byte) error {
	return os.WriteFile(outputFileName, data, 0644)
}

func preview(fileName string) error {
	commandName := ""
	commandParams := []string{}

	switch runtime.GOOS {
	case "linux":
		commandName = "xdg-open"
	case "windows":
		commandName = "cmd.exe"
		commandParams = []string{"/C", "start"}
	case "darwin":
		commandName = "open"
	default:
		return fmt.Errorf("OS not supported")
	}

	commandParams = append(commandParams, fileName)

	commandPath, err := exec.LookPath(commandName)
	if err != nil {
		return err
	}

	err = exec.Command(commandPath, commandParams...).Run()

	time.Sleep(2 * time.Second)
	return err
}
