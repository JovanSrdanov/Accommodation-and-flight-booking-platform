import React from 'react'
import Register from '../../components/register/register'
import "../page.css"
import "./register-page.css"

const RegisterPage = () => {
  return (
    <div className='App'>
      <h1 style={{marginTop: '4%', marginBottom: '20%'}}>Register</h1>
      <Register/>
    </div>
  )
}

export default RegisterPage