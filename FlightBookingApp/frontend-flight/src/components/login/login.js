import React, {useEffect, useRef, useState} from "react";
import useAuth from "../../hooks/useAuth";
import {Link, useLocation, useNavigate} from "react-router-dom";
import jwt_decode from "jwt-decode";
import useToggle from "../../hooks/useToggle";
import useInput from "../../hooks/useInput";
import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";
import "./login.css"

import axios from "../../api/axios";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faCheck, faTimes} from "@fortawesome/free-solid-svg-icons";

const LOGIN_URL = "/api/account/login";

const Login = () => {
    const navigate = useNavigate();
    const location = useLocation();
    const from = location.state?.from?.pathname || "/"

    const { setAuth } = useAuth();
    const userRef = useRef();
    const errRef = useRef();

    const [user, resetUser, userAttributes] = useInput('user', '')
    const [pwd, setPwd] = useState("");
    const [errMsg, setErrMsg] = useState("");

    const [check, toggleCheck] = useToggle('persist', false);

    useEffect(() => {
        userRef.current?.focus();
    }, []);

    useEffect(() => {
        setErrMsg("");
    }, [user, pwd]);

    const handleSubmit = async (e) => {
        e.preventDefault();

        try {
            const response = await axios.post(
                LOGIN_URL,
                JSON.stringify({username: user, password: pwd}),
                {
                    headers: {"Content-Type": "application/json"},
                    withCredentials: true,
                }
            );
            console.log(JSON.stringify(response?.data));
            //console.log(JSON.stringify(response));

            const accessToken = response?.data?.accessToken;
            //const refreshToken = response?.data?.refreshToken;

            const decodedToken = jwt_decode(accessToken);
            console.log('decoded token: ', decodedToken)

            const roles = decodedToken.roles
            console.log('roles: ', roles)

            setAuth({user, pwd, roles, accessToken});
            //setUser("");
            resetUser();
            setPwd("");
            roles.every(element => element === 1) ? navigate("/customer-info", {replace: true}) : navigate("/admin-info", {replace: true})
        } catch (err) {
            if (!err?.response) {
                setErrMsg("No Server Response");
            } else if (err.response?.status === 400) {
                setErrMsg("Missing Username or Password");
            } else if (err.response?.status === 401) {
                setErrMsg("Invalid username/password, or account not activated");
            } else {
                setErrMsg("Login Failed");
            }
            errRef.current.focus();
        }
    };

    return (
      <section style={{
          width: '100%',
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

        <form
            style={
                {
                    display: 'flex',
                    flexDirection: 'column',
                    justifyContent: 'space-evenly',
                    flexGrow: '1',
                    paddingBottom: '1rem'
                }}>
            <label style={{
                marginTop: '1rem',
                marginBottom: '1rem',
                fontSize: 'large'
            }}
                   htmlFor="username">
                Username:
            </label>
          <TextField
            style={{marginRight: '3%', fontSize: 'x-large', fontFamily: 'Nunito, sans-serif',
                padding: '0.25rem',
                borderRadius: '0.5rem',}}
            type="text"
            id="username"
            variant="standard"
            inputRef={userRef}
            autoComplete="off"
            {...userAttributes}
            required
          />
            <label style={{
                marginTop: '1rem',
                marginBottom: '1rem',
                fontSize: 'large'
            }}
                   htmlFor="password">
                Password:
            </label>
          <TextField
              style={{marginRight: '3%', fontSize: 'x-large', fontFamily: 'Nunito, sans-serif',
                  padding: '0.25rem',
                  borderRadius: '0.5rem',}}
            type="password"
            variant="standard"
            onChange={(e) => setPwd(e.target.value)}
            value={pwd}
            required
          />
          <div className="persistCheck">
              <input
                  style={{marginTop:'2%'}}
                  className="persistCheckbox"
                  type="checkbox"
                  id="persist"
                  onChange={toggleCheck}
                  checked={check}
              />
              <label style={{fontSize: 'large'}} htmlFor="persist">Trust This Device</label>
          </div>
          <Button
              style={{fontSize: 'x-large', marginLeft: '2%', fontFamily: 'Nunito, sans-serif',
                  borderRadius: '0.5rem',
                  marginTop: '1rem',
                  padding: '0.5rem'}}
              variant={user === "" || pwd === "" ? "outlined" : "contained" }
              disabled={user === "" || pwd === ""}
              onClick={handleSubmit}
          >
              Sign In
          </Button>
        </form>
        <p style={{fontSize: 'large'}}>
          Need an Account?
          <br />
          <span style={{fontSize: 'large'}} className="line">
            {/*put router link here*/}
            <Link to="/register" style={{color: 'aquamarine'}}>Sign up</Link>
          </span>
        </p>
      </section>
    );
}

export default Login;
