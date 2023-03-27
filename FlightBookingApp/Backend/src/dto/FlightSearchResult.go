package dto

import "FlightBookingApp/model"

type FlightSearchResult struct {
	Flight     model.Flight
	TotalPrice float32
}

func NewFlightSearchResult(flight *model.Flight, desiredNumberOfSeats int) *FlightSearchResult {
	flightSearchResults := new(FlightSearchResult)
	flightSearchResults.TotalPrice = float32(flight.Price * float32(desiredNumberOfSeats))
	flightSearchResults.Flight = *flight
	return flightSearchResults
}
