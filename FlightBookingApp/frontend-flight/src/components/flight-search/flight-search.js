import React, {useState} from "react";
import {useNavigate} from 'react-router-dom';
import axios from "axios";
import "./flight-search.css";
import {DatePicker, LocalizationProvider} from "@mui/x-date-pickers";
import {AdapterDayjs} from "@mui/x-date-pickers/AdapterDayjs";
import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';
import SearchSharpIcon from '@mui/icons-material/SearchSharp';
import {styled} from '@mui/material/styles';
import TableCell, {tableCellClasses} from '@mui/material/TableCell';
import TableRow from '@mui/material/TableRow';
import {Paper, Table, TableBody, TableContainer, TableHead} from "@mui/material";
import moment from "moment";
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';
import ArrowDownwardIcon from '@mui/icons-material/ArrowDownward';
import ArrowUpwardIcon from '@mui/icons-material/ArrowUpward';
import ImportExportIcon from '@mui/icons-material/ImportExport';
import dayjs from "dayjs";

import useAxiosPrivate from "../../hooks/useAxiosPrivate";
import useAuth from "../../hooks/useAuth";

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
const FlightSearch = ({LoggedIn}) => {

    const [data, setData] = useState([]);
    const [entityCount, setEntityCount] = useState(0);
    const axiosPrivate = useAxiosPrivate();

    const [searchParams, setSearchParams] = useState({
        departureDate: dayjs((new Date())),
        destinationCountry: "",
        destinationCity: "",
        startPointCountry: "",
        startPointCity: "",
        desiredNumberOfSeats: "1",
    });
    const [pagination, setPagination] = useState({
        pageNumber: 1,
        resultsPerPage: 4,
        sortDirection: "asc",
        sortType: "departureDateTime",
    });
    const [selectDesiredNumberOfSeats, setSelectDesiredNumberOfSeats] = React.useState(1);
    const [openNoFlightsDialog, setOpenNoFlightsDialog] = React.useState(false);
    const [buyTicketsDialog, setBuyTicketsDialog] = React.useState(false);
    const [selectedFlight, setSelectedFlight] = React.useState({});
    const [purchaseDialog, setPurchaseDialog] = React.useState(false);


    const fetchData = async () => {


        try {
            const {data} = await axios.get(process.env.REACT_APP_FLIGHT_APP_API + "search-flights", {
                params: {
                    ...searchParams,
                    departureDate: searchParams.departureDate.toISOString().substring(0, 10), // convert to ISO string format
                    ...pagination,
                },
            });
            if (data.Data != null)
                setData(data.Data);
            else {
                setOpenNoFlightsDialog(true);
                setData([])
            }
            setEntityCount(data.EntityCount)
        } catch (e) {
            alert("Unexpected error")
        }
    };

    const handleSearchParamsChange = (e) => {
        setSearchParams({...searchParams, [e.target.name]: e.target.value});
    };
    const handleSelectDesiredNumberOfSeatsChange = (e) => {
        setSelectDesiredNumberOfSeats(parseInt(e.target.value ?? 0));
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
        fetchData();

    };


    const handleCloseNoFlightsDialog = () => {
        setOpenNoFlightsDialog(false);
    };

    const handleOpenBuyTicketsDialog = (flight) => {

        setSelectedFlight(flight)
        setBuyTicketsDialog(true);

    };
    const handleCloseBuyTicketsDialog = () => {
        setBuyTicketsDialog(false);
        setSelectedFlight({})
    };

    const buyTickets = async () => {
        await axiosPrivate.post("/api/ticket/buy", {
            numberOfTickets: selectDesiredNumberOfSeats,
            flightId: selectedFlight.id,
        }).then((res) => {

            setPurchaseDialog(true);
            handleCloseBuyTicketsDialog();
        })
            .catch(err => {
                console.error(err)
                alert("unexpected error")
            });
    };
    const navigate = useNavigate();
    const {auth} = useAuth()


    return (
        <div className="flight-search">
            <div className="searchParams">
                <table>
                    <tbody>
                    <tr>
                        <td>
                            <TextField
                                fullWidth
                                variant="filled"
                                label="Start Point Country"
                                type="text"
                                name="startPointCountry"
                                value={searchParams.startPointCountry}
                                onChange={handleSearchParamsChange}
                            />
                        </td>
                        <td>
                            <TextField
                                fullWidth
                                variant="filled"
                                label="Destination country"
                                type="text"
                                name="destinationCountry"
                                value={searchParams.destinationCountry}
                                onChange={handleSearchParamsChange}
                            />
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <TextField
                                fullWidth
                                variant="filled"
                                label="Start Point city"
                                type="text"
                                name="startPointCity"
                                value={searchParams.startPointCity}
                                onChange={handleSearchParamsChange}
                            />
                        </td>
                        <td>
                            <TextField
                                fullWidth
                                variant="filled"
                                label="Destination city"
                                type="text"
                                name="destinationCity"
                                value={searchParams.destinationCity}
                                onChange={handleSearchParamsChange}
                            />
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <TextField type="number"
                                       fullWidth
                                       variant="filled"
                                       InputProps={{
                                           inputProps: {
                                               min: 1
                                           }
                                       }}
                                       value={searchParams.desiredNumberOfSeats}
                                       name="desiredNumberOfSeats"
                                       label="Desired number of seats:"
                                       onChange={handleSearchParamsChange}/>
                        </td>
                        <td>
                            <LocalizationProvider dateAdapter={AdapterDayjs}>
                                <DatePicker label="Departure date"
                                            value={searchParams.departureDate}
                                            minDate={dayjs((new Date()))}
                                            onChange={
                                                (newValue) => {
                                                    const result = new Date(newValue);
                                                    setSearchParams({
                                                            ...searchParams,
                                                            departureDate: newValue
                                                        }
                                                    )
                                                }
                                            }
                                />
                            </LocalizationProvider>
                        </td>

                    </tr>
                    </tbody>
                </table>
                <Button
                    variant="contained" endIcon={<SearchSharpIcon/>}
                    onClick={() => {
                        setPagination({...pagination, pageNumber: 1})
                        fetchData();
                    }}>Search
                </Button>
                <Dialog
                    open={openNoFlightsDialog}
                    onClose={handleCloseNoFlightsDialog}
                >
                    <DialogTitle id="alert-dialog-title">
                        {"No flights with this parameters found"}
                    </DialogTitle>
                    <DialogContent>
                        <DialogContentText id="alert-dialog-description">
                            Try new parameters!
                        </DialogContentText>
                    </DialogContent>
                    <DialogActions>
                        <Button onClick={handleCloseNoFlightsDialog}>Close</Button>
                    </DialogActions>
                </Dialog>
            </div>
            {entityCount > 0 &&
                <div className="flightResults">
                    <div className="pageButtons">
                        <Button
                            variant="contained" onClick={() => handlePaginationChange("prev")}
                            disabled={pagination.pageNumber < 2}>
                            Prev
                        </Button>
                        <span> Page {pagination.pageNumber} of {Math.ceil(entityCount / pagination.resultsPerPage)} </span>
                        <Button
                            variant="contained" onClick={() => handlePaginationChange("next")}
                            disabled={pagination.pageNumber >= entityCount / pagination.resultsPerPage}>
                            Next
                        </Button>
                    </div>
                    <TableContainer component={Paper}>
                        <Table>
                            <TableHead>
                                <StyledTableRow>
                                    <StyledTableCell className="cursor" align="center" style={{width: "20%"}}
                                                     onClick={() => {
                                                         setPagination({
                                                             ...pagination,
                                                             sortDirection: pagination.sortDirection === "asc" ? "dsc" : "asc",
                                                             sortType: "departureDateTime"
                                                         });
                                                     }
                                                     }>Departure
                                        Time {pagination.sortDirection === "asc" && pagination.sortType === "departureDateTime" &&
                                            <ArrowDownwardIcon> </ArrowDownwardIcon>
                                        }
                                        {pagination.sortDirection === "dsc" && pagination.sortType === "departureDateTime" &&
                                            <ArrowUpwardIcon> </ArrowUpwardIcon>
                                        }
                                        {pagination.sortType === "price" &&
                                            <ImportExportIcon></ImportExportIcon>
                                        }</StyledTableCell>
                                    <StyledTableCell align="center" style={{width: "20%"}}>Point of
                                        departure</StyledTableCell>
                                    <StyledTableCell align="center" style={{width: "20%"}}>Destination</StyledTableCell>
                                    <StyledTableCell align="center" style={{width: "5%"}}>Seats</StyledTableCell>
                                    <StyledTableCell sortable="down" align="center" style={{width: "20%"}}
                                                     className="cursor"
                                                     onClick={() => {
                                                         setPagination({
                                                             ...pagination,
                                                             sortDirection: pagination.sortDirection === "asc" ? "dsc" : "asc",
                                                             sortType: "price"
                                                         });
                                                     }}
                                    >Ticket price
                                        {pagination.sortDirection === "asc" && pagination.sortType === "price" &&
                                            <ArrowDownwardIcon> </ArrowDownwardIcon>
                                        }
                                        {pagination.sortDirection === "dsc" && pagination.sortType === "price" &&
                                            <ArrowUpwardIcon> </ArrowUpwardIcon>
                                        }
                                        {pagination.sortType === "departureDateTime" &&
                                            <ImportExportIcon></ImportExportIcon>
                                        }


                                    </StyledTableCell>
                                    <StyledTableCell align="center" style={{width: "10%"}}>Total price</StyledTableCell>
                                    {LoggedIn &&

                                        <StyledTableCell align="center" style={{width: "5%"}}>Buy
                                            tickets</StyledTableCell>

                                    }
                                </StyledTableRow>
                            </TableHead>
                            <TableBody>
                                {data.map((item, i) => (
                                    <StyledTableRow hover key={i}>
                                        <StyledTableCell
                                            align="center"
                                        > {moment(item.Flight.departureDateTime).format("MM.DD.YYYY HH:mm")}{" "}</StyledTableCell>
                                        <StyledTableCell align="center">
                                            <li>Airport name: {item.Flight.startPoint.name}</li>
                                            <li>City: {item.Flight.startPoint.address.city}</li>
                                            <li>Country {item.Flight.startPoint.address.country}</li>
                                            <li>Street: {item.Flight.startPoint.address.street}, {item.Flight.startPoint.address.streetNumber}</li>
                                        </StyledTableCell>
                                        <StyledTableCell align="center">
                                            <li>Airport name: {item.Flight.destination.name}</li>
                                            <li>City: {item.Flight.destination.address.city}</li>
                                            <li>Country {item.Flight.destination.address.country}</li>
                                            <li>Street: {item.Flight.destination.address.street}, {item.Flight.destination.address.streetNumber}</li>
                                        </StyledTableCell>
                                        <StyledTableCell align="center">
                                            <li>Total: {item.Flight.numberOfSeats}</li>
                                            <li>Vacant: {item.Flight.vacantSeats}</li>
                                        </StyledTableCell>
                                        <StyledTableCell align="center">{item.Flight.price}</StyledTableCell>
                                        <StyledTableCell align="center">{item.TotalPrice}</StyledTableCell>
                                        {LoggedIn &&
                                            <StyledTableCell align="center"> <Button
                                                onClick={() => handleOpenBuyTicketsDialog(item.Flight)}
                                                variant="contained">
                                                Buy now
                                            </Button>

                                            </StyledTableCell>}
                                    </StyledTableRow>
                                ))}
                            </TableBody>
                        </Table>
                    </TableContainer>
                </div>
            }
            <Dialog
                open={buyTicketsDialog}
                onClose={handleCloseBuyTicketsDialog}>
                <DialogTitle id="alert-dialog-title">
                    {"Select the number of tickets you want to buy"}
                </DialogTitle>
                <DialogContent>
                    <TextField type="number"
                               fullWidth
                               variant="filled"
                               InputProps={{
                                   inputProps: {
                                       min: 1,
                                       max: selectedFlight.vacantSeats
                                   }
                               }}
                               value={selectDesiredNumberOfSeats}
                               label="Desired number of seats:"
                               onChange={handleSelectDesiredNumberOfSeatsChange}/>
                </DialogContent>
                <DialogActions>
                    <Button onClick={handleCloseBuyTicketsDialog}>Close</Button>
                    <Button onClick={buyTickets}>Buy</Button>
                </DialogActions>
            </Dialog>
            <Dialog
                open={purchaseDialog}
                onClose={() => {
                    setPurchaseDialog(false)
                }}>
                <DialogTitle id="alert-dialog-title">
                    {"Successful purchase"}
                </DialogTitle>
                <DialogActions>
                    <Button onClick={() => {
                        setPurchaseDialog(false)
                        navigate('/bought-tickets');
                    }}>Close</Button>
                </DialogActions>
            </Dialog>
        </div>
    );
};
export default FlightSearch;
