package buspredictions

import (
	"encoding/xml"
	"errors"
	"github.com/awiede/wmata-go-sdk/wmata"
	"strings"
)

const busPredictionsServiceBaseUrl = "https://api.wmata.com/NextBusService.svc"

type GetNextBusResponse struct {
	XMLName            xml.Name            `json:"-" xml:"http://www.wmata.com NextBusResponse"`
	NextBusPredictions []NextBusPrediction `json:"Predictions" xml:"Predictions>NextBusPrediction"`
	StopName           string              `json:"StopName" xml:"StopName"`
}

type NextBusPrediction struct {
	DirectionNumber string `json:"DirectionNum" xml:"DirectionNum"`
	DirectionText   string `json:"DirectionText" xml:"DirectionText"`
	Minutes         int    `json:"Minutes" xml:"Minutes"`
	RouteID         string `json:"RouteID" xml:"RouteID"`
	TripID          string `json:"TripID" xml:"TripID"`
	VehicleID       string `json:"VehicleID" xml:"VehicleID"`
}

type BusPredictions interface {
	GetNextBuses(stopID string) (*GetNextBusResponse, error)
}

var _ BusPredictions = (*Service)(nil)

// NewService returns a new Incidents service with a reference to an existing wmata.Client
func NewService(client *wmata.Client, responseType wmata.ResponseType) *Service {
	return &Service{
		client:       client,
		responseType: responseType,
	}
}

type Service struct {
	client       *wmata.Client
	responseType wmata.ResponseType
}

func (service *Service) GetNextBuses(stopID string) (*GetNextBusResponse, error) {
	if stopID == "" {
		return nil, errors.New("stopID is required")
	}

	var requestUrl strings.Builder
	requestUrl.WriteString(busPredictionsServiceBaseUrl)

	switch service.responseType {
	case wmata.JSON:
		requestUrl.WriteString("/json/jPredictions")
	case wmata.XML:
		requestUrl.WriteString("/Predictions")
	}

	nextBus := GetNextBusResponse{}

	return &nextBus, service.client.BuildAndSendGetRequest(service.responseType, requestUrl.String(), map[string]string{"StopID": stopID}, &nextBus)

}