//User info is made of account and user info (see  backend model)
import cssClasses from './user-info.module.css'
import { useEffect, useState } from 'react'
import useAccountApi from '../../hooks/useAccountApi'
import useUserApi from '../../hooks/useUserApi'





function UserInfo(props){

    const InfoRow = (props) =>{
        return (
        <div className={cssClasses.infoRow}>
            <b><p>{props.info}:</p></b>
            <p>{props.value}</p>
        </div>
        )
    }

    const tempUser ={
        fullname : "Loading...",
        address : "Loading..."
    }

    const tempAcc ={
        username : "Loading...",
        email : "Loading..."
    }



    const [accountInfo, setAccountInfo] = useState(tempAcc)
    const [userInfo, setUserInfo] = useState(tempUser)

    const { GetAccountInfo } = useAccountApi();
    const { GetLoggedUserInfo } = useUserApi();

    useEffect( () => {
        console.log("USO")
        GetAccountInfo()
        .then(data => {
            setAccountInfo(data)
        })
        .catch(error =>{
            alert(error)
        })


        GetLoggedUserInfo()
            .then(data => {
                setUserInfo(data)
            })
            .catch(error =>{
                alert(error)
            })

    }, [])



    return (
        <div className={cssClasses.infoWrapper}>
            <h1>User info</h1>
            <InfoRow info="Username" value={accountInfo.username}/>
            <InfoRow info="Email" value={accountInfo.email}/>
            <InfoRow info="Fullname" value={userInfo.fullname}/>
            <InfoRow info="Address" value={userInfo.address}/>
        </div>
    )
}


export default UserInfo;