package main_test

import (
    "strconv"
    "log"
    "time"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "go-example"
    "github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"


	"testing"
)

var srv *http.Server


var _ = BeforeSuite(func() {
    startHTTPServer()
})

var _ = AfterSuite(func() {
})

func getRatesMock(apiID string) (int, []byte) {
    ResponseMock := []byte(`{
        "disclaimer": "Usage subject to terms: https://openexchangerates.org/terms",
        "license": "https://openexchangerates.org/license",
        "timestamp": `+strconv.FormatInt(time.Now().Unix(),10)+`,
        "base": "USD",
        "rates": {
            "AED": 3.673097,
            "AFN": 68.753,
            "ALL": 111.791078,
            "SCR": 13.631558,
            "SDG": 6.708748,
            "SEK": 7.986789,
            "SGD": 1.35839,
            "SHP": 0.774439,
            "SLL": 7524.999807,
            "SOS": 581.27,
            "SRD": 7.438,
            "SSP": 124.9444,
            "STD": 20583.157597,
            "SVC": 8.79201,
            "SYP": 515,
            "SZL": 13.09776,
            "USD": 1,
            "UYU": 28.764997,
            "UZS": 4202.05,
            "VEF": 10.05845,
            "VND": 22752.373857,
            "VUV": 105.078218,
            "WST": 2.511557,
            "ZAR": 13.032813,
            "ZMW": 9.078587,
            "ZWL": 322.355011
        }
    }`)
    return 200, ResponseMock
}

func startHTTPServer() {

    apiID := "" 
    var CACHETIME int64 = 2

    router := mux.NewRouter().StrictSlash(true)

    router.HandleFunc("/current_rates", main.Cached(getRatesMock, string(apiID), CACHETIME))

    srv = &http.Server{Addr: ":8080"}
    http.Handle("/", router)
    go func() {
        log.Fatal(srv.ListenAndServe())
    }()
}

type APIResponse struct {
    Expiration string
    Rates map[string]float64
}

func HTTPGetJson(url string) APIResponse {
    resp, err1 := http.Get(url)
    if err1 != nil {
        panic(err1)
    }
    body, err2 := ioutil.ReadAll(resp.Body)
    if err2 != nil {
        panic(err2)
    }
    var apiResponse APIResponse
    err3 := json.Unmarshal(body, &apiResponse)
    if err3 != nil {
        panic(err3)
    }
    err4 := resp.Body.Close()
    if err4 != nil {
        panic(err4)
    }
    return apiResponse
}

func TestGoExample(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GoExample Suite")
}

var _ = Describe("GoExample", func() {
    var expiration1 string
    var expiration2 string
    var expiration3 string
    Describe("First response", func() {
        Context("When querying the api for the first time", func() {
            It("Should return exchange rates", func() {
                apiResponse := HTTPGetJson("http://127.0.0.1:8080/current_rates")
                expiration1 = apiResponse.Expiration
                Expect(apiResponse.Rates["VUV"]).To(Equal(105.078218))
            })
        })
    })
    Describe("Second response", func() {
        Context("When querying the api for the second time", func() {
            It("Should return cached exchange rates", func() {
                time.Sleep(1*time.Second)
                apiResponse := HTTPGetJson("http://127.0.0.1:8080/current_rates")
                expiration2 = apiResponse.Expiration
                Expect(expiration2).To(Equal(expiration1))
            })
        })
    })
    Describe("Third response", func() {
        Context("When querying the api after the caching period", func() {
            It("Should return with new records with a new expiration time", func() {
                time.Sleep(2*time.Second)
                apiResponse := HTTPGetJson("http://127.0.0.1:8080/current_rates")
                expiration3 = apiResponse.Expiration
                Expect(expiration3).To(Not(Equal(expiration1)))
            })
        })
    })
})
