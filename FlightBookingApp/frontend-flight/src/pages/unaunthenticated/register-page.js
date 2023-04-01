import React from 'react'
import Register from '../../components/register/register'
import "../page.css"
import "./register-page.css"

const RegisterPage = () => {
  return (
    <div>
      <h1 style={{marginTop: '65%', marginBottom: '-45%'}}>Register</h1>
      <div className="App">
          <Register/>
      </div>
    </div>
  )
}

export default RegisterPage