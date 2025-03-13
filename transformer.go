package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/jung-kurt/gofpdf"
)

func main() {
	data, err := os.ReadDir("../../../Documents/Obsidian Vault/AWS")

	pdf := gofpdf.New("P", "mm", "A4", "")
  pdf.AddUTF8Font("OpenSans", "", "OpenSans-Regular.ttf")
  pdf.AddUTF8Font("OpenSans", "B", "OpenSans-Bold.ttf")

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range data {
		properName := file.Name()[:len(file.Name())-3]

		fileData, err := os.ReadFile("../../../Documents/Obsidian Vault/AWS/" + file.Name())
		if err != nil {
			log.Fatal(err)
		}

		ext := filepath.Ext(file.Name())
		if ext != ".md" {
			continue
		}

		newFileData := string(fileData)
		imagesRegex := regexp.MustCompile(`!\[\[([^\]]+)\]\]`)
		boldTextRegex := regexp.MustCompile(`\*\*(.*?)\*\*`)
    bulletPointsRegex := regexp.MustCompile(`- `)

		imagesMatches := imagesRegex.FindAllStringSubmatch(newFileData, -1)

    bullet := " • "

    boldIdx := boldTextRegex.FindAllStringIndex(newFileData, -1)
    lastBoldIndex := 0

    // heading
    pdf.AddPage()
    pdf.SetFont("OpenSans", "B", 20)
    pdf.MultiCell(190, 10, properName, "", "L", false)
    pdf.Ln(10)

    // all bold text
    for _, boldIndex := range boldIdx {
      fmt.Println(boldIndex[0])
      pdf.SetFont("OpenSans", "", 12)
      formattedBulets := bulletPointsRegex.ReplaceAllString(newFileData[lastBoldIndex:boldIndex[0]], bullet)
      pdf.Write(10,formattedBulets)

      pdf.SetFont("OpenSans", "B", 12)
      pdf.Write(10, " " + newFileData[boldIndex[0]+2:boldIndex[1]-2])
      lastBoldIndex = boldIndex[1]
    }

    // images
		for _, match := range imagesMatches {
			if len(match) > 1 {
				fmt.Println("Unsupported:", match[1])
			}
		}

    pdf.SetFont("OpenSans", "", 12)
    formattedBulets := bulletPointsRegex.ReplaceAllString(newFileData[lastBoldIndex:], bullet)
    pdf.Write(10,formattedBulets)
    pdf.Ln(30)
	}

	errPdf := pdf.OutputFileAndClose("output.pdf")

	if errPdf != nil {
		log.Fatal(errPdf)
	}
}
