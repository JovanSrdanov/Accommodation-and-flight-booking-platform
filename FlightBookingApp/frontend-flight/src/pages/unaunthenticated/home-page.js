import Login from '../../components/login/login';
import FlightSearch from "../../components/flight-search/flight-search";

import "../page.css"

function HomePage() {

    return (
        <div className='page'>
            <Login/>
            <FlightSearch LoggedIn={false}></FlightSearch>
        </div>
    );
}

export default HomePage;