import React from 'react';
import 'forum-frontend/src/components/Layout/Navbar.js';
const MainLayout = ({ children }) => {
    return (
        <>
            <Navbar />
            <Content>{children}</Content>
        </>
    );
};

export default MainLayout;