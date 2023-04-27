import React from 'react';
import "../../pages/wrapper.css"
import {Alert, Button, TextField} from "@mui/material";
import LoginIcon from '@mui/icons-material/Login';

import {useNavigate} from "react-router-dom";
import interceptor from "../../interceptor/interceptor";


function Login() {

    const navigate = useNavigate();
    const [email, setEmail] = React.useState("");
    const [password, setPassword] = React.useState("");
    const [showAlert, setShowAlert] = React.useState(false);

    const handleEmailChange = (event) => {
        setEmail(event.target.value);
    };

    const handlePasswordChange = (event) => {
        setPassword(event.target.value);
    };
    const handleLogin = async () => {
        interceptor.post('auth/login', {
            email: email,
            password: password
        }).then(res => {
            localStorage.setItem('jwt', res.data.jwt);
            const decoded = JSON.parse(atob(res.data.jwt.split('.')[1]));
            const role = decoded.role;
            if (role === 'ROLE_PKI_ADMIN') {
                navigate('/all-certificates')
                return
            }
            if (role === 'ROLE_CERTIFICATE_USER') {
                navigate('/my-certificates')
                return
            }
            if (role === 'ROLE_CERTIFICATE_USER_CHANGE_PASSWORD') {
                navigate('/change-password')
                return
            }
            navigate('/login')


        }).catch(err => {

            setShowAlert(true);
        })

    };
    const handleAlertClose = () => {
        setShowAlert(false);
    };

    return (
        <div>
            <div className="wrapper">

                <TextField
                    fullWidth
                    variant="filled"
                    label="E-mail"
                    type={"email"}
                    value={email}
                    onChange={handleEmailChange}
                />
                <TextField
                    fullWidth
                    variant="filled"
                    label="Password"
                    type="password"
                    value={password}
                    onChange={handlePasswordChange}
                />
                <Button
                    variant="contained" color="warning" endIcon={<LoginIcon/>}
                    onClick={handleLogin}
                >LOGIN
                </Button>


            </div>
            {showAlert && (
                <Alert sx={{width: "fit-content", margin: "10px auto"}} severity="error" onClose={handleAlertClose}>
                    Invalid email or password, please try again.
                </Alert>
            )}
        </div>
    );
}

export default Login;