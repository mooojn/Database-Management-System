import { useState, useEffect } from "react";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import './styles/ButtonNav.css'

function ButtonNav() {
    const buttons = [
        'Open Shell',
        'Create Database',
        'Start Transaction',
    ]
    return (
        <>
            <div className="buttonNav">
                {buttons.map((val, idx) => (
                    <button key={idx}>{val}</button>
                ))}
            </div>
        </>
    );
}

export default ButtonNav;
