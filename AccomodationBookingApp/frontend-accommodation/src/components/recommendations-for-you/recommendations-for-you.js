import React, {useEffect, useState} from 'react';

import interceptor from "../../interceptor/interceptor";
import {
    Accordion,
    AccordionDetails,
    AccordionSummary,
    Box,
    Paper,
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


function RecommendationsForYou(props) {
    const [recommended, setRecommended] = useState(null);
    const getRecommended = () => {
        interceptor.get("/api-2/accommodation/recommend").then((res) => {

            setRecommended(res.data)
        }).catch((err) => {
            console.log(err)
        })
    };
    useEffect(() => {
        getRecommended()
    }, []);

    return (
        <>
            <div>
                {recommended != null && recommended.length > 0 && (
                    <div className="wrapper">
                        <TableContainer component={Paper} sx={{maxHeight: 500, height: 500, overflowY: 'scroll'}}>
                            <Table>
                                <TableBody>
                                    {recommended.map((item) => (
                                        <React.Fragment key={`${item.id}-row`}>
                                            <StyledTableRow>
                                                <StyledTableCell>
                                                    <Box m={1} sx={{
                                                        width: 250,
                                                        height: 125,
                                                        overflowy: 'auto',
                                                        overflowX: 'auto'
                                                    }}>
                                                        <li>Name: {item.name}</li>
                                                        <li>Total price: {item.price}$</li>
                                                        <li>{item.address.city}, {item.address.country}</li>
                                                        <li>{item.address.street}, {item.address.streetNumber}</li>
                                                        <li>Guests: {item.minGuests} - {item.maxGuests}</li>
                                                        <li>Rating: {item.rating}</li>
                                                    </Box>
                                                </StyledTableCell>
                                                <StyledTableCell>
                                                    <Box m={1} sx={{
                                                        width: 350,
                                                        overflowy: 'auto',
                                                        overflowX: 'auto'
                                                    }}>
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
                                                        </Box>
                                                        <Box m={1}>
                                                            <Accordion sx={{border: '1px solid black'}}>
                                                                <AccordionSummary
                                                                    expandIcon={
                                                                        <ExpandMoreIcon/>}>Images</AccordionSummary>
                                                                <AccordionDetails>
                                                                    {item.images && item.images.length > 0 && (
                                                                        <ImageList
                                                                            variant="masonry"
                                                                            sx={{
                                                                                width: 250,
                                                                                height: 200,
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
                                                    </Box>

                                                </StyledTableCell>

                                            </StyledTableRow>
                                        </React.Fragment>
                                    ))}
                                </TableBody>
                            </Table>
                        </TableContainer>
                    </div>
                )}

            </div>


        </>
    );
}

export default RecommendationsForYou;