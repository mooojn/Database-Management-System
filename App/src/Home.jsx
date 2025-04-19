import { useState, useEffect } from "react";
import Transactions from "./components/Transactions";
import SideNav from "./components/SideNav";
import Testing from "./components/Testing";
import ButtonNav from "./components/ButtonNav";

import './styles/Home.css'

function Home() {

    return (
        <>
            <div className="main-container"> 
                {/* <SideNav />
                <ButtonNav /> */}
                <Testing/>
            </div>
        </>
    );
}

export default Home;
