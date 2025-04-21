import React, { useState } from 'react';
// import 'forum-frontend/src/css/chat.css';
const Chat = () => {
    const [message, setMessage] = useState('');
    const [messages, setMessages] = useState([]);

    const handleSend = () => {
        if (message) {
            setMessages([...messages, message]);
            setMessage('');
        }
    };

    return (
        <div className="chat">
            <h3>Чат</h3>
            <div className="chat-messages">
                {messages.map((msg, index) => (
                    <div key={index}>{msg}</div>
                ))}
            </div>
            <input
                type="text"
                placeholder="Введите сообщение..."
                value={message}
                onChange={(e) => setMessage(e.target.value)}
            />
            <button onClick={handleSend}>Отправить</button>
        </div>
    );
};

export default Chat;
