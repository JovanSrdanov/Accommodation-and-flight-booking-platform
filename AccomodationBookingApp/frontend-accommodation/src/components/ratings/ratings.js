import React, {useEffect, useState} from 'react';
import interceptor from "../../interceptor/interceptor";
import {Flex} from "reflexbox";
import {
    Accordion,
    AccordionDetails,
    AccordionSummary,
    Box,
    Button,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    Paper,
    Rating,
    styled,
    Table,
    TableBody,
    TableCell,
    tableCellClasses,
    TableContainer,
    TableRow
} from "@mui/material";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import ImageList from "@mui/material/ImageList";
import ImageListItem from "@mui/material/ImageListItem";

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

function Ratings(props) {
    const [ratableAccommodations, setRatableAccommodations] = useState(null);
    const [ratableHosts, setRatableHosts] = useState(null);
    const [selectedHost, setSelectedHost] = useState(null);
    const [selectedAccommodation, setSelectedAccommodation] = useState(null);
    const [rating, setRating] = React.useState(5);
    const [rateHostDialogShow, setRateHostDialogShow] = useState(false);
    const [rateAccommodationDialogShow, setRateAccommodationDialogShow] = useState(false);


    const getRatableHosts = () => {
        interceptor.get("api-2/accommodation/ratable/hosts").then((res) => {
            setRatableHosts(res.data)
        }).catch((err) => {
            console.log(err)
        })
    }

    const getRatableAccommodations = () => {
        interceptor.get("api-2/accommodation/ratable/accommodations").then((res) => {
            setRatableAccommodations(res.data)
        }).catch((err) => {
            console.log(err)
        })
    }

    useEffect(() => {
        getRatableAccommodations();
        getRatableHosts();
    }, []);


    const handleRateHostClick = (item) => {
        setSelectedHost(item)
        setRateHostDialogShow(true);
    };
    const handleRateHostDialogClose = () => {
        setRateHostDialogShow(false)
        setRating(5);
        setSelectedHost(null)
    };


    const handleRateHost = () => {
        var sendData = {
            rating: {
                hostId: selectedHost.hostId,
                rating: rating
            }
        }
        interceptor.post("api-1/rating/host", sendData).then((res) => {

            getRatableAccommodations();
            getRatableHosts();
            handleRateHostDialogClose()

        }).catch((err) => {
            console.log(err)
        })
    };
    const handleRateAccommodationDialogClose = () => {
        setRateAccommodationDialogShow(false)
        setRating(5);
        setSelectedAccommodation(null)
    };
    const handleRateAccommodation = () => {

        var sendData = {
            rating: {
                accommodationId: selectedAccommodation.id,
                rating: rating
            }
        }
        interceptor.post("api-1/rating/accommodation", sendData).then((res) => {
            getRatableAccommodations();
            handleRateAccommodationDialogClose()

        }).catch((err) => {
            console.log(err)
        })
    };
    const handleRateAccommodationClick = (item) => {
        setSelectedAccommodation(item)
        setRateAccommodationDialogShow(true)
    };
    const handleDeleteHostRating = (item) => {
        interceptor.delete("api-1/rating/host/" + item.hostId).then((res) => {
            getRatableHosts();
        }).catch((err) => {
            console.log(err)
        })

    };
    const handleDeleteAccommodationRating = (item) => {
        interceptor.delete("api-1/rating/accommodation/" + item.id).then((res) => {
            getRatableAccommodations();
        }).catch((err) => {
            console.log(err)
        })
    };
    return (
        <>
            <Dialog onClose={handleRateAccommodationDialogClose} open={rateAccommodationDialogShow}>
                {selectedAccommodation &&
                    (<DialogTitle>Rate the accommodation: {selectedAccommodation.name}
                    </DialogTitle>)}
                <DialogContent>
                    <Flex flexDirection="row" justifyContent="center" alignItems="center">
                        <Rating
                            size="large"
                            name="simple-controlled"
                            value={rating}
                            onChange={(event, newValue) => {
                                setRating(newValue);
                            }}
                        />
                    </Flex>
                </DialogContent>
                <DialogActions>
                    <Button
                        fullWidth
                        onClick={handleRateAccommodationDialogClose}
                        variant="contained"
                    >
                        Close
                    </Button>
                    <Button
                        fullWidth
                        variant="contained" color="success"
                        onClick={handleRateAccommodation}
                    >Rate</Button>
                </DialogActions>
            </Dialog>


            <Dialog onClose={handleRateHostDialogClose} open={rateHostDialogShow}>
                {selectedHost &&
                    (<DialogTitle>Rate the host: {selectedHost.username}
                    </DialogTitle>)}
                <DialogContent>
                    <Flex flexDirection="row" justifyContent="center" alignItems="center">
                        <Rating
                            size="large"
                            name="simple-controlled"
                            value={rating}
                            onChange={(event, newValue) => {
                                setRating(newValue);
                            }}
                        />
                    </Flex>
                </DialogContent>
                <DialogActions>
                    <Button
                        fullWidth
                        onClick={handleRateHostDialogClose}
                        variant="contained"
                    >
                        Close
                    </Button>
                    <Button
                        fullWidth
                        variant="contained" color="success"
                        onClick={handleRateHost}
                    >Rate</Button>
                </DialogActions>
            </Dialog>

            <div className="wrapper">
                <Flex flexDirection="rows">
                    <Flex flexDirection="column" alignItems="center" m={2}>
                        <Box m={1}>
                            Hosts
                        </Box>
                        {ratableHosts != null && ratableHosts.length > 0 && (
                            <div>
                                <TableContainer component={Paper}
                                                sx={{maxHeight: 700, overflowY: 'scroll'}}>
                                    <Table>
                                        <TableBody>
                                            {ratableHosts.map((item, idx) =>
                                                (
                                                    <React.Fragment key={`${idx}-row`}>
                                                        <StyledTableRow>
                                                            <StyledTableCell>
                                                                <Box m={1} sx={{
                                                                    overflowX: 'auto',
                                                                    width: 250,
                                                                    height: 100,
                                                                    overflowy: 'auto'
                                                                }}>
                                                                    <li>Username: {item.username}</li>
                                                                    <li>Email: {item.email}</li>


                                                                    <li>Surname: {item.surname}</li>
                                                                    <li>Name: {item.name}</li>
                                                                    <li>Your
                                                                        rating: {item.rating === 0 ? 'NOT YET RATED' : item.rating}</li>
                                                                </Box>
                                                            </StyledTableCell>
                                                            <StyledTableCell>
                                                                <Box m={1} sx={{
                                                                    overflowX: 'auto',
                                                                    width: 125,
                                                                    height: 125,
                                                                    overflowy: 'auto'
                                                                }}>
                                                                    <Box m={1}>
                                                                        <Button
                                                                            fullWidth
                                                                            color="success"
                                                                            variant="contained"
                                                                            onClick={() => {
                                                                                handleRateHostClick(item)
                                                                            }}
                                                                        >
                                                                            Rate
                                                                        </Button>
                                                                    </Box>
                                                                    {item.rating !== 0 && (
                                                                        <Box m={1}>
                                                                            <Button
                                                                                fullWidth
                                                                                color="error"
                                                                                variant="outlined"
                                                                                onClick={() => {
                                                                                    handleDeleteHostRating(item)
                                                                                }}
                                                                            >
                                                                                Delete rating
                                                                            </Button>
                                                                        </Box>
                                                                    )}
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
                    <Flex flexDirection="column" alignItems="center" m={2}>
                        <Box m={1}>
                            Accommodations
                        </Box>
                        {ratableAccommodations != null && ratableAccommodations.length > 0 && (
                            <div>
                                <TableContainer component={Paper}
                                                sx={{maxHeight: 700, overflowY: 'scroll'}}>
                                    <Table>
                                        <TableBody>
                                            {ratableAccommodations.map((item, idx) =>
                                                (
                                                    <React.Fragment key={`${idx}-row`}>
                                                        <StyledTableRow>
                                                            <StyledTableCell>
                                                                <Box m={1} sx={{
                                                                    width: 250,
                                                                    height: 125,
                                                                    overflowy: 'auto',
                                                                    overflowX: 'auto'
                                                                }}>
                                                                    <li>Name: {item.name}</li>
                                                                    <li>Minimum number of guests: {item.minGuests}</li>
                                                                    <li>Maximum number of guests: {item.maxGuests}</li>
                                                                    <li>{item.address.city}, {item.address.country}</li>
                                                                    <li>{item.address.street}, {item.address.streetNumber}</li>
                                                                    <li>Your
                                                                        rating: {item.rating === 0 ? 'NOT YET RATED' : item.rating}</li>
                                                                </Box>
                                                            </StyledTableCell>
                                                            <StyledTableCell>
                                                                <Box m={1}>
                                                                    <Accordion sx={{border: "2px solid black"}}>
                                                                        <AccordionSummary
                                                                            expandIcon={<ExpandMoreIcon/>}>
                                                                            List Of Ameneties
                                                                        </AccordionSummary>
                                                                        <AccordionDetails
                                                                            sx={{height: 200, overflowY: 'scroll'}}>
                                                                            {item.amenities.map((a) => (
                                                                                <Box m={1} key={a}>- {a}</Box>
                                                                            ))}
                                                                        </AccordionDetails>
                                                                    </Accordion>

                                                                    <Accordion sx={{border: '2px solid black'}}>
                                                                        <AccordionSummary
                                                                            expandIcon={
                                                                                <ExpandMoreIcon/>}>Images</AccordionSummary>
                                                                        <AccordionDetails>
                                                                            {item.images && item.images.length > 0 && (
                                                                                <ImageList
                                                                                    variant="masonry"
                                                                                    sx={{
                                                                                        width: 150,

                                                                                        border: '1px solid #f57c00',
                                                                                        margin: '0 auto' // Center horizontally
                                                                                    }}
                                                                                    cols={1}
                                                                                    gap={1}
                                                                                >
                                                                                    {item.images.map((item1, index) => (
                                                                                        <ImageListItem key={item1}>
                                                                                            <img src={item1} alt=""
                                                                                                 loading="lazy"/>
                                                                                        </ImageListItem>
                                                                                    ))}
                                                                                </ImageList>
                                                                            )}
                                                                        </AccordionDetails>
                                                                    </Accordion>
                                                                </Box>
                                                            </StyledTableCell>
                                                            <StyledTableCell>
                                                                <Box m={1} sx={{
                                                                    overflowX: 'auto',
                                                                    width: 125,
                                                                    height: 125,
                                                                    overflowy: 'auto'
                                                                }}>
                                                                    <Box m={1}>
                                                                        <Button
                                                                            fullWidth
                                                                            color="success"
                                                                            variant="contained"
                                                                            onClick={() => {
                                                                                handleRateAccommodationClick(item)
                                                                            }}
                                                                        >

                                                                            Rate
                                                                        </Button>
                                                                    </Box>
                                                                    {item.rating !== 0 && (
                                                                        <Box m={1}>
                                                                            <Button
                                                                                fullWidth
                                                                                color="error"
                                                                                variant="outlined"
                                                                                onClick={() => {
                                                                                    handleDeleteAccommodationRating(item)
                                                                                }}
                                                                            >
                                                                                Delete rating
                                                                            </Button>
                                                                        </Box>
                                                                    )}
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
                </Flex>
            </div>
        </>
    );
}

export default Ratings;