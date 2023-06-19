import React, {useEffect, useState} from 'react';
import {
    Box,
    Button,
    Paper,
    styled,
    Table,
    TableBody,
    TableCell,
    tableCellClasses,
    TableContainer,
    TableRow
} from "@mui/material";
import interceptor from "../../interceptor/interceptor";

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
    const [reservations, setReservations] = useState(false);

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

    return (
        <>
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
                                                disabled={r.status === 'rejected' || r.status === 'canceled' || new Date(r.dateRange.from * 1000).setHours(0, 0, 0, 0) <= new Date(new Date().getTime() + 1 * 24 * 60 * 60 * 1000).setHours(0, 0, 0, 0)}
                                                variant="contained" color="warning">Cancel reservation</Button>
                                        </Box>
                                        <Box m={1}>
                                            <Button
                                                onClick={() => {
                                                    handleCancel(r.Id)
                                                }}
                                                fullWidth
                                                disabled={r.status === 'rejected' || r.status === 'canceled' || new Date(r.dateRange.from * 1000).setHours(0, 0, 0, 0) <= new Date(new Date().getTime() + 1 * 24 * 60 * 60 * 1000).setHours(0, 0, 0, 0)}
                                                variant="outlined" color="info">Flights from accommodation</Button>
                                        </Box>
                                        <Box m={1}>
                                            <Button
                                                onClick={() => {
                                                    handleCancel(r.Id)
                                                }}
                                                fullWidth
                                                disabled={r.status === 'rejected' || r.status === 'canceled' || new Date(r.dateRange.from * 1000).setHours(0, 0, 0, 0) <= new Date(new Date().getTime() + 1 * 24 * 60 * 60 * 1000).setHours(0, 0, 0, 0)}
                                                variant="outlined" color="info">Flights to accommodation</Button>
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