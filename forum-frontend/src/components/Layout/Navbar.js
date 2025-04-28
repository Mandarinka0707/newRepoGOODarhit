import React from 'react';
 import '../MainLayout.css'; // Импортируйте CSS здесь
import { Link, useNavigate } from 'react-router-dom';

const Navbar = () => {
    const navigate = useNavigate();
    const isLoggedIn = localStorage.getItem('token') !== null;

    const handleLogout = () => {
        localStorage.removeItem('token');
        navigate('/login');
    };

    return (
        <nav className="navbar">
            <h1 className="navTitle">Forum</h1>
            <div className="navLinks">
                {isLoggedIn ? (
                    <>
                        <Link to="/posts">Posts</Link>
                        <Link to="#" onClick={handleLogout}>Logout</Link>
                    </>
                ) : (
                    <>
                        <Link to="/login">Login</Link>
                        <Link to="/register">Register</Link>
                    </>
                )}
            </div>
        </nav>
    );
};

export default Navbar;
