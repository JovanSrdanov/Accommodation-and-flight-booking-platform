import { Link, useNavigate } from "react-router-dom";
import useLogout from "../../hooks/useLogout";
import {Link} from "react-router-dom"

import classes from './MainNavigation.module.css'

//TODO Stefan: napraviti posebne navbar-ove za neautentifikovane, regularne korisnike i admine

const MainNavigation = () => {
  const navigate = useNavigate();
  const logout = useLogout();

  const signOut = async () => {
    await logout();
    navigate('/');
  }

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
            <Link to="/all-flights">All flights</Link>
          </li>
          <li>
                        <Link to='/create-flight'>Create Flight</Link>
                    </li>
                    <li>
                        <Link to='/bought-tickets'>Bought tickets</Link>
                    </li>
          <li>
             <button onClick={signOut}>Logout</button>
          </li>
        </ul>
      </nav>
    </header>
  );
}

export default MainNavigation