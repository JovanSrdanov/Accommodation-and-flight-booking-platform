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
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';

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


    const [selectedImages, setSelectedImages] = useState(null);
    const [selectedImage, setSelectedImage] = useState(null);

    const [viewImagesDialog, setViewImagesDialog] = useState(false);
    const [viewImageDialog, setViewImageDialog] = useState(false);

    const [errorDialogShow, setErrorDialogShow] = useState(false)
    const [openCreateAvailabilityDialog, setOpenCreateAvailabilityDialog] = useState(false);

    const [openUpdateAvailabilityDialog, setOpenUpdateAvailabilityDialog] = useState(false);

    const [idForCreatingAvailability, setIdForCreatingAvailability] = useState(null);
    const [idForUpdatingAvailability, setIdForUpdatingAvailability] = useState(null);
    const handleErrorClose = () => {
        setErrorDialogShow(false)
    };

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

    const getMyPlaces = async () => {
        try {
            const res = await interceptor.get("api-1/accommodation/all-my");
            const accommodations = res.data.accommodation;

            const res2 = await interceptor.get("/api-1/availability/all");
            const availabilities = res2.data.availabilities;

            const updatedAccommodations = accommodations.map((accommodation) => {
                availabilities.forEach((availability) => {
                    if (accommodation.id === availability.accommodationId) {
                        accommodation.listOfMyAvailabilities = availability;
                    }
                });
                return accommodation;
            });

            setMyAccommodations(updatedAccommodations);
        } catch (err) {
            console.log(err);
        }
    };

    const [myAccommodations, setMyAccommodations] = useState(null);


    useEffect(() => {
        getMyPlaces();

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

    const handleOpenUpdateAvailabilityDialog = (item) => {
        setOpenCreateAvailabilityDialog(true)
        setIdForUpdatingAvailability(parseObjectId(item.id))
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
        console.log(availability.priceWithDate.dateRange)

        // Convert start and end dates to Date objects
        const startDate = new Date(availability.priceWithDate.dateRange.from);
        const utcStartDate = new Date(startDate.getTime() - startDate.getTimezoneOffset() * 60000);
        const formattedStartDate = utcStartDate.toLocaleString("en-US", {timeZone: "GMT"}) + " GMT+0000";
        sendData.priceWithDate.dateRange.from = Date.parse(formattedStartDate) / 1000;

        const endDate = new Date(availability.priceWithDate.dateRange.to);
        const utcEndDate = new Date(endDate.getTime() - endDate.getTimezoneOffset() * 60000);
        const formattedEndDate = utcEndDate.toLocaleString("en-US", {timeZone: "GMT"}) + " GMT+0000";
        sendData.priceWithDate.dateRange.to = Date.parse(formattedEndDate) / 1000;


        interceptor.post("api-1/availability", {availability: sendData}).then(res => {
            getMyPlaces();

            handleCloseCreateAvailabilityDialog();
        }).catch(err => {
            setErrorDialogShow(true)
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


    const handleCloseUpdateAvailabilityDialog = () => {
        setOpenUpdateAvailabilityDialog(false)
        setIdForUpdatingAvailability(null)
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

    const [selectedAccommodationForAVCHANGE, setSelectedAccommodationForAVCHANGE] = useState(false);
    const [selectedAVCHANGE, setSelectedAVCHANGE] = useState(false);
    const handleChangeAvailabilityDialog = (item, a) => {
        setOpenUpdateAvailabilityDialog(true)
        setSelectedAccommodationForAVCHANGE(item)
        setSelectedAVCHANGE(a)
    };
    const handleChangeAvailability = () => {
        var sendData = {};
        sendData.updatedPriceWithDate = {}
        sendData.accommodationId = parseObjectId(selectedAccommodationForAVCHANGE.id);
        sendData.updatedPriceWithDate.Id = parseObjectId(selectedAVCHANGE.Id);
        sendData.updatedPriceWithDate.isPricePerPerson = availability.priceWithDate.isPricePerPerson;
        sendData.updatedPriceWithDate.price = availability.priceWithDate.price;
        sendData.updatedPriceWithDate.dateRange = {}

        // Convert start and end dates to Date objects
        const startDate = new Date(availability.priceWithDate.dateRange.from);
        const utcStartDate = new Date(startDate.getTime() - startDate.getTimezoneOffset() * 60000);
        const formattedStartDate = utcStartDate.toLocaleString("en-US", {timeZone: "GMT"}) + " GMT+0000";
        sendData.updatedPriceWithDate.dateRange.from = Date.parse(formattedStartDate) / 1000;

        const endDate = new Date(availability.priceWithDate.dateRange.to);
        const utcEndDate = new Date(endDate.getTime() - endDate.getTimezoneOffset() * 60000);
        const formattedEndDate = utcEndDate.toLocaleString("en-US", {timeZone: "GMT"}) + " GMT+0000";
        sendData.updatedPriceWithDate.dateRange.to = Date.parse(formattedEndDate) / 1000;
        console.log(utcEndDate)

        interceptor.put("api-1/availability", {priceWithDate: sendData}).then(res => {
            getMyPlaces();

            setOpenUpdateAvailabilityDialog(false)
        }).catch(err => {
            setErrorDialogShow(true)
            console.log(err)
        })

    };
    return (
        <>
            <Dialog onClose={handleCloseUpdateAvailabilityDialog} open={openUpdateAvailabilityDialog}>
                <DialogContent>
                    <Flex flexDirection="column">
                        <Flex flexDirection="row" justifyContent="center" alignItems="center">
                            <Box m={1}>
                                Update availability
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
                                    onClick={handleChangeAvailability}
                                    variant="contained" color="warning">Change availability</Button>
                            </Box>
                        </Flex>
                    </Flex>
                </DialogContent>

                <DialogActions>
                    <Button onClick={handleCloseCreateAvailabilityDialog} variant="contained">Close
                    </Button>
                </DialogActions>
            </Dialog>

            <Dialog onClose={handleErrorClose} open={errorDialogShow}>
                <DialogTitle>Error</DialogTitle>
                <DialogContent>
                    This availability can not be made because it overlaps with other availabilities or someone already
                    made a reservation on it
                </DialogContent>
                <DialogActions>
                    <Button onClick={handleErrorClose}
                            variant="contained"
                    >
                        Close
                    </Button>
                </DialogActions>
            </Dialog>

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
                        <TableContainer component={Paper}
                                        sx={{maxHeight: 500, height: 500, overflowY: 'scroll'}}>
                            <Table>
                                {myAccommodations != null &&
                                    <TableBody>
                                        {myAccommodations.map((item) => (
                                            <React.Fragment key={`${item.id}-row`}>
                                                <StyledTableRow hover>
                                                    <StyledTableCell>
                                                        <li>Name: {item.Name}</li>

                                                        <li>Number of guest: {item.MinGuests} - {item.MaxGuests}</li>
                                                        <li>{item.Address.city}, {item.Address.country}</li>
                                                        <li>{item.Address.street}, {item.Address.streetNumber}</li>
                                                    </StyledTableCell>
                                                    <StyledTableCell>
                                                        <Accordion sx={{border: "1px solid black"}}>
                                                            <AccordionSummary
                                                                expandIcon={<ExpandMoreIcon/>}>
                                                                List Of Ameneties
                                                            </AccordionSummary>
                                                            <AccordionDetails sx={{height: 200, overflowY: 'scroll'}}>
                                                                {item.Amenities.map((a) => (
                                                                    <Box m={1} key={a}>- {a}</Box>
                                                                ))}
                                                            </AccordionDetails>
                                                        </Accordion>
                                                    </StyledTableCell>
                                                    <StyledTableCell>
                                                        <Accordion sx={{border: "1px solid black"}}>
                                                            <AccordionSummary expandIcon={<ExpandMoreIcon/>}>
                                                                List Of Availabilities
                                                            </AccordionSummary>
                                                            <AccordionDetails>
                                                                <TableContainer sx={{height: 200, overflowY: 'scroll'}}
                                                                                component={Paper}
                                                                >
                                                                    <Table>
                                                                        <tbody>
                                                                        {item.listOfMyAvailabilities &&
                                                                            item.listOfMyAvailabilities.availableDates &&
                                                                            item.listOfMyAvailabilities.availableDates.map((a) => (
                                                                                <StyledTableRow key={a.Id} hover>
                                                                                    <StyledTableCell>
                                                                                        <li>From: {new Date(a.dateRange.from * 1000).toLocaleDateString("en-GB")}</li>
                                                                                        <li>To: {new Date(a.dateRange.to * 1000).toLocaleDateString("en-GB")}</li>
                                                                                        <li>Price: {a.price}</li>
                                                                                        <li>Is price per
                                                                                            person: {a.isPricePerPerson ? 'Yes' : 'No'}</li>
                                                                                    </StyledTableCell>
                                                                                    <StyledTableCell align="center">
                                                                                        <Box m={1}>
                                                                                            <Button fullWidth
                                                                                                    color="warning"
                                                                                                    variant="contained"
                                                                                                    onClick={() => handleChangeAvailabilityDialog(item, a)}>

                                                                                                Change
                                                                                            </Button>
                                                                                        </Box>
                                                                                    </StyledTableCell>
                                                                                </StyledTableRow>
                                                                            ))}
                                                                        </tbody>
                                                                    </Table>
                                                                </TableContainer>

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

                </Flex>
            </div>
        </>
    );
}

export default MyPlaces;