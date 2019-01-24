package lib

type BaseViewModel struct {
	Title    string
	DataSets []*DataSet
}

type OverviewViewModel struct {
	BaseViewModel
}

type ViewDataSetViewModel struct {
	BaseViewModel
	DataSet *DataSet
	Pager *Pager
	Content string
}

type ListDataSetViewModel struct {
	BaseViewModel
	DataSet *DataSet
	Table *Table
}

type ErrorsViewModel struct {
	BaseViewModel
	DataSet *DataSet
}