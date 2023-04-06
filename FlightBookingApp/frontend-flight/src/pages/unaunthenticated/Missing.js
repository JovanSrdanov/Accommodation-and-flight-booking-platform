import { Link } from "react-router-dom";
import "./Unauthorized.css"

const Missing = () => {
  return (
    <article style={{ padding: "100px" }}>
      <h1>Oops!</h1>
      <p>Page Not Found</p>
      <div className="flexGrow">
        <Link to="/">Go Back</Link> {/*TODO Stavi da vraca na profil korisnika*/}
      </div>
    </article>
  );
};

export default Missing;