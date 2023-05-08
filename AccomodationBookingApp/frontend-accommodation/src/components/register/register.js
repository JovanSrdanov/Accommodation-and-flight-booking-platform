import React, {useEffect, useState} from 'react';
import {Alert, Box, Button, FormControl, InputLabel, MenuItem, Select, TextField} from "@mui/material";
import LoginIcon from "@mui/icons-material/Login";
import HowToRegIcon from "@mui/icons-material/HowToReg";
import {useNavigate} from "react-router-dom";
import {Flex} from 'reflexbox'
import interceptor from "../../interceptor/interceptor";

function Register() {
    const navigate = useNavigate();
    const [user, setUser] = useState({
        username: '',
        password: '',
        role: 'Guest',
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
    const [showAlert, setShowAlert] = useState(false);
    const [errorMessage, setErrorMessage] = useState("No error");
    const [isDisabled, setIsDisabled] = useState(true);
    const handleInputChange = (event) => {

        const {name, value} = event.target;
        if (name.startsWith("address.")) {
            setUser((prevState) => {
                const address = {...prevState.address, [name.split(".")[1]]: value};
                return {...prevState, address};
            });
        } else {
            setUser((prevState) => ({...prevState, [name]: value}));
        }
    };

    useEffect(() => {
        console.log(user.password.length)
        console.log(user.password)
        // Check if all required fields are non-empty and password is at least 8 characters long
        const isValid =
            user.username !== "" &&
            user.password.length >= 8 &&
            user.name !== "" &&
            user.surname !== "" &&
            user.email.match(/^([\w.%+-]+)@([\w-]+\.)+([\w]{2,})$/i) &&
            user.address.country !== "" &&
            user.address.city !== "" &&
            user.address.street !== "" &&
            user.address.streetNumber !== "";
        setIsDisabled(!isValid);
    }, [user]);


    const handleAlertClose = () => {
        setShowAlert(false);
    };

    const handleRegisterClick = () => {
        // Send a POST request to the /user endpoint
        interceptor
            .post('/api-2/user', user)
            .then((response) => {
                navigate('/login');
            })
            .catch((error) => {
                console.error(error);

                setShowAlert(true);
            });
    };

    return (
        <div>
            <div className="wrapper">

                <Flex flexDirection="column">
                    <Flex>
                        <Box width={1 / 3} m={1}>
                            <TextField
                                fullWidth
                                variant="filled"
                                label="Username"
                                name="username"
                                value={user.username}
                                onChange={handleInputChange}
                            />
                        </Box>
                        <Box width={1 / 3} m={1}>
                            <TextField
                                fullWidth
                                variant="filled"
                                label="Password"
                                type="password"
                                name="password"
                                value={user.password}
                                onChange={handleInputChange}
                            />
                        </Box>
                        <Box width={1 / 3} m={1}>
                            <FormControl variant="filled" fullWidth>
                                <InputLabel id="role">Role</InputLabel>
                                <Select
                                    fullWidth
                                    variant="filled"
                                    id="role"
                                    name="role"
                                    value={user.role}
                                    label="Age"
                                    onChange={handleInputChange}
                                >
                                    <MenuItem value="Guest">Guest</MenuItem>
                                    <MenuItem value="Host">Host</MenuItem>
                                </Select>
                            </FormControl>
                        </Box>
                    </Flex>
                    <Flex>
                        <Box width={1 / 3} m={1}>
                            <TextField
                                fullWidth
                                variant="filled"
                                label="Name"
                                name="name"
                                value={user.name}
                                onChange={handleInputChange}
                            />
                        </Box>
                        <Box width={1 / 3} m={1}>
                            <TextField
                                fullWidth
                                variant="filled"
                                label="Surname"
                                name="surname"
                                value={user.surname}
                                onChange={handleInputChange}
                            />
                        </Box>
                        <Box width={1 / 3} m={1}>
                            <TextField
                                fullWidth
                                variant="filled"
                                label="E-mail"
                                type="email"
                                name="email"
                                value={user.email}
                                onChange={handleInputChange}
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
                                value={user.address.country}
                                onChange={handleInputChange}
                            />
                        </Box>
                        <Box width={1 / 4} m={1}>
                            <TextField
                                fullWidth
                                variant="filled"
                                label="City"
                                name="address.city"
                                value={user.address.city}
                                onChange={handleInputChange}
                            />
                        </Box>
                        <Box width={1 / 4} m={1}>
                            <TextField
                                fullWidth
                                variant="filled"
                                label="Street"
                                name="address.street"
                                value={user.address.street}
                                onChange={handleInputChange}
                            />
                        </Box>
                        <Box width={1 / 4} m={1}>
                            <TextField
                                fullWidth
                                variant="filled"
                                label="Street Number"
                                name="address.streetNumber"
                                value={user.address.streetNumber}
                                onChange={handleInputChange}
                            />
                        </Box>

                    </Flex>
                    <Flex flexDirection="row" justifyContent="center">

                        <p>All fields must be filled, password must be 8 characters or longer and email must be
                            in
                            valid
                            form</p>

                    </Flex>
                </Flex>
                <Flex flexDirection="row" justifyContent="space-between" alignItems="center">
                    <Box m={1}>
                        <Button variant="contained" color="warning" endIcon={<LoginIcon/>} onClick={() => {
                            navigate('/login')
                        }}>BACK TO LOGIN</Button>
                    </Box>
                    <Box m={1}>
                        <Button disabled={isDisabled} variant="contained" color="success" onClick={handleRegisterClick}
                                endIcon={<HowToRegIcon/>}>REGISTER</Button>
                    </Box>

                </Flex>

            </div>
            {showAlert && (
                <Alert sx={{width: "fit-content", margin: "10px auto"}} severity="error" onClose={handleAlertClose}>
                    {errorMessage}
                </Alert>
            )}
        </div>
    );
}

export default Register;