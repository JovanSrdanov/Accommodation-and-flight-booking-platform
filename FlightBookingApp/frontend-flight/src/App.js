import {Route, Routes} from "react-router-dom";

import HomePage from "./pages/unaunthenticated/home-page";
import FlightSearchPage from "./pages/customer/flight-search-page";
import MainNavigation from "./components/layout/MainNavigation";
import {createTheme, ThemeProvider} from '@mui/material/styles';

import Header from "./Header";

const darkTheme = createTheme({
    palette: {
        mode: 'dark',
    },
});


function App() {
    return (
        <ThemeProvider theme={darkTheme}>
            <Header />
            <MainNavigation/>
            <div className="plane">
                <img src="https://media.tenor.com/qsdblRVNZysAAAAC/flying-airplane.gif"/>
            </div>
            <div className="planeToLeft">
                <img src="https://media.tenor.com/qsdblRVNZysAAAAC/flying-airplane.gif"/>
            </div>
            <Routes>
                <Route path="/" element={<HomePage/>}/>
                <Route path="/flight-search" element={<FlightSearchPage/>}/>
            </Routes>
        </ThemeProvider>
    );
}

export default App;


