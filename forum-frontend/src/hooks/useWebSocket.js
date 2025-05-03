import { useState, useEffect, useCallback, useRef } from 'react';

const useWebSocket = (webSocketUrl, apiUrl, token) => {
  const [messages, setMessages] = useState([]);
  const [status, setStatus] = useState('disconnected');
  const [error, setError] = useState(null);
  const ws = useRef(null);
  const reconnectAttempts = useRef(0);

  const loadHistory = useCallback(async () => {
    try {
      const response = await fetch(apiUrl, {
        headers: { 'Authorization': `Bearer ${token}` }
      });
      
      if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
      
      const history = await response.json();
      setMessages(history);
    } catch (err) {
      setError(err.message);
      console.error('History load error:', err);
    }
  }, [apiUrl, token]);

  const connectWebSocket = useCallback(() => {
    if (!token || status === 'connected') return;

    const wsUrl = new URL(webSocketUrl);
    wsUrl.searchParams.set('token', token);

    ws.current = new WebSocket(wsUrl.href);
    setStatus('connecting');

    ws.current.onopen = () => {
      console.log('WebSocket connected');
      reconnectAttempts.current = 0;
      setStatus('connected');
      setError(null);
    };

    ws.current.onmessage = (event) => {
      try {
        const newMessage = JSON.parse(event.data);
        setMessages(prev => [...prev, newMessage]);
      } catch (err) {
        console.error('Message parse error:', err);
      }
    };

    ws.current.onclose = (event) => {
      console.log('WebSocket closed:', event.code, event.reason);
      setStatus('disconnected');
      
      if (event.code === 4001 || event.code === 4002) {
        setError('Authentication required');
        return;
      }

      // Exponential backoff reconnect
      const timeout = Math.min(1000 * 2 ** reconnectAttempts.current, 30000);
      setTimeout(() => {
        reconnectAttempts.current++;
        connectWebSocket();
      }, timeout);
    };

    ws.current.onerror = (err) => {
      console.error('WebSocket error:', err);
      setStatus('error');
      setError('Connection error');
    };
  }, [webSocketUrl, token, status]);

  useEffect(() => {
    if (!token) {
      setError('Authentication token is required');
      return;
    }

    loadHistory();
    connectWebSocket();

    return () => {
      if (ws.current?.readyState === WebSocket.OPEN) {
        ws.current.close();
      }
    };
  }, [token, loadHistory, connectWebSocket]);

  const sendMessage = useCallback((message) => {
    if (ws.current?.readyState === WebSocket.OPEN) {
      ws.current.send(JSON.stringify({
        ...message,
        timestamp: new Date().toISOString()
      }));
    }
  }, []);

  return {
    messages,
    sendMessage,
    status,
    error,
    retry: loadHistory
  };
};

export default useWebSocket;
// // useWebSocket.js
// import { useState, useEffect } from 'react';
// const useWebSocket = (url) => {
//     const [socket, setSocket] = useState(null);
//     const [messages, setMessages] = useState([]);
//     const [connectionStatus, setConnectionStatus] = useState('disconnected');

//     useEffect(() => {
//         const ws = new WebSocket(url);
        
//         ws.onopen = () => {
//             console.log('WebSocket connected');
//             setSocket(ws);
//             setConnectionStatus('connected');
//         };

//         ws.onmessage = (event) => {
//             try {
//                 const newMessage = JSON.parse(event.data);
//                 setMessages(prev => [...prev, newMessage]);
//             } catch (err) {
//                 console.error('Error parsing WebSocket message:', err);
//             }
//         };

//         ws.onclose = () => {
//             console.log('WebSocket disconnected');
//             setSocket(null);
//             setConnectionStatus('disconnected');
//             // Attempt to reconnect after 5 seconds
//             setTimeout(() => {
//                 setConnectionStatus('reconnecting');
//             }, 5000);
//         };

//         ws.onerror = (error) => {
//             console.error('WebSocket error:', error);
//             setConnectionStatus('error');
//         };

//         return () => {
//             if (ws && ws.readyState === WebSocket.OPEN) {
//                 ws.close();
//             }
//         };
//     }, [url]);

//     const sendMessage = (message) => {
//         if (socket && socket.readyState === WebSocket.OPEN) {
//             socket.send(JSON.stringify(message));
//         } else {
//             console.error('WebSocket is not connected');
//             // Optionally queue messages when disconnected
//         }
//     };

//     return { 
//         socket, 
//         messages, 
//         sendMessage, 
//         connectionStatus 
//     };
// };

// export default useWebSocket; // Добавьте эту строку

// import { useEffect, useRef, useState, useCallback } from 'react';

// export function useWebSocket(url, { manual = false } = {}) {
//   const socketRef = useRef(null);
//   const reconnectTimeoutRef = useRef(null);
//   const reconnectInterval = useRef(1000);
//   const [isConnected, setIsConnected] = useState(false);
//   const [connectionStatus, setConnectionStatus] = useState('disconnected');
//   const [messages, setMessages] = useState([]);
//   const [error, setError] = useState(null);
//   const messageQueue = useRef([]);

//   const getWebSocketUrl = useCallback(() => {
//     try {
//       const token = localStorage.getItem('token');
//       if (!token) {
//         setError('Authentication token not found');
//         return null;
//       }
      
//       const wsUrl = new URL(url);
//       wsUrl.searchParams.set('token', token);
//       return wsUrl.toString();
//     } catch (e) {
//       console.error('Invalid WebSocket URL:', e);
//       setError('Invalid WebSocket URL');
//       return null;
//     }
//   }, [url]);

//   const handleIncomingMessage = useCallback((data) => {
//     switch (data.type) {
//       case 'MESSAGE':
//         setMessages(prev => [...prev, data.data]);
//         break;
//       case 'HISTORY':
//         setMessages(data.data);
//         break;
//       case 'AUTH_ERROR':
//         setError(`Authentication error: ${data.message}`);
//         localStorage.removeItem('token');
//         window.location.reload();
//         break;
//       default:
//         console.warn('Unhandled message type:', data.type);
//     }
//   }, []);

//   const processMessageQueue = () => {
//     while (messageQueue.current.length > 0 && socketRef.current?.readyState === WebSocket.OPEN) {
//       const message = messageQueue.current.shift();
//       socketRef.current.send(JSON.stringify(message));
//     }
//   };

//   // 2. Обновленный useWebSocket.js (React)
// const connect = useCallback(() => {
//   const wsUrl = getWebSocketUrl();
//   if (!wsUrl) return;

//   if (socketRef.current) {
//     if (socketRef.current.readyState === WebSocket.OPEN) {
//       console.log('Already connected');
//       return;
//     }
//     socketRef.current.close();
//   }

//   console.log('Attempting WebSocket connection...');
//   setConnectionStatus('connecting');
  
//   const ws = new WebSocket(wsUrl);
//   socketRef.current = ws;

//   ws.onopen = () => {
//     console.log('WebSocket connected');
//     setIsConnected(true);
//     setConnectionStatus('connected');
//     setError(null);
//     reconnectInterval.current = 1000;
//     processMessageQueue();
    
//     // Запрос истории
//     const historyRequest = JSON.stringify({
//       type: 'GET_HISTORY',
//       timestamp: Date.now()
//     });
//     ws.send(historyRequest);
//   };

//   ws.onmessage = (event) => {
//     try {
//       const parsedData = JSON.parse(event.data);
//       handleIncomingMessage(parsedData);
//     } catch (e) {
//       console.error('Message parse error:', e);
//     }
//   };

//   ws.onerror = (error) => {
//     console.error('WebSocket error:', error);
//     setError('Connection error');
//     setConnectionStatus('error');
//   };

//   ws.onclose = (event) => {
//     console.log(`WebSocket closed: ${event.code}`, event.reason);
//     setIsConnected(false);
//     setConnectionStatus('disconnected');

//     if (event.code === 4002) { // Аутентификация
//       handleIncomingMessage({
//         type: 'AUTH_ERROR',
//         message: event.reason || 'Authentication failed'
//       });
//       return;
//     }

//     if (!event.wasClean && event.code !== 1000) {
//       const timeout = Math.min(reconnectInterval.current * 2, 30000);
//       reconnectInterval.current = timeout;
//       console.log(`Reconnecting in ${timeout}ms...`);
//       reconnectTimeoutRef.current = setTimeout(connect, timeout);
//     }
//   };
// }, [getWebSocketUrl, handleIncomingMessage]);

//   const disconnect = useCallback((permanent = false) => {
//     if (socketRef.current) {
//       if (permanent) {
//         socketRef.current.onclose = null;
//       }
//       socketRef.current.close(
//         permanent ? 1000 : 1001,
//         permanent ? 'Normal closure' : 'Reconnecting'
//       );
//     }
//     if (reconnectTimeoutRef.current) {
//       clearTimeout(reconnectTimeoutRef.current);
//     }
//   }, []);

//   const sendMessage = useCallback((message) => {
//     const userId = parseInt(localStorage.getItem('userId'), 10);
//     if (isNaN(userId)) {
//       setError('Invalid user ID');
//       return;
//     }
    
//     const fullMessage = {
//       ...message,
//       timestamp: Date.now(),
//       userId,
//       username: localStorage.getItem('username') || 'unknown'
//     };

//     if (socketRef.current?.readyState === WebSocket.OPEN) {
//       try {
//         socketRef.current.send(JSON.stringify(fullMessage));
//       } catch (e) {
//         console.error('Send error:', e);
//         setError('Failed to send message');
//       }
//     } else {
//       console.warn('Queueing message - connection not ready');
//       messageQueue.current.push(fullMessage);
//     }
//   }, []);

//   useEffect(() => {
//     if (!manual) {
//       connect();
//     }

//     return () => {
//       disconnect(true);
//       messageQueue.current = [];
//     };
//   }, [connect, disconnect, manual]);

//   return {
//     isConnected,
//     connectionStatus,
//     messages,
//     sendMessage,
//     connect,
//     disconnect,
//     error
//   };
// }

// export default useWebSocket;
// import { useEffect, useRef, useState, useCallback } from 'react';

// export function useWebSocket(url, { manual = false } = {}) {
//   const socketRef = useRef(null);
//   const reconnectTimeoutRef = useRef(null);
//   const [isConnected, setIsConnected] = useState(false);
//   const [connectionStatus, setConnectionStatus] = useState('disconnected');
//   const [messages, setMessages] = useState([]);
//   const [error, setError] = useState(null);

//   const getWebSocketUrl = useCallback(() => {
//     const token = localStorage.getItem('token');
//     if (!token) {
//       setError('Authentication token not found');
//       return null;
//     }

//     try {
//       const wsUrl = new URL(url);
//       wsUrl.searchParams.set('token', token);
//       return wsUrl.toString();
//     } catch (e) {
//       console.error('Invalid WebSocket URL:', e);
//       setError('Invalid WebSocket URL');
//       return null;
//     }
//   }, [url]);

//   const handleIncomingMessage = useCallback((data) => {
//     if (data.type === 'AUTH_ERROR') {
//       console.error('Authentication error:', data.message);
//       setError(data.message);
//       disconnect();
//       localStorage.removeItem('token');
//       window.location.reload();
//       return;
//     }
//     setMessages(prev => [...prev, data]);
//   }, []);

//   const connect = useCallback(() => {
//     const wsUrl = getWebSocketUrl();
//     if (!wsUrl) return;

//     if (socketRef.current && 
//       [WebSocket.OPEN, WebSocket.CONNECTING].includes(socketRef.current.readyState)) {
//       console.warn('WebSocket already connecting or connected');
//       return;
//     }

//     setConnectionStatus('connecting');
//     console.log('Connecting to WebSocket...');

//     socketRef.current = new WebSocket(wsUrl);

//     socketRef.current.onopen = () => {
//       console.log('WebSocket connected');
//       setIsConnected(true);
//       setConnectionStatus('connected');
//       setError(null);
//     };

//     socketRef.current.onmessage = (event) => {
//       try {
//         const parsedData = JSON.parse(event.data);
//         handleIncomingMessage(parsedData);
//       } catch (e) {
//         console.warn('Non-JSON message:', event.data);
//         handleIncomingMessage({ content: event.data });
//       }
//     };

//     socketRef.current.onerror = (event) => {
//       console.error('WebSocket error:', event);
//       setError('WebSocket connection error');
//       setConnectionStatus('error');
//     };

//     socketRef.current.onclose = (event) => {
//       console.log(`WebSocket closed: ${event.code} ${event.reason}`);
//       setIsConnected(false);
//       setConnectionStatus('disconnected');

//       if (!event.wasClean && event.code !== 1000) {
//         console.log('Reconnecting in 3 seconds...');
//         reconnectTimeoutRef.current = setTimeout(() => {
//           connect();
//         }, 3000);
//       }
//     };
//   }, [getWebSocketUrl, handleIncomingMessage]);

//   const disconnect = useCallback((permanent = false) => {
//     if (socketRef.current) {
//       if (permanent) {
//         socketRef.current.onclose = () => {};
//       }
//       socketRef.current.close(
//         permanent ? 1000 : 1001,
//         permanent ? 'Normal closure' : 'Reconnecting'
//       );
//     }
//     if (reconnectTimeoutRef.current) {
//       clearTimeout(reconnectTimeoutRef.current);
//     }
//   }, []);

//   const sendMessage = useCallback((message) => {
//     if (socketRef.current?.readyState === WebSocket.OPEN) {
//       const messageWithAuth = {
//         ...message,
//         timestamp: new Date().toISOString(),
//         user_id: parseInt(localStorage.getItem('userId'), 10),
//         username: localStorage.getItem('username') || 'unknown',
//       };
//       const raw = localStorage.getItem('username');
//       console.log('Stored username:', raw); // должно быть нормальное имя
      
//       try {
//         socketRef.current.send(JSON.stringify(messageWithAuth));
//       } catch (e) {
//         console.error('Error sending message:', e);
//         setError('Failed to send message');
//       }
//     } else {
//       console.error('Cannot send message - WebSocket not open');
//       setError('Connection not ready');
//     }
//   }, []);

//   useEffect(() => {
//     if (!manual) {
//       connect();
//     }

//     return () => {
//       disconnect(true);
//     };
//   }, [connect, disconnect, manual]);

//   return {
//     isConnected,
//     connectionStatus,
//     messages,
//     sendMessage,
//     connect,
//     disconnect,
//     error,
//   };
// }

// export default useWebSocket;