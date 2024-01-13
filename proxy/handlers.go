package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// На DADATA кончились запросы, их там всего 130, поэтому перевёл всё на Yandex GEO API

type Address struct {
	Result string `json:"result,omitempty"`
	GeoLat string `json:"lat,omitempty"`
	GeoLon string `json:"lon,omitempty"`
}

type AddressSearchRequest struct {
	Query string `json:"query"`
}
type SearchResponse struct {
	Addresses []*Address `json:"addresses"`
}

// geoFromAddressHandler getting coordinates from address via Yandex Geo API
//
//	@Summary		Get coordinates from address
//	@Description	Get coordinates from address
//	@Tags			search
//	@Accept			json
//	@Produce		json
//	@Param			input	body		AddressSearchRequest	true	"AddressSearchRequest"
//	@Success		200		{object}	SearchResponse
//	@Failure		503
//	@Failure		400
//	@Router			/address/search [post]
func geoFromAddressHandler(rw http.ResponseWriter, r *http.Request) {
	var searchRequest AddressSearchRequest
	var searchResponse SearchResponse

	err := json.NewDecoder(r.Body).Decode(&searchRequest)
	checkError(err)

	//Check for correct request
	if len(searchRequest.Query) == 0 {
		http.Error(rw, "Incorrect request", http.StatusBadRequest)
		return
	}

	controller := NewYandexAPI(`a4d0e5c5-b51f-4e0e-956f-bfeaeaffa363`)
	lat, lon := controller.GetGeoFromAddress(searchRequest.Query)

	//Check for correct response
	if len(lat) == 0 || len(lon) == 0 {
		http.Error(rw, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	searchResponse.Addresses = []*Address{{
		Result: "",
		GeoLat: lat,
		GeoLon: lon,
	}}
	err = json.NewEncoder(rw).Encode(searchResponse)
	checkError(err)

	//DADATA CODE
	/*
		api := dadata.NewCleanApi()
		addresses, err := api.Address(context.Background(), searchRequest.Query)
		log.Println(searchRequest.Query)
		checkError(err)

		if addresses == nil {
			http.Error(rw, "No answer from Dadata API", http.StatusServiceUnavailable)
			return
		}

		searchResponse.Addresses = []*Address{{GeoLat: addresses[0].GeoLat, GeoLon: addresses[0].GeoLon, Result: addresses[0].Result}}
		err = json.NewEncoder(rw).Encode(searchResponse)
		checkError(err)

	*/
}

// addressFromGeoHandler getting address from coordinates via Yandex GEO API
//
//	@Summary		Get address from coordinates
//	@Description	Get address from coordinates
//	@Tags			geocode
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	SearchResponse
//	@Param			input	body		Address	true	"Address"
//	@Failure		503
//	@Failure		400
//	@Router			/address/geocode [post]
func addressFromGeoHandler(rw http.ResponseWriter, r *http.Request) {
	var searchResponse SearchResponse
	var address Address
	err := json.NewDecoder(r.Body).Decode(&address)
	checkError(err)

	//Check for correct request
	if len(address.GeoLat) == 0 || len(address.GeoLon) == 0 {
		http.Error(rw, "Incorrect request", http.StatusBadRequest)
		return
	}

	controller := NewYandexAPI("a4d0e5c5-b51f-4e0e-956f-bfeaeaffa363")
	result := controller.GetAddrFromGeo(address.GeoLat, address.GeoLon)

	//Check for correct response
	if len(result) == 0 {
		http.Error(rw, "Service unavailable", http.StatusServiceUnavailable)
		return
	}
	address.Result = result
	searchResponse.Addresses = []*Address{&address}
	err = json.NewEncoder(rw).Encode(searchResponse)

	//DADATA CODE
	/*
		newReq, err := http.NewRequest("POST", AddressFromGeoAPIURL, r.Body)
		checkError(err)

		newReq.Header.Set("Content-Type", "application/json")
		newReq.Header.Set("Accept", "application/json")
		newReq.Header.Set("Authorization", "Token "+dadataAPIKey)
		newReq.Header.Set("X-Secret", dadataSecretKey)

		response, err := http.DefaultClient.Do(newReq)
		checkError(err)
		defer response.Body.Close()

		err = json.NewDecoder(response.Body).Decode(&addressFromGeo)
		checkError(err)
		query := addressFromGeo.Addresses[0].Value
		log.Println(query)

		api := dadata.NewCleanApi()
		addresses, err := api.Address(context.Background(), query)
		checkError(err)

		if addresses == nil {
			http.Error(rw, "No answer from Dadata API", http.StatusServiceUnavailable)
			return

		}

		searchResponse.Addresses = []*Address{{Result: addresses[0].Result, GeoLon: addresses[0].GeoLon, GeoLat: addresses[0].GeoLat}}
		err = json.NewEncoder(rw).Encode(searchResponse)
		checkError(err)

	*/
}

func checkError(err error) {
	if err != nil {
		log.Println(err)
	}
}

// NOT IN USE:
//type AddressFromGeoResponse struct {
//	Addresses []*struct {
//		Value string `json:"value"`
//	} `json:"suggestions"`
//}
//const (
//	AddressFromGeoAPIURL = "http://suggestions.dadata.ru/suggestions/api/4_1/rs/geolocate/address"
//	dadataAPIKey         = "2cb6ce15db3cb2b52b47f9c39d250875b89d0723"
//	dadataSecretKey      = "617a4664fbe290c49ea12e5500a53e5e69995246"
//)
