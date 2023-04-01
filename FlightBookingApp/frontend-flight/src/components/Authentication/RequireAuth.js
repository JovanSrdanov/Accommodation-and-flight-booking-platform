import { useLocation, Navigate, Outlet } from "react-router-dom";
import useAuth from "../../hooks/useAuth";

const RequireAuth = ({ allowedRoles }) => {
  const { auth } = useAuth();
  const location = useLocation();

  let foundRole = auth?.roles?.find((role) => allowedRoles?.includes(role));
  if (foundRole !== undefined) {
    foundRole = true
  }
  else {
    foundRole = false
  }

  return foundRole ? <Outlet /> : <Navigate to="/unauthorized" state={{ from: location }} replace />;
}

export default RequireAuth;