import { Outlet } from "react-router-dom";
import { useState, useEffect } from "react";
import useRefreshToken from "../../hooks/useRefreshToken";
import useAuth from "../../hooks/useAuth";

const PersistLogin = () => {
  const [isLoading, setIsloading] = useState(true);
  const refresh = useRefreshToken();
  const { auth } = useAuth();

  useEffect(() => {
    const verifyRegreshToken = async () => {
      try {
        await refresh();
      }
      catch (err) {
        console.error(err);
      }
      finally {
        setIsloading(false);
      }
    }

    !auth?.accessToken ? verifyRegreshToken() : setIsloading(false);
  }, [])

  useEffect(() => {
    console.log(`isLoading: ${isLoading}`)
    console.log(`aT: ${JSON.stringify(auth?.accessToken)}`)
  }, [isLoading])

  return (
    <>
      {isLoading
        ? <p>Loading...</p>
        : <Outlet/>
      }
    </>
  )
}

export default PersistLogin