import ParticlesBg from 'particles-bg'
import "./particles.css"
import {Navigate, Route, Routes, useNavigate} from "react-router-dom";
import {
    Alert,
    AppBar,
    Box,
    Button,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    FormControlLabel,
    Snackbar,
    Switch,
    Toolbar,
    Tooltip
} from "@mui/material";
import HotelIcon from '@mui/icons-material/Hotel'
import PersonOutlineOutlinedIcon from '@mui/icons-material/PersonOutlineOutlined';
import LogoutOutlinedIcon from '@mui/icons-material/LogoutOutlined';
import MyReservationsPage from "./pages/guest-pages/my-reservations-page";
import RecommendationsForYouPage from "./pages/guest-pages/recommendations-for-you-page";
import MyPlacesPage from "./pages/host-pages/my-places-page";
import HostAPlacePage from "./pages/host-pages/host-a-place-page";
import ReservationsAndRequestsPage from "./pages/host-pages/reservations-and-requests-page";
import ProfilePage from "./pages/guest-pages/profile-page";
import SearchAndFilterAccommodationsPage from "./pages/all-roles-pages/search-and-filter-accommodations-page";
import React, {useEffect, useRef, useState} from "react";
import HistoryIcon from '@mui/icons-material/History';
import RecommendOutlinedIcon from '@mui/icons-material/RecommendOutlined';
import OtherHousesOutlinedIcon from '@mui/icons-material/OtherHousesOutlined';
import AddCircleOutlineOutlinedIcon from '@mui/icons-material/AddCircleOutlineOutlined';
import ChecklistOutlinedIcon from '@mui/icons-material/ChecklistOutlined';
import LoginPage from "./pages/unauthenticated-pages/login-page";
import RegisterPage from "./pages/unauthenticated-pages/register-page";
import LoginIcon from '@mui/icons-material/Login';
import HowToRegIcon from '@mui/icons-material/HowToReg';
import TravelExploreIcon from '@mui/icons-material/TravelExplore';
import {Flex} from "reflexbox";
import interceptor from "./interceptor/interceptor";
import EditNotificationsIcon from '@mui/icons-material/EditNotifications';

function App() {

    const navigate = useNavigate();
    const [notificationSnackBar, setNotificationSnackBar] = useState(false);
    const [message, setMessage] = useState('');
    const openWebSocket = () => {

        const paseto = localStorage.getItem('paseto');
        if (!paseto) {
            localStorage.removeItem('paseto');
            return null
        }

        if (websocketOpen) {
            return
         
        }
        setWebsocketOpen(true);

        const ws = new WebSocket(`ws://localhost:8000/ws?authorization=${encodeURIComponent(paseto)}`);

        ws.onopen = () => {
            console.log('WEBSOCKET CONNECTION ESTABLISHED');

        };
        ws.onmessage = (event) => {

            setMessage(event.data.toUpperCase());
            setNotificationSnackBar(true)
        };

        ws.onclose = () => {
            console.log('WEBSOCKET CONNECTION CLOSED');
            setWebsocketOpen(false);
        };

        return () => {
            ws.close();
        };
    };


    const pasetoExpirationRole = () => {
        const paseto = localStorage.getItem('paseto');
        if (!paseto) {
            localStorage.removeItem('paseto');
            return null
        }
        const footer = paseto.split(".")[3];
        const decodedFooter = JSON.parse(atob(footer));
        const roleAndExp = decodedFooter.RoleAndExp;


        const regex = /role:(.*), expiration date: (.*)/;
        const matches = roleAndExp.match(regex);
        const role = matches[1];
        const expirationDateStr = matches[2];


        if (expirationDateStr) {

            const currentTime = new Date()
            const localOffset = currentTime.getTimezoneOffset() // in minutes

            let expirationDate = new Date(expirationDateStr.split(" ")[0] + "T" + expirationDateStr.split(" ")[1]);
            const utcOffset = expirationDateStr.substring(23, 28);
            const sign = utcOffset.substring(0, 1);
            const hours = parseInt(utcOffset.substring(1, 3));

            let expOffset = 0;
            if (sign === "+") {
                expOffset = hours;
            } else if (sign === "-") {
                expOffset = -hours;
            }
            const correctingOffset = expOffset - localOffset
            expirationDate = new Date(expirationDate.getTime() + correctingOffset * 60 * 1000)

            if (expirationDate < currentTime) {
                console.log("Token expired")
                localStorage.removeItem('paseto');
                return null;

            } else {

                if (role === "0") {

                    return "Host";

                } else if (role === "1") {

                    return "Guest";

                } else {
                    localStorage.removeItem('paseto');
                    return null;
                }
            }
        } else {
            localStorage.removeItem('paseto');
            return null;
        }
    };

    const ROLE = pasetoExpirationRole();
    if (ROLE === null) {
        interceptor
          .get("api-2/user/is-unique-visitor")
          .then((res) => {
          })
          .catch((err) => {
            console.log(err);
          });
    }

    const handleLogout = () => {
        localStorage.removeItem('paseto');
        window.location.href = "/login";
    };

    const [selectedItem, setSelectedItem] = useState({
        RequestMade: false,
        ReservationCanceled: true,
        HostRatingGiven: false,
        AccommodationRatingGiven: false,
        ProminentHost: false,
        HostResponded: false,
    });

    const handleClose = () => {
        setOpen(false);
    };

    const handleSwitchChange = (event) => {
        setSelectedItem((prevState) => ({
            ...prevState,
            [event.target.name]: event.target.checked,
        }));
    };

    const [open, setOpen] = useState(false);
    const isFirstRender = useRef(true);
    const isClickOpen = useRef(false);

    const handleClickOpen = () => {
        isClickOpen.current = true;
        interceptor
            .get("/api-1/notification/get-my")
            .then((res) => {
                setSelectedItem(res.data);
                setOpen(true);
            })
            .catch((err) => {
                console.log(err);
            });
    };
    const [websocketOpen, setWebsocketOpen] = useState(false);
    useEffect(() => {

        openWebSocket();


    }, [navigate]);


    useEffect(() => {

        if (isFirstRender.current || isClickOpen.current) {
            isFirstRender.current = false;
            isClickOpen.current = false;
            return;
        }

        interceptor
            .put("/api-1/notification/update-my", selectedItem)
            .then((res) => {
                // Handle the response
            })
            .catch((err) => {
                // Handle the error
            });
    }, [selectedItem]);

    return (
        <>
            <Snackbar
                anchorOrigin={{
                    vertical: 'bottom',
                    horizontal: 'right',
                }}
                open={notificationSnackBar}
                autoHideDuration={5000}
                onClose={() => setNotificationSnackBar(false)}
                message={message}
            >
                <Alert onClose={handleClose} severity="success" sx={{width: '100%'}}>
                    {message}
                </Alert>
            </Snackbar>

            <Dialog open={open} onClose={handleClose}>
                <DialogTitle>Choose about what you want to be notified</DialogTitle>
                <DialogContent>

                    <Flex flexDirection="column">
                        {ROLE === 'Host' && (
                            <>
                                <FormControlLabel
                                    control={
                                        <Switch
                                            checked={selectedItem.RequestMade}
                                            onChange={handleSwitchChange}
                                            name="RequestMade"
                                            color="success"
                                        />
                                    }
                                    label="Request Made"
                                />
                                <FormControlLabel
                                    control={
                                        <Switch
                                            checked={selectedItem.ReservationCanceled}
                                            onChange={handleSwitchChange}
                                            name="ReservationCanceled"
                                            color="success"
                                        />
                                    }
                                    label="Reservation Canceled"
                                />
                                <FormControlLabel
                                    control={
                                        <Switch
                                            checked={selectedItem.HostRatingGiven}
                                            onChange={handleSwitchChange}
                                            name="HostRatingGiven"
                                            color="success"
                                        />
                                    }
                                    label="Host Rating Given"
                                />
                                <FormControlLabel
                                    control={
                                        <Switch
                                            checked={selectedItem.AccommodationRatingGiven}
                                            onChange={handleSwitchChange}
                                            name="AccommodationRatingGiven"
                                            color="success"
                                        />
                                    }
                                    label="Accommodation Rating Given"
                                />
                                <FormControlLabel
                                    control={
                                        <Switch
                                            checked={selectedItem.ProminentHost}
                                            onChange={handleSwitchChange}
                                            name="ProminentHost"
                                            color="success"
                                        />
                                    }
                                    label="Prominent Host"
                                />
                            </>
                        )}
                        {ROLE === 'Guest' && (
                            <FormControlLabel
                                control={
                                    <Switch
                                        checked={selectedItem.HostResponded}
                                        onChange={handleSwitchChange}
                                        name="HostResponded"
                                        color="success"
                                    />
                                }
                                label="Host Responded"
                            />
                        )}
                    </Flex>
                </DialogContent>
                <DialogActions>
                    <Button onClick={handleClose} color="info" variant="outlined">
                        Close
                    </Button>

                </DialogActions>
            </Dialog>


            <ParticlesBg color="#FF9021" type="cobweb" num={100} bg={true}/>
            <Box>
                <AppBar position="static">
                    <Toolbar>
                        <Tooltip title="Search for your desired accommodation" arrow>
                            <Button sx={{color: "gold", mr: 5}}
                                    onClick={() => navigate('/search-and-filter-accommodations')}
                                    endIcon={<TravelExploreIcon/>} startIcon={<HotelIcon/>}>

                                Restful Reserve

                            </Button>
                        </Tooltip>

                        {ROLE === 'Guest' && (
                            <>

                                <Tooltip title="View your reservations" arrow>
                                    <Button startIcon={<HistoryIcon/>} sx={{color: 'inherit'}}
                                            onClick={() => navigate('/my-reservations')}>

                                        My Reservations
                                    </Button>
                                </Tooltip>

                                <Tooltip title="Our recommendations for you based on your preferences" arrow>
                                    <Button startIcon={<RecommendOutlinedIcon/>} sx={{color: 'inherit'}}
                                            onClick={() => navigate('/recommendations-for-you')}>

                                        Recommendations for you
                                    </Button>
                                </Tooltip>

                            </>
                        )}

                        {ROLE === 'Host' && (
                            <>

                                <Tooltip title="Places you host" arrow>
                                    <Button startIcon={<OtherHousesOutlinedIcon/>} sx={{color: 'inherit'}}
                                            onClick={() => navigate('/my-places')}>

                                        My places
                                    </Button>
                                </Tooltip>
                                <Tooltip title="Host a new place that you want to rent" arrow>
                                    <Button startIcon={<AddCircleOutlineOutlinedIcon/>} sx={{color: 'inherit'}}
                                            onClick={() => navigate('/host-a-place')}>

                                        Host a place
                                    </Button>
                                </Tooltip>
                                <Tooltip title="View all reservations and request for reservations" arrow>
                                    <Button startIcon={<ChecklistOutlinedIcon/>} sx={{color: 'inherit'}}
                                            onClick={() => navigate('/reservations-and-requests')}>

                                        Reservations and requests
                                    </Button>
                                </Tooltip>
                            </>
                        )}

                        {(ROLE === 'Guest' || ROLE === 'Host') && (
                            <>
                                <Tooltip title="Your notifications" arrow>
                                    <Button startIcon={<EditNotificationsIcon/>} sx={{marginLeft: 'auto'}}
                                            color="success"
                                            onClick={handleClickOpen}>
                                        Notifications
                                    </Button>
                                </Tooltip>


                                <Tooltip title="Your informations" arrow>
                                    <Button color="info"
                                            startIcon={<PersonOutlineOutlinedIcon/>}
                                            onClick={() => {
                                                navigate('/profile');
                                            }}>

                                        My profile
                                    </Button>
                                </Tooltip>
                                <Tooltip title="Log out of the system" arrow>
                                    <Button color="error" onClick={handleLogout} startIcon={<LogoutOutlinedIcon/>}>

                                        Log out
                                    </Button>
                                </Tooltip>
                            </>
                        )}

                        {(!(ROLE === 'Guest') && !(ROLE === 'Host')) && (
                            <>   <Tooltip title="View all reservations and request for reservations" arrow>
                                <Button color="warning" sx={{marginLeft: 'auto'}} startIcon={<LoginIcon/>}
                                        onClick={() => navigate('/login')}>

                                    Log in
                                </Button>
                            </Tooltip>
                                <Tooltip title="View all reservations and request for reservations" arrow>
                                    <Button color="success" startIcon={<HowToRegIcon/>}
                                            onClick={() => navigate('/register')}>

                                        Register
                                    </Button>
                                </Tooltip>
                            </>
                        )}
                    </Toolbar>

                </AppBar>
                <Routes>


                    {ROLE === 'Guest' && (
                        <>

                            <Route path="/my-reservations" element={<MyReservationsPage/>}/>

                            <Route path="/recommendations-for-you" element={<RecommendationsForYouPage/>}/>
                            <Route path="/profile" element={<ProfilePage/>}/>
                            <Route path="/search-and-filter-accommodations"
                                   element={<SearchAndFilterAccommodationsPage canBuy={true}/>}/>
                            <Route path="/*" element={<Navigate to="/search-and-filter-accommodations"/>}/>
                        </>
                    )}

                    {ROLE === 'Host' && (
                        <>
                            <Route path="/my-places" element={<MyPlacesPage/>}/>
                            <Route path="/host-a-place" element={<HostAPlacePage/>}/>
                            <Route path="/reservations-and-requests" element={<ReservationsAndRequestsPage/>}/>
                            <Route path="/profile" element={<ProfilePage/>}/>
                            <Route path="/search-and-filter-accommodations"
                                   element={<SearchAndFilterAccommodationsPage canBuy={false}/>}/>
                            <Route path="/*" element={<Navigate to="/search-and-filter-accommodations"/>}/>
                        </>
                    )}

                    {ROLE === null && (
                        <>
                            <Route path="/login" element={<LoginPage/>}/>
                            <Route path="/register" element={<RegisterPage/>}/>
                            <Route path="/search-and-filter-accommodations"
                                   element={<SearchAndFilterAccommodationsPage canBuy={false}/>}/>
                            <Route path="/*" element={<Navigate to="/search-and-filter-accommodations"/>}/>
                        </>
                    )}

                </Routes>
            </Box>
        </>


    );
}

export default App;
