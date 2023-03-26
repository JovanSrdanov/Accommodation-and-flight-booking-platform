package dto

import (
	"FlightBookingApp/errors"
	"strconv"
	"time"
)

type FlightSearchParameters struct {
	Date time.Time

	DestinationCountry string
	DestinationCity    string

	StartPointCountry string
	StartPointCity    string

	DesiredNumberOfSeats int
}

func NewFlightSearchParameters(Date string, destinationCountry string, destinationCity string,
	startPointCountry string, startPointCity string, desiredNumberOfSeats string) (*FlightSearchParameters, error) {

	flightSearchParameters := new(FlightSearchParameters)

	flightSearchParameters.DestinationCountry = destinationCountry
	flightSearchParameters.DestinationCity = destinationCity
	flightSearchParameters.StartPointCountry = startPointCountry
	flightSearchParameters.StartPointCity = startPointCity

	var err error
	layout := "2006-01-02"

	flightSearchParameters.Date, err = time.Parse(layout, Date)
	if err != nil {
		return nil, err
	}

	flightSearchParameters.DesiredNumberOfSeats, err = strconv.Atoi(desiredNumberOfSeats)
	if err != nil {
		return nil, err
	}

	parameters, err := flightSearchParameters.validate()
	if err != nil {
		return parameters, err
	}

	return flightSearchParameters, err
}

func (flightSearchParameters *FlightSearchParameters) validate() (*FlightSearchParameters, error) {
	if flightSearchParameters.DesiredNumberOfSeats <= 0 {
		return nil, &errors.DesiredNumberOfSeatsMustBeGreaterThanZeroError{}
	}

	if flightSearchParameters.Date.Before(time.Now()) {
		return nil, &errors.SearchDateMustBeInFutureError{}
	}
	return nil, nil
}
