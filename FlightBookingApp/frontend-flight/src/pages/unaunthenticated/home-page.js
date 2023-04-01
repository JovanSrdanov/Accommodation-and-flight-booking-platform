import Login from '../../components/login/login';
import FlightSearch from "../../components/flight-search/flight-search";

import "../page.css"
import React from "react";

function HomePage() {
    return (
        <div className='page'>
            <h1 style={{marginTop: '4%', marginBottom: '6%'}}>Welcome, please log in</h1>
            <Login/>
            <div style={{marginBottom: '4%'}}/>
            <FlightSearch LoggedIn={false}></FlightSearch>
        </div>
    );
}

export default HomePage;