package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/Kudzeri/go-mongo-api/models"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type UserController struct {
	client *mongo.Client
}

func NewUserController(client *mongo.Client) *UserController {
	return &UserController{client: client}
}

func isValidObjectId(id string) bool {
	_, err := primitive.ObjectIDFromHex(id)
	return err == nil
}

func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	if !isValidObjectId(id) {
		http.NotFound(w, r)
		return
	}

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ObjectId", http.StatusBadRequest)
		return
	}

	u := models.User{}
	f := bson.M{"_id": oid}

	err = uc.client.Database("mydb").Collection("users").FindOne(r.Context(), f)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	uj, err := json.Marshal(u)
	if err != nil {
		http.Error(w, "Error marshalling user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", uj)
}

func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u := models.User{}

	json.NewDecoder(r.Body).Decode(&u)

	u.ID = primitive.NewObjectID()

	uc.client.Database("mydb").Collection("users").InsertOne(r.Context(), u)

	uj, err := json.Marshal(u)
	if err != nil {
		http.Error(w, "Error marshalling user", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", uj)
}

func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	if !isValidObjectId(id) {
		http.NotFound(w, r)
		return
	}

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ObjectId", http.StatusBadRequest)
		return
	}

	result, err := uc.client.Database("mydb").Collection("users").DeleteOne(r.Context(), bson.M{"_id": oid})
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Deleted User %v\n", oid)
}
