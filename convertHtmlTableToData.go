package converthtmltabletodata

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Row represents a single row in a table. It is a map from header string to data string.
type Row map[string]string

// Table represents a table as a slice of Rows.
type Table []Row

// Convert takes an io.ReadCloser and returns Tables and an error.
// The Tables are Rows, where each Row is a hashmap with the header string has the key.
// If there is an error reading from the input or parsing the HTML, Convert returns the Tables processed so far and the error.
func Convert(reader io.ReadCloser) ([]Table, error) {
	tables := []Table{}

	document, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return tables, err
	}

	document.Find("table").Each(func(tableIndex int, tableSelection *goquery.Selection) {
		table := Table{}
		var headers []string

		tableSelection.Find("tr").Each(func(index int, rowSelection *goquery.Selection) {
			row := Row{}
			if index == 0 {
				rowSelection.Find("th").Each(func(i int, headerSection *goquery.Selection) {
					headers = append(headers, strings.TrimSpace(headerSection.Text()))
				})
			} else {
				rowSelection.Find("td").Each(func(dataSectionIndex int, dataSection *goquery.Selection) {
					if dataSectionIndex < len(headers) {
						header := headers[dataSectionIndex]
						row[header] = strings.TrimSpace(dataSection.Text())
					}
				})
				if len(row) > 0 {
					table = append(table, row)
				}
			}
		})

		if len(table) > 0 {
			tables = append(tables, table)
		}
	})

	return tables, nil
}

func ConvertUrl(filePath string) ([]Table, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	tables, err := Convert(file)
	if err != nil {
		return nil, err
	}

	return tables, nil
}

func ConvertReaderToJSON(reader io.ReadCloser) ([]byte, error) {
	tables, err := Convert(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	json, err := json.Marshal(tables)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return json, nil
}

func ConvertURLToJSON(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	return ConvertReaderToJSON(file)
}
