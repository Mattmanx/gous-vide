package controller
import (
   "net/http"
   "encoding/json"
)

type Controller interface {
   GetRoutes() Routes
}

// Uses the provided net/http/ResponseWriter to return an error 500 with the provided error message.
func respondError(w http.ResponseWriter, message string) {
   w.Header().Set("Content-Type", "application/json; charset=UTF-8")
   w.WriteHeader(http.StatusInternalServerError)

   if err := json.NewEncoder(w).Encode(message); err != nil {
      panic(err)
   }
}