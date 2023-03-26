package utils

import (
	"FlightBookingApp/errors"
	"strconv"
)

type PageInfo struct {
	PageNumber     int
	ResultsPerPage int
	SortDirection  string
	SortType       string
}

func NewPageInfo(_pageNumber string, _resultsPerPage string, sortDirection string, sortType string) (*PageInfo, error) {

	pageNumber, err := strconv.Atoi(_pageNumber)
	if err != nil {
		return nil, err
	}
	resultsPerPage, err := strconv.Atoi(_resultsPerPage)
	if err != nil {
		return nil, err
	}
	if pageNumber <= 0 {
		return nil, &errors.PageNumberMustBeGreaterThanZeroError{}
	}
	if resultsPerPage <= 0 {
		return nil, &errors.ResultsPerPageMustBeGreaterThanZeroError{}
	}
	pageInfo := PageInfo{pageNumber, resultsPerPage, sortDirection, sortType}

	return &pageInfo, nil

}
