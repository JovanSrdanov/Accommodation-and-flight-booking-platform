import React, {useEffect, useRef, useState} from "react";
import useAuth from "../../hooks/useAuth";
import {Link, useLocation, useNavigate} from "react-router-dom";
import jwt_decode from "jwt-decode";


import axios from "../../api/axios";

const LOGIN_URL = "/api/account/login";

const Login = () => {
    const navigate = useNavigate();
    const location = useLocation();
    const from = location.state?.from?.pathname || "/"

    const {setAuth} = useAuth();
    const userRef = useRef();
    const errRef = useRef();

    const [user, setUser] = useState("");
    const [pwd, setPwd] = useState("");
    const [errMsg, setErrMsg] = useState("");

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
            setUser("");
            setPwd("");
            navigate(from, {replace: true});
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
                <label htmlFor="username">Username:</label>
                <input
                    type="text"
                    id="username"
                    ref={userRef}
                    autoComplete="off"
                    onChange={(e) => setUser(e.target.value)}
                    value={user}
                    required
                />

                <label htmlFor="password">Password:</label>
                <input
                    type="password"
                    id="password"
                    onChange={(e) => setPwd(e.target.value)}
                    value={pwd}
                    required
                />
                <button>Sign In</button>
            </form>
            <p>
                Need an Account?
                <br/>
                <span className="line">
              {/*put router link here*/}
                    <Link to="/register">Sign up</Link>
            </span>
            </p>
        </section>
    )
}

export default Login;
