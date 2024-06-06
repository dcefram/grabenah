package main

import (
	"encoding/csv"
	"github.com/otiai10/gosseract/v2"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	client := gosseract.NewClient()
	defer client.Close()

	err := client.SetImage("./test/sample-2.jpeg")

	if err != nil {
		log.Fatal("Image is not valid")
	}

	text, _ := client.Text()

	chunks := strings.Split(text, "\n")

	var validLines []string
	var foundHeader bool
	for _, line := range chunks {
		if line == "" || strings.Contains(line, "Rate") {
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

	// For each valid line, we can assume that the expression idx+1 % 2 should give us an idea which are pairs
	var pairs [][]string
	for idx := 0; idx < len(validLines)/2; idx++ {
		realIdx := idx * 2
		name, price := getNameAndPrice(validLines[realIdx])
		pairs = append(pairs, []string{getDate(validLines[realIdx+1]), name, price})
	}

	// We save to csv
	file, err := os.Create("result.csv")
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
	re, err := regexp.Compile("([a-zA-Z]+\\s\\d{4})")
	if err != nil {
		return ""
	}
	return re.FindString(str)
}

func getNameAndPrice(str string) (string, string) {
	// We strip the known artifacts resulting from the icons prefixed in grab screenshots
	re, err := regexp.Compile("^(fs|fe\\.)\\s([\\S\\s]+)\\s[pP#]?(\\d+\\.\\d{2})")
	if err != nil {
		return str, "0"
	}

	m := re.FindStringSubmatch(str)

	if len(m) == 0 {
		return str, "0"
	}

	return m[2], m[3]
}
