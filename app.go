package main

import (
    "os"

    "net/http"

    "github.com/urfave/negroni"
    "github.com/gorilla/mux"
    "encoding/json"
    "github.com/tobyjsullivan/ues-v2/events"
    "github.com/tobyjsullivan/ues-v2/service"
)

const (
    svcEntityId = "48ab2171-bd06-4646-9809-108b56449353"
)

func main() {
    r := buildRoutes()

    n := negroni.New()
    n.UseHandler(r)

    port := os.Getenv("PORT")
    if port == "" {
        port = "3000"
    }

    n.Run(":" + port)
}

func buildRoutes() http.Handler {
    r := mux.NewRouter()
    r.HandleFunc("/", statusHandler).Methods("GET")
    r.HandleFunc("/commands/create-account", createAccountHandler).Methods("POST")
    r.HandleFunc("/commands/commit-event", commitEventHandler).Methods("POST")

    return r
}

type statusResponse struct {
    Status string `json:"status"`
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
    encoder := json.NewEncoder(w)
    resp := &statusResponse{
        Status: "ok",
    }
    encoder.Encode(resp)
}

func createAccountHandler(w http.ResponseWriter, r *http.Request) {
    requestId := r.Header.Get("x-request-id")
    if requestId == "" {
        http.Error(w, "All requests must include x-request-id header.", http.StatusBadRequest)
        return 
    }

    // Parse request
    vars := mux.Vars(r)
    accountId := vars["account-id"]
    email := vars["email"]
    password := vars["password"]

    exists, err := accountExists(accountId)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if exists {
        http.Error(w, "Account ID already exists.", http.StatusBadRequest)
        return
    }

    // TODO Validate identity (only email/password for now)

    // TODO Create the account

    // TODO Associate identity

    // TODO Generate auth token


    http.Error(w, "Not implemented", http.StatusNotImplemented)
}


func accountExists(accountId string) (bool, error) {
    svc, err := service.LoadAggregate(svcEntityId)
    if err != nil {
        return false, err
    }

    for _, id := range svc.AccountIDs {
        if id == accountId {
            return true, nil
        }
    }

    return false, nil
}

func emailInUse(email string) bool {
    // TODO
}


func commitEventHandler(w http.ResponseWriter, r *http.Request) {
    // TODO Parse request

    // TODO Load account associated with

    // TODO Check if entity is claimed by


    http.Error(w, "Not implemented", http.StatusNotImplemented)
}
