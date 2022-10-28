package utils

import (
	"bytes"
	"encoding/csv"
	"errors"
	"io"
	"io/ioutil"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/extrame/xls"
	"github.com/gin-gonic/gin"
)

// mType excel类型
var mType = map[string]Excel{
	"csv":  &Csv{},
	"xls":  &Xls{},
	"xlsx": &Xlsx{},
}

// CheckFileType excel类型检查
func CheckFileType(fileName string) error {
	temp := strings.Split(fileName, ".")
	if len(temp) < 2 {
		return errors.New("仅支持xlsx/xls/csv文件")
	}
	if _, ok := mType[temp[len(temp)-1]]; !ok {
		return errors.New("仅支持xlsx/xls/csv文件")
	}
	return nil
}

// excel文件读取方式
type Excel interface {
	Read(r io.Reader) ([][]string, error)
}

// 初始化excel文件实例
func NewExcel(fileType string) Excel {
	return mType[fileType]
}

// csv
type Csv struct{}

func (c *Csv) Read(r io.Reader) ([][]string, error) {
	// 初始化csv-reader
	reader := csv.NewReader(r)
	// 设置返回记录中每行数据期望的字段数，-1 表示返回所有字段
	reader.FieldsPerRecord = -1
	// 允许懒引号
	reader.LazyQuotes = true
	// 返回csv 中的所有内容
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, err
}

// xls
type Xls struct{}

func (x *Xls) Read(r io.Reader) ([][]string, error) {
	data, err := ioutil.ReadAll(r)
	xlFile, err := xls.OpenReader(bytes.NewReader(data), "utf-8")
	if err != nil {
		return nil, err
	}

	sheet := xlFile.GetSheet(0)
	if sheet.MaxRow <= 0 {
		return [][]string{}, nil
	}
	records := make([][]string, 0, sheet.MaxRow+1)
	for i := 0; i < int(sheet.MaxRow+1); i++ {
		row := sheet.Row(i)
		data := make([]string, 0)
		if row.LastCol() <= 0 {
			continue
		}
		for j := 0; j < row.LastCol(); j++ {
			col := row.Col(j)
			data = append(data, col)
		}
		records = append(records, data)
	}
	return records, nil
}

// xlsx
type Xlsx struct{}

func (x *Xlsx) Read(r io.Reader) ([][]string, error) {
	xlsx, err := excelize.OpenReader(r)
	if err != nil {
		return nil, err
	}
	records, err := xlsx.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}
	return records, nil
}

// Export 导出文件
func Export(c *gin.Context) {
	f := excelize.NewFile()
	index := f.NewSheet("Sheet1")
	f.SetCellValue("Sheet1", "A1", "Hello world.")
	f.SetCellValue("Sheet1", "B1", 100)
	f.SetActiveSheet(index)
	filename := "test" + ".xlsx"
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	//回写到web 流媒体 形成下载
	_ = f.Write(c.Writer)
}
