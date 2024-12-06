class DrawService {
    constructor() {
        this.socket = null;
        this.url = 'ws://localhost:8080/draw/'; // WebSocket endpoint
        this.isManuallyClosed = false;
        this.eventListeners = {}; // Custom event listeners
    }

    // Connect to the WebSocket server
    connect() {
        if (this.socket && (this.socket.readyState === WebSocket.OPEN || this.socket.readyState === WebSocket.CONNECTING)) {
            console.log('WebSocket is already connected or connecting');
            return;
        }

        console.log('Attempting to connect to WebSocket:', this.url);
        this.socket = new WebSocket(this.url);

        this.socket.onopen = () => {
            console.log('WebSocket connection established');
            this.emitEvent('open');
        };

        this.socket.onmessage = (event) => {
            console.log('Received message:', event.data);
            this.emitEvent('message', event.data);
        };

        this.socket.onerror = (error) => {
            console.error('WebSocket error:', error);
            this.emitEvent('error', error);
        };

        this.socket.onclose = (event) => {
            console.log('WebSocket connection closed:', event.reason);
            this.emitEvent('close', event);

            // If the connection is lost, log and do not try reconnecting.
            if (!this.isManuallyClosed) {
                console.log('WebSocket connection lost, no automatic reconnection.');
            }
        };
    }

    // Send a message to the WebSocket server
    sendMessage(message) {
        if (this.socket && this.socket.readyState === WebSocket.OPEN) {
            console.log('Sending message:', message);
            this.socket.send(message);
        } else {
            console.error('WebSocket is not open. Cannot send message.');
        }
    }

    // Close the WebSocket connection
    close() {
        if (this.socket) {
            this.isManuallyClosed = true;
            console.log('Manually closing WebSocket connection');
            this.socket.close();
        }
    }

    // Add an event listener for custom events
    on(eventType, callback) {
        if (!this.eventListeners[eventType]) {
            this.eventListeners[eventType] = [];
        }
        this.eventListeners[eventType].push(callback);
    }

    // Remove an event listener
    off(eventType, callback) {
        if (this.eventListeners[eventType]) {
            this.eventListeners[eventType] = this.eventListeners[eventType].filter(cb => cb !== callback);
        }
    }

    // Emit a custom event
    emitEvent(eventType, data) {
        if (this.eventListeners[eventType]) {
            this.eventListeners[eventType].forEach(callback => callback(data));
        }
    }
}

// Instantiate the draw service
const drawService = new DrawService();
export default drawService;
