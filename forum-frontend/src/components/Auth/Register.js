import React, { useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
// import './auth.css'; // Импортируйте CSS здесь

const Register = () => {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const navigate = useNavigate();

    const handleRegister = async () => {
        try {
            const response = await axios.post('http://localhost:8080/register', {
                username: username,
                password: password,
            });
            //console.log(response.data);
             // localStorage.setItem('token', response.data.token); // Раскомментировать если backend возвращает токен при регистрации
            navigate('/Login');
        } catch (error) {
            console.error('Registration failed:', error);
            alert('Registration failed. Please try again.');
        }
    };

    return (
        <div className="auth-container">
            <h2>Register</h2>
            <input // Заменил Input на input и добавил className
                type="text"
                className="form-control"
                placeholder="Username"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
            />
            <input // Заменил Input на input и добавил className
                type="password"
                className="form-control"
                placeholder="Password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
            />
            <button className="btn btn-primary" onClick={handleRegister}>Register</button>
        </div>
    );
};

export default Register;