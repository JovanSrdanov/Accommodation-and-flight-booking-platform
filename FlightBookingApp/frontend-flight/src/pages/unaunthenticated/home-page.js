import Login from '../../components/login/login';
import FlightSearch from "../../components/flight-search/flight-search";

import "./home-page.css"

function HomePage() {
    function clickHandler() {

    }

    return (
        <div className='App'>
            <Login/>
            <FlightSearch LoggedIn={false}></FlightSearch>
        </div>
    );
}

export default HomePage;