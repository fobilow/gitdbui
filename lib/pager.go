package lib

import (
	"fmt"
	"strconv"
)

type Pager struct {
	BlockPage   int
	RecordPage  int
	TotalBlocks int
	TotalRecords int
}

func (p *Pager) Reset(){
	fmt.Println("Resetting pager")
	p.BlockPage = 0
	p.RecordPage = 0
}

func (p *Pager) Set(blockFlag string, recordFlag string){
	fmt.Println("Setting pager: "+ blockFlag +","+ recordFlag)
	p.BlockPage, _ = strconv.Atoi(blockFlag)
	p.RecordPage, _ = strconv.Atoi(recordFlag)
}

func (p *Pager) NextRecordUri() string {
	recordPage := p.RecordPage
	if p.RecordPage < p.TotalRecords - 1 {
		recordPage = p.RecordPage + 1
	}

	return fmt.Sprintf("b%d/r%d", p.BlockPage, recordPage)
}

func (p *Pager) PrevRecordUri() string {
    recordPage := p.RecordPage
	if p.RecordPage > 0 {
		recordPage = p.RecordPage - 1
	}

	return fmt.Sprintf("b%d/r%d", p.BlockPage, recordPage)
}

func (p *Pager) NextBlockUri() string {
	blockPage := p.BlockPage
	if p.BlockPage < p.TotalBlocks - 1 {
		blockPage = p.BlockPage + 1
	}

	return fmt.Sprintf("b%d/r%d", blockPage, 0)
}

func (p *Pager) PrevBlockUri() string {
	blockPage := p.BlockPage
	if p.BlockPage > 0 {
		blockPage = p.BlockPage - 1
	}

	return fmt.Sprintf("b%d/r%d", blockPage, 0)
}