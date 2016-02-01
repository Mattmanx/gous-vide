package controller
import (
   "net/http"
   "encoding/json"
)

type Controller interface {
   GetRoutes() Routes
}

// Uses the provided net/http/ResponseWriter to return an simple json response with a string message
func respondMessage(statusCode int, w http.ResponseWriter, message string) {
   w.Header().Set("Content-Type", "application/json; charset=UTF-8")
   w.WriteHeader(statusCode)

   if err := json.NewEncoder(w).Encode(message); err != nil {
      panic(err)
   }
}