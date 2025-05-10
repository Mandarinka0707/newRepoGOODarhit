import React, { useState, useEffect, useRef } from 'react';
import useWebSocket from 'react-use-websocket';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import '/Users/darinautalieva/Desktop/GOProject/forum-frontend/src/components/MainLayout.css';

const Chat = () => {
    const [message, setMessage] = useState('');
    const [messages, setMessages] = useState([]);
    const [connectionStatus, setConnectionStatus] = useState('connecting');
    const [error, setError] = useState(null);
    const [isLoading, setIsLoading] = useState(true);
    const navigate = useNavigate();
    const messagesEndRef = useRef(null);

    const token = localStorage.getItem('token');
    const username = localStorage.getItem('username');
    const isAuthenticated = !!token;

    const { sendMessage, lastMessage } = useWebSocket(
        'ws://localhost:8082/ws',
        {
            onOpen: () => {
                console.log('WebSocket connection established');
                setConnectionStatus('connected');
                setError(null);
            },
            onClose: () => {
                console.log('WebSocket connection closed');
                setConnectionStatus('disconnected');
                setError('Connection lost. Reconnecting...');
            },
            onError: (event) => {
                console.error('WebSocket error:', event);
                setConnectionStatus('error');
                setError('Failed to connect to chat server');
            },
            shouldReconnect: () => true,
            reconnectAttempts: 10,
            reconnectInterval: 3000,
        }
    );

    useEffect(() => {
        const fetchMessages = async () => {
            try {
                const response = await axios.get('http://localhost:8082/messages');
                // Обрабатываем timestamp при загрузке сообщений
                const processedMessages = (response.data || []).map(msg => ({
                    ...msg,
                    timestamp: msg.timestamp || new Date().toISOString()
                }));
                setMessages(processedMessages);
                setIsLoading(false);
            } catch (error) {
                console.error('Error fetching messages:', error);
                setError('Failed to load chat history');
                setIsLoading(false);
            }
        };

        fetchMessages();
    }, []);

    useEffect(() => {
        if (lastMessage !== null) {
            try {
                const newMessage = JSON.parse(lastMessage.data);
                // Добавляем timestamp, если его нет
                if (!newMessage.timestamp) {
                    newMessage.timestamp = new Date().toISOString();
                }
                setMessages(prev => [...prev, newMessage]);
            } catch (err) {
                console.error('Error parsing message:', err);
            }
        }
    }, [lastMessage]);

    useEffect(() => {
        messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
    }, [messages]);

    const handleSendMessage = () => {
        if (!isAuthenticated) {
            navigate('/login');
            return;
        }

        if (!message.trim()) return;

        const msg = {
            username: username || 'Anonymous',
            message: message.trim(),
            timestamp: new Date().toISOString() // Всегда добавляем timestamp при отправке
        };

        sendMessage(JSON.stringify(msg));
        setMessage('');
    };

    // Функция для форматирования времени сообщения
    const formatMessageTime = (timestamp) => {
        try {
            const date = new Date(timestamp);
            if (isNaN(date.getTime())) {
                return new Date().toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'});
            }
            return date.toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'});
        } catch (e) {
            console.error('Error formatting date:', e);
            return new Date().toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'});
        }
    };

    if (isLoading) {
        return (
            <div className="chat-container">
                <div className="loading-message">Loading chat...</div>
            </div>
        );
    }

    return (
        <div className="chat-wrapper">
            <div className="chat-container">
                <div className="chat-header">
                    <h2>Community Chat</h2>
                    <div className={`connection-status ${connectionStatus}`}>
                        {connectionStatus.toUpperCase()}
                    </div>
                </div>

                {error && <div className="error-message">{error}</div>}

                <div className="messages-window">
                    {messages.length > 0 ? (
                        messages.map((msg, index) => (
                            <div key={index} className={`message ${msg.username === username ? 'own-message' : ''}`}>
                                <div className="message-header">
                                    <span className="message-username">{msg.username}</span>
                                    <span className="message-time">
                                        {formatMessageTime(msg.timestamp)}
                                    </span>
                                </div>
                                <div className="message-content">{msg.message}</div>
                            </div>
                        ))
                    ) : (
                        <div className="no-messages">No messages yet. Be the first to say hello!</div>
                    )}
                    <div ref={messagesEndRef} />
                </div>

                <div className="message-input-area">
                    <input
                        type="text"
                        value={message}
                        onChange={(e) => setMessage(e.target.value)}
                        onKeyPress={(e) => e.key === 'Enter' && handleSendMessage()}
                        placeholder={isAuthenticated ? "Type your message..." : "Please login to chat"}
                        disabled={!isAuthenticated || connectionStatus !== 'connected'}
                    />
                    <button
                        onClick={handleSendMessage}
                        disabled={!message.trim() || !isAuthenticated || connectionStatus !== 'connected'}
                    >
                        Send
                    </button>
                </div>

                {!isAuthenticated && (
                    <div className="login-prompt">
                        <p>You need to <a href="/login">login</a> to send messages</p>
                    </div>
                )}
            </div>
        </div>
    );
};

export default Chat;
// import React, { useState, useEffect, useRef } from 'react';
// import useWebSocket from 'react-use-websocket';
// import axios from 'axios';
// import { useNavigate } from 'react-router-dom';
// import '/Users/darinautalieva/Desktop/GOProject/forum-frontend/src/components/MainLayout.css';

// const Chat = () => {
//     const [message, setMessage] = useState('');
//     const [messages, setMessages] = useState([]);
//     const [connectionStatus, setConnectionStatus] = useState('connecting');
//     const [error, setError] = useState(null);
//     const [isLoading, setIsLoading] = useState(true);
//     const navigate = useNavigate();
//     const messagesEndRef = useRef(null);

//     const token = localStorage.getItem('token');
//     const username = localStorage.getItem('username');
//     const isAuthenticated = !!token;

//     const { sendMessage, lastMessage } = useWebSocket(
//         'ws://localhost:8082/ws',
//         {
//             onOpen: () => {
//                 console.log('WebSocket connection established');
//                 setConnectionStatus('connected');
//                 setError(null);
//             },
//             onClose: () => {
//                 console.log('WebSocket connection closed');
//                 setConnectionStatus('disconnected');
//                 setError('Connection lost. Reconnecting...');
//             },
//             onError: (event) => {
//                 console.error('WebSocket error:', event);
//                 setConnectionStatus('error');
//                 setError('Failed to connect to chat server');
//             },
//             shouldReconnect: () => true,
//             reconnectAttempts: 10,
//             reconnectInterval: 3000,
//         }
//     );

//     useEffect(() => {
//         const fetchMessages = async () => {
//             try {
//                 const response = await axios.get('http://localhost:8082/messages');
//                 setMessages(response.data || []);
//                 setIsLoading(false);
//             } catch (error) {
//                 console.error('Error fetching messages:', error);
//                 setError('Failed to load chat history');
//                 setIsLoading(false);
//             }
//         };

//         fetchMessages();
//     }, []);

//     useEffect(() => {
//         if (lastMessage !== null) {
//             try {
//                 const newMessage = JSON.parse(lastMessage.data);
//                 setMessages(prev => [...prev, newMessage]);
//             } catch (err) {
//                 console.error('Error parsing message:', err);
//             }
//         }
//     }, [lastMessage]);

//     useEffect(() => {
//         messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
//     }, [messages]);

//     const handleSendMessage = () => {
//         if (!isAuthenticated) {
//             navigate('/login');
//             return;
//         }

//         if (!message.trim()) return;

//         const msg = {
//             username: username || 'Anonymous',
//             message: message.trim(),
//             timestamp: new Date().toISOString()
//         };

//         sendMessage(JSON.stringify(msg));
//         setMessage('');
//     };

//     if (isLoading) {
//         return (
//             <div className="chat-container">
//                 <div className="loading-message">Loading chat...</div>
//             </div>
//         );
//     }

//     return (
//         <div className="chat-wrapper">
//             <div className="chat-container">
//                 <div className="chat-header">
//                     <h2>Community Chat</h2>
//                     <div className={`connection-status ${connectionStatus}`}>
//                         {connectionStatus.toUpperCase()}
//                     </div>
//                 </div>

//                 {error && <div className="error-message">{error}</div>}

//                 <div className="messages-window">
//                     {messages.length > 0 ? (
//                         messages.map((msg, index) => (
//                             <div key={index} className={`message ${msg.username === username ? 'own-message' : ''}`}>
//                                 <div className="message-header">
//                                     <span className="message-username">{msg.username}</span>
//                                     <span className="message-time">
//                                         {new Date(msg.timestamp).toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'})}
//                                     </span>
//                                 </div>
//                                 <div className="message-content">{msg.message}</div>
//                             </div>
//                         ))
//                     ) : (
//                         <div className="no-messages">No messages yet. Be the first to say hello!</div>
//                     )}
//                     <div ref={messagesEndRef} />
//                 </div>

//                 <div className="message-input-area">
//                     <input
//                         type="text"
//                         value={message}
//                         onChange={(e) => setMessage(e.target.value)}
//                         onKeyPress={(e) => e.key === 'Enter' && handleSendMessage()}
//                         placeholder={isAuthenticated ? "Type your message..." : "Please login to chat"}
//                         disabled={!isAuthenticated || connectionStatus !== 'connected'}
//                     />
//                     <button
//                         onClick={handleSendMessage}
//                         disabled={!message.trim() || !isAuthenticated || connectionStatus !== 'connected'}
//                     >
//                         Send
//                     </button>
//                 </div>

//                 {!isAuthenticated && (
//                     <div className="login-prompt">
//                         <p>You need to <a href="/login">login</a> to send messages</p>
//                     </div>
//                 )}
//             </div>
//         </div>
//     );
// };

// export default Chat;
 // src/Chat.js
// import React, { useState, useEffect, useRef } from 'react';
// import useWebSocket from 'react-use-websocket';
// import axios from 'axios';
// import '/Users/darinautalieva/Desktop/GOProject/forum-frontend/src/components/MainLayout.css';

// const Chat = () => {
//   const [message, setMessage] = useState('');
//   const [messages, setMessages] = useState([]);
//   const { sendMessage, lastMessage } = useWebSocket('ws://localhost:8082/ws');
//   const usernames = localStorage.getItem("username");
//   const messagesEndRef = useRef(null);

//   useEffect(() => {
//     // Загрузка сообщений из базы данных при инициализации
//     const fetchMessages = async () => {
//       try {
//         const response = await axios.get('http://localhost:8082/messages');
//         setMessages(response.data);
//       } catch (error) {
//         console.error('Error fetching messages:', error);
//       }
//     };

//     fetchMessages();
//   }, []);

//   useEffect(() => {
//     if (lastMessage !== null) {
//       setMessages((prevMessages) => {
//         if (Array.isArray(prevMessages)) {
//           return [...prevMessages, JSON.parse(lastMessage.data)];
//         } else {
//           return [JSON.parse(lastMessage.data)];
//         }
//       });
//     }
//   }, [lastMessage]);

//   useEffect(() => {
//     // Прокрутка вниз при добавлении нового сообщения
//     messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
//   }, [messages]);

//   const handleSendMessage = () => {
//     const msg = { username: usernames, message };
//     sendMessage(JSON.stringify(msg));
//     setMessage('');
//   };

//   return (
//     <div className="chat-container">
//       <div className="connection-status connected">
//         Чатик
//       </div>
//       <div className="messages">
//         {messages && messages.length > 0 ? (
//           messages.map((msg, index) => (
//             <div key={index} className="message">
//               <span className="user">{msg.username}:</span>
//               <span>{msg.message}</span>
//             </div>
//           ))
//         ) : (
//           <div className="message system">No messages yet.</div>
//         )}
//         <div ref={messagesEndRef} />
//       </div>
//       <div className="message-form">
//         <input
//           type="text"
//           value={message}
//           onChange={(e) => setMessage(e.target.value)}
//         />
//         <button onClick={handleSendMessage} disabled={!message.trim()}>Send</button>
//       </div>
//     </div>
//   );
// };

// export default Chat;