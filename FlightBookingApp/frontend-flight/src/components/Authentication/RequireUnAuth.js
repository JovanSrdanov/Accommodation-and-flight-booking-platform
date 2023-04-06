import useAuth from "../../hooks/useAuth";
import {Navigate, Outlet, useLocation} from "react-router-dom";

const RequireUnAuth = () => {
    const { auth } = useAuth();
    const location = useLocation();

    let isAuthenticated = auth?.roles?.length > 0

    return (
        !isAuthenticated ?
            <Outlet/> :
            <Navigate to= {auth?.roles.includes(0) ? "/admin-info" : "/customer-info"} state={{ from: location }} replace />
    )
}

export default  RequireUnAuth