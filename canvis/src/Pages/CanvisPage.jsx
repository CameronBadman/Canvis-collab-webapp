// src/components/CanvisPage.jsx
import React, { useState } from 'react';
import { useParams } from 'react-router-dom';
import CustomSVGDrawing from '../Componants/Canvas/Canvas';
import Sidebar from '../Componants/Sidebar/Sidebar';

function CanvisPage() {
  const { randomId } = useParams();
  const [currentTool, setCurrentTool] = useState('pen');

  return (
    <div className="canvas-page" style={{ 
      height: '100vh', 
      width: '100vw', 
      display: 'flex', 
      flexDirection: 'row'
    }}>
      <Sidebar 
        randomId={randomId} 
        currentTool={currentTool} 
        setCurrentTool={setCurrentTool} 
      />
      <div style={{ flex: 1, position: 'relative' }}>
        <CustomSVGDrawing strokeWidth={2} strokeColor="black" currentTool={currentTool} />
      </div>
    </div>
  );
}

export default CanvisPage;