import React, { useState } from 'react';
import useWebSocket from '/Users/darinautalieva/Desktop/GOProject/forum-frontend/src/hooks/useWebSocket.js';
import '/Users/darinautalieva/Desktop/GOProject/forum-frontend/src/components/MainLayout.css';

const Chat = () => {
    const [message, setMessage] = useState('');
    const userId = localStorage.getItem('userId');
    const { messages, sendMessage } = useWebSocket('ws://localhost:8081/websocket');

    if (!userId) {
        return <div>Please login to use chat</div>;
    }

    const handleSubmit = (e) => {
        e.preventDefault();
        if (message.trim()) {
            sendMessage({
                content: message,
                userId: userId,
                timestamp: new Date().toISOString()
            });
            setMessage('');
        }
    };

    return (
        <div className="chat-container">
            <div className="messages">
                {messages.map((msg, index) => (
                    <div key={index} className="message">
                        <span className="user">{msg.userId}</span>
                        <span className="content">{msg.content}</span>
                    </div>
                ))}
            </div>
            <form onSubmit={handleSubmit}>
                <input
                    type="text"
                    value={message}
                    onChange={(e) => setMessage(e.target.value)}
                    placeholder="Type your message..."
                />
                <button type="submit">Send</button>
            </form>
        </div>
    );
};
export default Chat;
