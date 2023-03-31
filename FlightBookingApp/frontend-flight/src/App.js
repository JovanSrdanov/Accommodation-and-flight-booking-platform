import {Route, Routes} from "react-router-dom";

import HomePage from "./pages/unaunthenticated/home-page";
import FlightSearchPage from "./pages/customer/flight-search-page";
import MainNavigation from "./components/layout/MainNavigation";
import RegisterPage from "./pages/unaunthenticated/register-page"
import AdminInfoPage from "./pages/admin/admin-info-page";

import {createTheme, ThemeProvider} from '@mui/material/styles';
import {Layout} from "./components/layout/Layout";
import Unauthorized from "./pages/unaunthenticated/Unauthorized";
import Missing from "./pages/unaunthenticated/Missing";
import AllFlightsPage from "./pages/admin/all-flights-page";
import CreateFlightPage from "./pages/admin/create-flight-page";
import Planes from "./components/planes/planes";
import HackerHeaders from "./components/hackerHeaders/hackerHeaders";
import BoughtTicketsPage from "./pages/customer/bought-tickets-page";
import RequireAuth from "./components/Authentication/RequireAuth";


const darkTheme = createTheme({
    palette: {
        mode: 'dark',
    }
});

const ROLES = {
    'ADMIN': 0,
    'REGULAR': 1
}

function App() {
    HackerHeaders();
    return (
        <div>
            <ThemeProvider theme={darkTheme}>
                <Planes/>
                <MainNavigation/>
                <Routes>
                    <Route path="/" element={<Layout/>}>

                        <Route element={<RequireAuth allowedRoles={[ROLES.REGULAR]}/>}>
                            <Route path="flight-search" element={<FlightSearchPage/>}/>
                        </Route>


                        <Route path="/" element={<HomePage/>}/>
                        <Route path="register" element={<RegisterPage/>}/>
                        <Route path="unauthorized" element={<Unauthorized/>}/>
                        <Route path="flight-search" element={<FlightSearchPage/>}/>
                        <Route path="/all-flights" element={<AllFlightsPage/>}/>
                        <Route path="/create-flight" element={<CreateFlightPage/>}/>
                        <Route path="/bought-tickets" element={<BoughtTicketsPage/>}/>
                        <Route path="admin-info" element={<AdminInfoPage/>}/>


                        <Route path="*" element={<Missing/>}/>
                    </Route>
                </Routes>
            </ThemeProvider>
        </div>
    );
}

export default App;


