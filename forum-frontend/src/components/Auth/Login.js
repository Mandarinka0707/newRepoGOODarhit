import React, { useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import '/Users/darinautalieva/Desktop/GOProject/forum-frontend/src/components/MainLayout.css';

const Login = () => {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const navigate = useNavigate();

    const handleLogin = async () => {
        try {
            const response = await axios.post('http://localhost:8080/login', { username, password });
            localStorage.setItem('token', response.data.token);
            navigate('/posts');
        } catch (error) {
            console.error('Login failed:', error);
            alert('Login failed. Please check your credentials.');
        }
    };

    return (
        <div className="auth-container"> {/* Добавьте класс для стилизации */}
            <h2>Login</h2>
            <input  //Заменил Input на input и добавил className
                type="text"
                className="form-control"
                placeholder="Username"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
            />
            <input //Заменил Input на input и добавил className
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