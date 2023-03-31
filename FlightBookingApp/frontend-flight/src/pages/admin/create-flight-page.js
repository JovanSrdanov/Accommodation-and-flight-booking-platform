import React from 'react';
import "../page.css"
import CreateFlight from "../../components/create-flight/create-flight";

function CreateFlightPage() {
    return (
        <div className="page">
            <h1>Create Flight</h1>
            <CreateFlight></CreateFlight>
        </div>
    );
}

export default CreateFlightPage;