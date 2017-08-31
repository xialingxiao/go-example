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
    "go-example/cache"
    "go-example/cache/memory"
    // "github.com/xialingxiao/go-example/cache"
    // "github.com/xialingxiao/go-example/cache/memory"
)


// PORT is declared as a top level variable
var PORT string // the port to listen to

// initialize the in-memory storage
var storage cache.Storage = memory.NewStorage()


// OpenExResponse only care about two fields in the reponse
// will ignore all other fields and errors
// all exchange rates use USD as base
type OpenExResponse struct {
    Rates map[string]float64
    Timestamp int64
}

// SuccessResponse returns a time string for expiration time instead of epoch
type SuccessResponse struct {
    Rates map[string]float64 `json:"rates"`
    Expiration string `json:"expiration"`
}

func getRates(apiID string) (int, []byte) {
    resp, err1 := http.Get("https://openexchangerates.org/api/latest.json?app_id="+apiID)
    if err1 != nil {
        log.Fatal("Error querying openexchangerates.org API")
    }
    defer func (){
        if cerr := resp.Body.Close(); cerr != nil && err1 == nil {
            err1 = cerr
        }
    }()
    body, err2 := ioutil.ReadAll(resp.Body)
    if err2 != nil {
        log.Fatal("Error querying openexchangerates.org API")
    }
    fmt.Println(string(body))
    return resp.StatusCode, body
}

// Cached wraps arround the query to openexchangerates.org and cache the results in memory
func Cached(getLatest func(string) (int, []byte), apiID string, cacheTime int64) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("content-type", "application/json")
        rates, expiration := storage.Get()

        if rates == nil {
            // Exchange rates not cached or have expired
            log.Println(fmt.Sprintf("Exchange rates not cached or have expired, needs to retrieve latest from openexchangerates.org"))
            var openExResponse OpenExResponse
            statusCode, temp := getLatest(apiID)
            err := json.Unmarshal(temp, &openExResponse)

            if err != nil {
                log.Println(fmt.Sprintf("Error %q", err.Error()))
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            if openExResponse.Timestamp == 0 {
                log.Println(fmt.Sprintf("Error querying openexchangerates.org, make sure to start the server with the correct api_id"))
                w.WriteHeader(statusCode)
                if err := json.NewEncoder(w).Encode(map[string]string{"error":"Error querying openexchangerates.org, make sure to start the server with the correct api_id"}); err != nil {
                    panic(err)
                }
                return
            }

            // initialize rates and expiration
            rates = openExResponse.Rates
            expiration = openExResponse.Timestamp+cacheTime
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

        w.Header().Set("content-type", "application/json")
        if err := json.NewEncoder(w).Encode(successResponse); err != nil {
            panic(err)
        }
        log.Println(fmt.Sprintf("Query ended successfully"))
    }
}

func main() {
    // use environment variable for setting the server port, defaults to 8080
    PORT = os.Getenv("PORT")
    if PORT == "" {
        PORT = "8080"
    }

    const CACHETIME int64 = 3600

    // You would need to manually input your api_id upon start
    fmt.Printf("Enter your api_id: ")
    apiID, _ := gopass.GetPasswdMasked() // Masked

    // Initialise the router
    router := mux.NewRouter().StrictSlash(true)

    // The only api endpoint we will implement
    router.HandleFunc("/current_rates", Cached(getRates, string(apiID), CACHETIME))

    // start server
    log.Fatal(http.ListenAndServe(":"+PORT, router))
}
