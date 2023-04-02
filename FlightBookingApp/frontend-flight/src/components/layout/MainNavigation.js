import {Link, useNavigate} from "react-router-dom";
import useLogout from "../../hooks/useLogout";
import useAuth from "../../hooks/useAuth";

import classes from './MainNavigation.module.css'

const MainNavigation = () => {
    const navigate = useNavigate();
    const logout = useLogout();
    const { auth } = useAuth();

    const canLoad = (role) => {
        // proverava da li je korisnik ulogovan, i da li je odgovarajuce role (ovako stoji posto trenutno svaki
        // korisnik ima jednu rolu, moze se prosiriti sa includes pa rola po potrebi
        return auth?.roles !== undefined && auth?.roles.every(element => element === role)
    }

    const signOut = async () => {
        await logout();
        navigate('/');
    }

    return (
        <header className={classes.header}>
            <div className={classes.logo}>FTN Airlines</div>
            <nav>
                { canLoad(0) &&
                    <ul>
                        <li>
                            <Link to="/all-flights">All flights</Link>
                        </li>
                        <li>
                            <Link to='/create-flight'>Create Flight</Link>
                        </li>
                        <li>
                            <Link to="/admin-info">Admin Info</Link>
                        </li>
                        <li>
                            <Link style={{ color: 'azure' }} to="/" onClick={signOut}>Logout</Link>
                        </li>
                    </ul>
                }
                { canLoad(1) &&
                    <ul>
                        <li>
                            <Link to="/flight-search">Flight Search</Link>
                        </li>
                        <li>
                            <Link to='/bought-tickets'>Bought tickets</Link>
                        </li>
                        <li>
                            <Link to="/customer-info">Customer Info</Link>
                        </li>
                        <li>
                            <Link style={{ color: 'azure' }} to="/" onClick={signOut}>Logout</Link>
                        </li>
                    </ul>
                }
            </nav>
        </header>
    );
}

export default MainNavigation