package internal

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/vaibhavxlr/KongTakeHomeAssignment/internal/DTOs"
	dbclient "github.com/vaibhavxlr/KongTakeHomeAssignment/internal/dbClient"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBConn interface {
	Find(context.Context, interface{}, ...*options.FindOptions) (Cursor, error)
}

type ActualCollection struct {
	collection *mongo.Collection
}

func (ac *ActualCollection) Find(ctx context.Context, filter interface{}, dbOptions ...*options.FindOptions) (Cursor, error) {
	cursor, err := ac.collection.Find(ctx, filter, dbOptions...)
	if err != nil {
		return nil, err
	}
	return cursor, nil
}

type Cursor interface {
	Next(context.Context) bool
	Decode(interface{}) error
	Close(context.Context) error
}

func listServices(w http.ResponseWriter, r *http.Request, coll DBConn) {
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

	services := make([]DTOs.Service, 0)

	findOpt := options.Find()
	if sortOrder == 1 {
		findOpt.SetSort(bson.D{{Key: "id", Value: -1}})
	} else {
		findOpt.SetSort(bson.D{{Key: "id", Value: 1}})
	}

	var cursor Cursor
	var err error
	if searchQuery == "" {
		cursor, err = coll.Find(ctx, bson.D{}, findOpt)
		if err != nil {
			log.Println("Failed to fetch data from DB in ListServices API, err: ", err)
			errResp := DTOs.ErrorResp{
				ErrorCode:   http.StatusText(http.StatusInternalServerError),
				ErrorString: err.Error(),
			}
			w.WriteHeader(http.StatusInternalServerError)
			respByte, _ := json.Marshal(errResp)
			w.Write([]byte(respByte))
			return
		}
		defer cursor.Close(ctx)
	} else {
		filter := bson.M{
			"$or": []bson.M{
				{"name": bson.M{"$regex": searchQuery, "$options": "i"}},
				{"info": bson.M{"regex": searchQuery, "$options": "i"}},
			},
		}
		cursor, err = coll.Find(ctx, filter, findOpt)
		if err != nil {
			log.Println("Failed to fetch data from DB in ListServices API, err: ", err)
			errResp := DTOs.ErrorResp{
				ErrorCode:   http.StatusText(http.StatusInternalServerError),
				ErrorString: err.Error(),
			}
			w.WriteHeader(http.StatusInternalServerError)
			respByte, _ := json.Marshal(errResp)
			w.Write([]byte(respByte))
			return
		}
		defer cursor.Close(ctx)
	}

	for cursor.Next(ctx) {
		var service DTOs.Service
		cursor.Decode(&service)
		services = append(services, service)
		totalPgCount++
	}

	var response DTOs.ListServicesResp
	sortodr := DTOs.SortOrder{
		AZ: 0,
		ZA: 1,
	}
	response.SortOrder = sortodr
	pagedetails := DTOs.PageDetails{
		Curr:  currPage,
		Count: count,
		Total: totalPgCount / count,
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
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}

// Handler func for service listing, sorting and filter/search
func ListServices(w http.ResponseWriter, r *http.Request) {
	coll := dbclient.DB_OBJ.Collection("serviceList")
	acColl := ActualCollection{
		collection: coll,
	}
	listServices(w, r, &acColl)
}

// Handler func for service and version details API
func ServiceDetails(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	serviceId := r.PathValue("id")
	collServ := dbclient.DB_OBJ.Collection("serviceList")
	collVer := dbclient.DB_OBJ.Collection("versions")

	filter := bson.M{"id": serviceId}
	var serviceData DTOs.Service
	err := collServ.FindOne(ctx, filter).Decode(&serviceData)
	if err != nil {
		log.Println("Failed to fetch service data, err: ", err)
		errResp := DTOs.ErrorResp{
			ErrorCode:   http.StatusText(http.StatusInternalServerError),
			ErrorString: err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		respByte, _ := json.Marshal(errResp)
		w.Write([]byte(respByte))
		return
	}

	filter = bson.M{"serviceId": serviceId}
	cursor, err := collVer.Find(ctx, filter)
	if err != nil {
		log.Println("Failed to fetch version data, err: ", err)
		errResp := DTOs.ErrorResp{
			ErrorCode:   http.StatusText(http.StatusInternalServerError),
			ErrorString: err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		respByte, _ := json.Marshal(errResp)
		w.Write([]byte(respByte))
		return
	}

	versions := make([]DTOs.Version, 0)
	for cursor.Next(ctx) {
		var version DTOs.Version
		cursor.Decode(&version)
		versions = append(versions, version)
	}
	serviceData.Versions = versions
	respBytes, _ := json.Marshal(serviceData)
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}
