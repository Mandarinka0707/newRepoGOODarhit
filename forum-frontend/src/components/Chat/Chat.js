// src/Chat.js
import React, { useState, useEffect, useRef } from 'react';
import useWebSocket from 'react-use-websocket';
import axios from 'axios';
import '/Users/darinautalieva/Desktop/GOProject/forum-frontend/src/components/MainLayout.css';

const Chat = () => {
  const [message, setMessage] = useState('');
  const [messages, setMessages] = useState([]);
  const { sendMessage, lastMessage } = useWebSocket('ws://localhost:8082/ws');
  const usernames = localStorage.getItem("username");
  const messagesEndRef = useRef(null);

  useEffect(() => {
    // Загрузка сообщений из базы данных при инициализации
    const fetchMessages = async () => {
      try {
        const response = await axios.get('http://localhost:8082/messages');
        setMessages(response.data);
      } catch (error) {
        console.error('Error fetching messages:', error);
      }
    };

    fetchMessages();
  }, []);

  useEffect(() => {
    if (lastMessage !== null) {
      setMessages((prevMessages) => {
        if (Array.isArray(prevMessages)) {
          return [...prevMessages, JSON.parse(lastMessage.data)];
        } else {
          return [JSON.parse(lastMessage.data)];
        }
      });
    }
  }, [lastMessage]);

  useEffect(() => {
    // Прокрутка вниз при добавлении нового сообщения
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  const handleSendMessage = () => {
    const msg = { username: usernames, message };
    sendMessage(JSON.stringify(msg));
    setMessage('');
  };

  return (
    <div className="chat-container">
      <div className="connection-status connected">
        Чатик
      </div>
      <div className="messages">
        {messages && messages.length > 0 ? (
          messages.map((msg, index) => (
            <div key={index} className="message">
              <span className="user">{msg.username}:</span>
              <span>{msg.message}</span>
            </div>
          ))
        ) : (
          <div className="message system">No messages yet.</div>
        )}
        <div ref={messagesEndRef} />
      </div>
      <div className="message-form">
        <input
          type="text"
          value={message}
          onChange={(e) => setMessage(e.target.value)}
        />
        <button onClick={handleSendMessage} disabled={!message.trim()}>Send</button>
      </div>
    </div>
  );
};

export default Chat;