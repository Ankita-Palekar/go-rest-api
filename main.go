package main

import (
  "net/http"
  "encoding/json"
  "github.com/gorilla/handlers"
  "github.com/gorilla/mux"
  "github.com/dgrijalva/jwt-go"
  "os"
  jwtmiddleware "github.com/auth0/go-jwt-middleware"
  "time"
  "github.com/auth0-community/auth0"
)

type Product struct {
  Id int
  Name string
  Slug string
  Description string
}

var mySigningKey = []byte("a super secret phrase")

var productList = []Product {
  {Id: 0, Name: "hover board", Slug: "slug-one", Description: "Some Description"},
  {Id: 1, Name: "Hover Shooters", Slug: "hover-shooters", Description: "Shoot your way to the top on 14 different hoverboards"},
  {Id: 2, Name: "Ocean Explorer", Slug: "ocean-explorer", Description: "Explore the depths of the sea in this one of a kind underwater experience"},
  {Id: 3, Name: "Dinosaur Park", Slug: "dinosaur-park", Description: "Go back 65 million years in the past and rIde a T-Rex"},
  {Id: 4, Name: "Cars VR", Slug: "cars-vr", Description: "Get behind the wheel of the fastest cars in the world."},
  {Id: 5, Name: "Robin Hood", Slug: "robin-hood", Description: "Pick up the bow and arrow and master the art of archery"},
  {Id: 6, Name: "Real World VR", Slug: "real-world-vr", Description: "Explore the seven wonders of the world in VR"},
}

var GetTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  token := jwt.New(jwt.SigningMethodHS256)
  claims := token.Claims.(jwt.MapClaims)
  claims["admin"] = true
  claims["name"] = "Bob Bobberson"
  claims["exp"] = time.Now().Add(time.Second*120).Unix()
  tokenString, _ := token.SignedString(mySigningKey)
  w.Write([]byte(tokenString))
})

func main() {
  r := mux.NewRouter()
  r.Handle("/", http.FileServer(http.Dir("./views")))
  r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
 
  // r.Handle("/status", StatusHandler).Methods("GET")
  // r.Handle("/products", jwtMiddleware.Handler(ProductsHandler)).Methods("GET")
  // r.Handle("/products/{slug}/feedback", jwtMiddleware.Handler(AddFeedbackHandler)).Methods("POST")
  r.Handle("/status", StatusHandler).Methods("GET")
  r.Handle("/products", ProductsHandler).Methods("GET")
  r.Handle("/products/{slug}/feedback", AddFeedbackHandler).Methods("POST")

  r.Handle("/get_token", GetTokenHandler).Methods("GET")
  http.ListenAndServe(":3000", handlers.LoggingHandler(os.Stdout, r))
}



var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
  ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
    return mySigningKey, nil
  },
  SigningMethod: jwt.SigningMethodHS256,
})

var ProductsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  payload, _ := json.Marshal(productList)
  w.Header().Set("Content-Type", "application/json")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Write([]byte(payload))
})

var StatusHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
  w.Write([]byte("API is up and running"))
})


var AddFeedbackHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
  var product Product
  vars := mux.Vars(r)
  slug := vars["slug"]

  for _, p := range productList {
    if p.Slug == slug {
        product = p
    }
  }

  w.Header().Set("Content-Type", "application/json")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  if product.Slug != "" {
    payload, _ := json.Marshal(product)
    w.Write([]byte(payload))
  } else {
    w.Write([]byte("Product Not Found"))
  }
})
