import RequireAuth from "./components/Authentication/RequireAuth";
import PersistLogin from "./components/Authentication/PersistLogin"
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
        <main>
            <ThemeProvider theme={darkTheme}>
                <MainNavigation/>
                <Planes/>
                <Routes>
                    <Route path="/" element={<Layout/>}>
                        {/* public rotues*/}
                        <Route path="/" element={<HomePage/>}/>
                        <Route path="register" element={<RegisterPage/>}/>
                        <Route path="unauthorized" element={<Unauthorized/>}/>
                        {/* protected routes*/}
                        {/* Ovako se stite rute - stavis rutu sa required auth i prosledis role koje su
              dozvoljene u allowerRoles */}
                        {/* Za zasticene rute ne koristiti axios, vec axiosPrivate, u njega su ugradjeni interceptori */}
                        {/*dodati const axiosPrivate = useAxiosePrivate() za svaki request koji zahteva Auth i sa time praviti pozive*/}
                        {/*za vise detalja Admin info page*/}
                        <Route element={<PersistLogin/>}>
                            <Route element={<RequireAuth allowedRoles={[ROLES.REGULAR]}/>}>
                                <Route path="flight-search" element={<FlightSearchPage/>}/>
                                <Route path="/bought-tickets" element={<BoughtTicketsPage/>}/>
                            </Route>

                <Route element={<RequireAuth allowedRoles={[ROLES.ADMIN]} />}>
                  <Route path="/create-flight" element={<CreateFlightPage/>}/>
                  <Route path="/all-flights" element={<AllFlightsPage />} />
                  <Route path="admin-info" element={<AdminInfoPage />} />
                </Route>
              </Route>
              {/* catch all */}
              <Route path="*" element={<Missing />} />
            </Route>
          </Routes>
        </ThemeProvider>
      </main>
    );
}

export default App;


