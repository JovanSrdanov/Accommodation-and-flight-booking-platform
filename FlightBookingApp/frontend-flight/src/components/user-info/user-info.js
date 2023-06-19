import React, {useEffect, useState} from 'react';
import {Box, Button, Checkbox, Dialog, DialogActions, DialogContent, DialogTitle} from '@mui/material';

import {useTheme} from '@mui/material/styles';

import cssClasses from './user-info.module.css';
import useAccountApi from '../../hooks/useAccountApi';
import useUserApi from '../../hooks/useUserApi';
import useAxiosPrivate from '../../hooks/useAxiosPrivate';
import {DatePicker, LocalizationProvider} from "@mui/x-date-pickers";
import {AdapterDayjs} from "@mui/x-date-pickers/AdapterDayjs";
import dayjs from "dayjs";

function UserInfo(props) {
    const axios = useAxiosPrivate();
    const theme = useTheme();

    const InfoRow = (props) => {
        return (
            <div className={cssClasses.infoRow}>
                <b>
                    <p>{props.info}:</p>
                </b>
                <p>{props.value}</p>
            </div>
        );
    };

    const tempUser = {
        fullname: 'Loading...',
        address: 'Loading...',
    };

    const tempAcc = {
        username: 'Loading...',
        email: 'Loading...',
    };

    const [accountInfo, setAccountInfo] = useState(tempAcc);
    const [userInfo, setUserInfo] = useState(tempUser);
    const [apiKey, setApiKey] = useState(null);
    const [openDialog, setOpenDialog] = useState(false);
    const [validForever, setValidForever] = useState(false);
    const [selectedDate, setSelectedDate] = useState(dayjs().add(2, 'day'));

    const {GetAccountInfo} = useAccountApi();
    const {GetLoggedUserInfo} = useUserApi();

    const handleGenerateClick = () => {
        if (validForever === true) {
            axios.post("api/account/api-key", {}).then((res) => {
                getApiKey();
            }).catch((err) => {
                console.log(err)
            })
            setOpenDialog(false)
            getApiKey();
        } else {
            let date = new Date(selectedDate)
            axios.post("api/account/api-key", {expirationDate: date}).then((res) => {
                setOpenDialog(false)
                getApiKey();
            }).catch((err) => {
                console.log(err)
            })
        }
    };

    const handleDialogClose = () => {
        setOpenDialog(false);
    };

    const getApiKey = () => {
        axios
            .get('api/account/api-key')
            .then((res) => {
                setApiKey(res.data);
            })
            .catch((err) => {
                // Handle error
            });
    };

    useEffect(() => {
        GetAccountInfo()
            .then((data) => {
                setAccountInfo(data);
            })
            .catch((error) => {
                alert(error);
            });

        GetLoggedUserInfo()
            .then((data) => {
                setUserInfo(data);
            })
            .catch((error) => {
                alert(error);
            });

        getApiKey();
    }, []);

    const handleCheckboxChange = (event) => {
        setValidForever(event.target.checked);
    };

    const handleDateTimeChange = (date) => {
        setSelectedDate(date);
    };

    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);

    return (
        <div className={cssClasses.infoWrapper}>
            <h1>User info</h1>
            <InfoRow info="Username" value={accountInfo.username}/>
            <InfoRow info="Email" value={accountInfo.email}/>
            <InfoRow info="Fullname" value={userInfo.fullname}/>
            <InfoRow info="Address" value={userInfo.address}/>


            <Box m={2}>
                <Button variant="contained" onClick={() => setOpenDialog(true)}>
                    Generate API key
                </Button>
            </Box>
            <Box m={2}>
                {apiKey !== null && (<li>API KEY: {apiKey.value}</li>)}
            </Box>
            <Box m={2}>
                {apiKey !== null && (
                    apiKey.expirationDate.slice(0, 10) !== '0001-01-01' ? (

                        <li>Expiration date: {apiKey.expirationDate.slice(0, 10)}</li>
                    ) : (
                        <li>Expiration date: VALID UNTIL NEW API KEY IS GENERATED</li>
                    )
                )}
            </Box>

            <Dialog open={openDialog} onClose={handleDialogClose}>
                <DialogTitle>Generate API Key</DialogTitle>
                <DialogContent>
                    <Box m={2}>
                        <Checkbox
                            checked={validForever}
                            onChange={handleCheckboxChange}
                            inputProps={{'aria-label': 'Valid forever checkbox'}}
                        />
                        Valid forever
                    </Box>
                    <LocalizationProvider dateAdapter={AdapterDayjs}>
                        <DatePicker
                            label="Select Date and Time"
                            value={selectedDate}
                            onChange={handleDateTimeChange}
                            disabled={validForever}
                            minDate={dayjs().add(1, 'day')}

                        />
                    </LocalizationProvider>
                </DialogContent>
                <DialogActions>
                    <Button onClick={handleDialogClose}>Close</Button>
                    <Button onClick={handleGenerateClick}>Generate</Button>
                </DialogActions>
            </Dialog>
        </div>
    );
}

export default UserInfo;
