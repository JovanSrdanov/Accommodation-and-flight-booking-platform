import RequireAuth from "./components/Authentication/RequireAuth";
import PersistLogin from "./components/Authentication/PersistLogin"
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
import Planes from "./components/Planes/Planes";
import {useEffect} from "react";
import AllFlightsPage from "./pages/admin/all-flights-page";

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
    useEffect(() => {
        function handleHover(event) {
            const uppercaseLetters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ';
            const lowercaseLetters = 'abcdefghijklmnopqrstuvwxyz';

            const {target} = event;
            if (target.tagName === 'H1') {
                let interval = null;
                const originalText = target.dataset.originalText || target.innerText.trim();

                const text = originalText;
                let iteration = 0;
                let animationComplete = false;
                let animationStarted = false;
                let animationReversed = false;

                const originalTextWithoutSpaces = originalText.replace(/\s/g, '');

                clearInterval(interval);

                interval = setInterval(() => {
                    let newText = '';
                    for (let i = 0; i < text.length; i++) {
                        const letter = text[i];
                        if (letter === ' ') {
                            newText += letter;
                            continue;
                        }

                        if (i < iteration) {
                            newText += text[i];
                        } else {
                            if (!animationStarted) {
                                animationStarted = true;
                                target.dataset.originalText = originalText;
                            }

                            let newLetter;
                            if (uppercaseLetters.includes(letter)) {
                                newLetter = uppercaseLetters[Math.floor(Math.random() * 26)];
                            } else if (lowercaseLetters.includes(letter)) {
                                newLetter = lowercaseLetters[Math.floor(Math.random() * 26)];
                            } else {
                                newLetter = letter;
                            }

                            if (animationComplete) {
                                if (newLetter !== originalText[i]) {
                                    animationReversed = true;
                                }
                                newLetter = originalText[i];
                            }
                            newText += newLetter;
                        }
                    }

                    target.innerText = newText;

                    if (iteration >= originalTextWithoutSpaces.length) {
                        clearInterval(interval);
                        animationComplete = true;
                        target.innerText = originalText;
                    }

                    if (animationComplete && animationStarted && !animationReversed) {
                        iteration -= 0.5;
                    } else {
                        iteration += 0.5;
                    }
                }, 50);
            }
        }

        document.addEventListener('mouseover', handleHover);

        return () => {
            document.removeEventListener('mouseover', handleHover);
        };
    }, []);
    return (
      <main>
        <ThemeProvider theme={darkTheme}>
          <MainNavigation />
          <Planes />
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
              <Route element={<PersistLogin />}>
                <Route element={<RequireAuth allowedRoles={[ROLES.REGULAR]} />}>
                  <Route path="flight-search" element={<FlightSearchPage />} />
                </Route>

                <Route element={<RequireAuth allowedRoles={[ROLES.ADMIN]} />}>
                  <Route path="/all-flights" element={<AllFlightsPage />} />
                  <Route path="admin-info" element={<AdminInfoPage />} />
                </Route>
              </Route>
              {/* catch all */}
              <Route path="*" element={<Missing />} />
            </Route>
          </Routes>{" "}
          origin/develop
        </ThemeProvider>
      </main>
    );
}

export default App;


