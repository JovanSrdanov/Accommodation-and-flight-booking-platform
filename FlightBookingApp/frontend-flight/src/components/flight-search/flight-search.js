import React, {useEffect, useState} from "react";
import axios from "axios";
import {usePagination, useTable} from "react-table";
import "./flight-search.css";
import {DatePicker, LocalizationProvider} from "@mui/x-date-pickers";
import {AdapterDayjs} from "@mui/x-date-pickers/AdapterDayjs";
import TextField from '@mui/material/TextField';

const FlightSearch = () => {
    const [data, setData] = useState([]);
    const [entityCount, setEntityCount] = useState(0);
    const [loading, setLoading] = useState(false);
    const [searchParams, setSearchParams] = useState({
        departureDate: "",
        destinationCountry: "",
        destinationCity: "",
        startPointCountry: "",
        startPointCity: "",
        desiredNumberOfSeats: "",
    });
    const [pagination, setPagination] = useState({
        pageNumber: 1,
        resultsPerPage: 2,
        sortDirection: "asc",
        sortType: "departureDateTime",
    });

    useEffect(() => {
        const fetchData = async () => {
            setLoading(true);
            const {data} = await axios.get("http://localhost:4200/api/flight/search", {
                params: {...searchParams, ...pagination},
            });
            console.log(data)
            if (data.Data != null)
                setData(data.Data);
            else {

                setData([])
            }
            setEntityCount(data.EntityCount)
            setLoading(false);
        };
        fetchData();
    }, [pagination]);

    const columns = React.useMemo(
        () => [
            {
                Header: "Departure Time",
                accessor: "departureDate",
                sortType: "departureDateTime",
                onClick: () => {
                    setPagination({
                        ...pagination,
                        sortType: "departureDateTime",
                        sortDirection:
                            pagination.sortType === "departureDateTime" &&
                            pagination.sortDirection === "asc"
                                ? "desc"
                                : "asc",
                    });
                }

            },
            {
                Header: "Point of departure",
                accessor: "pointOfDeparture",
            },
            {
                Header: "Destination",
                accessor: "destination",
            },
            {
                Header: "Number of seats",
                accessor: "numberOfSeats",
            },
            {
                Header: "Number of vacant seats",
                accessor: "NumberOfVacantSeats",
            },
            {
                Header: "Ticket price",
                accessor: "ticketPrice",
                sortType: "price",
            },
            {
                Header: "Total price",
                accessor: "totalPrice",
            },
        ],
        []
    );

    const {
        getTableProps,
        headerGroups,
    } = useTable(
        {
            columns,
            data,
            initialState: {pageIndex: pagination.pageNumber - 1, pageSize: pagination.resultsPerPage},
        },
        usePagination
    );

    const handleSearchParamsChange = (e) => {
        setSearchParams({...searchParams, [e.target.name]: e.target.value});
    };

    const handlePaginationChange = (type) => {
        switch (type) {
            case "next":
                setPagination({...pagination, pageNumber: pagination.pageNumber + 1})
                break;
            case "prev":
                setPagination({...pagination, pageNumber: pagination.pageNumber - 1})
                break;
            default:
                break;
        }
    };

    return (
        <div className="flight-search">

            <div>
                <LocalizationProvider dateAdapter={AdapterDayjs}>
                    <DatePicker label="departureDate"
                                onChange={
                                    (newValue) => {
                                        var result = new Date(newValue)
                                        result.setDate(result.getDate() + 1)
                                        setSearchParams({
                                                ...searchParams,
                                                departureDate: result.toISOString().substring(0, 10)
                                            }
                                        )
                                    }
                                }
                    />
                </LocalizationProvider>


            </div>
            <div>
                <label htmlFor="destinationCountry">Destination Country:</label>
                <TextField
                    type="text"
                    name="destinationCountry"
                    onChange={handleSearchParamsChange}

                />
            </div>
            <div>
                <label htmlFor="destinationCity">Destination City:</label>
                <input
                    type="text"
                    name="destinationCity"
                    onChange={handleSearchParamsChange}

                />
            </div>
            <div>
                <label htmlFor="startPointCountry">Start Point Country:</label>
                <input
                    type="text"
                    name="startPointCountry"
                    onChange={handleSearchParamsChange}

                />
            </div>
            <div>
                <label htmlFor="startPointCity">Start Point City:</label>
                <input
                    type="text"
                    name="startPointCity"
                    onChange={handleSearchParamsChange}

                />
            </div>
            <div>
                <label htmlFor="desiredNumberOfSeats">Desired Number of Seats:</label>
                <input
                    type="text"
                    name="desiredNumberOfSeats"
                    onChange={handleSearchParamsChange}

                />
            </div>
            <button onClick={() => {
                setPagination({...pagination, pageNumber: 1})
            }}>
                Search
            </button>
            <table {...getTableProps()}>
                {entityCount > 0 &&
                    <thead>
                    {headerGroups.map((headerGroup) => (
                        <tr {...headerGroup.getHeaderGroupProps()}>
                            {headerGroup.headers.map((column) => (
                                <th {...column.getHeaderProps()}>{column.render("Header")}</th>
                            ))}
                        </tr>
                    ))}
                    </thead>}
                <tbody>
                {
                    data.map((item, i) => (
                        <tr key={i}>
                            <td>{item.Flight.departureDateTime} </td>

                            <td>
                                <ul>
                                    <li>{item.Flight.startPoint.name} </li>
                                    <li>{item.Flight.startPoint.address.city} </li>
                                    <li>{item.Flight.startPoint.address.country} </li>
                                    <li>{item.Flight.startPoint.address.street} </li>
                                    <li>{item.Flight.startPoint.address.streetNumber} </li>
                                </ul>
                            </td>

                            <td>
                                <ul>
                                    <li>{item.Flight.destination.name} </li>
                                    <li>{item.Flight.destination.address.city} </li>
                                    <li>{item.Flight.destination.address.country} </li>
                                    <li>{item.Flight.destination.address.street} </li>
                                    <li>{item.Flight.destination.address.streetNumber} </li>
                                </ul>
                            </td>

                            <td>{item.Flight.numberOfSeats} </td>
                            <td>{item.Flight.vacantSeats} </td>
                            <td>{item.Flight.price} </td>
                            <td>{item.TotalPrice} </td>
                        </tr>
                    ))}
                </tbody>
            </table>
            <div>
                <button onClick={() => handlePaginationChange("prev")} disabled={pagination.pageNumber < 2}>
                    Prev
                </button>
                <span>Page {pagination.pageNumber} of {Math.floor(entityCount / pagination.resultsPerPage)}</span>
                <button onClick={() => handlePaginationChange("next")}
                        disabled={pagination.pageNumber >= entityCount / pagination.resultsPerPage}>
                    Next
                </button>
            </div>
        </div>
    );
};
export default FlightSearch;
