import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import '../MainLayout.css'; // Импортируйте CSS здесь

const Login = () => {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const navigate = useNavigate();

    useEffect(() => {
        const token = localStorage.getItem('token');
        if (token) {
            navigate('/posts');
        }
    }, [navigate]);

    const handleLogin = async () => {
        try {
          const response = await axios.post('http://localhost:8080/login', { 
            username, 
            password 
          });
          localStorage.setItem('token', response.data.token);
          navigate('/chat');
        } catch (error) {
          console.error('Login failed:', error);
        }
      };

    return (
        <div className="auth-container"> {/* Добавьте класс для стилизации */}
            <h2>Login</h2>
            <input
                type="text"
                className="form-control"
                placeholder="Username"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
            />
            <input
                type="password"
                className="form-control"
                placeholder="Password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
            />
            <button className="btn btn-primary" onClick={handleLogin}>Login</button>
        </div>
    );
};

export default Login;
