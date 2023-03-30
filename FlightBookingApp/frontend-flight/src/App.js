import {Route, Routes} from "react-router-dom";

import HomePage from "./pages/unaunthenticated/home-page";
import FlightSearchPage from "./pages/customer/flight-search-page";
import MainNavigation from "./components/layout/MainNavigation";
import {createTheme, ThemeProvider} from '@mui/material/styles';
import Planes from "./components/Planes/Planes";
import {useEffect} from "react";




const darkTheme = createTheme({
    palette: {
        mode: 'dark',
    },
});




function App() {
    useEffect(() => {
        const uppercaseLetters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ';
        const lowercaseLetters = 'abcdefghijklmnopqrstuvwxyz';

        function handleHover(event) {
            const { target } = event;
            if (target.tagName === "H1") {
                let interval = null;
                const originalText = target.dataset.originalText || target.innerText.trim();

                const text = originalText;
                let iteration = 0;

                clearInterval(interval);

                interval = setInterval(() => {
                    target.innerText = text
                        .split('')
                        .map((letter, index) => {
                            if (letter === ' ') {
                                return letter;
                            }

                            if (index < iteration) {
                                return text[index];
                            }

                            if (uppercaseLetters.includes(letter)) {
                                return uppercaseLetters[Math.floor(Math.random() * 26)];
                            }

                            if (lowercaseLetters.includes(letter)) {
                                return lowercaseLetters[Math.floor(Math.random() * 26)];
                            }

                            return letter;
                        })
                        .join('');

                    if (iteration >= text.replace(/\s/g, '').length) {
                        clearInterval(interval);
                    }

                    iteration += 1 / 3;
                }, 30);

                // Store the original text in the dataset
                target.dataset.originalText = originalText;
            }
        }

        document.addEventListener("mouseover", handleHover);

        return () => {
            document.removeEventListener("mouseover", handleHover);
        };
    }, []);

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


