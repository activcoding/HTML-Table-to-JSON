package converthtmltabletodata

import (
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Row map[string]string

type Table []Row

func Convert(reader io.ReadCloser) ([]Table, error) {
	tables := []Table{}

	document, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return tables, err
	}

	document.Find("table").Each(func(tableIndex int, tableSelection *goquery.Selection) {
		if tableIndex == 0 {
			return
		}

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
				table = append(table, row)
			}
		})

		tables = append(tables, table)
	})

	return tables, nil
}
