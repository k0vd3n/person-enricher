package models

type AgifyResponse struct {
    Age *int `json:"age"`
}

type GenderizeResponse struct {
    Gender *string `json:"gender"`
}

type NationalizeResponse struct {
    Country []struct {
        CountryID   string  `json:"country_id"`
        Probability float64 `json:"probability"`
    } `json:"country"`
}