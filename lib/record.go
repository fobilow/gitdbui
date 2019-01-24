package lib

type Record struct {
	ID      string
	Content string
}

type Table struct {
	Headers []string
	Rows [][]string
}