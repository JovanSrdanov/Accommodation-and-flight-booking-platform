import React from 'react';
import AllFlights from "../../components/all-flights/all-flights";
import "../page.css"

function AllFlightsPage() {
    return (
        <div className="page">
            <h1>All flights</h1>
            <AllFlights></AllFlights>
        </div>
    );
}

export default AllFlightsPage;