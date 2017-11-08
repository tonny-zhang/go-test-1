package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"path"

	"io/ioutil"

	"strings"

	"strconv"

	"github.com/tealeg/xlsx"
)

var dirCurrent, err = os.Getwd()

func errPrint(msg string) {
	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", msg)
}
func getNumStr(numStr string) (string, bool) {
	i := strings.LastIndex(numStr, "%")
	isHavePercent := i > 0 && i == strings.Count(numStr, "")-1-1
	if isHavePercent {
		numStr = numStr[0:i]
	}
	return numStr, isHavePercent
}
func convert(excelFileName string) {
	excelFileName = strings.Replace(excelFileName, "\\", "/", -1)
	if strings.Index(path.Base(excelFileName), "~") == 0 {
		fmt.Println(excelFileName + " 文件名不合法，或不是一个完整的excel文件")
		return
	}
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		if r := recover(); r != nil {
			errPrint(fmt.Sprintf("%s 解析错误 panic的内容%v\n", excelFileName, r))
		}
		// fmt.Printf("%s 解析错误，头列数不一致，请检查！\n", excelFileName)
	}()
	for _, sheet := range xlFile.Sheets {
		if strings.Index(sheet.Name, "_") != 0 || len(sheet.Rows) == 0 {
			continue
		}
		rows := sheet.Rows
		// fmt.Printf("rows = %d %s\n", len(rows), excelFileName)
		nameZH := rows[0].Cells
		types := rows[1].Cells
		nameEN := rows[2].Cells

		lenZHCell := len(nameZH)
		lenEnCell := len(nameEN)
		var lenCellActual = 0
		var headerRow []map[string]string
		for i, cell := range nameZH {
			name, _ := cell.String()
			if len(name) == 0 || i >= lenEnCell {
				break
			}
			enCell := nameEN[i]
			if nil != enCell {
				cellHeader := make(map[string]string)
				en, _ := enCell.String()
				cellHeader["name"] = name
				cellHeader["en"] = en

				typeCell := types[i]
				if nil != typeCell {
					t, _ := types[i].String()
					cellHeader["type"] = t
				}
				headerRow = append(headerRow, cellHeader)
				lenCellActual++
			} else {
				break
			}

			// headerRow = append(headerRow, &headerCell{name, t, en})
		}
		// fmt.Println(headerRow)
		// b, err := json.Marshal(headerRow)
		// fmt.Println(string(b), err, len(b))

		var data []map[string]interface{}
		for _, row := range rows[4:] {
			len := len(row.Cells)
			if len > 0 && len <= lenZHCell {
				var dMap = make(map[string]interface{})
				var lenNull = 0
				var isEmpty = false
				for index, cell := range row.Cells {
					if index < lenCellActual {
						d, _ := cell.String()
						if d == "" {
							if index == 0 {
								isEmpty = true
								break
							}
							lenNull++
						}
						en, _ := nameEN[index].String()
						t, _ := types[index].String()
						t = strings.ToLower(t)

						numStr, isHavePercent := getNumStr(d)
						if t == "int" {
							valNumber, _ := strconv.Atoi(numStr)
							dMap[en] = valNumber
							if isHavePercent {
								dMap[en+"_isp"] = true
							}
						} else if t == "float" {
							valNumber, _ := strconv.ParseFloat(numStr, 64)
							dMap[en] = valNumber
							if isHavePercent {
								dMap[en+"_isp"] = true
							}
						} else if t == "bool" {
							dMap[en] = strings.ToUpper(d) == "T"
						} else {
							dMap[en] = d
						}

						// fmt.Printf("%s\t", d)
					}
				}
				// fmt.Printf("lenNull = %d, len = %d, %t, %v\n", lenNull, len, lenNull < len, dMap);
				// 过滤全空行
				if !isEmpty && lenNull < len {
					data = append(data, dMap)
				}
				// fmt.Printf("\n")
			}
		}

		// b1, err1 := json.Marshal(data)
		// fmt.Println(string(b1), err1, len(b1))

		result := make(map[string]interface{})
		result["header"] = headerRow
		result["root"] = data
		bResult, _ := json.Marshal(result)
		// fmt.Println(string(bResult), errResult, len(bResult))

		outputdir := path.Join(dirCurrent, "output")
		os.MkdirAll(outputdir, os.ModePerm)

		regPostfix := regexp.MustCompile("\\..+$")
		filenameNew := regPostfix.ReplaceAllString(path.Base(excelFileName), ".json")
		outfilename := path.Join(outputdir, filenameNew)

		f, err := os.Create(outfilename)
		defer f.Close()
		if err == nil {
			f.Write(bResult)
			fmt.Printf("%s save!\n", outfilename)
		} else {
			fmt.Println(err)
		}
	}
}
func walk(dir string) {
	files, err := ioutil.ReadDir(dir)
	if err == nil {
		for _, file := range files {
			filepath := path.Join(dir, file.Name())
			convert(filepath)
		}
	} else {
		fmt.Println(err)
	}
}
func main() {
	// convert("E:\\source\\nodejs\\tool\\police\\data\\data\\skillEquip.xlsx")
	// walk("E:\\source\\nodejs\\tool\\police\\data\\data\\")
	// file := "E:\\source\\nodejs\\tool\\police\\data\\data\\activity.xlsx"
	// fmt.Println(file)
	// file = strings.Replace(file, "\\", "/", -1)
	// fmt.Println(file)
	// fmt.Println(path.Base(file))
	// reg := regexp.MustCompile("\\..+$")
	// fmt.Println(reg.ReplaceAllString(path.Base(file), ".json"))

	dirExcel := path.Join(dirCurrent, "data")
	if info, err := os.Stat(dirExcel); !os.IsNotExist(err) && info.IsDir() {
		walk(dirExcel)
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("\n\n回车退出...")
		reader.ReadByte()
		os.Exit(0)
	} else {
		errPrint("当前目录下没有用于存放excel文件的data目录")
	}
}
