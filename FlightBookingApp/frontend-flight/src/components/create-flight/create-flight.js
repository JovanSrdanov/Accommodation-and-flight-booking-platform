import React, {useEffect, useState} from "react";
import axios from "axios";
import {Paper, Table, TableBody, TableContainer, TableHead} from "@mui/material";
import {styled} from "@mui/material/styles";
import TableCell, {tableCellClasses} from "@mui/material/TableCell";
import TableRow from "@mui/material/TableRow";
import "./create-flight.css"
import TextField from "@mui/material/TextField";
import {LocalizationProvider, MobileDateTimePicker} from "@mui/x-date-pickers";
import {AdapterDayjs} from "@mui/x-date-pickers/AdapterDayjs";
import dayjs from "dayjs";
import Button from "@mui/material/Button";

import AirplanemodeActiveIcon from '@mui/icons-material/AirplanemodeActive';
import Dialog from "@mui/material/Dialog";
import DialogTitle from "@mui/material/DialogTitle";
import DialogContent from "@mui/material/DialogContent";
import DialogActions from "@mui/material/DialogActions";
import {useNavigate} from "react-router-dom";


const StyledTableCell = styled(TableCell)(({theme}) => ({
    [`&.${tableCellClasses.head}`]: {
        backgroundColor: theme.palette.common.black,
        color: theme.palette.common.white,
    },
    [`&.${tableCellClasses.body}`]: {
        fontSize: 14,
    },
}));

const StyledTableRow = styled(TableRow)(({theme}) => ({
    '&:nth-of-type(odd)': {
        backgroundColor: theme.palette.action.focusOpacity,
    }
}));

function CreateFlight() {
    const [airports, setAirports] = useState([]);
    const [selectedStartPointAirport, setSelectedStartPointAirport] = useState(null);
    const [selectedDestinationAirport, setSelectedDestinationAirport] = useState(null);
    const [selectCorrectAirportsDialog, setSelectCorrectAirportsDialog] = useState(false);
    const [flightCreatedDialog, setFlightCreatedDialog] = useState(false);
    const [numberOfSeats, setNumberOfSeats] = useState(250);
    const [ticketPrice, setTicketPrice] = useState(175);
    const [dateTime, setDateTime] = useState(dayjs((new Date())));


    useEffect(() => {
        axios
            .get(process.env.REACT_APP_FLIGHT_APP_API + "airport")
            .then((response) => setAirports(response.data))
            .catch((error) => console.error(error));
    }, []);

    const handleSelectedStartPointAirport = (airport) => {
        setSelectedStartPointAirport(airport);
    };
    const handleSelectedDestinationAirport = (airport) => {
        setSelectedDestinationAirport(airport);
    };
    const isSelectedStartPointAirport = (airport) => {
        return selectedStartPointAirport && selectedStartPointAirport.id === airport.id;
    };

    const isSelectedDestinationAirport = (airport) => {
        return selectedDestinationAirport && selectedDestinationAirport.id === airport.id;
    };

    const createFlight = () => {
        if (!selectedStartPointAirport || !selectedDestinationAirport || selectedStartPointAirport.id === selectedDestinationAirport.id) {
            setSelectCorrectAirportsDialog(true)
        }
        let body = {
            destination: selectedDestinationAirport,
            startPoint: selectedStartPointAirport,
            departureDateTime: dateTime,
            price: ticketPrice,
            numberOfSeats: numberOfSeats
        }
        axios.post(process.env.REACT_APP_FLIGHT_APP_API + "flight", body).then(r => {
            setFlightCreatedDialog(true)

        }).catch(e => {
            console.log(e);
            alert("Unexpected error,please try again latter")
        })
    }

    const handleNumberOfSeatsChange = (event) => {
        setNumberOfSeats(parseInt(event.target.value));
    };


    const handleTicketPriceChange = (event) => {
        setTicketPrice(parseInt(event.target.value));
    }
    const handleDateTimeChange = (newValue) => {
        if (newValue.isAfter(dayjs())) {
            setDateTime(newValue);
        }
    }
    const navigate = useNavigate();
    return (
        <div className="createFlight">
            <div className="airports">
                <h1>Choose Start Point</h1>
                <TableContainer component={Paper} sx={{maxHeight: 400}}>
                    <Table stickyHeader>
                        <TableHead>
                            <StyledTableRow>
                                <StyledTableCell align="center">Airports</StyledTableCell>
                            </StyledTableRow>
                        </TableHead>
                        <TableBody>
                            {airports.map((airport) => (
                                <StyledTableRow

                                    key={airport.id}
                                    style={{
                                        backgroundColor: isSelectedStartPointAirport(airport) ? "var(--outlines)" : "inherit",

                                    }}
                                    onClick={() => handleSelectedStartPointAirport(airport)}>
                                    <StyledTableCell align="center">
                                        <li>Airport name: {airport.name}</li>
                                        <li>City: {airport.address.city}</li>
                                        <li>Country: {airport.address.country}</li>
                                        <li>Street: {airport.address.street},{airport.address.streetNumber} </li>

                                    </StyledTableCell>
                                </StyledTableRow>
                            ))}
                        </TableBody>
                    </Table>
                </TableContainer>
                {
                    selectedStartPointAirport && (
                        <div className="wrapper">
                            <p>Selected start point airport:</p>
                            <p>{selectedStartPointAirport.name}</p>
                        </div>
                    )
                }
            </div>
            <div className="flightInfo">
                <h1>Fill Flight Info</h1>
                <div className="flightInfoWrapper">

                    <TextField type="number"
                               variant="filled"
                               fullWidth
                               InputProps={{
                                   inputProps: {
                                       min: 1
                                   }
                               }}
                               value={numberOfSeats}
                               onChange={handleNumberOfSeatsChange}
                               label="Number of seats:"
                    />
                    <TextField type="number"
                               variant="filled"
                               fullWidth
                               InputProps={{
                                   inputProps: {
                                       min: 1
                                   }
                               }}
                               value={ticketPrice}
                               onChange={handleTicketPriceChange}
                               label="Ticket price:"
                    />

                    <LocalizationProvider dateAdapter={AdapterDayjs}>
                        <MobileDateTimePicker label="Departure date"
                                              value={dateTime}
                                              minDate={dayjs((new Date()))}
                                              onChange={handleDateTimeChange}
                        />
                    </LocalizationProvider>


                    <Button
                        variant="contained" endIcon={<AirplanemodeActiveIcon color="error"/>}
                        onClick={createFlight}
                    >Create flight
                    </Button>

                </div>

            </div>
            <div className="airports">
                <h1>Choose Destination</h1>
                <TableContainer component={Paper} sx={{maxHeight: 400}}>
                    <Table stickyHeader>
                        <TableHead>
                            <StyledTableRow>
                                <StyledTableCell align="center">Airports</StyledTableCell>
                            </StyledTableRow>
                        </TableHead>
                        <TableBody>
                            {airports.map((airport) => (
                                <StyledTableRow
                                    hover
                                    key={airport.id}
                                    style={{backgroundColor: isSelectedDestinationAirport(airport) ? "var(--outlines)" : "inherit"}}
                                    onClick={() => handleSelectedDestinationAirport(airport)}>
                                    <StyledTableCell align="center">
                                        <li>Airport name: {airport.name}</li>
                                        <li>City: {airport.address.city}</li>
                                        <li>Country: {airport.address.country}</li>
                                        <li>Street: {airport.address.street},{airport.address.streetNumber} </li>
                                    </StyledTableCell>
                                </StyledTableRow>
                            ))}
                        </TableBody>
                    </Table>
                </TableContainer>
                {
                    selectedDestinationAirport && (
                        <div className="wrapper">
                            <p>Selected destination airport:</p>
                            <p>{selectedDestinationAirport.name}</p>
                        </div>
                    )
                }

                <Dialog
                    open={selectCorrectAirportsDialog}
                    onClose={() => {
                        setSelectCorrectAirportsDialog(false)
                    }}
                >
                    <DialogTitle id="alert-dialog-title">
                        {"Both airports must be selected and they must be different from each other"}
                    </DialogTitle>
                    <DialogContent>

                    </DialogContent>
                    <DialogActions>
                        <Button onClick={() => {
                            setSelectCorrectAirportsDialog(false)
                        }}>Close</Button>
                    </DialogActions>
                </Dialog>

                <Dialog
                    open={flightCreatedDialog}
                    onClose={() => {
                        setFlightCreatedDialog(false)
                    }}
                >
                    <DialogTitle id="alert-dialog-title">
                        {"Flight created!"}
                    </DialogTitle>
                    <DialogContent>
                    </DialogContent>
                    <DialogActions>
                        <Button onClick={() => {
                            setFlightCreatedDialog(false)
                            navigate('/all-flights');
                        }}>Close</Button>
                    </DialogActions>
                </Dialog>

            </div>
        </div>
    );
}

export default CreateFlight;