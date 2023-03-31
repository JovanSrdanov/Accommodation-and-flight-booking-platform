import React from 'react';
import plane from "../../assets/flying-airplane.gif"
import "../../index.css"

const Planes = () => (
    <div className="plane-container">
        <div className="planeToRight">
            <img src={plane} alt=""/>
        </div>
        <div className="planeToLeft">
            <img src={plane} alt=""/>
        </div>
        <div className="planeToRight">
            <img src={plane} alt=""/>
        </div>
        <div className="planeToLeft">
            <img src={plane} alt=""/>
        </div>
        <div className="planeToRight">
            <img src={plane} alt=""/>
        </div>
        <div className="planeToLeft">
            <img src={plane} alt=""/>
        </div>
        <div className="planeToRight">
            <img src={plane} alt=""/>
        </div>
    </div>

);

export default Planes;