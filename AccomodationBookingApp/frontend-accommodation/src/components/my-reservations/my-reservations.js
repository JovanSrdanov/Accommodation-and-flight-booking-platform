import React, {useEffect, useState} from 'react';
import {
    Box,
    Button,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    Paper,
    styled,
    Table,
    TableBody,
    TableCell,
    tableCellClasses,
    TableContainer,
    TableRow,
    TextField
} from "@mui/material";
import interceptor from "../../interceptor/interceptor";
import {Flex} from "reflexbox";
import axios from "axios";

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

function MyReservations() {
    const [reservations, setReservations] = useState(null);
    const [flightResults, setFlightResults] = useState(null);
    const [selectedReservation, setSelectedReservation] = useState(null);
    const [toDialogShow, setToDialogShow] = useState(false);
    const [fromDialogShow, setFromDialogShow] = useState(false);
    const [ticketBoughtDialogShow, setTicketBoughtDialogShow] = useState(false);
    const [APIKey, setAPIkey] = useState("");
    const [selectedFlight, setSelectedFlight] = useState(null);
    const [enterAPIkeyDialogShow, setEnterAPIkeyDialogShow] = useState(false);
    const [errorAPIDialogShow, setErrorAPIDialogShow] = useState(false);

    const getAllReservations = () => {
        interceptor.get("api-2/accommodation/reservation/all/guest").then(res => {
            setReservations(res.data)
        }).catch(err => {
            console.log(err)
        })
    }
    useEffect(() => {
        getAllReservations();
    }, []);
    const handleCancel = (Id) => {
        interceptor.get("api-1/reservation/cancel/" + Id).then(res => {
            getAllReservations();
        }).catch(err => {
            console.log(err)
        })
    };
    const [toInputData, setToInputData] = useState({
        startPointCountry: '',
        startPointCity: '',
        desiredNumberOfSeats: 1,
    });

    const [fromInputData, setFromInputData] = useState({
        destinationCountry: '',
        destinationCity: '',
        desiredNumberOfSeats: 1,
    });
    const handleChangeTo = (event) => {
        setToInputData((prevState) => ({
            ...prevState,
            [event.target.name]: event.target.value,
        }));
    };
    const handleChangeFrom = (event) => {
        setFromInputData((prevState) => ({
            ...prevState,
            [event.target.name]: event.target.value,
        }));
    };


    const handleSearchTo = () => {
        const params = {
            departureDate: selectedReservation.dateRange.from.substring(0, 10),
            destinationCountry: selectedReservation.address.country,
            destinationCity: selectedReservation.address.city,
            startPointCountry: toInputData.startPointCountry,
            startPointCity: toInputData.startPointCity,
            desiredNumberOfSeats: toInputData.desiredNumberOfSeats,
            pageNumber: 1,
            resultsPerPage: 25,
            sortDirection: "asc",
            sortType: "departureDateTime"
        };
        axios.get("http://localhost:4200/api/search-flights", {params}).then((res) => {
            setFlightResults(res.data.Data)
        }).catch((err) => {
            console.log(err)
        })
    };

    const flightsFromAccommodationClickHandle = (r) => {
        setSelectedReservation(r)
        setFromDialogShow(true)
    };
    const handleFromDialogShowClose = () => {
        setFromDialogShow(false)
        setFlightResults(null)
    };
    const flightsToAccommodationClickHandle = (r) => {
        setSelectedReservation(r)
        setToDialogShow(true)
    };

    const handleToDialogShowClose = () => {
        setToDialogShow(false)
        setFlightResults(null)
    };
    const handleBuyTicketsClick = (item) => {
        setEnterAPIkeyDialogShow(true)
        setSelectedFlight(item)
    };


    const handleCloseBuyApiKeyDialog = () => {
        setEnterAPIkeyDialogShow(false)
        setAPIkey("")
    };
    const handleBuyTickets = () => {
        let sendData = {
            apiKey: APIKey,
            flightId: selectedFlight.Flight.id,
            numberOfTickets: parseInt(toInputData.desiredNumberOfSeats)
        }
        axios.post("http://localhost:4200/api/ticket/api-key", sendData)
            .then((res) => {
                setTicketBoughtDialogShow(true)
            }).catch((err) => {
            console.log(err)
            setErrorAPIDialogShow(true)
        })

    };
    const handleCloseTickedDialogBought = () => {
        setTicketBoughtDialogShow(false)
        handleCloseBuyApiKeyDialog()
    };
    const handleSearchFrom = () => {
        const params = {
            departureDate: selectedReservation.dateRange.to.substring(0, 10),
            destinationCountry: fromInputData.destinationCountry,
            destinationCity: fromInputData.destinationCity,
            startPointCountry: selectedReservation.address.country,
            startPointCity: selectedReservation.address.city,
            desiredNumberOfSeats: fromInputData.desiredNumberOfSeats,
            pageNumber: 1,
            resultsPerPage: 25,
            sortDirection: "asc",
            sortType: "departureDateTime"
        };
        console.log(params)
        axios.get("http://localhost:4200/api/search-flights", {params}).then((res) => {

            setFlightResults(res.data.Data)
        }).catch((err) => {
            console.log(err)
        })
    };
    const handleCloseErrorApiDialog = () => {
        setErrorAPIDialogShow(false)
    };
    return (
        <>

            <Dialog open={fromDialogShow} onClose={handleFromDialogShowClose} fullWidth maxWidth="lg">
                <DialogTitle>Search for fligths from accommodation:</DialogTitle>
                <DialogContent>
                    <Flex flexDirection="column">
                        <Flex m={1}>
                            <Box width={1 / 3} m={1}>
                                <TextField
                                    label="Destination Country"
                                    name="destinationCountry"
                                    value={fromInputData.destinationCountry}
                                    onChange={handleChangeFrom}
                                    fullWidth

                                />
                            </Box>
                            <Box width={1 / 3} m={1}>
                                <TextField
                                    label="Destination City"
                                    name="destinationCity"
                                    value={fromInputData.destinationCity}
                                    onChange={handleChangeFrom}
                                    fullWidth

                                />
                            </Box>
                            <Box width={1 / 3} m={1}>
                                <TextField
                                    label="Desired Number of Seats"
                                    name="desiredNumberOfSeats"
                                    type="number"
                                    value={fromInputData.desiredNumberOfSeats}
                                    onChange={handleChangeFrom}
                                    fullWidth
                                    InputProps={{
                                        inputProps: {
                                            min: 1
                                        }
                                    }}

                                />
                            </Box>
                        </Flex>
                        {flightResults != null && flightResults.length > 0 && (
                            <div>
                                <TableContainer component={Paper}
                                                sx={{maxHeight: 700, overflowY: 'scroll'}}>
                                    <Table>
                                        <TableBody>
                                            {flightResults.map((item, idx) =>
                                                (
                                                    <React.Fragment key={`${idx}-row`}>
                                                        <StyledTableRow>
                                                            <StyledTableCell>
                                                                <Box m={1} sx={{
                                                                    overflowX: 'auto',
                                                                    width: 200,
                                                                    height: 100,
                                                                    overflowy: 'auto'
                                                                }}>
                                                                    <li>Airport name: {item.Flight.startPoint.name}</li>
                                                                    <li>City: {item.Flight.startPoint.address.city}</li>
                                                                    <li>Country {item.Flight.startPoint.address.country}</li>
                                                                    <li>Street: {item.Flight.startPoint.address.street}, {item.Flight.startPoint.address.streetNumber}</li>
                                                                </Box>
                                                            </StyledTableCell>
                                                            <StyledTableCell>
                                                                <Box m={1} sx={{
                                                                    overflowX: 'auto',
                                                                    width: 200,
                                                                    height: 100,
                                                                    overflowy: 'auto'
                                                                }}>
                                                                    <li>Airport
                                                                        name: {item.Flight.destination.name}</li>
                                                                    <li>City: {item.Flight.destination.address.city}</li>
                                                                    <li>Country {item.Flight.destination.address.country}</li>
                                                                    <li>Street: {item.Flight.destination.address.street}, {item.Flight.destination.address.streetNumber}</li>
                                                                </Box>
                                                            </StyledTableCell>
                                                            <StyledTableCell>
                                                                <li>Total: {item.Flight.numberOfSeats}</li>
                                                                <li>Vacant: {item.Flight.vacantSeats}</li>
                                                                <li>Price per person: {item.Flight.price}</li>
                                                                <li>Total price: {item.TotalPrice}</li>
                                                            </StyledTableCell>
                                                            <StyledTableCell>
                                                                <Box m={1}>
                                                                    <Button fullWidth color="success"
                                                                            variant="contained"
                                                                            onClick={() => {
                                                                                handleBuyTicketsClick(item)
                                                                            }}
                                                                    >
                                                                        Buy tickets
                                                                    </Button>
                                                                </Box>
                                                            </StyledTableCell>
                                                        </StyledTableRow>
                                                    </React.Fragment>
                                                )
                                            )}
                                        </TableBody>
                                    </Table>
                                </TableContainer>
                            </div>
                        )}


                    </Flex>
                </DialogContent>
                <DialogActions>
                    <Button variant="outlined"
                            color="warning"
                            onClick={handleFromDialogShowClose}>
                        Close
                    </Button>
                    <Button variant="contained"
                            color="info"
                            onClick={handleSearchFrom}>
                        Search
                    </Button>
                </DialogActions>
            </Dialog>

            <Dialog onClose={handleCloseErrorApiDialog} open={errorAPIDialogShow}>
                <DialogTitle>You can not buy ticket with that api key</DialogTitle>
                <DialogActions>
                    <Button onClick={handleCloseErrorApiDialog} variant="contained">Close</Button>
                </DialogActions>
            </Dialog>


            <Dialog onClose={handleCloseTickedDialogBought} open={ticketBoughtDialogShow}>
                <DialogTitle>Ticket bought successfully!</DialogTitle>
                <DialogActions>
                    <Button onClick={handleCloseTickedDialogBought} variant="contained">Close</Button>
                </DialogActions>
            </Dialog>

            <Dialog open={enterAPIkeyDialogShow} onClose={handleCloseBuyApiKeyDialog}>
                <DialogTitle>Enter your API key from your <b><i>FTN Airlines</i></b> account:</DialogTitle>
                <DialogContent>
                    <Box m={1}>
                        <TextField
                            label="API Key"
                            value={APIKey}
                            onChange={(e) => setAPIkey(e.target.value)}
                            fullWidth
                        />
                    </Box>
                </DialogContent>
                <DialogActions>
                    <Button variant="outlined"
                            color="warning"
                            fullWidth
                            onClick={handleCloseBuyApiKeyDialog}>
                        Close
                    </Button>
                    <Button variant="contained"
                            color="success"
                            fullWidth
                            onClick={handleBuyTickets}
                            disabled={APIKey === ""}
                    >
                        Buy tickets
                    </Button>
                </DialogActions>
            </Dialog>

            <Dialog open={toDialogShow} onClose={handleToDialogShowClose} fullWidth maxWidth="lg">
                <DialogTitle>Search for fligths to your accommodation:</DialogTitle>
                <DialogContent>
                    <Flex flexDirection="column">
                        <Flex m={1}>
                            <Box width={1 / 3} m={1}>
                                <TextField
                                    label="Start Point Country"
                                    name="startPointCountry"
                                    value={toInputData.startPointCountry}
                                    onChange={handleChangeTo}
                                    fullWidth

                                />
                            </Box>
                            <Box width={1 / 3} m={1}>
                                <TextField
                                    label="Start Point City"
                                    name="startPointCity"
                                    value={toInputData.startPointCity}
                                    onChange={handleChangeTo}
                                    fullWidth

                                />
                            </Box>
                            <Box width={1 / 3} m={1}>
                                <TextField
                                    label="Desired Number of Seats"
                                    name="desiredNumberOfSeats"
                                    type="number"
                                    value={toInputData.desiredNumberOfSeats}
                                    onChange={handleChangeTo}
                                    fullWidth
                                    InputProps={{
                                        inputProps: {
                                            min: 1
                                        }
                                    }}

                                />
                            </Box>
                        </Flex>
                        {flightResults != null && flightResults.length > 0 && (
                            <div>
                                <TableContainer component={Paper}
                                                sx={{maxHeight: 700, overflowY: 'scroll'}}>
                                    <Table>
                                        <TableBody>
                                            {flightResults.map((item, idx) =>
                                                (
                                                    <React.Fragment key={`${idx}-row`}>
                                                        <StyledTableRow>
                                                            <StyledTableCell>
                                                                <Box m={1} sx={{
                                                                    overflowX: 'auto',
                                                                    width: 200,
                                                                    height: 100,
                                                                    overflowy: 'auto'
                                                                }}>
                                                                    <li>Airport name: {item.Flight.startPoint.name}</li>
                                                                    <li>City: {item.Flight.startPoint.address.city}</li>
                                                                    <li>Country {item.Flight.startPoint.address.country}</li>
                                                                    <li>Street: {item.Flight.startPoint.address.street}, {item.Flight.startPoint.address.streetNumber}</li>
                                                                </Box>
                                                            </StyledTableCell>
                                                            <StyledTableCell>
                                                                <Box m={1} sx={{
                                                                    overflowX: 'auto',
                                                                    width: 200,
                                                                    height: 100,
                                                                    overflowy: 'auto'
                                                                }}>
                                                                    <li>Airport
                                                                        name: {item.Flight.destination.name}</li>
                                                                    <li>City: {item.Flight.destination.address.city}</li>
                                                                    <li>Country {item.Flight.destination.address.country}</li>
                                                                    <li>Street: {item.Flight.destination.address.street}, {item.Flight.destination.address.streetNumber}</li>
                                                                </Box>
                                                            </StyledTableCell>
                                                            <StyledTableCell>
                                                                <li>Total: {item.Flight.numberOfSeats}</li>
                                                                <li>Vacant: {item.Flight.vacantSeats}</li>
                                                                <li>Price per person: {item.Flight.price}</li>
                                                                <li>Total price: {item.TotalPrice}</li>
                                                            </StyledTableCell>
                                                            <StyledTableCell>
                                                                <Box m={1}>
                                                                    <Button fullWidth color="success"
                                                                            variant="contained"
                                                                            onClick={() => {
                                                                                handleBuyTicketsClick(item)
                                                                            }}
                                                                    >
                                                                        Buy tickets
                                                                    </Button>
                                                                </Box>
                                                            </StyledTableCell>
                                                        </StyledTableRow>
                                                    </React.Fragment>
                                                )
                                            )}
                                        </TableBody>
                                    </Table>
                                </TableContainer>
                            </div>
                        )}


                    </Flex>
                </DialogContent>
                <DialogActions>
                    <Button variant="outlined"
                            color="warning"
                            onClick={handleToDialogShowClose}>
                        Close
                    </Button>
                    <Button variant="contained"
                            color="info"
                            onClick={handleSearchTo}>
                        Search
                    </Button>
                </DialogActions>
            </Dialog>


            <div className="wrapper">
                <TableContainer component={Paper} sx={{maxHeight: 500, height: 500}}>
                    <Table>
                        <TableBody>
                            {reservations && reservations.length > 0 && reservations.map((r, idx) => (
                                <StyledTableRow key={idx} hover>
                                    <StyledTableCell>
                                        <Box m={1} sx={{
                                            width: 300,
                                            height: 200,
                                            overflowy: 'auto',
                                            overflowX: 'auto'
                                        }}>
                                            <li>Total price: {r.price}$</li>
                                            <li>Number of guests: {r.numberOfGuests}</li>
                                            <li>From: {new Date(r.dateRange.from).toLocaleDateString()}</li>
                                            <li>To: {new Date(r.dateRange.to).toLocaleDateString()}</li>
                                            <li>Status: {r.status}</li>
                                            <li>Accommodation name: {r.accommodationName}</li>
                                            <li>Location:</li>
                                            <li>{r.address.country}, {r.address.city} </li>
                                            <li>{r.address.street}, {r.address.streetNumber} </li>
                                        </Box>
                                    </StyledTableCell>
                                    <StyledTableCell>
                                        <Box m={1}>
                                            <Button
                                                onClick={() => {
                                                    handleCancel(r.Id)
                                                }}
                                                fullWidth
                                                disabled={r.status === 'rejected' || r.status === 'canceled' || new Date(r.dateRange.from) <= new Date(new Date().setDate(new Date().getDate() - 1))}
                                                variant="contained" color="warning">Cancel reservation</Button>
                                        </Box>
                                        <Box m={1}>
                                            <Button
                                                onClick={() => {
                                                    flightsToAccommodationClickHandle(r)
                                                }}
                                                fullWidth
                                                disabled={r.status === 'rejected' || r.status === 'canceled' || r.status === 'pending' || new Date(r.dateRange.to) <= new Date(new Date().setDate(new Date().getDate() - 1))}
                                                variant="outlined" color="info">Flights to accommodation</Button>
                                        </Box>
                                        <Box m={1}>
                                            <Button
                                                onClick={() => {
                                                    flightsFromAccommodationClickHandle(r)
                                                }}

                                                fullWidth
                                                disabled={r.status === 'rejected' || r.status === 'canceled' || r.status === 'pending' || new Date(r.dateRange.from) <= new Date(new Date().setDate(new Date().getDate() - 1))}
                                                variant="outlined" color="info">Flights from accommodation</Button>
                                        </Box>
                                    </StyledTableCell>
                                </StyledTableRow>
                            ))}
                        </TableBody>
                    </Table>
                </TableContainer>
                <Box m={1}>
                    Pending and accepted reservations can be canceled up to 24h before start date
                </Box>
            </div>
        </>
    );
}

export default MyReservations;