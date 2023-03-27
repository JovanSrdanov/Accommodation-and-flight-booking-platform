function Login() {
  return (
    <div className="card">
      <div>
        <input type="text" className="text-field" />
        <input type="text" className="text-field" />
      </div>
      <div className="actions">
        <button className="btn">Login</button>
        <button className="btn">Register</button>
      </div>
    </div>
  );
}

export default Login;