import React from 'react';
import {Flex} from "reflexbox";
import {Box, Paper, Table, TableContainer} from "@mui/material";

function ReservationsAndRequests() {
    return (
        <>
            <div className="wrapper">
                <Flex flexDirection="rows">
                    <Flex flexDirection="column" alignItems="center" m={2}>
                        <Box m={1}>
                            Reservations
                        </Box>
                        <TableContainer component={Paper} sx={{maxHeight: 500, height: 500}}>
                            <Table>

                            </Table>
                        </TableContainer>
                    </Flex>
                    <Flex flexDirection="column" alignItems="center" m={2}>
                        <Box m={1}>
                            Requests
                        </Box>

                        <TableContainer component={Paper} sx={{maxHeight: 500, height: 500}}>
                            <Table>

                            </Table>
                        </TableContainer>
                    </Flex>
                </Flex>
            </div>
        </>
    );
}

export default ReservationsAndRequests;