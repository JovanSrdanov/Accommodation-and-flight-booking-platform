import React from 'react';
import "../../pages/wrapper.css"
import {Alert, Button, TextField} from "@mui/material";
import LoginIcon from '@mui/icons-material/Login';
import HowToRegIcon from '@mui/icons-material/HowToReg';
import {useNavigate} from "react-router-dom";
import interceptor from "../../interceptor/interceptor";

function Login() {

    const navigate = useNavigate();
    const [username, setUsername] = React.useState("");
    const [password, setPassword] = React.useState("");
    const [showAlert, setShowAlert] = React.useState(false);

    const handleEmailChange = (event) => {
        setUsername(event.target.value);
    };

    const handlePasswordChange = (event) => {
        setPassword(event.target.value);
    };
    const handleLogin = async () => {
        interceptor.post('api-1/account-credentials/login', {
            username: username,
            password: password
        }).then(res => {
            const paseto = res.data.accessToken;
            const role = res.data.role;
            const expirationDate = res.data.expirationDate;
            localStorage.setItem('paseto', paseto);
            localStorage.setItem('role', role);
            localStorage.setItem('expirationDate', expirationDate);
            if (role === 'Host') {
                navigate('/host-profile')
                return
            }
            if (role === 'Guest') {
                navigate('/guest-profile')
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
                    label="Username"
                    type={"email"}
                    value={username}
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
                <Button
                    variant="contained" color="success" endIcon={<HowToRegIcon/>}
                    onClick={() => {
                        navigate('/register')
                    }}
                >REGISTER
                </Button>


            </div>
            {showAlert && (
                <Alert sx={{width: "fit-content", margin: "10px auto"}} severity="error" onClose={handleAlertClose}>
                    Invalid username or password, please try again.
                </Alert>
            )}
        </div>
    );
}

export default Login;