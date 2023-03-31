import axios from "../api/axios"
import useAuth from "./useAuth"

const useRefreshToken = () => {
  const { setAuth } = useAuth();

  const refresh = async () => {
    const response = await axios.get('api/account/refresh-token', {
      withCredentials: true
    });
    setAuth(prev => {
      console.log("Token refreshed")
      console.log(JSON.stringify(prev));
      console.log(response.data.accessToken);
      return { ...prev, accessToken: response.data.accessToken }
    });
    return response.data.accessToken;
  }

  return refresh
}

export default useRefreshToken