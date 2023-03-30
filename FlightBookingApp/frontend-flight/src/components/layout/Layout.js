import { Outlet } from "react-router-dom"
import "./Layout.css"

export const Layout = () => {
  return (
    <main className="App">
      <Outlet/>
    </main>
  )
}
