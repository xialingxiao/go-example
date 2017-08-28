#### Task:
Create a stand alone HTTP service written in Golang that returns the current foreign exchange rates within the hour.

#### Specifics:
Openexchange to be used as the source of forex data and retrieved rates should be cached for an hour at a time. You are free to use any version of Golang above 1.5.

#### Endpoints that must be implemented:
*All endpoints should return JSON*

`GET /current_rates`
- Returns the list of latest rates retrieved from Openexchange.

`GET /current_rates?currency=<currency_code>`
- Returns the current rate for the desired currency code e.g /current_rates?currency=USD

#### Things you may need:
- https://tour.golang.org/welcome/1
- https://docs.openexchangerates.org/ (you would need to sign up for a free API key)

#### To run
`make run` should start the web application at port `8080`

#### To test
`make test` should test the web application


