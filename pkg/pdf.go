package pkg

import (
	"os"
	"os/exec"
)

func ConvertHTMLToPDF(outputFile string, data []byte) error {
	inputFile := "/tmp/document.html"
	if err := os.WriteFile(inputFile, data, 0644); err != nil {
		return err
	}

	return exec.Command("wkhtmltopdf", "--title", "Transaction Report",
		"--zoom", "0.98", "--page-size", "A4", "--orientation", "Portrait", inputFile, outputFile).Run()
}
