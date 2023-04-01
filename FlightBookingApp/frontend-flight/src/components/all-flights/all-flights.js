import React, {useEffect, useState} from 'react';
import {Paper, Table, TableBody, TableContainer, TableHead} from "@mui/material";
import TableRow from "@mui/material/TableRow";
import {styled} from "@mui/material/styles";
import TableCell, {tableCellClasses} from "@mui/material/TableCell";
import moment from "moment/moment";
import Button from "@mui/material/Button";
import "./all-flights.css"
import Dialog from '@mui/material/Dialog';
import DialogTitle from "@mui/material/DialogTitle";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import DialogContentText from "@mui/material/DialogContentText";
import useAxiosPrivate from "../../hooks/useAxiosPrivate";

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


function AllFlights() {
    const [data, setData] = useState(null);
    const [cancelFlightDialog, setCancelFlightDialog] = React.useState(false);
    const [successfulCancelFlightDialog, setSuccessfulCancelFlightDialog] = React.useState(false);
    const [flightId, setFlightId] = React.useState(false);
    const axiosPrivate = useAxiosPrivate();
    const openCancelFlightDialog = (flightID) => {
        setCancelFlightDialog(true);
        setFlightId(flightID)
    };
    const closeCancelFlightDialog = () => {
        setCancelFlightDialog(false);
    };
    const openSuccessfulCancelFlightDialog = () => {
        setSuccessfulCancelFlightDialog(true);
    };
    const closeSuccessfulCancelFlightDialog = () => {
        setSuccessfulCancelFlightDialog(false);
        fetchData();
    };


    const fetchData = async () => {
        try {
            const {data} = await axiosPrivate.get("/api/flight");
            setData(data)

        } catch (e) {
            alert("Unexpected error")
        }
    }

    const cancelFlight = async () => {
        try {
            const {res} = await axiosPrivate.patch("/api/flight/" + flightId);
            openSuccessfulCancelFlightDialog()

        } catch (e) {
            alert("Unexpected error")
        }
    }

    useEffect(() => {
        fetchData();
    }, []);


    return (<div>
            <Dialog
                onClose={closeCancelFlightDialog} open={cancelFlightDialog}>
                <DialogTitle id="alert-dialog-title1">
                    {"Are you sure you want to cancel this flight?"}
                </DialogTitle>
                <DialogContent>
                    <DialogContentText id="alert-dialog-description">

                    </DialogContentText>
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => {
                        closeCancelFlightDialog();
                        cancelFlight();
                    }
                    }>Cancel flight </Button>
                    <Button onClick={closeCancelFlightDialog}>Close</Button>
                </DialogActions>
            </Dialog>

            <Dialog
                onClose={closeSuccessfulCancelFlightDialog} open={successfulCancelFlightDialog}>
                <DialogTitle id="alert-dialog-title">
                    {"Flight canceled"}
                </DialogTitle>
                <DialogActions>

                    <Button onClick={closeSuccessfulCancelFlightDialog}>Close</Button>
                </DialogActions>
            </Dialog>

            <TableContainer className="all-flights" component={Paper} sx={{maxHeight: 600}}>
                <Table stickyHeader>
                    <TableHead>
                        <TableRow>
                            <StyledTableCell align="center" style={{width: "10%"}}>Departure Time</StyledTableCell>
                            <StyledTableCell align="center" style={{width: "20%"}}>Point of departure</StyledTableCell>
                            <StyledTableCell align="center" style={{width: "20%"}}>Destination</StyledTableCell>
                            <StyledTableCell align="center" style={{width: "10%"}}>Seats</StyledTableCell>
                            <StyledTableCell align="center" style={{width: "20%"}}>Price</StyledTableCell>
                            <StyledTableCell align="center" style={{width: "20%"}}>Cancel Flight</StyledTableCell>
                        </TableRow>
                    </TableHead>
                    {
                        data != null &&
                        <TableBody>
                            {data.map((item, i) => (
                                <StyledTableRow hover key={i}>
                                    <StyledTableCell
                                        align="center"
                                    > {moment(item.departureDateTime).format("MM.DD.YYYY HH:mm")}{" "}</StyledTableCell>
                                    <StyledTableCell align="center">
                                        <li>Airport name: {item.startPoint.name}</li>
                                        <li>City: {item.startPoint.address.city}</li>
                                        <li>Country {item.startPoint.address.country}</li>
                                        <li>Street: {item.startPoint.address.street}, {item.startPoint.address.streetNumber}</li>
                                    </StyledTableCell>
                                    <StyledTableCell align="center" style={{width: "35%"}}>
                                        <li>Airport name: {item.destination.name}</li>
                                        <li>City: {item.destination.address.city}</li>
                                        <li>Country {item.destination.address.country}</li>
                                        <li>Street: {item.destination.address.street}, {item.destination.address.streetNumber}</li>
                                    </StyledTableCell>
                                    <StyledTableCell align="center">
                                        <li>Total: {item.numberOfSeats}</li>
                                        <li>Vacant: {item.vacantSeats}</li>
                                    </StyledTableCell>
                                    <StyledTableCell align="center">{item.price}</StyledTableCell>
                                    <StyledTableCell align="center">
                                        {
                                            item.canceled === false &&
                                            <Button variant="outlined" color="error"
                                                    onClick={() => openCancelFlightDialog(item.id)}>Cancel
                                                now</Button>
                                        }
                                        {
                                            item.canceled === true &&
                                            <span>CANCELED</span>
                                        }
                                    </StyledTableCell>
                                </StyledTableRow>))}
                        </TableBody>
                    }
                </Table>
            </TableContainer>
        </div>
    );
}

export default AllFlights;