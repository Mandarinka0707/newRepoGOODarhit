import { useEffect, useRef, useState, useCallback } from 'react';

export function useWebSocket(url, { manual = false } = {}) {
  const socketRef = useRef(null);
  const reconnectTimeoutRef = useRef(null);
  const [isConnected, setIsConnected] = useState(false);
  const [connectionStatus, setConnectionStatus] = useState('disconnected');
  const [messages, setMessages] = useState([]);
  const [error, setError] = useState(null);

  const getWebSocketUrl = useCallback(() => {
    const token = localStorage.getItem('token');
    if (!token) {
      setError('Authentication token not found');
      return null;
    }

    try {
      const wsUrl = new URL(url);
      wsUrl.searchParams.set('token', token);
      return wsUrl.toString();
    } catch (e) {
      console.error('Invalid WebSocket URL:', e);
      setError('Invalid WebSocket URL');
      return null;
    }
  }, [url]);

  const handleIncomingMessage = useCallback((data) => {
    if (data.type === 'AUTH_ERROR') {
      console.error('Authentication error:', data.message);
      setError(data.message);
      disconnect();
      localStorage.removeItem('token');
      window.location.reload();
      return;
    }
    setMessages(prev => [...prev, data]);
  }, []);

  const connect = useCallback(() => {
    const wsUrl = getWebSocketUrl();
    if (!wsUrl) return;

    if (socketRef.current && 
      [WebSocket.OPEN, WebSocket.CONNECTING].includes(socketRef.current.readyState)) {
      console.warn('WebSocket already connecting or connected');
      return;
    }

    setConnectionStatus('connecting');
    console.log('Connecting to WebSocket...');

    socketRef.current = new WebSocket(wsUrl);

    socketRef.current.onopen = () => {
      console.log('WebSocket connected');
      setIsConnected(true);
      setConnectionStatus('connected');
      setError(null);
    };

    socketRef.current.onmessage = (event) => {
      try {
        const parsedData = JSON.parse(event.data);
        handleIncomingMessage(parsedData);
      } catch (e) {
        console.warn('Non-JSON message:', event.data);
        handleIncomingMessage({ content: event.data });
      }
    };

    socketRef.current.onerror = (event) => {
      console.error('WebSocket error:', event);
      setError('WebSocket connection error');
      setConnectionStatus('error');
    };

    socketRef.current.onclose = (event) => {
      console.log(`WebSocket closed: ${event.code} ${event.reason}`);
      setIsConnected(false);
      setConnectionStatus('disconnected');

      if (!event.wasClean && event.code !== 1000) {
        console.log('Reconnecting in 3 seconds...');
        reconnectTimeoutRef.current = setTimeout(() => {
          connect();
        }, 3000);
      }
    };
  }, [getWebSocketUrl, handleIncomingMessage]);

  const disconnect = useCallback((permanent = false) => {
    if (socketRef.current) {
      if (permanent) {
        socketRef.current.onclose = () => {};
      }
      socketRef.current.close(
        permanent ? 1000 : 1001,
        permanent ? 'Normal closure' : 'Reconnecting'
      );
    }
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current);
    }
  }, []);

  const sendMessage = useCallback((message) => {
    if (socketRef.current?.readyState === WebSocket.OPEN) {
      const messageWithAuth = {
        ...message,
        timestamp: new Date().toISOString(),
        user_id: parseInt(localStorage.getItem('userId'), 10),
        username: localStorage.getItem('username') || 'unknown',
      };

      try {
        socketRef.current.send(JSON.stringify(messageWithAuth));
      } catch (e) {
        console.error('Error sending message:', e);
        setError('Failed to send message');
      }
    } else {
      console.error('Cannot send message - WebSocket not open');
      setError('Connection not ready');
    }
  }, []);

  useEffect(() => {
    if (!manual) {
      connect();
    }

    return () => {
      disconnect(true);
    };
  }, [connect, disconnect, manual]);

  return {
    isConnected,
    connectionStatus,
    messages,
    sendMessage,
    connect,
    disconnect,
    error,
  };
}

export default useWebSocket;