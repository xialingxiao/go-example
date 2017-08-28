package main

import (
    "github.com/howeyc/gopass"
    "fmt"
    "net/http"
    "encoding/json"
    "io/ioutil"
    "log"
    "time"
    "os"

    "github.com/gorilla/mux"
    "github.com/xialingxiao/go-example/cache"
    "github.com/xialingxiao/go-example/cache/memory"
)


// declare the top level variable(s)
var PORT string // the port to listen to

// initialize the in-memory storage
var storage cache.Storage = memory.NewStorage()


// we only care about two fields in the reponse
// will ignore all other fields and errors
// all exchange rates use USD as base
type OpenExResponse struct {
    Rates map[string]float64
    Timestamp int64
}

type SuccessResponse struct {
    Rates map[string]float64 `json:"rates"`
    Expiration string `json:"expiration"`
}

func get_rates(api_id string) (int, []byte) {
    resp, _ := http.Get("https://openexchangerates.org/api/latest.json?app_id="+api_id)
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    return resp.StatusCode, body
}

func cached(get_latest func(string) (int, []byte), api_id string) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("content-type", "application/json")
        rates, expiration := storage.Get()

        if rates == nil {
            // Exchange rates not cached or have expired
            log.Println(fmt.Sprintf("Exchange rates not cached or have expired, needs to retrieve latest from openexchangerates.org"))
            var openExResponse OpenExResponse
            statusCode, temp := get_latest(api_id)
            err := json.Unmarshal(temp, &openExResponse)

            if err != nil {
                log.Println(fmt.Sprintf("Error %q", err.Error()))
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            if openExResponse.Timestamp == 0 {
                log.Println(fmt.Sprintf("Error querying openexchangerates.org, make sure to start the server with the correct api_id"))
                w.WriteHeader(statusCode)
                w.Write([]byte(`{"error": "Error querying openexchangerates.org, make sure to start the server with the correct api_id"}`))
                return
            }

            // initialize rates and expiration
            rates = openExResponse.Rates
            expiration = openExResponse.Timestamp+3600
            storage.Set(rates, expiration)
        }

        currency := r.URL.Query().Get("currency")
        successResponse := SuccessResponse{}
        successResponse.Expiration = time.Unix(expiration, 0).Format("2006-01-02 15:04:05")
        if len(currency)>0 {
            log.Println(fmt.Sprintf("Query for exchange rate of %q", currency))
            if rate, ok := rates[currency]; ok {
                successResponse.Rates = make(map[string]float64)
                successResponse.Rates[currency] = rate
            }
        } else {
            log.Println(fmt.Sprintf("Query for all exchange rates"))
            successResponse.Rates = rates 
        }
        successResponseBytes, _ := json.Marshal(successResponse)

        w.Header().Set("content-type", "application/json")
        w.Write(successResponseBytes)
        log.Println(fmt.Sprintf("Query ended successfully"))
    }
}

func main() {
    // use environment variable for setting the server port, defaults to 8080
    PORT = os.Getenv("PORT")
    if PORT == "" {
        PORT = "8080"
    }

    // You would need to manually input your api_id upon start
    fmt.Printf("Enter your api_id: ")
    api_id, _ := gopass.GetPasswdMasked() // Masked

    // Initialise the router
    router := mux.NewRouter().StrictSlash(true)

    // The only api endpoint we will implement
    router.HandleFunc("/current_rates", cached(get_rates, string(api_id)))

    // start server
    log.Fatal(http.ListenAndServe(":"+PORT, router))
}
