import FlightSearch from "../../components/flight-search/flight-search";
import "../page.css"

function FlightSearchPage() {
    return (
        <div className="page">
            <h1>Flight Search</h1>
            <FlightSearch LoggedIn={true}></FlightSearch>
        </div>
    )
}

export default FlightSearchPage