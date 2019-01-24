package lib

import (
	"path/filepath"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"strings"
	"sort"
)

type Block struct {
	DataSet *DataSet
	Name    string
	Size    int64
	Records []*Record
	BadRecords []string
}

func (b *Block) HumanSize() string {
	return formatBytes(uint64(b.Size))
}

func (b *Block) RecordCount() int {
	return len(b.Records)
}

func (b *Block) LoadRecords(){
	b.Records = b.records()
}

func (b *Block) Show() []string {
	var blockContents []string
	for _, record := range b.Records {
		blockContents = append(blockContents, record.Content)

	}

	return blockContents
}

func (b *Block) readBlock() ([]string, error) {

	var result []string
	var fmtErr error
	var jsonErr error
	var formattedJson []byte

	blockFile := filepath.Join(b.DataSet.DbPath, b.DataSet.Name, b.Name+".json")
	data, err := ioutil.ReadFile(blockFile)
	if err != nil {
		return result, err
	}

	var dataBlock map[string]interface{}
	var record map[string]interface{}

	fmtErr = json.Unmarshal(data, &dataBlock)
	if fmtErr != nil {
		return result, &badBlockError{fmtErr.Error()+" - "+blockFile, blockFile}
	}

	recordKeys := orderMapKeys(dataBlock)

	for _, k := range recordKeys {
		recordJson := dataBlock[k].(string)
		jsonErr = json.Unmarshal([]byte(recordJson), &record)
		formattedJson, jsonErr = json.MarshalIndent(record, "", "\t")
		if jsonErr != nil {
			return result, &badRecordError{jsonErr.Error()+" - "+k, k}
		}

		result = append(result, string(formattedJson))
	}

	return result, err
}

func (b *Block) records() []*Record {

	var records []*Record
	b.DataSet.BadBlocks = []string{}
	b.DataSet.BadRecords = []string{}

	recs, err := b.readBlock()

	if err != nil {
		if  be, ok := err.(*badBlockError); ok {
			b.DataSet.BadBlocks = append(b.DataSet.BadBlocks, be.blockFile)
		}else if  re, ok := err.(*badRecordError); ok {
			b.DataSet.BadRecords = append(b.DataSet.BadRecords, re.recordId)
		}

		return records
	}

	for _, rec := range recs {
		records = append(records, &Record{Content: rec})
	}

	return records
}

func (b *Block) Table() *Table {
	b.LoadRecords()

	table := &Table{}
	var jsonMap map[string]interface{}
	for i, record := range b.Records {
		e := json.Unmarshal([]byte(record.Content), &jsonMap)
		if e != nil {
			fmt.Println(e.Error())
		}

		var row []string
		if i == 0 {
			table.Headers = orderMapKeys(jsonMap)
		}
		for _, key := range table.Headers {
			val := fmt.Sprintf("%v",jsonMap[key])
			if len(val) > 40 {
				val = val[0:40]
			}
			row = append(row, val)
		}

		table.Rows = append(table.Rows, row)
	}

	return table
}


const (
	BYTE = 1.0 << (10 * iota)
	KILOBYTE
	MEGABYTE
	GIGABYTE
	TERABYTE
)

func formatBytes(bytes uint64) string {
	unit := ""
	value := float32(bytes)

	switch {
	case bytes >= TERABYTE:
		unit = "TB"
		value = value / TERABYTE
	case bytes >= GIGABYTE:
		unit = "GB"
		value = value / GIGABYTE
	case bytes >= MEGABYTE:
		unit = "MB"
		value = value / MEGABYTE
	case bytes >= KILOBYTE:
		unit = "KB"
		value = value / KILOBYTE
	case bytes >= BYTE:
		unit = "B"
	case bytes == 0:
		return "0"
	}

	stringValue := fmt.Sprintf("%.1f", value)
	stringValue = strings.TrimSuffix(stringValue, ".0")
	return fmt.Sprintf("%s%s", stringValue, unit)
}

func orderMapKeys(_map map[string]interface{}) []string{
	// To store the keys in slice in sorted order
	var keys []string
	for k := range _map {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i, v := range keys {
		if v == "ID" {
			//swap ID to the front of the array
			keys[i], keys[0] = keys[0], keys[i]
			break
		}
	}

	return keys
}
