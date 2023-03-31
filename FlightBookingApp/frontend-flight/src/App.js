import RequireAuth from "./components/Authentication/RequireAuth";
import {Route, Routes} from "react-router-dom";

import HomePage from "./pages/unaunthenticated/home-page";
import FlightSearchPage from "./pages/customer/flight-search-page";
import MainNavigation from "./components/layout/MainNavigation";
import RegisterPage from "./pages/unaunthenticated/register-page"
import AdminInfoPage from "./pages/admin/admin-info-page";

import {createTheme, ThemeProvider} from '@mui/material/styles';
import { Layout } from "./components/layout/Layout";
import Unauthorized from "./pages/unaunthenticated/Unauthorized";
import Missing from "./pages/unaunthenticated/Missing";

const darkTheme = createTheme({
    palette: {
        mode: 'dark',
    },
});

const ROLES = {
  'ADMIN': 0, 
  'REGULAR': 1
}

function App() {
    return (
      <main>
        <ThemeProvider theme={darkTheme}>
          <MainNavigation />
          <Routes>
            <Route path="/" element={<Layout />}>
              {/* public rotues*/}
              <Route path="/" element={<HomePage />} />
              <Route path="register" element={<RegisterPage />} />
              <Route path="unauthorized" element={<Unauthorized />} />

              {/* protected routes*/}
              {/* Ovako se stite rute - stavis rutu sa required auth i prosledis role koje su 
              dozvoljene u allowerRoles */}
              {/* Za zasticene rute ne koristiti axios, vec axiosPrivate, u njega su ugradjeni interceptori */}
              <Route
                element={
                  <RequireAuth allowedRoles={[ROLES.REGULAR, ROLES.ADMIN]} />
                }
              >
                <Route path="flight-search" element={<FlightSearchPage />} />
              </Route>

              <Route
                element={
                  <RequireAuth allowedRoles={[ROLES.ADMIN]} />
                }
              >
                <Route path="admin-info" element={<AdminInfoPage />} />
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
