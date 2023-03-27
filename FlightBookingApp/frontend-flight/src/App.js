import { Route, Routes } from "react-router-dom";

import HomePage from "./pages/unaunthenticated/home-page";
import FlightSearchPage from "./pages/customer/flight-search-page";
import MainNavigation from "./components/layout/MainNavigation";

function App() {
  return (
    <div>
      <MainNavigation/>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/flight-search" element={<FlightSearchPage />}/>
      </Routes>
    </div>
  );
}

export default App;
