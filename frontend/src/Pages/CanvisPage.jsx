// src/components/CanvisPage.jsx
import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import Sidebar from '../Componants/Sidebar/Sidebar.jsx';
import DrawService from '../services/DrawService.js';

function CanvisPage() {
    const { randomId } = useParams();
    const [currentTool, setCurrentTool] = useState('pen');
    const [receivedMessage, setReceivedMessage] = useState(''); // To store received messages

    useEffect(() => {
        // Establish the WebSocket connection when the component mounts
        DrawService.connect();

        // Handle incoming WebSocket messages
        DrawService.socket.onmessage = (event) => {
            console.log('Received message from WebSocket:', event.data);
            setReceivedMessage(event.data);  // Store the received message in state
        };

        // Cleanup WebSocket connection on component unmount
        return () => {
            DrawService.close();
        };
    }, []);

    // Send a test message to the WebSocket server
    const handleSendMessage = () => {
        const message = 'Hello from CanvisPage!';
        DrawService.sendMessage(message);
    };

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
                {/* Canvas component will go here, but we're not modifying it yet */}
            </div>

            <div>
                <h2>Received WebSocket Message:</h2>
                <p>{receivedMessage}</p>  {/* Display received message */}
            </div>

            <div>
                <button onClick={handleSendMessage}>Send Test Message</button>  {/* Send message to WebSocket */}
            </div>
        </div>
    );
}

export default CanvisPage;
