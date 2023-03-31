import Login from '../../components/login/login';
import FlightSearch from "../../components/flight-search/flight-search";

import "../page.css"
import React from "react";

function HomePage() {
    return (
        <div className='page'>
            <h1>Home Page</h1>
            <Login/>
            <FlightSearch LoggedIn={false}></FlightSearch>
        </div>
    );
}

export default HomePage;