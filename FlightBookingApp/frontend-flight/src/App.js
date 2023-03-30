import {Route, Routes} from "react-router-dom";

import HomePage from "./pages/unaunthenticated/home-page";
import FlightSearchPage from "./pages/customer/flight-search-page";
import MainNavigation from "./components/layout/MainNavigation";
import {createTheme, ThemeProvider} from '@mui/material/styles';
import Planes from "./components/Planes/Planes";
import useHoverAnimation from "./HackerHeader";



const darkTheme = createTheme({
    palette: {
        mode: 'dark',
    },
});




function App() {
    useHoverAnimation()
    return (
        <ThemeProvider theme={darkTheme}>
            <MainNavigation/>
            <Planes/>
            <Routes>
                <Route path="/" element={<HomePage/>}/>
                <Route path="/flight-search" element={<FlightSearchPage/>}/>
            </Routes>
        </ThemeProvider>
    );
}

export default App;


