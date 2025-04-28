import React, { useState } from 'react';
import { BrowserRouter as Router, Route, Routes, Navigate } from 'react-router-dom';
import Login from './components/Auth/Login';
import Register from './components/Auth/Register';
import PostList from './components/Posts/PostList';
import CreatePost from './components/Posts/CreatePost';
import MainLayout from './components/Layout/MainLayout';
import Chat from './components/Chat/Chat'; // Добавьте импорт
import './components/MainLayout.css';

const PrivateRoute = ({ children }) => {
    const token = localStorage.getItem('token');
    return token ? children : <Navigate to="/login" />;
};

const App = () => {
    const [refreshPosts, setRefreshPosts] = useState(false);

    const onPostCreated = () => {
        setRefreshPosts(prev => !prev);
    };

    return (
        <Router>
            <MainLayout>
                <Routes>
                    <Route path="/register" element={<Register />} />
                    <Route path="/login" element={<Login />} />
                    <Route path="/posts" element={
                        <PrivateRoute>
                            <>
                                <CreatePost onPostCreated={onPostCreated} />
                                <PostList key={refreshPosts} />
                            </>
                        </PrivateRoute>
                    } />
                    <Route path="/chat" element={
                        <PrivateRoute>
                            <Chat />
                        </PrivateRoute>
                    } />
                    <Route path="/" element={<Navigate to="/login" />} />
                </Routes>
            </MainLayout>
        </Router>
    );
};

export default App;