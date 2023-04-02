import { useRef, useState, useEffect } from "react";
import {Link, useNavigate} from "react-router-dom";
import {
  faCheck,
  faTimes,
  faInfoCircle,
} from "@fortawesome/free-solid-svg-icons";
import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import axios from "../../api/axios"
import "./register.css";

const USER_REGEX = /^[A-z][A-z0-9]{3,23}$/;
const PWD_REGEX = /^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[!@#$%]).{8,24}$/;
const EMAIL_REGEX = /^[^\s@]+@[^\s@]+\.[^\s@]+$/i
const REGISTER_URL = "/api/account/register";

const Register = () => {
  const userRef = useRef();
  const errRef = useRef();
  const navigate = useNavigate();

  // User
  const [user, setUser] = useState("");
  const [validName, setValidName] = useState(false);
  const [userFocus, setUserFocus] = useState(false);

  // Password
  const [pwd, setPwd] = useState("");
  const [validPwd, setValidPwd] = useState(false);
  const [pwdFocus, setPwdFocus] = useState(false);

  // Confirm password
  const [matchPwd, setMatchPwd] = useState("");
  const [validMatch, setValidMatch] = useState(false);
  const [matchFocus, setMatchFocus] = useState(false);

  // Email
  const [email, setEmail] = useState("")
  const [validEmail, setValidEmail] = useState(false)
  const [emailFocus, setEmailFocus] = useState(false)

  // Error message
  const [errMsg, setErrMsg] = useState("");
  const [success, setSuccess] = useState(false);

  const [name, setName] = useState("");
  const [surname, setSurname] = useState("");
  const [country, setCountry] = useState("");
  const [city, setCity] = useState("");
  const [street, setStreet] = useState("");
  const [streetNumber, setStreetNumber] = useState("");

  const [userInfoDialogVisible, setUserInfoDialogVisible] = useState(false)

  const [isSecondFormPartDisabled, setIsSecondFormPartDisabled] = useState(true)

  const validSecondPart = () => {
    return name === "" || surname === "" || country === "" || city === "" || street === "" || streetNumber === ""
  }

  useEffect(() => {
    setValidName(USER_REGEX.test(user));
  }, [user]);

  useEffect(() => {
    setValidPwd(PWD_REGEX.test(pwd));
    setValidMatch(pwd === matchPwd);
  }, [pwd, matchPwd]);

  useEffect(() => {
    setValidEmail(EMAIL_REGEX.test(email))
  }, [email])

  useEffect(() => {
    setErrMsg("");
  }, [user, pwd, matchPwd]);

  useEffect(() => {
    setIsSecondFormPartDisabled(validSecondPart())
  }, [name, surname, country, city, street, streetNumber])

  const handleSubmit = async (e) => {
    e.preventDefault();
    // if button enabled with JS hack
    const v1 = USER_REGEX.test(user);
    const v2 = PWD_REGEX.test(pwd);
    if (!v1 || !v2) {
      setErrMsg("Invalid Entry");
      return;
    }

    await axios.post(REGISTER_URL, JSON.stringify({ username:user, password:pwd, email:email, name:name, surname:surname,
      address: {country:country, city:city, street:street, streetNumber:streetNumber}}))
        .then(res => {
          console.log(res.data)
          setSuccess(true);
          //clear state and controlled inputs
          //need value attrib on inputs for this
          setUser("");
          setPwd("");
          setMatchPwd("");

          //navigate("/")
        })
        .catch(err => {
          if (!err?.response) {
            setErrMsg("No Server Response");
          } else if (err.response?.status === 400) {    // TODO Stefan: promeni
            setErrMsg("Username or Email Taken");
          } else {
            setErrMsg("Registration Failed");
          }
          errRef.current?.focus();
        })
  };

  return (
    <>
      {success ? (
        <section>
          <h1>An activation link has been sent to {email}</h1>
          <p>
            <Link to="/" style={{color: 'aquamarine'}}>Sign in</Link>
          </p>
        </section>
      ) : (
        <section style={{width: '100%',
                         maxWidth: '420px',
                         minHeight: '400px',
                         display: 'flex',
                         flexDirection: 'column',
                         justifyContent: 'flex-start',
                         padding: '1rem',
                         margin: 'auto',
                         backgroundColor: '#282e3b'}}>
          <p
            ref={errRef}
            className={errMsg ? "errmsg" : "offscreen"}
            aria-live="assertive"
          >
            {errMsg}
          </p>
          <div style={{color: 'red'}}>All fields are required!</div>
          {!userInfoDialogVisible ?
              <form style={
                {
                  display: 'flex',
                  flexDirection: 'column',
                  justifyContent: 'space-evenly',
                  flexGrow: '1',
                  paddingBottom: '1rem'
                }
              }>
                <label style={{
                  marginTop: '1rem',
                  marginBottom: '1rem'
                }}
                       htmlFor="username">
                  Username:
                  <FontAwesomeIcon
                      icon={faCheck}
                      className={validName ? "valid" : "hide"}
                  />
                  <FontAwesomeIcon
                      icon={faTimes}
                      className={validName || !user ? "hide" : "invalid"}
                  />
                </label>
                <TextField
                    style={
                      {
                        fontFamily: 'Nunito, sans-serif',
                        fontSize: '22px',
                        padding: '0.25rem',
                        borderRadius: '0.5rem',
                      }
                    }
                    type="text"
                    id="username"
                    variant="standard"
                    inputRef={userRef}
                    autoComplete="off"
                    onChange={(e) => setUser(e.target.value)}
                    value={user}
                    required
                    aria-invalid={validName ? "false" : "true"}
                    aria-describedby="uidnote"
                    onFocus={() => setUserFocus(true)}
                    onBlur={() => setUserFocus(false)}
                />
                <p
                    id="uidnote"
                    className={
                      userFocus && user && !validName ? "instructions" : "offscreen"
                    }
                >
                  <FontAwesomeIcon icon={faInfoCircle} />
                  4 to 24 characters.
                  <br />
                  Must begin with a letter.
                  <br />
                  Letters, numbers, underscores, hyphens allowed.
                </p>

                <label style={{
                  marginTop: '1rem',
                  marginBottom: '1rem'
                }}
                       htmlFor="password">
                  Password:
                  <FontAwesomeIcon
                      icon={faCheck}
                      className={validPwd ? "valid" : "hide"}
                  />
                  <FontAwesomeIcon
                      icon={faTimes}
                      className={validPwd || !pwd ? "hide" : "invalid"}
                  />
                </label>
                <TextField
                    style={
                      {
                        fontFamily: 'Nunito, sans-serif',
                        fontSize: '22px',
                        padding: '0.25rem',
                        borderRadius: '0.5rem',
                      }
                    }
                    type="password"
                    id="password"
                    variant="standard"
                    onChange={(e) => setPwd(e.target.value)}
                    value={pwd}
                    required
                    aria-invalid={validPwd ? "false" : "true"}
                    aria-describedby="pwdnote"
                    onFocus={() => setPwdFocus(true)}
                    onBlur={() => setPwdFocus(false)}
                />
                <p
                    id="pwdnote"
                    className={pwdFocus && !validPwd ? "instructions" : "offscreen"}
                >
                  <FontAwesomeIcon icon={faInfoCircle} />
                  8 to 24 characters.
                  <br />
                  Must include uppercase and lowercase letters, a number and a
                  special character.
                  <br />
                  Allowed special characters:{" "}
                  <span aria-label="exclamation mark">!</span>{" "}
                  <span aria-label="at symbol">@</span>{" "}
                  <span aria-label="hashtag">#</span>{" "}
                  <span aria-label="dollar sign">$</span>{" "}
                  <span aria-label="percent">%</span>
                </p>

                <label style={{
                  marginTop: '1rem',
                  marginBottom: '1rem'
                }}
                       htmlFor="confirm_pwd">
                  Confirm Password:
                  <FontAwesomeIcon
                      icon={faCheck}
                      className={validMatch && matchPwd ? "valid" : "hide"}
                  />
                  <FontAwesomeIcon
                      icon={faTimes}
                      className={validMatch || !matchPwd ? "hide" : "invalid"}
                  />
                </label>
                <TextField
                    style={
                      {
                        fontFamily: 'Nunito, sans-serif',
                        fontSize: '22px',
                        padding: '0.25rem',
                        borderRadius: '0.5rem',
                      }
                    }
                    type="password"
                    id="confirm_pwd"
                    variant="standard"
                    onChange={(e) => setMatchPwd(e.target.value)}
                    value={matchPwd}
                    required
                    aria-invalid={validMatch ? "false" : "true"}
                    aria-describedby="confirmnote"
                    onFocus={() => setMatchFocus(true)}
                    onBlur={() => setMatchFocus(false)}
                />
                <p
                    id="confirmnote"
                    className={
                      matchFocus && !validMatch ? "instructions" : "offscreen"
                    }
                >
                  <FontAwesomeIcon icon={faInfoCircle} />
                  Must match the first password input field.
                </p>
                <label style={{
                  marginTop: '1rem',
                  marginBottom: '1rem'
                }}
                       htmlFor="confirm_pwd">
                  Email:
                  <FontAwesomeIcon
                      icon={faCheck}
                      className={validEmail ? "valid" : "hide"}
                  />
                  <FontAwesomeIcon
                      icon={faTimes}
                      className={validEmail || !email ? "hide" : "invalid"}
                  />
                </label>
                <TextField
                    style={
                      {
                        fontFamily: 'Nunito, sans-serif',
                        fontSize: '22px',
                        padding: '0.25rem',
                        borderRadius: '0.5rem',
                      }
                    }
                    type="email"
                    id="email"
                    variant="standard"
                    onChange={(e) => setEmail(e.target.value)}
                    value={email}
                    required
                    aria-invalid={validEmail ? "false" : "true"}
                    aria-describedby="emailnote"
                    onFocus={() => setEmailFocus(true)}
                    onBlur={() => setEmailFocus(false)}
                />
                <p
                    id="emailnote"
                    className={
                      emailFocus && !validEmail ? "instructions" : "offscreen"
                    }
                >
                  <FontAwesomeIcon icon={faInfoCircle} />
                  No whitespaces allowed
                  <br/>
                  Only one '@' symbol allowed
                  <br/>
                  must have a '.' symbol after the '@' symbol, followed by a domain name (.com, .edu.rs, etc.)
                </p>
                <Button
                    style={
                      {
                        fontFamily: 'Nunito, sans-serif',
                        fontSize: '22px',
                        borderRadius: '0.5rem',
                        marginTop: '1rem',
                        padding: '0.5rem'
                      }
                    }
                    //className={!validName || !validPwd || !validMatch !validEmail ? "" : "validButton"}
                    variant={!validName || !validPwd || !validMatch || !validEmail ? "outlined" : "contained"}
                    disabled={!validName || !validPwd || !validMatch || !validEmail}
                    onClick={(e) => {
                      e.preventDefault();
                      setUserInfoDialogVisible(true)
                    }}
                >
                  Next
                </Button>
              </form> :
              <form style={
                {
                  display: 'flex',
                  flexDirection: 'column',
                  justifyContent: 'space-evenly',
                  flexGrow: '1',
                  paddingBottom: '1rem',
                }}>
                <label style={{
                  marginTop: '1rem',
                  marginBottom: '1rem'
                }}
                       htmlFor="name">
                  Name:
                </label>
                <TextField
                    style={
                      {
                        fontFamily: 'Nunito, sans-serif',
                        fontSize: '22px',
                        padding: '0.25rem',
                        borderRadius: '0.5rem',
                      }
                    }
                    type="text"
                    id="name"
                    variant="standard"
                    autoComplete="on"
                    onChange={(e) => setName(e.target.value)}
                    value={name}
                    required
                />
                <label style={{
                  marginTop: '1rem',
                  marginBottom: '1rem'
                }}
                       htmlFor="name">
                  Surname:
                </label>
                <TextField
                    style={
                      {
                        fontFamily: 'Nunito, sans-serif',
                        fontSize: '22px',
                        padding: '0.25rem',
                        borderRadius: '0.5rem',
                      }
                    }
                    type="text"
                    id="surname"
                    variant="standard"
                    autoComplete="on"
                    onChange={(e) => setSurname(e.target.value)}
                    value={surname}
                    required
                />
                <label style={{
                  marginTop: '1rem',
                  marginBottom: '1rem'
                }}
                       htmlFor="name">
                  Country:
                </label>
                <TextField
                    style={
                      {
                        fontFamily: 'Nunito, sans-serif',
                        fontSize: '22px',
                        padding: '0.25rem',
                        borderRadius: '0.5rem',
                      }
                    }
                    type="text"
                    id="country"
                    variant="standard"
                    autoComplete="on"
                    onChange={(e) => setCountry(e.target.value)}
                    value={country}
                    required
                />
                <label style={{
                  marginTop: '1rem',
                  marginBottom: '1rem'
                }}
                       htmlFor="name">
                  City:
                </label>
                <TextField
                    style={
                      {
                        fontFamily: 'Nunito, sans-serif',
                        fontSize: '22px',
                        padding: '0.25rem',
                        borderRadius: '0.5rem',
                      }
                    }
                    type="text"
                    id="city"
                    variant="standard"
                    autoComplete="on"
                    onChange={(e) => setCity(e.target.value)}
                    value={city}
                    required
                />
                <label style={{
                  marginTop: '1rem',
                  marginBottom: '1rem'
                }}
                       htmlFor="name">
                  Street:
                </label>
                <TextField
                    style={
                      {
                        fontFamily: 'Nunito, sans-serif',
                        fontSize: '22px',
                        padding: '0.25rem',
                        borderRadius: '0.5rem',
                      }
                    }
                    type="text"
                    id="street"
                    variant="standard"
                    autoComplete="on"
                    onChange={(e) => setStreet(e.target.value)}
                    value={street}
                    required
                />
                <label style={{
                  marginTop: '1rem',
                  marginBottom: '1rem'
                }}
                       htmlFor="name">
                  Street Number:
                </label>
                <TextField
                    style={
                      {
                        fontFamily: 'Nunito, sans-serif',
                        fontSize: '22px',
                        padding: '0.25rem',
                        borderRadius: '0.5rem',
                      }
                    }
                    type="text"
                    id="streetNumber"
                    variant="standard"
                    autoComplete="on"
                    onChange={(e) => setStreetNumber(e.target.value)}
                    value={streetNumber}
                    required
                />
                <Button
                  type="submit"
                  style={{
                  marginTop: '6%',
                  fontSize: 'x-large',
                  textAlign: 'center'
                }}
                className={!isSecondFormPartDisabled ? "validButton" : ""}
                variant={!isSecondFormPartDisabled ? "contained" : "outlined"}
                color={!isSecondFormPartDisabled ? "success" : "primary"}
                disabled={isSecondFormPartDisabled}
                onClick={handleSubmit}
                >
                  Submit
                </Button>
               <Button style={ {
                 marginTop: '6%',
                 fontSize: 'x-large',
                 textAlign: 'center'
               }}
                variant="contained"
                onClick={(e) => {
                  e.preventDefault();
                  setUserInfoDialogVisible(false)
                }}
               >
                 Back
               </Button>
              </form>
          }
          <p>
            Already registered?
            <br />
            <span className="line">
              {/*put router link here*/}
              <Link to="/" style={{color: 'aquamarine'}}>Sign in</Link>
            </span>
          </p>
        </section>
      )}
    </>
  );
};

export default Register;
