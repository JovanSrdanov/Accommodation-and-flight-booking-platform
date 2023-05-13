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
    Paper,
    styled,
    Table,
    TableBody,
    TableCell,
    tableCellClasses,
    TableContainer,
    TableRow
} from "@mui/material";
import {Flex} from "reflexbox";
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

function MyPlaces() {

    const [myAccommodations, setMyAccommodations] = useState(null);
    const [myAvailabilities, setMyAvailabilities] = useState(null);

    const [selectedImages, setSelectedImages] = useState(null);
    const [selectedImage, setSelectedImage] = useState(null);

    const [viewImagesDialog, setViewImagesDialog] = useState(false);
    const [viewImageDialog, setViewImageDialog] = useState(false);

    const [idForCreatingAvailability, setIdForCreatingAvailability] = useState(null);
    const [idForUpdatingAvailability, setIdForUpdatingAvailability] = useState(null);

    const getMyPlaces = () => {
        interceptor.get("api-1/accommodation/all-my").then(res => {
            setMyAccommodations(res.data.accommodation)

        }).catch(err => {
            console.log(err)
        })
    }

    const getMyAvailabilities = () => {
        interceptor.get("/api-1/availability/all").then(res => {
            //  console.log(res.data)
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
    return (
        <>

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
                                                            <Button fullWidth variant="contained" color="warning">Create
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
                                            <React.Fragment key={`${item.id}-row`}>
                                                <StyledTableRow hover>
                                                    <StyledTableCell>

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