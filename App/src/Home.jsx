import { useState, useEffect } from "react";
import Transactions from "./components/Transactions";
import SideNav from "./components/SideNav";
import ButtonNav from "./components/ButtonNav";

import './styles/Home.css'

function Home() {

    return (
        <>
            <div className="main-container"> 
                <SideNav />
                <ButtonNav />
            </div>
        </>
    );
}

export default Home;
