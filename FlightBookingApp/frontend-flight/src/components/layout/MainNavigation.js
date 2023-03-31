import {Link} from "react-router-dom"

import classes from './MainNavigation.module.css'

function MainNavigation() {
    return (
        <header className={classes.header}>
            <div className={classes.logo}>TEMP NAVBAR</div>
            <nav>
                <ul>
                    <li>
                        <Link to="/flight-search">Flight Search</Link>
                    </li>
                    <li>
                        <Link to="/admin-info">Admin Info</Link>
                    </li>
                    <li>
                        <Link to="/">Home</Link>
                    </li>
                    <li>
                        <Link to='/all-flights'>All flights</Link>
                    </li>
                    <li>
                        <Link to='/create-flight'>Create Flight</Link>
                    </li>
                    <li>
                        <Link to='/bought-tickets'>Bought tickets</Link>
                    </li>
                </ul>
            </nav>
        </header>
    );
}

export default MainNavigation