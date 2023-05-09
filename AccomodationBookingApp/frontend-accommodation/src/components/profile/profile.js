import React, {useEffect, useState} from 'react';
import {Flex} from "reflexbox";
import {Box, Button, Dialog, DialogActions, DialogContent, DialogTitle, TextField} from "@mui/material";
import interceptor from "../../interceptor/interceptor";

function Profile() {


    const [username, setUsername] = useState("")
    const [oldPassword, setOldPassword] = useState("")
    const [newPassword, setNewPassword] = useState("")
    const [passwordDialogShow, setPasswordDialogShow] = useState(false)
    const [successDialogShow, setSuccessDialogShow] = useState(false)
    const [usernameTakenDialogShow, setUsernameTakenDialogShow] = useState(false)
    const [userInfo, setUserInfo] = useState({
        name: '',
        surname: '',
        email: '',
        address: {
            country: '',
            city: '',
            street: '',
            streetNumber: ''
        }
    });

    const handleOldPasswordChange = (event) => {
        setOldPassword(event.target.value);
    };

    const handleNewPasswordChange = (event) => {
        setNewPassword(event.target.value);
    };

    const closePasswordDialog = () => {
        interceptor.put('api-1/account-credentials/change-password', {oldPassword, newPassword})
            .then((response) => {
                setOldPassword("");
                setNewPassword("");
                setPasswordDialogShow(false);
                setSuccessDialogShow(true)
            })
            .catch((error) => {
                // Handle the error here, such as showing an error message
            });


    };


    const handleUpdateUsernameClick = () => {
        interceptor.put('api-1/account-credentials/change-username', {username})
            .then((response) => {
                setSuccessDialogShow(true)
            })
            .catch((error) => {
                setUsernameTakenDialogShow(true)
            });
    };


    const getAllUserInfo = () => {

        interceptor.get("api-2/user/logged-in-info").then(res => {
            let user = {
                name: res.data.name,
                surname: res.data.surname,
                email: res.data.email,
                address: {
                    country: res.data.address.country,
                    city: res.data.address.city,
                    street: res.data.address.street,
                    streetNumber: res.data.address.streetNumber
                }
            }

            setUserInfo(user)

            setUsername(res.data.username)
        }).catch(err => console.log(err));

    }

    useEffect(() => {
        getAllUserInfo();
    }, []);

    const handleBasicUserInfoInputChange = (event) => {

        const {name, value} = event.target;
        if (name.startsWith("address.")) {
            setUserInfo((prevState) => {
                const address = {...prevState.address, [name.split(".")[1]]: value};
                return {...prevState, address};
            });
        } else {
            setUserInfo((prevState) => ({...prevState, [name]: value}));
        }
    };


    const handleUsernameChange = (event) => {
        setUsername(event.target.value);

    }

    const UpdateBasicUserInfo = () => {
        interceptor.put('api-1/user-profile', userInfo)
            .then((response) => {
                setSuccessDialogShow(true)
            })
            .catch((error) => {
                // Handle the error here, such as showing an error message
            });
    }

    const DeleteProfile = () => {
        // interceptor.put('api-2/user-profile', userInfo)
        //     .then((response) => {
        //
        //     })
        //     .catch((error) => {
        //
        //     });
    };
    const handleClose = () => {
        setSuccessDialogShow(false)
    };

    const usernameTakenDialogClose = () => {
        setUsernameTakenDialogShow(false)
    };
    return (
        <>
            <Dialog onClose={handleClose} open={successDialogShow}>
                <DialogTitle>Update Successful!</DialogTitle>
                <DialogActions>
                    <Button onClick={handleClose}
                            variant="contained"
                    >
                        Close
                    </Button>
                </DialogActions>
            </Dialog>


            <Dialog onClose={usernameTakenDialogClose} open={usernameTakenDialogShow}>
                <DialogTitle>That username is already taken</DialogTitle>
                <DialogActions>
                    <Button onClick={usernameTakenDialogClose}
                            variant="contained"
                    >
                        Close
                    </Button>
                </DialogActions>
            </Dialog>

            <Dialog open={passwordDialogShow} onClose={() => setPasswordDialogShow(false)}>
                <DialogTitle>Change password</DialogTitle>
                <DialogContent>
                    <Box m={1}>
                        <TextField
                            fullWidth
                            variant="filled"
                            label="Old password"
                            type="password"
                            name="oldPassword"
                            value={oldPassword}
                            onChange={handleOldPasswordChange}
                        />
                    </Box>

                    <Box m={1}>
                        <TextField m={1}
                                   fullWidth
                                   variant="filled"
                                   label="New password"
                                   type="password"
                                   name="newPassword"
                                   value={newPassword}
                                   onChange={handleNewPasswordChange}
                        />
                    </Box>
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => {
                        setOldPassword("");
                        setNewPassword("");
                        setPasswordDialogShow(false);
                    }} color="error"
                            variant="outlined">
                        Close
                    </Button>
                    <Button onClick={closePasswordDialog} disabled={oldPassword.length < 8 || newPassword.length < 8}
                            color="warning"
                            variant="contained">
                        Change password
                    </Button>
                </DialogActions>
            </Dialog>

            <div className="wrapper">
                <Flex flexDirection="column">


                    <Flex flexDirection="column" justifyContent="center" alignItems="center">
                        <p>All fields must be filled, password must be 8 characters or longer and email must be
                            in
                            valid
                            form</p>
                        <hr
                            style={{
                                width: "100%",
                                borderColor: 'black',
                            }}
                        />
                        <Box m={1}>
                            Change username
                        </Box>
                        <Box width={1 / 3} m={1}>
                            <TextField
                                fullWidth
                                variant="filled"
                                label="Username"
                                name="username"
                                value={username}
                                onChange={handleUsernameChange}
                            />
                        </Box>


                        <Box width={1 / 3} m={1}
                             display="flex"
                             justifyContent="center"
                             alignItems="center">
                            <Button
                                onClick={handleUpdateUsernameClick}
                                fullWidth
                                color="warning"
                                variant="contained"
                                disabled={username === ''}>

                                Update username

                            </Button>
                        </Box>
                        <hr
                            style={{
                                width: "100%",
                                borderColor: 'black',
                            }}
                        />


                    </Flex>
                    <Flex flexDirection="column" justifyContent="center" alignItems="center">
                        <Box m={1}>
                            Change basic information
                        </Box>
                    </Flex>

                    <Flex>

                        <Box width={1 / 3} m={1}>
                            <TextField
                                fullWidth
                                variant="filled"
                                label="Name"
                                name="name"
                                value={userInfo.name}
                                onChange={handleBasicUserInfoInputChange}
                            />
                        </Box>
                        <Box width={1 / 3} m={1}>
                            <TextField
                                fullWidth
                                variant="filled"
                                label="Surname"
                                name="surname"
                                value={userInfo.surname}
                                onChange={handleBasicUserInfoInputChange}
                            />
                        </Box>
                        <Box width={1 / 3} m={1}>
                            <TextField
                                fullWidth
                                variant="filled"
                                label="E-mail"
                                type="email"
                                name="email"
                                value={userInfo.email}
                                onChange={handleBasicUserInfoInputChange}
                            />
                        </Box>
                    </Flex>
                    <Flex>
                        <Box width={1 / 4} m={1}>
                            <TextField
                                fullWidth
                                variant="filled"
                                label="Country"
                                name="address.country"
                                value={userInfo.address.country}
                                onChange={handleBasicUserInfoInputChange}
                            />
                        </Box>
                        <Box width={1 / 4} m={1}>
                            <TextField
                                fullWidth
                                variant="filled"
                                label="City"
                                name="address.city"
                                value={userInfo.address.city}
                                onChange={handleBasicUserInfoInputChange}
                            />
                        </Box>
                        <Box width={1 / 4} m={1}>
                            <TextField
                                fullWidth
                                variant="filled"
                                label="Street"
                                name="address.street"
                                value={userInfo.address.street}
                                onChange={handleBasicUserInfoInputChange}
                            />
                        </Box>
                        <Box width={1 / 4} m={1}>
                            <TextField
                                fullWidth
                                variant="filled"
                                label="Street Number"
                                name="address.streetNumber"
                                value={userInfo.address.streetNumber}
                                onChange={handleBasicUserInfoInputChange}
                            />
                        </Box>


                    </Flex>


                    <Flex flexDirection="column" justifyContent="center" alignItems="center">
                        <Box width={1 / 3} m={1}>
                            <Button
                                fullWidth
                                color="warning"
                                variant="contained"
                                onClick={UpdateBasicUserInfo}

                                disabled={!(Object.values(userInfo).every(val => val !== '') && /\S+@\S+\.\S+/.test(userInfo.email))}>


                                Update basic information

                            </Button>
                        </Box>
                    </Flex>
                    <hr
                        style={{
                            width: "100%",
                            borderColor: 'black',
                        }}
                    />
                    <Flex flexDirection="column" justifyContent="center" alignItems="center">
                        <Box m={1}>
                            Change password (old password must be re-entered so that it can be changed)
                        </Box>
                        <Box width={1 / 3} m={1}>
                            <Button
                                fullWidth
                                color="warning"
                                variant="contained"
                                onClick={() => {
                                    setPasswordDialogShow(true)
                                }}
                            >

                                Change Password

                            </Button>
                        </Box>


                        <hr
                            style={{
                                width: "100%",
                                borderColor: 'black',
                            }}
                        />
                    </Flex>
                    <Flex flexDirection="column" justifyContent="center" alignItems="center">
                        <Box width={1 / 4} m={1}>
                            <Button
                                onClick={DeleteProfile}
                                fullWidth
                                color="error"
                                variant="contained">

                                Delete account

                            </Button>
                        </Box>
                    </Flex>

                </Flex>


            </div>

        </>
    );
}

export default Profile;