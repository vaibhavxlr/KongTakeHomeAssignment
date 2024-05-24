package internal

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/vaibhavxlr/KongTakeHomeAssignment/internal/DTOs"
	dbclient "github.com/vaibhavxlr/KongTakeHomeAssignment/internal/dbClient"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ListServices(w http.ResponseWriter, r *http.Request)  {
	currPageStr := r.URL.Query().Get("curr")
	countStr := r.URL.Query().Get("count")
	sort := r.URL.Query().Get("sortOrder")
	searchQuery := r.URL.Query().Get("search")
	ctx := r.Context()
	currPage := 1
	count := 5
	totalPgCount := 0
	sortOrder := 0

	if currPageStr != "" {
		currPage, _ = strconv.Atoi(currPageStr)
	}
	if countStr != "" {
		count, _ = strconv.Atoi(countStr)
	}
	if sort != "" {
		sortOrder, _ = strconv.Atoi(sort)
	}

	coll := dbclient.DB_OBJ.Collection("serviceList")
	services := make([]DTOs.Service, 0)

	findOpt := options.Find() 
	
	if sortOrder == 1 {
		findOpt.SetSort(bson.D{{Key: "id", Value: -1}})
	} else {
		findOpt.SetSort(bson.D{{Key: "id", Value: 1}})
	}

	if searchQuery == "" {
		cursor, err := coll.Find(ctx, bson.D{}, findOpt)
		if err != nil {
			log.Println("Failed to fetch data from DB in ListServices API, err: ", err)
		}
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var service DTOs.Service
			if err = cursor.Decode(&service); err != nil {
				log.Println("Cursor couldn't find document, err: ", err)
			}
			services = append(services, service)
			totalPgCount++
		}
	} else {
		filter := bson.M{
			"$or" : []bson.M {
				{	"name":bson.M{"$regex": searchQuery, "$options": "i"}},
				{"info":bson.M{"regex": searchQuery, "$options": "i"}},
			},
		}
		cursor, err := coll.Find(ctx, filter, findOpt)
		if err != nil {
			log.Println("Failed to fetch data from DB in ListServices API, err: ", err)
		}
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var service DTOs.Service
			cursor.Decode(&service)
			services = append(services, service)
			totalPgCount++ 
		}
	}
	
	var response DTOs.ListServicesResp
	sortodr := DTOs.SortOrder{
		AZ: 0,
		ZA: 1,
	}
	response.SortOrder = sortodr
	pagedetails := DTOs.PageDetails {
		Curr:currPage,
		Count: count,
		Total: totalPgCount/count,
	}
	response.PageDetails = pagedetails

	startInd := (currPage - 1) * count
	endInd := startInd + count

	servicelist := make([]DTOs.Service, 0)
	for startInd < endInd && startInd < totalPgCount {
		servicelist = append(servicelist, services[startInd])
		startInd++
	}

	response.Services = servicelist

	respBytes, _ := json.Marshal(response)
	w.Write(respBytes)
}

func ServiceDetails(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	serviceId := r.PathValue("id")
	collServ := dbclient.DB_OBJ.Collection("serviceList")
	collVer := dbclient.DB_OBJ.Collection("versions")
	filter := bson.M{"id": serviceId}
	var serviceData DTOs.Service
	err := collServ.FindOne(ctx, filter).Decode(&serviceData)
	log.Println(err)
	filter = bson.M{"serviceId": serviceId}
	cursor, err := collVer.Find(ctx, filter)
	if err != nil {
		log.Println("Failed to fetch versions, err: ", err)
	}
	versions := make([]DTOs.Version, 0)
	for cursor.Next(ctx) {
		var version DTOs.Version
		cursor.Decode(&version)
		versions = append(versions, version)
	}	
	serviceData.Versions = versions
	respBytes, _ := json.Marshal(serviceData)
	w.Write(respBytes)
}