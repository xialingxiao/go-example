package main

import (
    "github.com/howeyc/gopass"
    "fmt"
    "net/http"
    "encoding/json"
    "io/ioutil"
    "log"
    "os"

    "github.com/gorilla/mux"
    "github.com/xialingxiao/go-example/cache"
    "github.com/xialingxiao/go-example/cache/memory"
)

var storage cache.Storage = memory.NewStorage()
var duration int64 = 3600

type Intermediate struct {
    Rates map[string]float64
    Timestamp int64
}

func get_rates(api_id string) []byte {
    resp, err := http.Get("https://openexchangerates.org/api/latest.json?app_id="+api_id)
    if err != nil {
        // handle error
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    return body
}


func cached(request func(string) []byte, api_id string) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        rates_map := storage.Get()
        if rates_map == nil {
            fmt.Println("Not cached!!!!!")
            var response Intermediate
            temp := request(api_id)
            err := json.Unmarshal(temp, &response)

            if err != nil {
                fmt.Println("err")
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            rates_map = response.Rates
            timestamp := response.Timestamp
            storage.Set(rates_map, timestamp+duration)
        }
        currency := r.URL.Query().Get("currency")
        return_map := make(map[string]float64)
        if len(currency)>0 {
            fmt.Println("single currency")
            if rate, ok := rates_map[currency]; ok {
                return_map[currency] = rate
            }
        } else {

            fmt.Println("return all")
            return_map = rates_map 
        }
        rates_bytes, err := json.Marshal(return_map)
        if err != nil {
            fmt.Println(err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.Header().Set("content-type", "application/json")
        w.Write(rates_bytes)
    }
}


func main() {
    fmt.Printf("Enter your api_id: ")
    api_id, _ := gopass.GetPasswdMasked() // Masked

    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/current_rates", cached(get_rates, string(api_id)))
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Fatal(http.ListenAndServe(":"+port, router))
}
