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

    function parseObjectId(str) {
        const regex = /^ObjectID\("(.+)"\)$/;
        const match = str.match(regex);
        return match ? match[1] : null;
    }


    useEffect(() => {
        getAllReservations();
    }, []);
    const handleCancel = (Id) => {
        interceptor.get("api-1/reservation/cancel/" + parseObjectId(Id)).then(res => {
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
                            {reservations && reservations.map((r) => (
                                <StyledTableRow key={r.Id} hover>
                                    <StyledTableCell>
                                        <li>Price: {r.price}</li>
                                        <li>Total price:</li>
                                        <li>Number of guests: {r.numberOfGuests}</li>
                                        <li>From: {new Date(r.dateRange.from * 1000).toLocaleDateString("en-GB")}</li>
                                        <li>To: {new Date(r.dateRange.to * 1000).toLocaleDateString("en-GB")}</li>
                                        <li>Reservation made by:</li>
                                        <li>Reserved place:</li>
                                        <li>Status: {r.status}</li>
                                    </StyledTableCell>
                                    <StyledTableCell>
                                        <Box m={1}>
                                            <Button
                                                onClick={() => {
                                                    handleCancel(r.Id)
                                                }}
                                                disabled={r.status === 'rejected' || r.status === 'canceled'}
                                                variant="contained" color="warning">Cancel reservation</Button>
                                        </Box>
                                    </StyledTableCell>
                                </StyledTableRow>
                            ))}
                        </TableBody>
                    </Table>
                </TableContainer>
            </div>
        </>
    );
}

export default MyReservations;