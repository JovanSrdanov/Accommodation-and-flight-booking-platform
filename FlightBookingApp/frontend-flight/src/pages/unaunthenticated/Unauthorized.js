import { useNavigate } from "react-router-dom";
import "./Unauthorized.css";
import "../page.css"
import Button from '@material-ui/core/Button';

const Unauthorized = () => {
  const navigate = useNavigate();

  const goBack = () => navigate(-1);

  return (
    <section className="page">
      <h1 style={{fontSize: "xxx-large", marginTop: '10%'}}>Unauthorized</h1>
      <br />
      <p style={{fontSize: "xx-large", marginTop: '2%'}}>You do not have access to the requested page.</p>
      <div>
        <Button style={{fontSize: "larger", marginTop: '4%'}} variant="contained" color="secondary" onClick={goBack}>Go Back</Button>
      </div>
    </section>
  );
};

export default Unauthorized;
