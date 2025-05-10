import React from 'react';
import Navbar from './Navbar';
import Chat from '../Chat/Chat';
import '/Users/darinautalieva/Desktop/GOProject/forum-frontend/src/components/MainLayout.css';

const MainLayout = ({ children }) => {
    return (
        <div className="app-container">
            <Navbar />
            <div className="main-content-wrapper">
                <div className="content-and-chat-container">
                    <div className="page-content">{children}</div>
                    <div className="chat-wrapper">
                        <Chat />
                    </div>
                </div>
            </div>
        </div>
    );
};

export default MainLayout;
// import React from 'react';
// import Navbar from './Navbar';
// import Chat from '../Chat/Chat';
// import '/Users/darinautalieva/Desktop/GOProject/forum-frontend/src/components/MainLayout.css';

// const MainLayout = ({ children }) => {
//     return (
//         <div className="app-container">
//             <Navbar />
//             <div className="main-content">
//                 <div className="page-content">{children}</div>
//                 <Chat /> {/* Чат виден всем */}
//             </div>
//         </div>
//     );
// };

// export default MainLayout;
// import React from 'react';
// import Navbar from './Navbar';
// import Chat from '../Chat/Chat';
// import '../MainLayout.css'; // Импортируйте CSS здесь

// const MainLayout = ({ children }) => {
//     const isLoggedIn = localStorage.getItem('token') !== null;

//     return (
//         <>
//             <Navbar />
//             <div className="main-layout">
//                 {isLoggedIn && <Chat />}
//                 <div className="content">{children}</div>
//             </div>
//         </>
//     );
// };

// export default MainLayout;
