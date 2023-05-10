import ParticlesBg from 'particles-bg'
import "./particles.css"
import {Navigate, Route, Routes, useNavigate} from "react-router-dom";
import {AppBar, Box, Button, Toolbar, Tooltip} from "@mui/material";
import HotelIcon from '@mui/icons-material/Hotel'
import PersonOutlineOutlinedIcon from '@mui/icons-material/PersonOutlineOutlined';
import LogoutOutlinedIcon from '@mui/icons-material/LogoutOutlined';
import BookedPlacesPage from "./pages/guest-pages/booked-places-page";
import VisitingHistoryPage from "./pages/guest-pages/visiting-history-page";
import RecommendationsForYouPage from "./pages/guest-pages/recommendations-for-you-page";
import MyPlacesPage from "./pages/host-pages/my-places-page";
import HostAPlacePage from "./pages/host-pages/host-a-place-page";
import ReservationsAndRequestsPage from "./pages/host-pages/reservations-and-requests-page";
import ProfilePage from "./pages/guest-pages/profile-page";
import SearchAccommodationPage from "./pages/all-roles-pages/search-accommodation-page";
import React, {useEffect, useState} from "react";
import HistoryIcon from '@mui/icons-material/History';
import CheckIcon from '@mui/icons-material/Check';
import RecommendOutlinedIcon from '@mui/icons-material/RecommendOutlined';
import OtherHousesOutlinedIcon from '@mui/icons-material/OtherHousesOutlined';
import AddCircleOutlineOutlinedIcon from '@mui/icons-material/AddCircleOutlineOutlined';
import ChecklistOutlinedIcon from '@mui/icons-material/ChecklistOutlined';
import LoginPage from "./pages/unauthenticated-pages/login-page";
import RegisterPage from "./pages/unauthenticated-pages/register-page";
import LoginIcon from '@mui/icons-material/Login';
import HowToRegIcon from '@mui/icons-material/HowToReg';
import TravelExploreIcon from '@mui/icons-material/TravelExplore';

function App() {


    const [role, setRole] = useState(null);
    const navigate = useNavigate();

    const pasetoExpirationRole = () => {
        const expirationDateStr = localStorage.getItem('expirationDate');
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
                localStorage.removeItem('role');
                localStorage.removeItem('expirationDate');
                setRole(null);
            } else {
                setRole(localStorage.getItem('role'))
            }
        } else {
            localStorage.removeItem('paseto');
            localStorage.removeItem('role');
            localStorage.removeItem('expirationDate');
            setRole(null);
        }
    };


    useEffect(() => {
        pasetoExpirationRole();
    }, [navigate]);

    const handleLogout = () => {
        localStorage.removeItem('paseto');
        localStorage.removeItem('role');
        localStorage.removeItem('expirationDate');
        setRole(null);
        navigate('/login');
    };


    const IS_HOST = role === 'Host';
    const IS_GUEST = role === 'Guest';


    return (


        <div>


            <ParticlesBg color="#FF9021" type="cobweb" num={100} bg={true}/>
            <Box>
                <AppBar position="static">
                    <Toolbar>
                        <Tooltip title="Search for your desired accommodation" arrow>
                            <Button sx={{color: "gold", mr: 5}} onClick={() => navigate('/search-accommodation')}
                                    endIcon={<TravelExploreIcon/>} startIcon={<HotelIcon/>}>

                                Restful Reserve

                            </Button>
                        </Tooltip>

                        {IS_GUEST && (
                            <>
                                <Tooltip title="View the places you have booked" arrow>
                                    <Button sx={{color: 'inherit'}}
                                            startIcon={<CheckIcon/>}

                                            onClick={() => navigate('/booked-places')}>

                                        Booked places
                                    </Button>
                                </Tooltip>

                                <Tooltip title="View the places you have visited" arrow>
                                    <Button startIcon={<HistoryIcon/>} sx={{color: 'inherit'}}
                                            onClick={() => navigate('/visiting-history')}>

                                        Visiting history
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

                        {IS_HOST && (
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

                        {(IS_GUEST || IS_HOST) && (
                            <>
                                <Tooltip title="Your informations" arrow>
                                    <Button color="info" sx={{marginLeft: 'auto'}}
                                            startIcon={<PersonOutlineOutlinedIcon/>}
                                            onClick={() => {
                                                navigate('/profile');
                                            }}>

                                        My profile
                                    </Button>
                                </Tooltip>
                                <Tooltip title="Log out of the sistem" arrow>
                                    <Button color="error" onClick={handleLogout} startIcon={<LogoutOutlinedIcon/>}>

                                        Log out
                                    </Button>
                                </Tooltip>
                            </>
                        )}

                        {(!IS_GUEST && !IS_HOST) && (
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

                    {IS_GUEST && (
                        <>
                            <Route path="/booked-places" element={<BookedPlacesPage/>}/>
                            <Route path="/visiting-history" element={<VisitingHistoryPage/>}/>
                            <Route path="/recommendations-for-you" element={<RecommendationsForYouPage/>}/>

                            <Route path="/profile" element={<ProfilePage/>}/>
                            <Route path="/search-accommodation" element={<SearchAccommodationPage/>}/>
                        </>
                    )}

                    {IS_HOST && (
                        <>
                            <Route path="/my-places" element={<MyPlacesPage/>}/>
                            <Route path="/host-a-place" element={<HostAPlacePage/>}/>
                            <Route path="/reservations-and-requests" element={<ReservationsAndRequestsPage/>}/>
                            <Route path="/profile" element={<ProfilePage/>}/>
                            <Route path="/search-accommodation" element={<SearchAccommodationPage/>}/>
                        </>
                    )}

                    {!IS_HOST && !IS_GUEST && (
                        <>
                            <Route path="/login" element={<LoginPage/>}/>
                            <Route path="/register" element={<RegisterPage/>}/>
                        </>
                    )}
                    <Route path="/search-accommodation" element={<SearchAccommodationPage/>}/>
                    <Route path="/" element={<Navigate to="/search-accommodation"/>}/>
                    <Route path="*" element={<Navigate to="/search-accommodation"/>}/>
                </Routes>
            </Box>
        </div>


    );
}

export default App;