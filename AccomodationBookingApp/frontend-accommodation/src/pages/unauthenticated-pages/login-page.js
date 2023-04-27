import React from 'react';

import "../page.css"
import Login from "../../components/login/login";

function LoginPage() {
    return (
        <div className="page">
            <h1>Login</h1>
            <Login></Login>
        </div>
    );
}

export default LoginPage;