import useAxiosPrivate from './useAxiosPrivate'
export default function useAccountApi(){
    const axios = useAxiosPrivate();
    const endpointPath = "/api/account/";

    //Moze i preko async/await, to je cak modernije vljd, ali svejedno je moze i preko promisa
    const GetAccountInfo =  () => {
        return new Promise((resolve, reject) =>{
            axios.get(endpointPath + "logged/info")
                .then(response => {
                    resolve(response.data)
                })
                .catch(error => {
                    reject(error)
                })
        })
    }

    //Vracas svaku funkciju koju si definisao
    return {
        GetAccountInfo
    }
}