import { Link } from "react-router-dom"

import  classes from './MainNavigation.module.css'

function MainNavigation() {
  return (
    <header className={classes.header}>
      <div className={classes.logo}>TEMP NAVBAR</div>
      <nav>
        <ul>
          <li>
            <Link to='/flight-search'>Flight Search</Link>

          </li>
            <li>

                <Link to='/all-flights'>All flights</Link>
            </li>
        </ul>
      </nav>
    </header>
  )
} 

export default MainNavigation