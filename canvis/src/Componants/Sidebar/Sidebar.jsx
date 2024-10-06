// src/components/Sidebar.jsx
import React from 'react';

function Sidebar({ randomId, currentTool, setCurrentTool }) {
  return (
    <div className="sidebar" style={{
      width: '250px',
      height: '100%',
      backgroundColor: '#f0f0f0',
      padding: '20px',
      boxShadow: '2px 0 5px rgba(0,0,0,0.1)',
      color: 'black'  // Set the default text color to black
    }}>
      <h1 style={{ marginBottom: '20px', color: 'black' }}>Canvas Tools</h1>
      <p style={{ marginBottom: '20px', color: 'black' }}>Page ID: <strong>{randomId}</strong></p>
      
      <div className="tool">
        <button 
          style={{
            width: '100%',
            padding: '10px',
            marginBottom: '10px',
            backgroundColor: currentTool === 'pen' ? '#007bff' : '#6c757d',
            color: 'white',  // Keep button text white for contrast
            border: 'none',
            borderRadius: '5px',
            cursor: 'pointer'
          }}
          onClick={() => setCurrentTool('pen')}
        >
          Pen Tool
        </button>
      </div>
      <div className="tool">
        <button 
          style={{
            width: '100%',
            padding: '10px',
            marginBottom: '10px',
            backgroundColor: currentTool === 'eraser' ? '#28a745' : '#6c757d',
            color: 'white',  // Keep button text white for contrast
            border: 'none',
            borderRadius: '5px',
            cursor: 'pointer'
          }}
          onClick={() => setCurrentTool('eraser')}
        >
          Eraser
        </button>
      </div>
    </div>
  );
}

export default Sidebar;