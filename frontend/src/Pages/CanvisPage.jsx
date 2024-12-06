import React, { useState } from 'react';
import { useParams } from 'react-router-dom';
import CanvasComponent from "../Componants/Canvas/CanvasComponant.jsx";
import Sidebar from '../Componants/Sidebar/Sidebar.jsx';

function CanvasPage() {
    const { randomId } = useParams();


    return (
        <div className="flex h-screen w-screen overflow-hidden">


            {/* Canvas Area */}
            <div className="flex-1 relative">
                <CanvasComponent />
            </div>
        </div>
    );
}

export default CanvasPage;

