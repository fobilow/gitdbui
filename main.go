package main

import (
	"github.com/gorilla/mux"
	"fmt"
	"net/http"
	"github.com/fobilow/gitdbui/lib"
	"flag"
	"os"
)

var DataSets []*lib.DataSet

type Endpoint struct {
	Path    string
	Handler func(w http.ResponseWriter, r *http.Request)
}

type GUI struct {
	dbPath string
	host string
	port int64
}

func (g *GUI ) start() {

	DataSets = lib.LoadDatasets(g.dbPath)
	eps := g.getGUIEndpoints()

	r := mux.NewRouter()
	for _, ep := range eps {
		r.HandleFunc(ep.Path, ep.Handler)
	}

	addr := fmt.Sprintf("%s:%d", g.host, g.port)
	fmt.Println("GitDB GUI will run at http://" + addr)

	err := http.ListenAndServe(":4120", r)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Server started!")
	}
}

func (g *GUI ) getGUIEndpoints() []*Endpoint {
	endpoints := []*Endpoint{
		{"/gitdb/css/app.css", appCss},
		{"/gitdb", overview},
		{"/gitdb/errors/{dataset}", viewErrors},
		{"/gitdb/list/{dataset}", list},
		{"/gitdb/view/{dataset}", view},
		{"/gitdb/view/{dataset}/b{b}/r{r}", view},
	}

	return endpoints
}

func main() {

	dbPathFlag := flag.String("dbPath", "", "Path to database")
	hostFlag := flag.String("host", "localhost", "Hostname for GUI server")
	portFlag := flag.Int64("port", 4120, "Port for GUI sever")
	flag.Parse()

	if *dbPathFlag == "" {
		println("Please specify dbPath i.e ./gitdbui -dbPath=/path/to/db")
		os.Exit(400)
	}

	ui := &GUI{*dbPathFlag, *hostFlag, *portFlag}
	ui.start()
}