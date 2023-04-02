import React from 'react'
import Register from '../../components/register/register'
import "../page.css"
import "./register-page.css"

const RegisterPage = () => {
  return (
      <div className="page">
          <h1 style={{marginTop: '3%', marginBottom: '4%'}}>Register</h1>
          <Register/>
      </div>
  )
}

export default RegisterPage