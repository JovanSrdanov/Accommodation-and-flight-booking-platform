
import useAxiosPrivate from './useAxiosPrivate'
export default function useUserApi(){
    const axios = useAxiosPrivate();
    const endpointPath = "/api/user/";

    const GetLoggedUserInfo = () => {
        return new Promise((resolve, reject) => {
            axios.get( endpointPath + "logged-in")
            .then(response => {
                const fullname = response.data.name + " " + response.data.surname;
                const address = response.data.address;
                const addressInline = address.street + " " + address.streetNumber + ", " + address.city + ", " + address.country;
                resolve({
                    fullname : fullname,
                    address : addressInline
                })
            })
            .catch(error => {
                reject(error)
            })
        })
    }


    return {
        GetLoggedUserInfo
    }
}