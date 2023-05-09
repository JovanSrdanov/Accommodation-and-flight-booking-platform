import React, {useEffect, useState} from 'react';
import {Flex} from "reflexbox";
import {Box, TextField} from "@mui/material";
import interceptor from "../../interceptor/interceptor";

function GuestProfile() {

    const [user, setUser] = useState({
        username: '',
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

    const [userName, setUsername] = useState()

    const getAllUserInfo = () => {

        interceptor.get("api-2/user/logged-in-info").then(res => {
            setUser(res.data)
        })

    }

    useEffect(() => {
        getAllUserInfo();
    }, []);

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


    return (
        <>
            <div className="wrapper">
                <Flex flexDirection="column">
                    <Flex flexDirection="row" justifyContent="center">
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


            </div>

        </>
    );
}

export default GuestProfile;