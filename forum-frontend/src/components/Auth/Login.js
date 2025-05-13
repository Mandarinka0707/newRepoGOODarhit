
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
            const response = await axios.post('http://localhost:8080/api/v1/auth/login', { username, password });
            console.log('Login response:', response.data);
            const token = response.data.token;
            localStorage.setItem('token', token);

            // Распарсить user_id из токена
            const base64Url = token.split('.')[1];
            const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
            const payload = JSON.parse(atob(base64));
            console.log(payload);
            localStorage.setItem('userId', String(payload.user_id));
            localStorage.setItem('username', response.data.username);
            localStorage.setItem('userRole', payload.role); 
            console.log('username', response.data.username);
            
        console.log('User role:', payload.role);
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