import React, {useState} from 'react';
import ImageList from '@mui/material/ImageList';
import ImageListItem from '@mui/material/ImageListItem';
import CancelIcon from '@mui/icons-material/Cancel';
import {Box, Button, Card, Grid, TextField} from "@mui/material";
import {Flex} from "reflexbox";
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import Checkbox from '@mui/material/Checkbox';
import AddPhotoAlternateIcon from '@mui/icons-material/AddPhotoAlternate';
import CardHeader from '@mui/material/CardHeader';

import Divider from '@mui/material/Divider';
import AddCircleOutlineOutlinedIcon from "@mui/icons-material/AddCircleOutlineOutlined";


function HostAPlace() {
    const [uploadedImages, setUploadedImages] = useState([]);

    const handleImageRemove = (index) => {
        const updatedItems = [...uploadedImages];
        updatedItems.splice(index, 1);
        setUploadedImages(updatedItems);
    };

    const handleImageUpload = (event) => {
        const files = Array.from(event.target.files);

        const uploadedItems = files.filter(file => file.type.startsWith('image/'))
            .map((file) => ({
                img: URL.createObjectURL(file),
                file,
                id: Date.now(), // add a unique identifier to each file object
            }));

        setUploadedImages([...uploadedImages, ...uploadedItems]);
    };

    ////////////////////////////////////////////
    function not(a, b) {
        return a.filter((value) => b.indexOf(value) === -1);
    }

    function intersection(a, b) {
        return a.filter((value) => b.indexOf(value) !== -1);
    }

    function union(a, b) {
        return [...a, ...not(b, a)];
    }


    const [checked, setChecked] = React.useState([]);
    const [left, setLeft] = React.useState(["Wi-Fi", "Heating", "AC", "Kitchen", "Kitchen", "Kitchen", "Kitchen", "Kitchen", "Kitchen", "Kitchen", "Kitchen"]);
    const [right, setRight] = React.useState([]);

    const leftChecked = intersection(checked, left);
    const rightChecked = intersection(checked, right);

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
                    height: 300,
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


    ////////////////////////////////////////////
    const [placeData, setPlaceData] = useState({
        name: '',
        minGuests: '1',
        maxGuests: '1',
        country: '',
        city: '',
        street: '',
        streetNumber: '',
    });

    const handleChange = (event) => {
        const {name, value} = event.target;
        setPlaceData((prevPlaceData) => ({
            ...prevPlaceData,
            [name]: value,
        }));
    };


    return (
        <div className="wrapper">
            <Flex flexDirection="rows">
                <Flex flexDirection="column" alignItems="center" m={2}>
                    <Box m={1}>
                        Provide basic information
                    </Box>
                    <Box m={1}>
                        <TextField
                            fullWidth
                            variant="filled"
                            label="Name of the place"
                            type="text"
                            name="name"
                            value={placeData.name}
                            onChange={handleChange}
                        />
                    </Box>

                    <Box m={1}>
                        <TextField
                            fullWidth
                            variant="filled"
                            type="number"
                            label="Minimum number of guests"
                            InputProps={{inputProps: {min: 1}}}
                            name="minGuests"
                            value={placeData.minGuests}
                            onChange={handleChange}
                        />
                    </Box>

                    <Box m={1}>
                        <TextField
                            fullWidth
                            variant="filled"
                            type="number"
                            label="Maximum number of guests"
                            InputProps={{inputProps: {min: 1}}}
                            name="maxGuests"
                            value={placeData.maxGuests}
                            onChange={handleChange}
                        />
                    </Box>

                    <Box m={1}>
                        <TextField
                            fullWidth
                            variant="filled"
                            label="Country"
                            type="text"
                            name="country"
                            value={placeData.country}
                            onChange={handleChange}
                        />
                    </Box>

                    <Box m={1}>
                        <TextField
                            fullWidth
                            variant="filled"
                            label="City"
                            type="text"
                            name="city"
                            value={placeData.city}
                            onChange={handleChange}
                        />
                    </Box>

                    <Box m={1}>
                        <TextField
                            fullWidth
                            variant="filled"
                            label="Street"
                            type="text"
                            name="street"
                            value={placeData.street}
                            onChange={handleChange}
                        />
                    </Box>

                    <Box m={1}>
                        <TextField
                            fullWidth
                            variant="filled"
                            label="Street number"
                            type="text"
                            name="streetNumber"
                            value={placeData.streetNumber}
                            onChange={handleChange}
                        />
                    </Box>

                </Flex>
                <hr
                    style={{
                        margin: "10px",
                        border: "1px solid grey",
                    }}
                />
                <Flex flexDirection="column" alignItems="center" m={2}>
                    <Box m={1}>
                        Select amenities
                    </Box>
                    <Grid container spacing={2} justifyContent="center">
                        <Grid item>{customList('Choices', left)}</Grid>
                        <Grid item>
                            <Grid container direction="column" alignItems="center">
                                <Button
                                    sx={{my: 0.5}}
                                    variant="outlined"
                                    size="small"
                                    onClick={handleCheckedRight}
                                    disabled={leftChecked.length === 0}
                                    aria-label="move selected right"
                                >
                                    &gt;
                                </Button>
                                <Button
                                    sx={{my: 0.5}}
                                    variant="outlined"
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
                </Flex>
                <hr
                    style={{
                        margin: "10px",
                        border: "1px solid grey",
                    }}
                />
                <Flex flexDirection="column" justifyContent="center" alignItems="center" m={2}>
                    <Box>
                        Pictures
                    </Box>
                    <ImageList
                        sx={{width: 400, height: 400, border: '1px solid #f57c00'}}
                        cols={3}
                        rowHeight={164}
                    >
                        {uploadedImages.map((item, index) => (
                            <ImageListItem key={item.img}>
                                <img src={item.img} alt="" loading="lazy"/>
                                <button
                                    className="remove-image-button"
                                    onClick={() => handleImageRemove(index)}
                                >
                                    <CancelIcon/>
                                </button>
                            </ImageListItem>
                        ))}
                    </ImageList>

                    <Button variant="contained" color="warning" component="label"
                            endIcon={<AddPhotoAlternateIcon/>}>
                        Add pictures
                        <input
                            type="file"
                            onChange={handleImageUpload}
                            multiple
                            accept="image/*"
                            style={{display: 'none'}}
                        />
                    </Button>

                </Flex>


            </Flex>
            <Flex m={2} flexDirection="column"
                  justifyContent="center"
                  alignItems="center">
                <Button
                    endIcon={<AddCircleOutlineOutlinedIcon/>}
                    color="warning"
                    variant="contained"
                    disabled={(uploadedImages.length === 0 ||
                        Object.values(placeData).some((val) => val === "") ||
                        parseInt(placeData.minGuests) > parseInt(placeData.maxGuests))}
                >
                    Create accommodation
                </Button>
                <Box m={1}>
                    At least one picture must be uploaded, all fields must be filled and number maximum number of guest
                    must be greater or equal to minimum number of guests
                </Box>

            </Flex>


        </div>
    );
}

export default HostAPlace;
