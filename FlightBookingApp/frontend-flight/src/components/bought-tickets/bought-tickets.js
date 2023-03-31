import React, {useEffect, useState} from 'react';
import {Paper, Table, TableBody, TableContainer, TableHead} from "@mui/material";
import TableRow from "@mui/material/TableRow";
import {styled} from "@mui/material/styles";
import TableCell, {tableCellClasses} from "@mui/material/TableCell";
import useAxiosPrivate from "../../hooks/useAxiosPrivate";
import moment from "moment";

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


function BoughtTickets() {
    const [data, setData] = useState(null);

    const axiosPrivate = useAxiosPrivate();
    const fetchData = async () => {
        await axiosPrivate.get("/api/ticket/myTickets").then(res => {
            console.log(res);
            setData(res.data);
        }).catch(e => {
            alert("Unexpected error")
        });
    }

    useEffect(() => {
        fetchData();
    }, [])

    return (
        <div><TableContainer className="all-flights" component={Paper} sx={{maxHeight: 600}}>
            <Table stickyHeader>
                <TableHead>
                    <TableRow>
                        <StyledTableCell align="center" style={{width: "10%"}}>Ticked ID</StyledTableCell>
                        <StyledTableCell align="center" style={{width: "10%"}}>Departure Time</StyledTableCell>
                        <StyledTableCell align="center" style={{width: "25%"}}>Point of departure</StyledTableCell>
                        <StyledTableCell align="center" style={{width: "25%"}}>Destination</StyledTableCell>
                        <StyledTableCell align="center" style={{width: "10%"}}>Seats</StyledTableCell>
                        <StyledTableCell align="center" style={{width: "10%"}}>Price</StyledTableCell>
                    </TableRow>
                </TableHead>
                {
                    data != null &&
                    <TableBody>
                        {data.map((item, i) => (
                            <StyledTableRow hover key={i}>
                                <StyledTableCell align="center">{item.id}</StyledTableCell>
                                <StyledTableCell
                                    align="center"
                                > {moment(item.flightInfo.departureDateTime).format("MM.DD.YYYY HH:mm")}{" "}</StyledTableCell>

                                <StyledTableCell align="center">
                                    <li>Airport name: {item.flightInfo.destination.name}</li>
                                    <li>City: {item.flightInfo.destination.address.city}</li>
                                    <li>Country {item.flightInfo.destination.address.country}</li>
                                    <li>Street: {item.flightInfo.destination.address.street}, {item.flightInfo.destination.address.streetNumber}</li>
                                </StyledTableCell>

                                <StyledTableCell align="center">
                                    <li>Airport name: {item.flightInfo.startPoint.name}</li>
                                    <li>City: {item.flightInfo.startPoint.address.city}</li>
                                    <li>Country {item.flightInfo.startPoint.address.country}</li>
                                    <li>Street: {item.flightInfo.startPoint.address.street}, {item.flightInfo.startPoint.address.streetNumber}</li>
                                </StyledTableCell>
                                <StyledTableCell align="center">
                                    <li>Total: {item.flightInfo.numberOfSeats}</li>
                                    <li>Vacant: {item.flightInfo.vacantSeats}</li>
                                </StyledTableCell>

                                <StyledTableCell align="center">{item.flightInfo.price}</StyledTableCell>


                            </StyledTableRow>))}
                    </TableBody>
                }
            </Table>
        </TableContainer></div>
    );
}

export default BoughtTickets;