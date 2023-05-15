import React, {useEffect, useState} from 'react';
import {Flex} from "reflexbox";
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

function ReservationsAndRequests() {

    const [accepted, setAccepted] = useState(null);
    const [pending, setPending] = useState(null);

    const getReservations = () => {
        interceptor.get("api-1/reservation/pending").then(res => {
            setPending(res.data.reservation)
        }).catch(err => {
            console.log(err)
        })

    }
    const getRequests = () => {
        interceptor.get("api-1/reservation/accepted").then(res => {
            setAccepted(res.data.reservation)
        }).catch(err => {
            console.log(err)
        })
    }
    const getData = () => {
        getRequests()
        getReservations();
    }
    useEffect(() => {
        getData();
    }, []);


    const handleAccept = (Id) => {
        interceptor.get("api-1/reservation/accept/" + Id).then(res => {
            getData();
        }).catch(err => {
            console.log(err)
        })
    };
    const handleReject = (Id) => {
        interceptor.get("api-1/reservation/reject/" + Id).then(res => {
            getData();
        }).catch(err => {
            console.log(err)
        })

    };
    return (
        <>
            <div className="wrapper">
                <Flex flexDirection="rows" alignItems="center" justifyContent="center">
                    <Flex flexDirection="column" alignItems="center" m={2} w={1 / 2}>
                        <Box m={1}>
                            Reservations
                        </Box>
                        <TableContainer component={Paper} sx={{maxHeight: 500, height: 500}}>
                            <Table>
                                <TableBody>
                                    {accepted && accepted.length > 0 && accepted.map((a) => (
                                        <StyledTableRow key={a.Id} hover>
                                            <StyledTableCell>
                                                <li>Total price: {a.price}$</li>
                                                <li>Number of guests: {a.numberOfGuests}</li>
                                                <li>From: {new Date(a.dateRange.from * 1000).toLocaleDateString("en-GB")}</li>
                                                <li>To: {new Date(a.dateRange.to * 1000).toLocaleDateString("en-GB")}</li>
                                            </StyledTableCell>
                                        </StyledTableRow>
                                    ))}
                                </TableBody>
                            </Table>
                        </TableContainer>
                    </Flex>
                    <Flex flexDirection="column" alignItems="center" m={2} w={1 / 2}>
                        <Box m={1}>
                            Requests
                        </Box>

                        <TableContainer component={Paper} sx={{maxHeight: 500, height: 500}}>
                            <Table>
                                <TableBody>
                                    {pending && pending.length > 0 && pending.map((a) => (
                                        <StyledTableRow key={a.Id} hover>
                                            <StyledTableCell>
                                                <li>Number of past cancellations by this
                                                    person: {a.numberOfCancellations}</li>
                                                <li>Total price: {a.price}$</li>
                                                <li>Number of guests: {a.numberOfGuests}</li>
                                                <li>From: {new Date(a.dateRange.from * 1000).toLocaleDateString("en-GB")}</li>
                                                <li>To: {new Date(a.dateRange.to * 1000).toLocaleDateString("en-GB")}</li>
                                            </StyledTableCell>

                                            <StyledTableCell>
                                                <Box m={1}>
                                                    <Button fullWidth variant="contained"
                                                            color="warning"
                                                            onClick={() => {
                                                                handleAccept(a.Id)
                                                            }}
                                                    >Accept
                                                    </Button>
                                                </Box>
                                                <Box m={1}>
                                                    <Button fullWidth variant="contained"
                                                            onClick={() => {
                                                                handleReject(a.Id)
                                                            }}
                                                            color="error"
                                                    >Reject
                                                    </Button>
                                                </Box>
                                            </StyledTableCell>
                                        </StyledTableRow>
                                    ))}
                                </TableBody>
                            </Table>
                        </TableContainer>
                    </Flex>
                </Flex>
            </div>
        </>
    );
}

export default ReservationsAndRequests;