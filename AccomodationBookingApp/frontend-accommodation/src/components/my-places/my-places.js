import React, {useEffect, useState} from 'react';
import interceptor from "../../interceptor/interceptor";
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
    FormControlLabel,
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
import {Flex} from "reflexbox";
import ImageList from "@mui/material/ImageList";
import ImageListItem from "@mui/material/ImageListItem";
import {DatePicker, LocalizationProvider} from "@mui/x-date-pickers";
import {AdapterDayjs} from "@mui/x-date-pickers/AdapterDayjs";
import dayjs from "dayjs";
import Checkbox from "@mui/material/Checkbox";

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

function MyPlaces() {

    const [myAccommodations, setMyAccommodations] = useState(null);
    const [myAvailabilities, setMyAvailabilities] = useState(null);

    const [selectedImages, setSelectedImages] = useState(null);
    const [selectedImage, setSelectedImage] = useState(null);

    const [viewImagesDialog, setViewImagesDialog] = useState(false);
    const [viewImageDialog, setViewImageDialog] = useState(false);

    const [openCreateAvailabilityDialog, setOpenCreateAvailabilityDialog] = useState(false);

    const [openUpdateAvailabilityDialog, setOpenUpdateAvailabilityDialog] = useState(false);

    const [idForCreatingAvailability, setIdForCreatingAvailability] = useState(null);
    const [idForUpdatingAvailability, setIdForUpdatingAvailability] = useState(null);


    const handleFormChange = (event) => {
        const {name, value, checked} = event.target;

        setAvailability((prevState) => ({
            ...prevState,
            priceWithDate: {
                ...prevState.priceWithDate,
                [name]: name === 'isPricePerPerson' ? checked : value
            }
        }));
    };

    const getMyPlaces = () => {
        interceptor.get("api-1/accommodation/all-my").then(res => {
            setMyAccommodations(res.data.accommodation)

        }).catch(err => {
            console.log(err)
        })
    }

    const getMyAvailabilities = () => {
        interceptor.get("/api-1/availability/all").then(res => {
            console.log(res.data.availabilities)
            setMyAvailabilities(res.data.availabilities)
        }).catch(err => {
            console.log(err)
        })
    }

    const getMyInfo = () => {
        getMyPlaces();
        getMyAvailabilities();
    }

    useEffect(() => {
        getMyInfo();

    }, []);

    const handleImages = (Images) => {
        setSelectedImages(Images);
        setViewImagesDialog(true);
    };

    const handleCloseViewImagesDialog = () => {
        setSelectedImages(null);
        setViewImagesDialog(false);
    };
    const handleImage = (item) => {
        setSelectedImage(item)
        setViewImageDialog(true);
    };
    const handleCloseViewImageDialog = () => {
        setSelectedImage(null)
        setViewImageDialog(false);
    };
    const handleCloseCreateAvailabilityDialog = () => {
        setOpenCreateAvailabilityDialog(false)
        setIdForCreatingAvailability(null)
        setAvailability({
            accommodationId: '',
            priceWithDate: {
                dateRange: {
                    from: null,
                    to: null
                },
                isPricePerPerson: true,
                price: 1,
            },
        });
    };

    function parseObjectId(str) {
        const regex = /^ObjectID\("(.+)"\)$/;
        const match = str.match(regex);
        return match ? match[1] : null;
    }

    const handleOpenCreateAvailabilityDialog = (item) => {
        setOpenCreateAvailabilityDialog(true)
        setIdForCreatingAvailability(parseObjectId(item.id))
    };

    const [availability, setAvailability] = useState({
        accommodationId: '',
        priceWithDate: {
            dateRange: {
                from: null,
                to: null
            },
            isPricePerPerson: true,
            price: 1,
        },
    });


    const handleCreateAvailability = () => {
        const sendData = {...availability};
        sendData.accommodationId = idForCreatingAvailability;

        // Convert start and end dates to Date objects
        const startDate = dayjs(availability.priceWithDate.dateRange.from).toDate();
        const endDate = dayjs(availability.priceWithDate.dateRange.to).toDate();
        sendData.priceWithDate.dateRange.from = startDate.getTime();
        sendData.priceWithDate.dateRange.to = endDate.getTime();

        console.log(sendData);

        interceptor.post("api-1/availability", {availability: sendData}).then(res => {
            getMyInfo();
            handleCloseCreateAvailabilityDialog();
        }).catch(err => {
            console.log(err)
        })
    };

    const handleStartDateChange = (date) => {
        const newAvailability = {...availability};
        newAvailability.priceWithDate.dateRange.from = date;
        setAvailability(newAvailability);
    };

    const handleEndDateChange = (date) => {
        const newAvailability = {...availability};
        newAvailability.priceWithDate.dateRange.to = date;
        setAvailability(newAvailability);
    };

    const handlePriceChange = (event) => {
        const newAvailability = {...availability};
        newAvailability.priceWithDate.price = parseInt(event.target.value);
        setAvailability(newAvailability);
    };

    const handleIsPricePerPersonChange = (event) => {
        const newAvailability = {...availability};
        newAvailability.priceWithDate.isPricePerPerson = event.target.checked;
        setAvailability(newAvailability);
    };


    return (
        <>
            <Dialog onClose={handleCloseCreateAvailabilityDialog} open={openCreateAvailabilityDialog}>
                <DialogContent>
                    <Flex flexDirection="column">
                        <Flex flexDirection="row" justifyContent="center" alignItems="center">
                            <Box m={1}>
                                Create availability
                            </Box>
                        </Flex>
                        <Flex flexDirection="row" justifyContent="center" alignItems="center">
                            <Box width={1 / 2} m={1}>
                                <LocalizationProvider dateAdapter={AdapterDayjs}>
                                    <DatePicker
                                        label="Start date"
                                        minDate={dayjs()}
                                        name="from"
                                        onChange={handleStartDateChange}
                                        value={availability.priceWithDate.dateRange.from}

                                    />
                                </LocalizationProvider>
                            </Box>
                            <Box width={1 / 2} m={1}>
                                <LocalizationProvider dateAdapter={AdapterDayjs}>
                                    <DatePicker
                                        label="End date"
                                        minDate={dayjs()}
                                        name="to"
                                        onChange={handleEndDateChange}

                                        value={availability.priceWithDate.dateRange.to}

                                    />
                                </LocalizationProvider>
                            </Box>
                        </Flex>
                        <Flex flexDirection="row" justifyContent="center" alignItems="center">
                            <Box m={1} width={1 / 2}>
                                <TextField
                                    fullWidth
                                    variant="filled"
                                    onChange={handlePriceChange}
                                    type="number"
                                    label="Price"
                                    InputProps={{inputProps: {min: 1}}}
                                    name="price"
                                    value={availability.priceWithDate.price}

                                />
                            </Box>
                            <Box m={1} width={1 / 2}>
                                <FormControlLabel
                                    control={
                                        <Checkbox
                                            name="isPricePerPerson"
                                            onChange={handleIsPricePerPersonChange}
                                            checked={availability.priceWithDate.isPricePerPerson}

                                        />
                                    }
                                    label="Price per person"
                                />
                            </Box>
                        </Flex>
                        <Flex flexDirection="row" justifyContent="center" alignItems="center">
                            <Box m={1}>
                                <Button
                                    disabled={(!availability.priceWithDate.dateRange.from || !availability.priceWithDate.dateRange.to || dayjs(availability.priceWithDate.dateRange.to).isBefore(availability.priceWithDate.dateRange.from))}
                                    onClick={handleCreateAvailability}
                                    variant="contained" color="warning">Create availability</Button>
                            </Box>
                        </Flex>
                    </Flex>
                </DialogContent>

                <DialogActions>
                    <Button onClick={handleCloseCreateAvailabilityDialog} variant="contained">Close
                    </Button>
                </DialogActions>
            </Dialog>


            <Dialog onClose={handleCloseViewImageDialog} open={viewImageDialog}>
                <DialogContent>
                    <img src={selectedImage} alt="" style={{width: '100%'}}/>
                </DialogContent>
                <DialogActions>
                    <Button onClick={handleCloseViewImageDialog}
                            variant="contained"
                    >
                        Close
                    </Button>
                </DialogActions>
            </Dialog>

            <Dialog onClose={handleCloseViewImagesDialog} open={viewImagesDialog}>
                <DialogTitle>Images</DialogTitle>
                <DialogContent>
                    {selectedImages && selectedImages.length > 0 && (
                        <ImageList variant="masonry"
                                   sx={{width: 500, height: 500, border: '1px solid #f57c00'}}
                                   cols={3}
                                   gap={5}>
                            {selectedImages.map((item, index) => (
                                <ImageListItem key={item}>
                                    <img src={item} alt="" loading="lazy" onClick={() => handleImage(item)}
                                         style={{cursor: 'pointer'}}/>
                                </ImageListItem>
                            ))}
                        </ImageList>
                    )}
                </DialogContent>
                <DialogActions>
                    <Button onClick={handleCloseViewImagesDialog}
                            variant="contained"
                    >
                        Close
                    </Button>
                </DialogActions>
            </Dialog>

            <div className="wrapper">
                <Flex flexDirection="rows">
                    <Flex flexDirection="column" alignItems="center" m={2}>
                        <Box m={1}>
                            My places
                        </Box>
                        <TableContainer component={Paper} sx={{maxHeight: 500, height: 500}}>
                            <Table>
                                {myAccommodations != null &&
                                    <TableBody>
                                        {myAccommodations.map((item) => (
                                            <React.Fragment key={`${item.id}-row`}>
                                                <StyledTableRow hover>
                                                    <StyledTableCell>
                                                        <li>Name: {item.Name}</li>
                                                        <li>Automatic reservation: {item.isAutomaticReservation}</li>
                                                        <li>Number of guest: {item.MinGuests} - {item.MaxGuests}</li>
                                                        <li>{item.Address.city}, {item.Address.country}</li>
                                                        <li>{item.Address.street}, {item.Address.streetNumber}</li>
                                                    </StyledTableCell>
                                                    <StyledTableCell align="center">
                                                        <Accordion>
                                                            <AccordionSummary>
                                                                Ameneties
                                                            </AccordionSummary>
                                                            <AccordionDetails>
                                                                {item.Amenities.map((a) => (
                                                                    <li key={a}>{a}</li>
                                                                ))}
                                                            </AccordionDetails>
                                                        </Accordion>
                                                    </StyledTableCell>
                                                    <StyledTableCell align="center">
                                                        <Box m={1}>
                                                            <Button fullWidth variant="contained"
                                                                    onClick={() => handleImages(item.Images)}>View
                                                                images</Button>
                                                        </Box>
                                                        <Box m={1}>
                                                            <Button fullWidth variant="contained" color="warning"
                                                                    onClick={() => handleOpenCreateAvailabilityDialog(item)}
                                                            >Create
                                                                availability</Button>
                                                        </Box>
                                                    </StyledTableCell>
                                                </StyledTableRow>
                                            </React.Fragment>
                                        ))}
                                    </TableBody>
                                }
                            </Table>
                        </TableContainer>
                    </Flex>
                    <Flex flexDirection="column" alignItems="center" m={2}>
                        <Box m={1}>
                            My availabilities
                        </Box>

                        <TableContainer component={Paper} sx={{maxHeight: 500, height: 500}}>
                            <Table>
                                {myAvailabilities != null &&
                                    <TableBody>
                                        {myAvailabilities.map((item) => (
                                            <React.Fragment key={`${item.id}-row-myAvailabilities`}>
                                                <StyledTableRow hover>
                                                    <StyledTableCell>
                                                        <li>From:</li>
                                                        <li>To</li>
                                                        <li>Price</li>
                                                        <li>Per person</li>
                                                    </StyledTableCell>
                                                    <StyledTableCell align="center">
                                                        <Box m={1}>
                                                            <Button fullWidth variant="contained" color="warning">Change
                                                                availability</Button>
                                                        </Box>
                                                    </StyledTableCell>
                                                </StyledTableRow>
                                            </React.Fragment>
                                        ))}
                                    </TableBody>
                                }
                            </Table>
                        </TableContainer>
                    </Flex>
                </Flex>
            </div>
        </>
    );
}

export default MyPlaces;