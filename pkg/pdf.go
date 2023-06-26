package pkg

import (
	"github.com/jung-kurt/gofpdf"
)

func ConvertToPDF(text string) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 16)

	// Add the text to the PDF document
	pdf.MultiCell(190, 10, string(text), "", "", false)

	filename := "/Users/masud/Documents/output.pdf"
	return filename, pdf.OutputFileAndClose(filename)
}

//
//func ConvertToPDFFromMarkdown(text string) (string, error) {
//	pdf := mdtopdf.NewPdfRenderer()
//
//	// Set options
//	pdf.Title = "Example PDF"
//	pdf.Author = "John Doe"
//	pdf.Creator = "mdtopdf"
//	pdf.Subject = "Example Subject"
//
//	// Convert markdown to PDF
//	err = pdf.Parse(string(md))
//	if err != nil {
//		return "", err
//	}
//
//	// Save PDF to file
//	return "", pdf.OutputFileAndClose("example.pdf")
//}
