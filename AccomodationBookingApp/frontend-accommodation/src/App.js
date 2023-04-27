import ParticlesBg from 'particles-bg'
import "./particles.css"
import {Navigate, Route, Routes, useNavigate} from "react-router-dom";
import {AppBar, Box, Button, Toolbar, Tooltip} from "@mui/material";
import HotelIcon from '@mui/icons-material/Hotel'
import HistoryIcon from '@mui/icons-material/History';
import CheckIcon from '@mui/icons-material/Check';
import RecommendOutlinedIcon from '@mui/icons-material/RecommendOutlined';
import OtherHousesOutlinedIcon from '@mui/icons-material/OtherHousesOutlined';
import AddCircleOutlineOutlinedIcon from '@mui/icons-material/AddCircleOutlineOutlined';
import ChecklistOutlinedIcon from '@mui/icons-material/ChecklistOutlined';
import PersonOutlineOutlinedIcon from '@mui/icons-material/PersonOutlineOutlined';
import LogoutOutlinedIcon from '@mui/icons-material/LogoutOutlined';
import LoginPage from "./pages/unauthenticated-pages/login-page";
import RegisterPage from "./pages/unauthenticated-pages/register-page";
import BookedPlacesPage from "./pages/guest-pages/booked-places-page";
import VisitingHistoryPage from "./pages/guest-pages/visiting-history-page";
import RecommendationsForYouPage from "./pages/guest-pages/recommendations-for-you-page";
import MyPlacesPage from "./pages/host-pages/my-places-page";
import HostAPlacePage from "./pages/host-pages/host-a-place-page";
import ReservationsAndRequestsPage from "./pages/host-pages/reservations-and-requests-page";
import GuestProfilePage from "./pages/guest-pages/guest-profile-page";
import HostProfilePage from "./pages/host-pages/host-profile-page";
import SearchAccommodationPage from "./pages/all-roles-pages/search-accommodation-page";

function App() {
    const navigate = useNavigate();
    return (
        <div className="App">
            <ParticlesBg color="#FF9021" type="cobweb" num={100} bg={true}/>
            <Box>
                <AppBar position="static">
                    <Toolbar>
                        <Tooltip title="Search for your desired accommodation" arrow>
                            <Button sx={{color: "gold", mr: 5}} onClick={() => navigate('/search-accommodation')}>
                                <HotelIcon sx={{display: {xs: 'none', md: 'flex'}, mr: 1}}/>
                                Restful Reserve
                            </Button>
                        </Tooltip>

                        <Tooltip title="View the places you have booked" arrow>
                            <Button sx={{color: 'inherit'}} onClick={() => navigate('/booked-places')}>
                                <CheckIcon sx={{display: {xs: 'none', md: 'flex'}, mr: 1}}/>
                                Booked places
                            </Button>
                        </Tooltip>

                        <Tooltip title="View the places you have visited" arrow>
                            <Button sx={{color: 'inherit'}} onClick={() => navigate('/visiting-history')}>
                                <HistoryIcon sx={{display: {xs: 'none', md: 'flex'}, mr: 1}}/>
                                Visiting history
                            </Button>
                        </Tooltip>

                        <Tooltip title="Our recommendations for you based on your preferences" arrow>
                            <Button sx={{color: 'inherit'}} onClick={() => navigate('/recommendations-for-you')}>
                                <RecommendOutlinedIcon sx={{display: {xs: 'none', md: 'flex'}, mr: 1}}/>
                                Recommendations for you
                            </Button>
                        </Tooltip>

                        <Tooltip title="Places you host" arrow>
                            <Button sx={{color: 'inherit'}} onClick={() => navigate('/my-places')}>
                                <OtherHousesOutlinedIcon sx={{display: {xs: 'none', md: 'flex'}, mr: 1}}/>
                                My places
                            </Button>
                        </Tooltip>
                        <Tooltip title="Host a new place that you want to rent" arrow>
                            <Button sx={{color: 'inherit'}} onClick={() => navigate('/host-a-place')}>
                                <AddCircleOutlineOutlinedIcon sx={{display: {xs: 'none', md: 'flex'}, mr: 1}}/>
                                Host a place
                            </Button>
                        </Tooltip>
                        <Tooltip title="View all reservations and request for reservations" arrow>
                            <Button sx={{color: 'inherit'}} onClick={() => navigate('/reservations-and-requests')}>
                                <ChecklistOutlinedIcon sx={{display: {xs: 'none', md: 'flex'}, mr: 1}}/>
                                Reservations and requests
                            </Button>
                        </Tooltip>
                        <Tooltip title="Your basic informations" arrow>
                            <Button color="info" sx={{marginLeft: 'auto'}}
                                    onClick={() => navigate('/guest-profile') /*TODO  OVDE DODATI PROVERU ZA KOJI JE TACNO PROFIL*/}>
                                <PersonOutlineOutlinedIcon sx={{display: {xs: 'none', md: 'flex'}, mr: 1}}/>
                                My profile
                            </Button>
                        </Tooltip>
                        <Tooltip title="Log out of the sistem" arrow>
                            <Button color="error">
                                <LogoutOutlinedIcon sx={{display: {xs: 'none', md: 'flex'}, mr: 1}}/>
                                Log out
                            </Button>
                        </Tooltip>
                    </Toolbar>
                </AppBar>
                <Routes>
                    <Route path="/login" element={<LoginPage/>}/>
                    <Route path="/register" element={<RegisterPage/>}/>
                    <Route path="/search-accommodation" element={<SearchAccommodationPage/>}/>
                    
                    <Route path="/booked-places" element={<BookedPlacesPage/>}/>
                    <Route path="/visiting-history" element={<VisitingHistoryPage/>}/>
                    <Route path="/recommendations-for-you" element={<RecommendationsForYouPage/>}/>

                    <Route path="/my-places" element={<MyPlacesPage/>}/>
                    <Route path="/host-a-place" element={<HostAPlacePage/>}/>
                    <Route path="/reservations-and-requests" element={<ReservationsAndRequestsPage/>}/>

                    <Route path="/guest-profile" element={<GuestProfilePage/>}/>
                    <Route path="/host-profile" element={<HostProfilePage/>}/>


                    <Route path="*" element={<Navigate to="/search-accommodation"/>}/>

                </Routes>
            </Box>
        </div>


    );
}

export default App;
