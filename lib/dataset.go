package lib

import (
	"strings"
	"time"
	"io/ioutil"
	"path/filepath"
	"log"
)

type DataSet struct {
	Name   string
	DbPath string
	Blocks []*Block
	BadBlocks []string
	BadRecords []string
	LastModified time.Time
}

func (d *DataSet) Size() int64 {
	size := int64(0)
	for _, block := range d.Blocks {
		size += block.Size
	}

	return size
}

func (d *DataSet) HumanSize() string {
	return formatBytes(uint64(d.Size()))
}

func (d *DataSet) BlockCount() int {
	return len(d.Blocks)
}

func (d *DataSet) RecordCount() int {
	count := 0
	for _, block := range d.Blocks {
		block.LoadRecords()
		count += block.RecordCount()
	}

	return count
}

func (d *DataSet) BadBlocksCount() int {
	return len(d.BadBlocks)
}

func (d *DataSet) BadRecordsCount() int {
	return len(d.BadRecords)
}

func (d *DataSet) LastModifiedDate() string {
	return d.LastModified.Format("02 Jan 2006 15:04:05")
}

func (d *DataSet) LoadBlocks(){
	var blocks []*Block
	blks, err := ioutil.ReadDir(filepath.Join(d.DbPath, d.Name))
	if err != nil {
		log.Print(err.Error())
	}
	for _, block := range blks {
		if !block.IsDir() && strings.HasSuffix(block.Name(), ".json") {
			blockName := strings.TrimSuffix(block.Name(), ".json")

			b := &Block{
				DataSet: d,
				Name:    blockName,
				Size:    block.Size(),
			}

			blocks = append(blocks, b)
		}
	}

	d.Blocks = blocks
}

func (d *DataSet) blocks() []*Block {
	var blocks []*Block
	blks, err := ioutil.ReadDir(filepath.Join(d.DbPath, d.Name))
	if err != nil {
		return blocks
	}

	for _, block := range blks {
		if !block.IsDir() && strings.HasSuffix(block.Name(), ".json") {
			blockName := strings.TrimSuffix(block.Name(), ".json")

			b := &Block{
				DataSet: d,
				Name:    blockName,
				Size:    block.Size(),
			}

			blocks = append(blocks, b)
		}
	}

	return blocks
}

func (d *DataSet) Indexes() []string {
	//grab indexes
	var indexes []string

	indexFiles, err := ioutil.ReadDir(filepath.Join(d.DbPath, ".gitdb/Index/", d.Name))
	if err != nil {
		return indexes
	}

	for _, indexFile := range indexFiles {
		indexes = append(indexes, strings.TrimSuffix(indexFile.Name(), ".json"))
	}

	return indexes
}

func LoadDatasets(dbPath string) []*DataSet {
	var dataSets []*DataSet

	dirs, err := ioutil.ReadDir(dbPath)
	if err != nil {
		log.Print(err.Error())
		return dataSets
	}

	for _, dir := range dirs {
		if !strings.HasPrefix(dir.Name(), ".") && dir.IsDir() {
			dataset := &DataSet{
				Name: dir.Name(),
				DbPath: dbPath,
				LastModified: dir.ModTime(),
			}

			dataset.LoadBlocks()
			dataSets = append(dataSets, dataset)
		}
	}
	return dataSets
}