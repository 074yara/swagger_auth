package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	net "net/url"
	"strings"
)

var yandexURL = "https://geocode-maps.yandex.ru/1.x/"

type YandexAPIer interface {
	GetGeoFromAddress(address string) (string, string)
	GetAddrFromGeo(lat, lon string) string
}

type YandexAPI struct {
	apikey string
	format string
}

func NewYandexAPI(apikey string) YandexAPIer {
	return &YandexAPI{apikey: apikey, format: "json"}
}

type YandexResponse struct {
	Response struct {
		GeoObjectCollection struct {
			MetaDataProperty struct {
				GeocoderResponseMetaData struct {
					Request string `json:"request"`
				} `json:"GeocoderResponseMetaData"`
			} `json:"metaDataProperty"`
			FeatureMember []struct {
				GeoObject struct {
					Point struct {
						Pos string `json:"pos"`
					} `json:"Point"`
					MetaDataProperty struct {
						GeocoderMetaData struct {
							Text string `json:"text"`
						} `json:"GeocoderMetaData"`
					} `json:"metaDataProperty"`
				} `json:"GeoObject"`
			} `json:"featureMember"`
		} `json:"GeoObjectCollection"`
	} `json:"response"`
}

func (y YandexAPI) GetGeoFromAddress(address string) (string, string) {
	var yandexResponse YandexResponse
	//без QueryEscape не работает на русском
	address = net.QueryEscape(address)
	url := fmt.Sprintf(`%v?apikey=%v&geocode=%v&format=%v&resultsresults=1`, yandexURL, y.apikey, address, y.format)
	req, err := http.NewRequest("GET", url, nil)
	checkError(err)

	resp, err := http.DefaultClient.Do(req)
	checkError(err)

	data, err := io.ReadAll(resp.Body)
	checkError(err)

	err = json.Unmarshal(data, &yandexResponse)
	checkError(err)

	lat, lon := ``, ``
	geoStr := yandexResponse.Response.GeoObjectCollection.FeatureMember[0].GeoObject.Point.Pos
	geoSlice := strings.Split(geoStr, ` `)
	if len(geoSlice) == 2 {
		lat, lon = geoSlice[1], geoSlice[0]
	}
	//fmt.Println(yandexResponse.Response.GeoObjectCollection.FeatureMember[0].GeoObject.MetaDataProperty.GeocoderMetaData.Text)
	return lat, lon
}

func (y YandexAPI) GetAddrFromGeo(lat, lon string) string {
	var yandexResponse YandexResponse
	var address string
	//без QueryEscape не работает на русском
	lat = net.QueryEscape(lat)
	lon = net.QueryEscape(lon)
	url := fmt.Sprintf(`%v?apikey=%v&geocode=%v,%v&format=%v&resultsresults=1`, yandexURL, y.apikey, lon, lat, y.format)
	req, err := http.NewRequest("GET", url, nil)
	checkError(err)

	resp, err := http.DefaultClient.Do(req)
	checkError(err)

	data, err := io.ReadAll(resp.Body)
	checkError(err)

	err = json.Unmarshal(data, &yandexResponse)
	checkError(err)

	address = yandexResponse.Response.GeoObjectCollection.FeatureMember[0].
		GeoObject.MetaDataProperty.GeocoderMetaData.Text
	return address
}
