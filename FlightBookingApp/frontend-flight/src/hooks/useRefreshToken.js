import axios from "../api/axios"
import useAuth from "./useAuth"
import jwt_decode from "jwt-decode";

const useRefreshToken = () => {
  const { setAuth } = useAuth();

  const refresh = async () => {
    const response = await axios.get('api/account/refresh-token', {
      withCredentials: true
    });
    setAuth(prev => {
      console.log("Token refreshed")

      const accessToken = response?.data?.accessToken;

      const decodedToken = jwt_decode(accessToken);
      console.log("decoded token: ", decodedToken);

      const roles = decodedToken.roles;
      console.log("roles: ", roles);

      return { ...prev, roles: roles, accessToken: accessToken };
    });
    return response.data.accessToken;
  }

  return refresh
}

export default useRefreshToken