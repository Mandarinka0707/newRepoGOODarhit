import React, { useState, useEffect } from 'react';
import axios from 'axios';
import useWebSocket from '/Users/darinautalieva/Desktop/GOProject/forum-frontend/src/hooks/useWebSocket.js'; // путь поменяй если другой
import '/Users/darinautalieva/Desktop/GOProject/forum-frontend/src/components/MainLayout.css'; // путь к стилям

const Chat = () => {
  const [message, setMessage] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [errorMessage, setErrorMessage] = useState('');
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [userId, setUserId] = useState(null);
  const [username, setUsername] = useState(null);

  const {
    messages,
    sendMessage,
    connectionStatus,
    connect,
    disconnect,
  } = useWebSocket('ws://localhost:8082/ws', { manual: true });

  // Load user data from localStorage on initial mount
  useEffect(() => {
    const storedUserId = localStorage.getItem('userId');
    const storedUsername = localStorage.getItem('username');

    if (storedUserId && storedUsername) {
      setUserId(storedUserId);
      setUsername(storedUsername);
      setIsLoggedIn(true);
    }
  }, []);

  // Connect when isLoggedIn changes to true
  useEffect(() => {
    if (isLoggedIn) {
      connect(); // Connect WebSocket when isLoggedIn is true
    } else {
      disconnect(); // Disconnect if isLoggedIn becomes false (e.g., logout)
    }

    return () => {
        disconnect(); // Disconnect when component unmounts
    }
  }, [isLoggedIn, connect, disconnect]);

  const handleLogin = async (e) => {
    e.preventDefault();

    try {
      const response = await axios.post('http://localhost:8080/login', {
        email,
        password,
      });

      if (response.status === 200) {
        const { userId, username, token } = response.data;

        if (!userId || !username || username === 'undefined') {
          console.warn('Invalid user session, cannot connect');
          setErrorMessage('Login failed: server did not provide user info.');
          return;
        }

        localStorage.setItem('userId', userId);
        localStorage.setItem('username', username);
        localStorage.setItem('token', token);

        setUserId(userId);
        setUsername(username);
        setIsLoggedIn(true);  // Set isLoggedIn AFTER setting userId and username

      }
    } catch (error) {
      console.error("Login error:", error.response || error.message || error); // Log detailed error
      setErrorMessage('Login failed. Please check your credentials.');
    }
  };

  const handleSubmit = (e) => {
    e.preventDefault();

    if (!message.trim() || !username) {
      console.warn('Cannot send empty message or missing username');
      return;
    }

    sendMessage({
      user_id: userId,
      username: username,
      content: message.trim(),
    });

    setMessage('');
  };

  if (!isLoggedIn) {
    return (
      <div>
        <h2>Login to use Chat</h2>
        <form onSubmit={handleLogin}>
          <div>
            <label>Email</label>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
          </div>
          <div>
            <label>Password</label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </div>
          <button type="submit">Login</button>
        </form>
        {errorMessage && <p style={{ color: 'red' }}>{errorMessage}</p>}
      </div>
    );
  }

  return (
    <div className="chat-container">
      <div className="connection-status">
        Status: {connectionStatus}
      </div>
      <div className="messages">
        {messages.map((msg, index) => (
          <div key={index} className="message">
            <span className="user">{msg.username}:</span>
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

// import React, { useState, useEffect } from 'react';
// import axios from 'axios';
// import useWebSocket from '/Users/darinautalieva/Desktop/GOProject/forum-frontend/src/hooks/useWebSocket.js'; // путь поменяй если другой
// import '/Users/darinautalieva/Desktop/GOProject/forum-frontend/src/components/MainLayout.css'; // путь к стилям

// const Chat = () => {
//   const [message, setMessage] = useState('');
//   const [email, setEmail] = useState('');
//   const [password, setPassword] = useState('');
//   const [errorMessage, setErrorMessage] = useState('');
//   const [isLoggedIn, setIsLoggedIn] = useState(false);

//   const [userId, setUserId] = useState(null);
//   const [username, setUsername] = useState(null);

//   const {
//     messages,
//     sendMessage,
//     connectionStatus,
//     connect, // добавляем ручное подключение
//     disconnect
//   } = useWebSocket('ws://localhost:8082/ws', { manual: true }); // manual mode!

//   useEffect(() => {
//     const storedUserId = localStorage.getItem('userId');
//     const storedUsername = localStorage.getItem('username');
//     if (storedUserId && storedUsername) {
//       setUserId(storedUserId);
//       setUsername(storedUsername);
//       setIsLoggedIn(true);
//       connect(); // подключаем WebSocket только если есть данные
//     }
//   }, [connect]);

//   const handleLogin = async (e) => {
//     e.preventDefault();

//     try {
//       const response = await axios.post('http://localhost:8080/login', {
//         email,
//         password,
//       });

//       if (response.status === 200) {
//         const { userId, username, token } = response.data;

//         if (!userId || !username || username === 'undefined') {
//           console.warn('Invalid user session, cannot connect');
//           setErrorMessage('Login failed: server did not provide user info.');
//           return;
//         }

//         localStorage.setItem('userId', userId);
//         localStorage.setItem('username', username);
//         localStorage.setItem('token', token);

//         setUserId(userId);
//         setUsername(username);
//         setIsLoggedIn(true);

//         connect(); // подключаем WebSocket после логина

//       }
//     } catch (error) {
//       console.error(error);
//       setErrorMessage('Login failed. Please check your credentials.');
//     }
//   };

//   const handleSubmit = (e) => {
//     e.preventDefault();

//     if (!message.trim() || !username) {
//       console.warn('Cannot send empty message or missing username');
//       return;
//     }

//     sendMessage({
//       user_id:  userId,
//       username: username,
//       content: message.trim(),
//     });

//     setMessage('');
//   };

//   if (!isLoggedIn) {
//     return (
//       <div>
//         <h2>Login to use Chat</h2>
//         <form onSubmit={handleLogin}>
//           <div>
//             <label>Email</label>
//             <input
//               type="email"
//               value={email}
//               onChange={(e) => setEmail(e.target.value)}
//               required
//             />
//           </div>
//           <div>
//             <label>Password</label>
//             <input
//               type="password"
//               value={password}
//               onChange={(e) => setPassword(e.target.value)}
//               required
//             />
//           </div>
//           <button type="submit">Login</button>
//         </form>

//         {errorMessage && <p style={{ color: 'red' }}>{errorMessage}</p>}
//       </div>
//     );
//   }

//   return (
//     <div className="chat-container">
//       <div className="connection-status">
//         Status: {connectionStatus}
//       </div>

//       <div className="messages">
//         {messages.map((msg, index) => (
//           <div key={index} className="message">
//             <span className="user">{msg.username}:</span>  
//             <span className="content">{msg.content}</span>
//           </div>
//         ))}
//       </div>

//       <form onSubmit={handleSubmit} className="chat-form">
//         <input
//           type="text"
//           value={message}
//           onChange={(e) => setMessage(e.target.value)}
//           placeholder="Type your message..."
//           disabled={connectionStatus !== 'connected'}
//           className="chat-input"
//         />
//         <button
//           type="submit"
//           disabled={connectionStatus !== 'connected' || !message.trim()}
//           className="chat-send-button"
//         >
//           Send
//         </button>
//       </form>
//     </div>
//   );
// };

// export default Chat;
