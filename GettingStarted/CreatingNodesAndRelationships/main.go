package main

import (
	"fmt"
	"github.com/jmcvetta/neoism"
	"github.com/verdverm/neo4j-tutorials/common/reset"
)

var (
	db *neoism.Database
)

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func resetDB() {
	reset.RemoveNeo4jDB()
	reset.StartNeo4jDB()
}

func init() {
	resetDB()
	var err error
	db, err = neoism.Connect("http://localhost:7474/db/data")
	if err != nil {
		panic(err)
	}
}

func main() {
	createNode()
	queryNode()
	createMovie()
	createUnique()
	setNodeProperty()
	queryMovies()

	queryAllNodes()
}

// Create a node with neoism function
func createNode() {
	actor := "Tom Hanks"
	// Create a node
	n, err := db.CreateNode(neoism.Props{"name": actor})
	if err != nil {
		panic(err)
	}
	// Add a label
	n.AddLabel("Actor")

	fmt.Println("createNode()", n.Data)
}

func queryNode() {
	// query statemunt
	stmt := `
		MATCH (actor:Actor)
		WHERE actor.name = {actorSub}
		RETURN actor
	`
	// query params
	actor := "Tom Hanks"
	params := neoism.Props{"actorSub": actor}

	// query results
	res := []struct {
		Actor neoism.Node
	}{}

	// construct query
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	// execute query
	err := db.Cypher(&cq)
	panicErr(err)

	// check results
	if len(res) != 1 {
		panic(fmt.Sprintf("Incorrect results len in query1()\n\tgot %d, expected 1\n", len(res)))
	}
	n := res[0].Actor // Only one row of data returned

	fmt.Printf("queryNode() -> %+v\n", n.Data)
}

func createMovie() {
	actor := "Tom Hanks"
	movie := "Sleepless in Seattle"

	// query statemunt
	stmt := `
		MATCH (actor:Actor)
		WHERE actor.name = {actorSub}
		CREATE (movie:Movie {title: {movieSub}})
		CREATE (actor)-[:ACTED_IN]->(movie);
	`
	// query params
	params := neoism.Props{
		"actorSub": actor,
		"movieSub": movie,
	}

	// construct query
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     nil,
	}
	// execute query
	err := db.Cypher(&cq)
	panicErr(err)

	fmt.Println("createMovie()")
}

func createUnique() {
	actor := "Tom Hanks"
	movie := "Forrest Gump"

	// query statemunt
	stmt := `
		MATCH (actor:Actor {name: {actorSub}})
		CREATE UNIQUE (actor)-[r:ACTED_IN]->(movie:Movie {title: {movieSub}})
		RETURN r;
	`
	// query params
	params := neoism.Props{
		"actorSub": actor,
		"movieSub": movie,
	}

	// query results
	res := []struct {
		A   string `json:"a.name"` // `json` tag matches column name in query
		Rel string `json:"type(r)"`
		B   string `json:"b.name"`
	}{}

	// construct query
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	// execute query
	err := db.Cypher(&cq)
	panicErr(err)

	r := res[0]
	fmt.Println("createUnique()", r.A, r.Rel, r.B)
}

func setNodeProperty() {
	actor := "Tom Hanks"
	dob := 1944

	// query statemunt
	stmt := `
		MATCH (actor:Actor {name: {actorSub}})
		SET actor.DoB = {dobSub}
		RETURN actor.name, actor.DoB;
	`
	// query params
	params := neoism.Props{
		"actorSub": actor,
		"dobSub":   dob,
	}

	// query results
	res := []struct {
		Name string `json:"actor.name"` // `json` tag matches column name in query
		DoB  string `json:"actor.DoB"`
	}{}

	// construct query
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	// execute query
	err := db.Cypher(&cq)
	panicErr(err)

	r := res[0]
	fmt.Println("setNodeProperty()", r.Name, r.DoB)
}

func queryMovies() {
	// query statemunt
	stmt := `
		MATCH (movie:Movie)
		RETURN movie;
	`
	// query params
	actor := "Tom Hanks"
	params := neoism.Props{"actorSub": actor}

	// query results
	res := []struct {
		Movie neoism.Node
	}{}

	// construct query
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	// execute query
	err := db.Cypher(&cq)
	panicErr(err)

	// check results
	if len(res) != 2 {
		panic(fmt.Sprintf("Incorrect results len in query1()\n\tgot %d, expected 2\n", len(res)))
	}

	fmt.Printf("queryMovies()\n")
	for i, _ := range res {
		n := res[i].Movie // Only one row of data returned
		fmt.Printf("  Node[%d] %+v\n", i, n.Data)
	}
}

func queryAllNodes() {
	// query results
	res := []struct {
		N neoism.Node // Column "n" gets automagically unmarshalled into field N
	}{}

	// construct query
	cq := neoism.CypherQuery{
		Statement: "MATCH (n) RETURN n",
		Result:    &res,
	}
	// execute query
	err := db.Cypher(&cq)
	panicErr(err)

	fmt.Printf("queryAllNodes(%d)\n", len(res))
	for i, _ := range res {
		n := res[i].N // Only one row of data returned
		fmt.Printf("  Node[%d] %+v\n", i, n.Data)
	}

}