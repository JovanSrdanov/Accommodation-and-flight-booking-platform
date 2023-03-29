import Login from '../../components/login/login';
import FlightSearch from "../../components/flight-search/flight-search";

function HomePage() {
    function clickHandler() {

    }

    return (
        <div>
            <div className="title">Login</div>
            <Login/>
            <FlightSearch LoggedIn={false}></FlightSearch>
        </div>
    );
}

export default HomePage;