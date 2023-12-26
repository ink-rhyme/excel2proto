package main

import (
	"log/slog"
	"os"
	"path/filepath"
	"strconv"

	"github.com/xuri/excelize/v2"
)

// excel Header Enum
const (
	HeaderName = iota
	HeaderType

	HeaderMax
)
const messageBodyRowStart = ""

// global
var ProtoPath string

func excelType2ProtoType(excelType string) string {
	return ""
}

func EachDirGetFiles(path string) ([]os.FileInfo, error) {
	dir, err := os.Open(path)
	if nil != err {
		slog.Error("dir open faild")
		return nil, err
	}
	return dir.Readdir(0)
}

func createProto(name string, table [][]string) bool {
	if "" == ProtoPath {
		return false
	}
	if HeaderMax > len(table) {
		slog.Error("table is empty")
		return false
	}

	names := table[HeaderName]
	types := table[HeaderType]
	messageBody := "syntax=\"proto3\"\npackage " + name + messageBodyRowStart + "\n{ // Start " + name
	rowStart := messageBodyRowStart + "\n\t"
	for i := 0; i < len(names) && i < len(types); i += 1 {
		rowBody := excelType2ProtoType(types[i]) + "\t\t" + names[i] + "\t\t = " + strconv.Itoa(i) + ";"
		messageBody += rowStart + rowBody
	}
	messageBody += messageBody + "\n} // End " + name
	filename := name + ".proto"
	file, err := os.Create(filepath.Join(ProtoPath, filename))
	if nil != err {
		slog.Error("create file faild", "file", name)
		return false
	}
	wSize, err := file.WriteString(messageBody)
	if nil != err || wSize != len(messageBody) {
		slog.Error("create file faild", "name", filename)
		return false
	}
	file.Close()
	slog.Debug("create file success", "name", filename)
	return true
}

func execExcel(excel *excelize.File) bool {
	if nil == excel {
		slog.Error("nil excel")
		return false
	}
	for _, sheetName := range excel.GetSheetList() {
		rows, err := excel.GetRows(sheetName)
		if nil != err {
			continue
		}
		createProto(sheetName, rows)
	}
	return true
}

func main() {
	files, err := EachDirGetFiles("")
	if nil != err {
		return
	}
	for _, file := range files {
		excel, err := excelize.OpenFile(file.Name())
		if nil != err {
			slog.Error("file open failed", "name", file.Name())
			continue
		}
		execExcel(excel)
		excel.Close()
	}
}
