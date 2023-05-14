import React, {useState} from 'react';
import {Flex} from "reflexbox";
import {
    Accordion,
    AccordionDetails,
    AccordionSummary,
    Box,
    Button,
    Card,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    FormControlLabel,
    Grid,
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
import {DatePicker, LocalizationProvider} from "@mui/x-date-pickers";
import {AdapterDayjs} from "@mui/x-date-pickers/AdapterDayjs";
import dayjs from "dayjs";
import Checkbox from "@mui/material/Checkbox";
import CardHeader from "@mui/material/CardHeader";
import Divider from "@mui/material/Divider";
import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import ListItemIcon from "@mui/material/ListItemIcon";
import ListItemText from "@mui/material/ListItemText";
import interceptor from "../../interceptor/interceptor";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import ImageList from "@mui/material/ImageList";
import ImageListItem from "@mui/material/ImageListItem";
import {useNavigate} from "react-router-dom";

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

function SearchAndFilterAccommodations(props) {

    const [left, setLeft] = React.useState(["Wi-Fi", "Heating", "AC", "Kitchen"]);
    const [right, setRight] = React.useState([]);
    const [checked, setChecked] = React.useState([]);
    const leftChecked = intersection(checked, left);
    const rightChecked = intersection(checked, right);
    const [resultData, setResultData] = React.useState(null);

    const [resultDialogShow, setResultDialogShow] = React.useState(false);


    function not(a, b) {
        return a.filter((value) => b.indexOf(value) === -1);
    }

    function intersection(a, b) {
        return a.filter((value) => b.indexOf(value) !== -1);
    }

    function union(a, b) {
        return [...a, ...not(b, a)];
    }

    const handleToggle = (value) => () => {
        const currentIndex = checked.indexOf(value);
        const newChecked = [...checked];

        if (currentIndex === -1) {
            newChecked.push(value);
        } else {
            newChecked.splice(currentIndex, 1);
        }

        setChecked(newChecked);
    };
    const numberOfChecked = (items) => intersection(checked, items).length;
    const handleToggleAll = (items) => () => {
        if (numberOfChecked(items) === items.length) {
            setChecked(not(checked, items));
        } else {
            setChecked(union(checked, items));
        }
    };
    const handleCheckedRight = () => {
        setRight(right.concat(leftChecked));
        setLeft(not(left, leftChecked));
        setChecked(not(checked, leftChecked));
    };
    const handleCheckedLeft = () => {
        setLeft(left.concat(rightChecked));
        setRight(not(right, rightChecked));
        setChecked(not(checked, rightChecked));
    };
    const customList = (title, items) => (
        <Card>
            <CardHeader

                avatar={
                    <Checkbox
                        onClick={handleToggleAll(items)}
                        checked={numberOfChecked(items) === items.length && items.length !== 0}
                        indeterminate={
                            numberOfChecked(items) !== items.length && numberOfChecked(items) !== 0
                        }
                        disabled={items.length === 0}
                        inputProps={{
                            'aria-label': 'all items selected',
                        }}
                    />
                }
                title={title}
                subheader={`${numberOfChecked(items)}/${items.length} selected`}
            />
            <Divider/>
            <List
                sx={{
                    width: 200,
                    height: 150,
                    overflow: 'auto',
                }}
                dense
                component="div"
                role="list"
            >
                {items.map((value) => {
                    const labelId = `transfer-list-all-item-${value}-label`;

                    return (
                        <ListItem
                            key={value}
                            role="listitem"
                            button
                            onClick={handleToggle(value)}
                        >
                            <ListItemIcon>
                                <Checkbox
                                    checked={checked.indexOf(value) !== -1}
                                    tabIndex={-1}
                                    disableRipple
                                    inputProps={{
                                        'aria-labelledby': labelId,
                                    }}
                                />
                            </ListItemIcon>
                            <ListItemText id={labelId} primary={` ${value}`}/>
                        </ListItem>
                    );
                })}
            </List>
        </Card>
    );

    const handleResultDialogShow = () => {
        setResultDialogShow(false)
    };

    const [formData, setFormData] = useState({
        location: '',
        minGuests: 1,
        startDate: dayjs(),
        endDate: dayjs(),
        minPrice: 1,
        maxPrice: 2,
        prominentHost: false,
        amenities: []
    });

    const handleInputChange = (event) => {
        const {name, value, type, checked} = event.target;
        const newValue = type === 'checkbox' ? checked : value;
        setFormData((prevState) => ({
            ...prevState,
            [name]: newValue
        }));
    };


    function handleSearch() {


        const searchAndFilterData = {...formData};
        searchAndFilterData.amenities = right;

        const startDate = new Date(formData.startDate);
        const utcStartDate = new Date(startDate.getTime() - startDate.getTimezoneOffset() * 60000);
        console.log(utcStartDate)
        searchAndFilterData.startDate = Math.round(utcStartDate.getTime() / 1000);

        const endDate = new Date(formData.endDate);
        const utcEndDate = new Date(endDate.getTime() - endDate.getTimezoneOffset() * 60000);
        searchAndFilterData.endDate = Math.round(utcEndDate.getTime() / 1000);
        searchAndFilterData.minGuests = parseInt(searchAndFilterData.minGuests)


        setResultDialogShow(true);
        interceptor.post("api-2/accommodation/search", searchAndFilterData).then(res => {
            setResultData(res.data)


        }).catch(err => {
                console.log(err)
            }
        );


    }

    const navigate = useNavigate();
    const handleReserve = (item) => {

        const startDate = new Date(formData.startDate);
        const utcStartDate = new Date(startDate.getTime() - startDate.getTimezoneOffset() * 60000);


        const endDate = new Date(formData.endDate);
        const utcEndDate = new Date(endDate.getTime() - endDate.getTimezoneOffset() * 60000);


        var sendData = {}
        sendData.accommodationId = item.id
        sendData.numberOfGuests = parseInt(formData.minGuests)
        sendData.dateRange = {}
        sendData.dateRange.from = Math.round(utcStartDate.getTime() / 1000);
        sendData.dateRange.to = Math.round(utcEndDate.getTime() / 1000);

        interceptor.post("api-1/reservation", {reservation: sendData}).then(res => {
            navigate("/profile")

        }).catch(err => {
                console.log(err)
            }
        );


    };
    return (
        <>
            <Dialog onClose={handleResultDialogShow} open={resultDialogShow}>
                <DialogTitle>Search and filter results</DialogTitle>

                <DialogContent>
                    {resultData != null && resultData.length > 0 &&
                        <TableContainer component={Paper}
                                        sx={{maxHeight: 500, height: 500, overflowY: 'scroll'}}>
                            <Table>

                                <TableBody>
                                    {resultData.map((item) => (
                                        <React.Fragment key={`${item.id}-row`}>
                                            <StyledTableRow>
                                                <StyledTableCell>
                                                    <li>Name: {item.name}</li>
                                                    <li>Total price: {item.price}</li>
                                                    <li>{item.address.city}, {item.address.country}</li>
                                                    <li>{item.address.street}, {item.address.streetNumber}</li>

                                                </StyledTableCell>
                                                <StyledTableCell>
                                                    <Accordion sx={{border: "1px solid black"}}>
                                                        <AccordionSummary
                                                            expandIcon={<ExpandMoreIcon/>}>
                                                            Images
                                                        </AccordionSummary>
                                                        <AccordionDetails>
                                                            {item.images && item.images.length > 0 && (
                                                                <ImageList variant="masonry"
                                                                           sx={{
                                                                               width: 250,
                                                                               height: 250,
                                                                               border: '1px solid #f57c00'
                                                                           }}
                                                                           cols={2}
                                                                           gap={1}>
                                                                    {item.images.map((item1, index) => (
                                                                        <ImageListItem key={item1}>
                                                                            <img src={item1} alt="" loading="lazy"
                                                                            />
                                                                        </ImageListItem>
                                                                    ))}
                                                                </ImageList>
                                                            )}
                                                        </AccordionDetails>
                                                    </Accordion>
                                                </StyledTableCell>

                                                {props.canBuy && props.canBuy == true &&
                                                    <StyledTableCell>
                                                        <Button
                                                            onClick={() => {
                                                                handleReserve(item)
                                                            }}
                                                            fullWidth
                                                            color="warning"
                                                            variant="contained">

                                                            Reserve
                                                        </Button>
                                                    </StyledTableCell>
                                                }

                                            </StyledTableRow>
                                        </React.Fragment>
                                    ))
                                    }
                                </TableBody>
                            </Table>
                        </TableContainer>
                    }

                </DialogContent>

                <DialogActions>
                    <Button onClick={handleResultDialogShow}
                            variant="contained"
                    >
                        Close
                    </Button>
                </DialogActions>
            </Dialog>

            <div className="wrapper">
                <Flex flexDirection="column">
                    <Flex flexDirection="row" justifyContent="center" alignItems="center">
                        Search parameters
                    </Flex>
                    <Flex flexDirection="row" justifyContent="center" alignItems="center">
                        <Box width={1 / 4} m={1}>
                            <TextField
                                fullWidth
                                variant="filled"
                                label="Location"
                                name="location"
                                value={formData.location}
                                onChange={handleInputChange}
                            />
                        </Box>
                        <Box width={1 / 4} m={1}>
                            <TextField
                                fullWidth
                                variant="filled"
                                type="number"
                                label="Number of guests"
                                InputProps={{inputProps: {min: 1}}}

                                name="minGuests"
                                value={formData.minGuests}
                                onChange={handleInputChange}
                            />
                        </Box>
                        <Box width={1 / 4} m={1}>
                            <LocalizationProvider dateAdapter={AdapterDayjs}>
                                <DatePicker label="Start date"
                                            value={formData.startDate}
                                            onChange={(date) =>
                                                setFormData((prevState) => ({
                                                    ...prevState,
                                                    startDate: date
                                                }))
                                            }
                                            minDate={dayjs()}/>
                            </LocalizationProvider>
                        </Box>
                        <Box width={1 / 4} m={1}>
                            <LocalizationProvider dateAdapter={AdapterDayjs}>
                                <DatePicker label="End date"
                                            value={formData.endDate}
                                            onChange={(date) =>
                                                setFormData((prevState) => ({
                                                    ...prevState,
                                                    endDate: date
                                                }))
                                            }
                                            minDate={dayjs()}/>
                            </LocalizationProvider>
                        </Box>
                    </Flex>
                    <hr style={{width: "100%", border: "1px solid grey"}}
                    />
                    <Flex flexDirection="row" justifyContent="center" alignItems="center">
                        <Box m={1}>
                            Filter parameters
                        </Box>
                    </Flex>
                    <Flex flexDirection="row" m={1} justifyContent="center" alignItems="center">

                        <Flex flexDirection="column" width={1 / 4} m={1}>
                            <Box m={1}>
                                <TextField
                                    fullWidth
                                    variant="filled"
                                    type="number"
                                    label="Minimum price"
                                    InputProps={{inputProps: {min: 1}}}
                                    name="minPrice"
                                    value={formData.minPrice}
                                    onChange={handleInputChange}
                                />
                            </Box>
                            <Box m={1}>
                                <TextField
                                    fullWidth
                                    variant="filled"
                                    type="number"
                                    label="Maximum price"
                                    name="maxPrice"
                                    value={formData.maxPrice}
                                    onChange={handleInputChange}
                                    InputProps={{inputProps: {min: 1}}}
                                />
                            </Box>
                            <Box m={1}>
                                <FormControlLabel
                                    control={
                                        <Checkbox
                                            checked={formData.prominentHost}
                                            onChange={handleInputChange}
                                            name="prominentHost"
                                        />
                                    }
                                    label="Prominent Host"
                                />
                            </Box>
                        </Flex>

                        <Box direction="column" m={1}
                             alignItems="center" justifyContent="center">


                            <Grid container spacing={1}>

                                <Grid item>{customList('Choices', left)}</Grid>
                                <Grid item>
                                    <Grid container direction="column" alignItems="center">
                                        <Flex alignItems="center" justifyContent="center" m={1}>
                                            Select desired amenities
                                        </Flex>
                                        <Button
                                            sx={{my: 0.5}}
                                            variant="contained"
                                            size="small"
                                            onClick={handleCheckedRight}
                                            disabled={leftChecked.length === 0}
                                            aria-label="move selected right"
                                        >
                                            &gt;
                                        </Button>
                                        <Button
                                            sx={{my: 0.5}}
                                            variant="contained"
                                            size="small"
                                            onClick={handleCheckedLeft}
                                            disabled={rightChecked.length === 0}
                                            aria-label="move selected left"
                                        >
                                            &lt;
                                        </Button>
                                    </Grid>
                                </Grid>
                                <Grid item>{customList('Chosen', right)}</Grid>
                            </Grid>

                        </Box>
                    </Flex>
                    <Flex flexDirection="column" justifyContent="center" alignItems="center">
                        <Box width={1 / 3} m={1}>
                            <Button
                                onClick={handleSearch}
                                fullWidth
                                color="warning"
                                variant="contained">

                                Search and filter
                            </Button>
                        </Box>
                    </Flex>
                </Flex></div>
        </>
    );
}

export default SearchAndFilterAccommodations;