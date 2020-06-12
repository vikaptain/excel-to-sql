package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/tealeg/xlsx"

)

var dirOut = "data"

var headerStr = ""

var footerStr="";

func main() {

	dirs, _ := ioutil.ReadDir(dirOut)

	for _, dir := range dirs {
		_, err := os.Stat(dirOut + "/" + dir.Name() + "/input.xlsx")
		if err != nil {
			if !os.IsExist(err) {
				continue
			}
		}
		_, err = os.Stat(dirOut + "/" + dir.Name() + "/template.sql")
		if err != nil {
			if !os.IsExist(err) {
				continue
			}
		}

		ToSQL(dirOut+"/"+dir.Name()+"/input.xlsx", dirOut+"/"+dir.Name()+"/template.sql", "result/"+dir.Name()+".sql",dirOut+"/"+dir.Name()+"/footer.sql")
	}

}

func ToSQL(input string, template string, output string,footer string) {
	xlFile, err := xlsx.OpenFile(input)
	if err != nil {
		fmt.Println(err.Error())
	}

	strArray := make([]string, 0)
	for _, sheet := range xlFile.Sheets {
		isHeader := true
		headers := make([]string, 0)
		values := make([]map[string]string, 0)

		for _, row := range sheet.Rows {
			if isHeader {
				for _, cell := range row.Cells {
					text := cell.String()
					headers = append(headers, text)
				}
				isHeader = false
			} else {
				length := len(row.Cells)
				if length > len(headers) {
					length = len(headers)
				}
				rowMap := make(map[string]string)
				for i := 0; i < length; i++ {
					if i == 0 && row.Cells[i].String() == "" {
						break
					}

					if row.Cells[i].IsTime() {
						t, _ := row.Cells[i].GetTime(false)
						t1 := xlsx.TimeToUTCTime(t)
						ts := t1.Format("2006-01-02 15:04:05")
						rowMap[headers[i]] = ts
						continue
					}
					text := row.Cells[i].String()
					rowMap[headers[i]] = text
				}
				if len(rowMap)!=0 {
					values = append(values, rowMap)
				}
			}
		}
		template := getTemplate(template)
		for _, value := range values {
			s := template
			for key := range value {
				s = strings.ReplaceAll(s, "{"+key+"}", value[key])
			}
			strArray = append(strArray, s)
		}
		break
	}

	sqlStr:=headerStr+"\r\n"+strings.Join(strArray, "\r\n");

	// ioutil.WriteFile(output, []byte(headerStr+"\r\n"+strings.Join(strArray, "\r\n")), 0644)

	footerStr := getTemplate(footer)

	sqlStr+="\r\n";
	sqlStr=sqlStr+footerStr

	ioutil.WriteFile(output, []byte(sqlStr), 0644)

}

func getTemplate(filename string) string {
	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	return string(fileBytes)
}
