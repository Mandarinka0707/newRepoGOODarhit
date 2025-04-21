import React from 'react';
import Navbar from './Navbar';
import Chat from '../Chat/Chat';
import '../MainLayout.css'; // Импортируйте CSS здесь

const MainLayout = ({ children }) => {
    const isLoggedIn = localStorage.getItem('token') !== null;

    return (
        <>
            <Navbar />
            <div className="main-layout">
                {isLoggedIn && <Chat />}
                <div className="content">{children}</div>
            </div>
        </>
    );
};

export default MainLayout;
