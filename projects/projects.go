package projects

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"gregoryalbouy-server-go/db"
	"gregoryalbouy-server-go/utl"
)

// Routes used by projects
var Routes = utl.RouteMap{
	"/projects":        handleRoute,
	"/projects/{test}": handleTest,
}

// type projectResult struct {
// 	Count   int        `json:"count"`
// 	Results []*Project `json:"results"`
// }

var collection *mongo.Collection

func handleTest(w http.ResponseWriter, r *http.Request) {
	resp := []byte("salut")
	w.Write(resp)
}

func todo(w http.ResponseWriter, r *http.Request) {

}

func getRequestedID(path string) string {
	ID := ""

	return ID
}

// HandleRoute for projects
func handleRoute(w http.ResponseWriter, r *http.Request) {
	collection = db.Collection("projects")

	projects := All()

	response, err := json.MarshalIndent(projects, "", "    ")
	utl.Check(err)

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// All projects
func All() (results []*Project) {
	opts := options.Find()
	// opts.SetLimit(3)
	opts.SetSort(bson.D{{Key: "addedOn", Value: -1}})

	cur, err := collection.Find(context.TODO(), bson.D{{}}, opts)
	utl.Check(err)
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var p Project
		err := cur.Decode(&p)
		utl.Check(err)

		results = append(results, &p)
	}

	err = cur.Err()
	utl.Check(err)

	return
}

// ByID func
func ByID(id primitive.ObjectID) (result *Project) {
	if err := collection.FindOne(context.TODO(), bson.M{"name": "Sarah Cornish"}).Decode(&result); err != nil {
		log.Fatal(err)
	}
	return
}
