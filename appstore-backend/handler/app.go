package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"appstore/model"
	"appstore/service"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/pborman/uuid"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {// the entrance for uploading app
    // Parse r *http.Request to build app
    fmt.Println("Received one upload request")
    user := r.Context().Value("user")
    claims := user.(*jwt.Token).Claims
    username := claims.(jwt.MapClaims)["username"]
 
    app := model.App{
        Id:          uuid.New(),// we new a id from strach
        User:        username.(string),
        Title:       r.FormValue("title"),
        Description: r.FormValue("description"),
    }


  price, err := strconv.Atoi(r.FormValue("price"))
    fmt.Printf("%v,%T", price, price)
    if err != nil {
        fmt.Println(err)
    }
    app.Price = price
 
 
    file, _, err := r.FormFile("media_file")
    if err != nil {
        http.Error(w, "Media file is not available", http.StatusBadRequest)
        fmt.Printf("Media file is not available %v\n", err)
        return
    }
 
    err = service.SaveApp(&app, file)
    if err != nil {
        http.Error(w, "Failed to save app to backend", http.StatusInternalServerError)
        fmt.Printf("Failed to save app to backend %v\n", err)
        return
    }
 
 
    fmt.Println("App is saved successfully.")
 
     fmt.Fprintf(w, "App is saved successfully: %s\n", app.Description)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one search request")
    w.Header().Set("Content-Type", "application/json")
    //get params from request
    title := r.URL.Query().Get("title")
    description := r.URL.Query().Get("description")
 
    //call service to handle request
    var apps []model.App
    var err error
    apps, err = service.SearchApps(title, description)
    if err != nil {
        http.Error(w, "Failed to read Apps from backend", http.StatusInternalServerError)
        return
    }
 
    //construct response
    js, err := json.Marshal(apps)
    //it serializes the Go data into JSON format
    if err != nil {
        http.Error(w, "Failed to parse Apps into JSON format", http.StatusInternalServerError)
        return
    }
    w.Write(js)
 }
 
 func checkoutHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one checkout request")
    w.Header().Set("Content-Type", "text/plain")
    //get params from request
    appID := r.FormValue("appID")
    //call service
    s, err := service.CheckoutApp(r.Header.Get("Origin"), appID)
    if err != nil {
        fmt.Println("Checkout failed.")
        w.Write([]byte(err.Error()))
        return
    }
    //construct URL
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(s.URL))
 
    fmt.Println("Checkout process started!")
 }
 
 