import { Route, Routes } from "react-router-dom";

import HomePage from "./pages/unaunthenticated/home-page";
import FlightSearchPage from "./pages/customer/flight-search-page";
import MainNavigation from "./components/layout/MainNavigation";
import RegisterPage from "./pages/unaunthenticated/register-page";

function App() {
  return (
    <div>
      <MainNavigation/>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/register" element={<RegisterPage/>}/>
        <Route path="/flight-search" element={<FlightSearchPage />}/>
      </Routes>
    </div>
  );
}

export default App;
