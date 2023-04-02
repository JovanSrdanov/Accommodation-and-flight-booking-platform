import React, {useEffect, useRef, useState} from "react";
import useAuth from "../../hooks/useAuth";
import {Link, useLocation, useNavigate} from "react-router-dom";
import jwt_decode from "jwt-decode";
import useToggle from "../../hooks/useToggle";
import useInput from "../../hooks/useInput";
import "./login.css"

import axios from "../../api/axios";

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
                setErrMsg("Unauthorized");
            } else {
                setErrMsg("Login Failed");
            }
            errRef.current.focus();
        }
    };

    return (
      <section>
        <p
          ref={errRef}
          className={errMsg ? "errmsg" : "offscreen"}
          aria-live="assertive"
        >
          {errMsg}
        </p>

        <form onSubmit={handleSubmit}>
          <label style={{fontSize: 'x-large'}} htmlFor="username">Username:</label>
          <input
            style={{marginRight: '3%', fontSize: 'x-large'}}
            type="text"
            id="username"
            ref={userRef}
            autoComplete="off"
            {...userAttributes}
            required
          />

          <label style={{fontSize: 'x-large'}} htmlFor="password">Password:</label>
          <input
            style={{fontSize: 'x-large'}}
            type="password"
            id="password"
            onChange={(e) => setPwd(e.target.value)}
            value={pwd}
            required
          />
          <button style={{fontSize: 'x-large'}}>Sign In</button>
          <div className="persistCheck">
            <input
              style={{marginTop:'2%'}}
              className="persistCheckbox"
              type="checkbox"
              id="persist"
              onChange={toggleCheck}
              checked={check}
            />
            <label style={{fontSize: 'x-large'}} htmlFor="persist">Trust This Device</label>
          </div>
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
