import { useState, useEffect } from "react";
import './styles/SideNav.css'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

function SideNav() {
    const dbs = [
        'StdManagement',
        'admin',
        'config',
        'local'
    ]
    return (
        <>
            <div className="sideNav">
                <div className="sideNav-head">
                    <h2>Compass</h2>
                    <h4>{'{ }'} My Queries</h4>
                </div>
                <hr />
                <div className="sideNav-body">

                    <h4>Connections</h4>
                    <h5>Moon</h5>

                    <div className="sideNav-databases">
                        <ul>
                            {dbs.map((val, idx) => (
                                <li id="sideNav-item" key={idx}>
                                    <FontAwesomeIcon icon={["fas", "database"]} style={{ marginRight: "8px" }} /> {val}
                                </li>

                            ))}
                        </ul>
                    </div>
                </div>
            </div>
        </>
    );
}

export default SideNav;
