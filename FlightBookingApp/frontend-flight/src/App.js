import {Route, Routes} from "react-router-dom";

import HomePage from "./pages/unaunthenticated/home-page";
import FlightSearchPage from "./pages/customer/flight-search-page";
import MainNavigation from "./components/layout/MainNavigation";
import RegisterPage from "./pages/unaunthenticated/register-page"

import {createTheme, ThemeProvider} from '@mui/material/styles';

const darkTheme = createTheme({
    palette: {
        mode: 'dark',
    },
});


function App() {
    return (
      <main>
        <ThemeProvider theme={darkTheme}>
          <MainNavigation />
          <Routes>
            <Route path="/" element={<HomePage />} />
            <Route path="/register" element={<RegisterPage />} />
            <Route path="/flight-search" element={<FlightSearchPage />} />
          </Routes>
        </ThemeProvider>
      </main>
    );
}

export default App;
