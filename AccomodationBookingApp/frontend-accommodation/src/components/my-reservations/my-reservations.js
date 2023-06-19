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
        interceptor.get("api-1/reservation/all/guest").then(res => {
            setReservations(res.data.reservations)
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
                            {reservations && reservations.length > 0 && reservations.map((r) => (
                                <StyledTableRow key={r.Id} hover>
                                    <StyledTableCell>
                                        <li>Total price: {r.price}$</li>
                                        <li>Number of guests: {r.numberOfGuests}</li>
                                        <li>From: {new Date(r.dateRange.from * 1000).toLocaleDateString("en-GB")}</li>
                                        <li>To: {new Date(r.dateRange.to * 1000).toLocaleDateString("en-GB")}</li>
                                        <li>Status: {r.status}</li>
                                    </StyledTableCell>
                                    <StyledTableCell>
                                        <Box m={1}>
                                            <Button
                                                onClick={() => {
                                                    handleCancel(r.Id)
                                                }}
                                                disabled={r.status === 'rejected' || r.status === 'canceled' || new Date(r.dateRange.from * 1000).setHours(0, 0, 0, 0) <= new Date(new Date().getTime() + 1 * 24 * 60 * 60 * 1000).setHours(0, 0, 0, 0)}
                                                variant="contained" color="warning">Cancel reservation</Button>
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