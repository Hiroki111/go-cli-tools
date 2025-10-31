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
	header = `<!DOCTYPE html>
	<html>
		<head>
			<meta http-equiv="content-type" content="text/html; charset=utf-8">
			<title>Markdown Preview Tool</title>
		</head>
		<body>`
	footer = `
	   </body>
	</html>
	`
)

func main() {
	fileName := flag.String("file", "", "Markdown file to preview")
	skipPreview := flag.Bool("s", false, "Skip auto-preview")
	flag.Parse()

	if *fileName == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*fileName, os.Stdout, *skipPreview); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(fileName string, out io.Writer, skipPreview bool) error {
	input, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	htmlData := parseContent(input)

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

	return preview(outName)
}

func parseContent(input []byte) []byte {
	output := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	var buffer bytes.Buffer
	buffer.WriteString(header)
	buffer.Write(body)
	buffer.WriteString(footer)

	return buffer.Bytes()
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

	return exec.Command(commandPath, commandParams...).Run()
}
