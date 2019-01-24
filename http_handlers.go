package main

import (
	"net/http"
	"io/ioutil"
	"log"
	"github.com/fobilow/gitdbui/lib"
	"github.com/gorilla/mux"
	"html/template"
)

func appCss(w http.ResponseWriter, r *http.Request) {
	css := readView("static/css/app.css")
	w.Header().Set("Content-Type", "text/css");
	w.Write([]byte(css))
}

func overview(w http.ResponseWriter, r *http.Request) {

	viewModel := &lib.OverviewViewModel{}
	viewModel.Title =  "DB GUI"
	viewModel.DataSets = DataSets

	t, _ := template.ParseFiles("static/index.html", "static/sidebar.html")
	t.Execute(w, viewModel)
}

func list(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	viewDs := vars["dataset"]

	var dataSet *lib.DataSet
	for _, ds := range DataSets {
		if ds.Name == viewDs {
			dataSet = ds
			break
		}
	}

	if dataSet != nil {

		block := dataSet.Blocks[0]

		table := block.Table()

		viewModel := &lib.ListDataSetViewModel{DataSet: dataSet, Table: table}
		viewModel.DataSets = DataSets

		t, _ := template.ParseFiles("static/list.html", "static/sidebar.html")
		t.Execute(w, viewModel)
	}else{
		w.Write([]byte("Dataset ("+viewDs+") does not exist"))
	}
}

func view(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	viewModel := &lib.ViewDataSetViewModel{}
	if viewModel.Pager == nil {
		viewModel.Pager = &lib.Pager{}
	}

	viewDs := vars["dataset"]
	blockFlag := vars["b"]
	recordFlag := vars["r"]

	if blockFlag == "" && recordFlag == "" {
		viewModel.Pager.Reset()
	}else{
		viewModel.Pager.Set(blockFlag, recordFlag)
	}

	var dataSet *lib.DataSet
	for _, ds := range DataSets {
		if ds.Name == viewDs {
			dataSet = ds
			break
		}
	}

	if dataSet != nil {

		block := dataSet.Blocks[viewModel.Pager.BlockPage]

		viewModel.Pager.TotalBlocks = dataSet.BlockCount()
		viewModel.Pager.TotalRecords = block.RecordCount()

		//TODO only load record once per block - use caching
		block.LoadRecords()
		content := "No record found"
		if len(block.Records) > viewModel.Pager.RecordPage {
			content = block.Records[viewModel.Pager.RecordPage].Content
		}

		viewModel.DataSet = dataSet
		viewModel.Content = content
		viewModel.DataSets = DataSets

		t, _ := template.ParseFiles("static/view.html", "static/sidebar.html")
		t.Execute(w, viewModel)
	}else{
		w.Write([]byte("Dataset ("+viewDs+") does not exist"))
	}
}

func viewErrors(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	viewDs := vars["dataset"]

	var dataSet *lib.DataSet
	for _, ds := range DataSets {
		if ds.Name == viewDs {
			dataSet = ds
			break
		}
	}

	if dataSet != nil {
		//TODO refactor this
		dataSet.RecordCount() //hack to get records loaded so errors can be populated in dataset
		viewModel := &lib.ErrorsViewModel{DataSet: dataSet}
		viewModel.Title = "Errors"
		viewModel.DataSets = DataSets

		t, _ := template.ParseFiles("static/errors.html", "static/sidebar.html")
		t.Execute(w, viewModel)
	}else{
		w.Write([]byte("Dataset ("+viewDs+") does not exist"))
	}
}


func readView(fileName string) string {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Print(err.Error())
		return ""
	}

	return string(data)
}
