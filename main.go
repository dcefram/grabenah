package main

import (
	"encoding/csv"
	"github.com/otiai10/gosseract/v2"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

var validExt = []string{".jpg", ".jpeg", ".png"}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please specify the path to the image: grabenah ~/Downloads/grab.jpeg")
	}

	imgPath := os.Args[1]

	if !slices.Contains(validExt, filepath.Ext(imgPath)) {
		log.Fatal("File is not a valid image")
	}

	client := gosseract.NewClient()
	defer client.Close()

	err := client.SetImage(imgPath)

	if err != nil {
		log.Fatal("Image is not valid")
	}

	text, _ := client.Text()
	chunks := strings.Split(text, "\n")

	var validLines []string
	var foundHeader bool
	for _, line := range chunks {
		if line == "" || strings.Contains(line, "Rate") || strings.Contains(line, "Reorder") || strings.Contains(line, "= =") || strings.Contains(line, "—_—") {
			continue
		}

		if foundHeader {
			validLines = append(validLines, line)
			continue
		}

		if strings.Contains(line, "Transport") {
			foundHeader = true
		}
	}

	pairs := [][]string{{"Date", "Name", "Price"}}
	for idx := 0; idx < len(validLines)/2; idx++ {
		realIdx := idx * 2
		name, price := getNameAndPrice(validLines[realIdx])
		pairs = append(pairs, []string{getDate(validLines[realIdx+1]), name, price})
	}

	// We save to csv
	outputPath := "./"
	if len(os.Args) > 2 {
		outputPath = os.Args[2]
	}

	outputPath, filename := path.Split(outputPath)
	ext := path.Ext(filename)

	if filename == "" || ext == "" {
		outputPath = path.Join(outputPath, filename)
		filename = "result.csv"
	}

	file, err := os.Create(path.Join(outputPath, filename))
	if err != nil {
		log.Fatal("Failed to create file")
	}
	defer file.Close()

	w := csv.NewWriter(file)
	err = w.WriteAll(pairs)
	if err != nil {
		log.Fatal("Failed to write to file")
	}
}

func getDate(str string) string {
	re, err := regexp.Compile("((\\d+\\s)?[a-zA-Z]+\\s\\d{4})")
	if err != nil {
		return ""
	}
	return re.FindString(str)
}

func getNameAndPrice(str string) (string, string) {
	// We strip the known artifacts resulting from the icons prefixed in grab screenshots
	re, err := regexp.Compile("^(fs|fe|fe\\.)\\s([\\S\\s]+)\\s([pP#]|\\*\\?)?([\\d,]+\\.\\d{2})")
	if err != nil {
		return str, "0"
	}

	m := re.FindStringSubmatch(str)

	if len(m) == 0 {
		return str, "0"
	}

	p := strings.ReplaceAll(m[4], ",", "")

	return m[2], p
}
