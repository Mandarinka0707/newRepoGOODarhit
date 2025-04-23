// useWebSocket.js
import { useState, useEffect } from 'react';
const useWebSocket = (url) => {
    const [socket, setSocket] = useState(null);
    const [messages, setMessages] = useState([]);
    const [connectionStatus, setConnectionStatus] = useState('disconnected');

    useEffect(() => {
        const ws = new WebSocket(url);
        
        ws.onopen = () => {
            console.log('WebSocket connected');
            setSocket(ws);
            setConnectionStatus('connected');
        };

        ws.onmessage = (event) => {
            try {
                const newMessage = JSON.parse(event.data);
                setMessages(prev => [...prev, newMessage]);
            } catch (err) {
                console.error('Error parsing WebSocket message:', err);
            }
        };

        ws.onclose = () => {
            console.log('WebSocket disconnected');
            setSocket(null);
            setConnectionStatus('disconnected');
            // Attempt to reconnect after 5 seconds
            setTimeout(() => {
                setConnectionStatus('reconnecting');
            }, 5000);
        };

        ws.onerror = (error) => {
            console.error('WebSocket error:', error);
            setConnectionStatus('error');
        };

        return () => {
            if (ws && ws.readyState === WebSocket.OPEN) {
                ws.close();
            }
        };
    }, [url]);

    const sendMessage = (message) => {
        if (socket && socket.readyState === WebSocket.OPEN) {
            socket.send(JSON.stringify(message));
        } else {
            console.error('WebSocket is not connected');
            // Optionally queue messages when disconnected
        }
    };

    return { 
        socket, 
        messages, 
        sendMessage, 
        connectionStatus 
    };
};

export default useWebSocket; // Добавьте эту строку