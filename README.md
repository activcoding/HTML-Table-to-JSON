# HTML-Table-to-JSON
The package allows you to extract the data from an html table

Example Usage:
```go
func main() {
	data, err := converthtmltabletodata.ConvertURLToJSON("src/index.html")
	if err != nil { return }
	jsonFile, err := os.Create("jsonData.json")
	if err != nil { return }
	defer jsonFile.Close()

	_, err = jsonFile.Write(data)
	if err != nil { return }
}
```
