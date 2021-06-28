package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/karriereat/dependinger/internal/pkg/database"
)

type Dependency struct {
	ID         string `json:"id"`
	ComponentA string `json:"componentA"`
	ComponentB string `json:"componentB"`
}

type Component struct {
	ID       string        `json:"id"`
	Name     string        `json:"name"`
	Colname  string        `json:"colname"`
	Children []interface{} `json:"children"`
}

type Comp struct {
	ID      int
	Name    string
	Parents []Comp
}

var dependencies []Dependency
var components []Component

var comp Comp

func main() {
	defer database.Close()
	router := mux.NewRouter()
	router.HandleFunc("/component", getComponents).Methods("GET")
	router.HandleFunc("/component", createComponent).Methods("POST")
	router.HandleFunc("/parents", getParentsForComponents).Methods("POST")
	router.HandleFunc("/component/{componentId}/parents", getComponentParents).Methods("GET")
	router.HandleFunc("/dependency", addDependency).Methods("POST")
	log.Fatal(http.ListenAndServe(":5000", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router)))
}

func getComponents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	result := database.GetComponents()
	components = nil

	for result.Next() {
		var err error
		var component Component
		err = result.Scan(&component.ID, &component.Name)

		components = append(components, component)

		if err != nil {
			panic(err.Error())
		}
	}

	json.NewEncoder(w).Encode(components)
}

func getParentsForComponents(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	body, _ := ioutil.ReadAll(request.Body)
	ids := strings.Split(string(body), ",")

	var result [][]string

	for _, id := range ids {

		dbResult := database.GetComponentParents(id)
		var tmp []string
		for dbResult.Next() {
			var err error
			var component Component

			err = dbResult.Scan(&component.ID, &component.Name)

			tmp = append(tmp, component.ID)

			if err != nil {
				panic(err.Error())
			}
		}

		result = append(result, tmp)

	}

	json.NewEncoder(writer).Encode(result)
}

func createComponent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var component Component
	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &component)

	result := database.CreateComponent(component.Name)

	if result != 0 && err == nil {
		component.ID = strconv.FormatInt(result, 10)
		json.NewEncoder(w).Encode(&component)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func addDependency(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var dependency Dependency
	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &dependency)

	componentA, err := strconv.Atoi(dependency.ComponentA)
	componentB, err := strconv.Atoi(dependency.ComponentB)

	result := database.CreateDependency(componentA, componentB)

	if result != 0 && err == nil {
		dependency.ID = strconv.FormatInt(result, 10)
		json.NewEncoder(w).Encode(&dependency)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func getComponentParents(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(request)
	requestComponent := vars["componentId"]

	result := database.GetComponentParents(requestComponent)
	components = nil

	var tmpComp Component

	tmpComp.ID = requestComponent
	tmpComp.Name = "Foo"

	for result.Next() {
		var err error
		var component Component

		err = result.Scan(&component.ID, &component.Name)

		component.Colname = "level2"
		tmpComp.Children = append(tmpComp.Children, component)

		if err != nil {
			panic(err.Error())
		}
	}

	json.NewEncoder(writer).Encode(&tmpComp)
}

/*
func fetchStuff() {

    result := database.GetComponents()
    var newParents []Comp
    components = nil

    for result.Next() {
        var err error
        var component Component
        err = result.Scan(&component.ID, &component.Name)

        var tmp Comp
        componentId, _ := strconv.Atoi(component.ID)
        tmp.ID = componentId
        tmp.Name = component.Name

        newParents = append(newParents, tmp)

        if err != nil {
            panic(err.Error())
        }
    }

    comp.Parents = newParents
}
*/
/*
func getComponentParents(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(request)
	requestComponent := vars["componentId"]

	result := database.GetComponentParents(requestComponent)
	components = nil

	for result.Next() {
		var err error
		var component Component

		err = result.Scan(&component.ID, &component.Name)

		components = append(components, component)

		if err != nil {
			panic(err.Error())
		}
	}

	//firstComponent := components[:0]

	comp.ID = 7
	comp.Name = "Finance API"
	fetchStuff()
}
*/
