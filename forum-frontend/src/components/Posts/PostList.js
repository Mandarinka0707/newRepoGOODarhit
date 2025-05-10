import React, { useEffect, useState, useCallback } from 'react';
import axios from 'axios';
import Comments from '/Users/darinautalieva/Desktop/GOProject/forum-frontend/src/components/Posts/Comment.js';
import '/Users/darinautalieva/Desktop/GOProject/forum-frontend/src/components/MainLayout.css';
import { useNavigate } from 'react-router-dom';

const PostList = ({ refreshTrigger }) => {
    const [posts, setPosts] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [editingPostId, setEditingPostId] = useState(null);
    const [editFormData, setEditFormData] = useState({
        title: '',
        content: ''
    });

    const currentUserId = localStorage.getItem('userId') ? parseInt(localStorage.getItem('userId'), 10) : null;
    const currentUserRole = localStorage.getItem('userRole');
    const token = localStorage.getItem('token');
    const isAuthenticated = !!token;
    const navigate = useNavigate();

    const fetchPosts = useCallback(async () => {
        try {
            setLoading(true);
            setError(null);
            setEditingPostId(null);

            const response = await axios.get('http://localhost:8081/api/v1/posts', {
                headers: { 
                    'Accept': 'application/json',
                    ...(token && { 'Authorization': `Bearer ${token}` })
                }
            });

            const rawPosts = response.data.data || response.data;

            if (!Array.isArray(rawPosts)) {
                throw new Error('Invalid posts data format');
            }

            const processedPosts = rawPosts.map(post => ({
                ...post,
                id: parseInt(post.id, 10),
                author_id: parseInt(post.author_id, 10),
                created_at: new Date(post.created_at).toISOString()
            }));

            setPosts(processedPosts);
        } catch (err) {
            setError(err.message || 'Failed to load posts');
        } finally {
            setLoading(false);
        }
    }, [token]);

    useEffect(() => {
        fetchPosts();
    }, [fetchPosts, refreshTrigger]);

    const handleDeletePost = async (postId, authorId) => {
        if (!isAuthenticated) {
            navigate('/login');
            return;
        }

        const confirmDelete = window.confirm("Are you sure you want to delete this post?");
        if (!confirmDelete) return;

        try {
            await axios.delete(`http://localhost:8081/api/v1/posts/${postId}`, {
                headers: { 
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                }
            });

            setPosts(prev => prev.filter(post => post.id !== postId));
        } catch (error) {
            console.error('Delete post error:', error);
            const errorMessage = error.response?.data?.error || 
                               error.response?.data?.message || 
                               'Failed to delete post';
            
            alert(errorMessage);
            await fetchPosts();
        }
    };

    const startEditing = (post) => {
        if (!isAuthenticated) {
            navigate('/login');
            return;
        }
        setEditingPostId(post.id);
        setEditFormData({
            title: post.title,
            content: post.content
        });
    };

    const cancelEditing = () => {
        setEditingPostId(null);
    };

    const handleEditChange = (e) => {
        const { name, value } = e.target;
        setEditFormData(prev => ({
            ...prev,
            [name]: value
        }));
    };

    const handleUpdatePost = async (postId) => {
        try {
            await axios.put(
                `http://localhost:8081/api/v1/posts/${postId}`,
                {
                    title: editFormData.title,
                    content: editFormData.content
                },
                {
                    headers: {
                        'Authorization': `Bearer ${token}`,
                        'Content-Type': 'application/json'
                    }
                }
            );
            
            setEditingPostId(null);
            await fetchPosts();
        } catch (error) {
            console.error('Update post error:', error);
            const errorMessage = error.response?.data?.error || 
                               error.response?.data?.message || 
                               'Failed to update post';
            alert(errorMessage);
        }
    };

    return (
        <div className="post-list-container">
            {loading && <div className="loading-indicator">Loading posts...</div>}
            {error && <div className="error-message">Error: {error}</div>}

            {posts.map(post => (
                <div key={post.id} className="post-item">
                    {editingPostId === post.id ? (
                        <div className="edit-form">
                            <input
                                type="text"
                                name="title"
                                value={editFormData.title}
                                onChange={handleEditChange}
                                className="edit-title-input"
                            />
                            <textarea
                                name="content"
                                value={editFormData.content}
                                onChange={handleEditChange}
                                className="edit-content-input"
                                rows={5}
                            />
                            <div className="edit-actions">
                                <button 
                                    onClick={() => handleUpdatePost(post.id)}
                                    className="save-button"
                                >
                                    Save
                                </button>
                                <button 
                                    onClick={cancelEditing}
                                    className="cancel-button"
                                >
                                    Cancel
                                </button>
                            </div>
                        </div>
                    ) : (
                        <>
                            <div className="post-header">
                                <h3>{post.title}</h3>
                                {isAuthenticated && (currentUserId === post.author_id || currentUserRole === 'admin') && (
                                    <div className="post-actions">
                                        <button
                                            onClick={() => startEditing(post)}
                                            className="edit-button"
                                            title="Edit post"
                                        >
                                            ✎
                                        </button>
                                        <button
                                            onClick={() => handleDeletePost(post.id, post.author_id)}
                                            className="delete-button"
                                            title={currentUserRole === 'admin' 
                                                ? "Delete post (admin)" 
                                                : "Delete your post"}
                                        >
                                            ✕
                                        </button>
                                    </div>
                                )}
                            </div>
                            <div className="post-content">
                                {post.content.split('\n').map((p, i) => (
                                    <p key={i}>{p}</p>
                                ))}
                            </div>
                        </>
                    )}
                    
                    <div className="post-meta">
                        <span className="author">{post.author_name}</span>
                        <span className="separator">•</span>
                        <span className="timestamp">
                            {new Date(post.created_at).toLocaleDateString('en-US', {
                                year: 'numeric',
                                month: 'long',
                                day: 'numeric',
                                hour: '2-digit',
                                minute: '2-digit'
                            })}
                        </span>
                        {currentUserRole === 'admin' && post.author_id !== currentUserId && (
                            <span className="admin-badge">(admin action)</span>
                        )}
                    </div>
                    
                    <Comments postId={post.id} />
                </div>
            ))}
        </div>
    );
};

export default PostList;
// import React, { useEffect, useState, useCallback } from 'react';
// import axios from 'axios';
// import Comments from '/Users/darinautalieva/Desktop/GOProject/forum-frontend/src/components/Posts/Comment.js';
// import '/Users/darinautalieva/Desktop/GOProject/forum-frontend/src/components/MainLayout.css';

// const PostList = ({ refreshTrigger }) => {
//     const [posts, setPosts] = useState([]);
//     const [loading, setLoading] = useState(true);
//     const [error, setError] = useState(null);
//     const [editingPostId, setEditingPostId] = useState(null);
//     const [editFormData, setEditFormData] = useState({
//         title: '',
//         content: ''
//     });

//     const currentUserId = parseInt(localStorage.getItem('userId'), 10);
//     const currentUserRole = localStorage.getItem('userRole');
//     const token = localStorage.getItem('token');

//     const fetchPosts = useCallback(async () => {
//         try {
//             setLoading(true);
//             setError(null);
//             setEditingPostId(null); // Сбрасываем режим редактирования
    
//             const response = await axios.get('http://localhost:8081/api/v1/posts', {
//                 headers: { 
//                     'Accept': 'application/json',
//                     'Authorization': `Bearer ${token}`
//                 }
//             });
    
//             const rawPosts = response.data.data || response.data;
    
//             if (!Array.isArray(rawPosts)) {
//                 throw new Error('Invalid posts data format');
//             }
    
//             const processedPosts = rawPosts.map(post => ({
//                 ...post,
//                 id: parseInt(post.id, 10),
//                 author_id: parseInt(post.author_id, 10),
//                 created_at: new Date(post.created_at).toISOString()
//             }));
    
//             setPosts(processedPosts);
//         } catch (err) {
//             setError(err.message || 'Failed to load posts');
//         } finally {
//             setLoading(false);
//         }
//     }, [token]);

//     useEffect(() => {
//         fetchPosts();
//     }, [fetchPosts, refreshTrigger]);

//     const handleDeletePost = async (postId, authorId) => {
//         const confirmDelete = window.confirm("Are you sure you want to delete this post?");
//         if (!confirmDelete) return;

//         try {
//             await axios.delete(`http://localhost:8081/api/v1/posts/${postId}`, {
//                 headers: { 
//                     'Authorization': `Bearer ${token}`,
//                     'Content-Type': 'application/json'
//                 }
//             });

//             setPosts(prev => prev.filter(post => post.id !== postId));
//         } catch (error) {
//             console.error('Delete post error:', error);
//             const errorMessage = error.response?.data?.error || 
//                                error.response?.data?.message || 
//                                'Failed to delete post';
            
//             alert(errorMessage);
//             await fetchPosts();
//         }
//     };

//     const startEditing = (post) => {
//         setEditingPostId(post.id);
//         setEditFormData({
//             title: post.title,
//             content: post.content
//         });
//     };

//     const cancelEditing = () => {
//         setEditingPostId(null);
//     };

//     const handleEditChange = (e) => {
//         const { name, value } = e.target;
//         setEditFormData(prev => ({
//             ...prev,
//             [name]: value
//         }));
//     };

//     const handleUpdatePost = async (postId) => {
//         try {
//             await axios.put(
//                 `http://localhost:8081/api/v1/posts/${postId}`,
//                 {
//                     title: editFormData.title,
//                     content: editFormData.content
//                 },
//                 {
//                     headers: {
//                         'Authorization': `Bearer ${token}`,
//                         'Content-Type': 'application/json'
//                     }
//                 }
//             );
            
//             setEditingPostId(null);
//             await fetchPosts();
//         } catch (error) {
//             console.error('Update post error:', error);
//             const errorMessage = error.response?.data?.error || 
//                                error.response?.data?.message || 
//                                'Failed to update post';
//             alert(errorMessage);
//         }
//     };

//     return (
//         <div className="post-list-container">
//             {loading && <div className="loading-indicator">Loading posts...</div>}
//             {error && <div className="error-message">Error: {error}</div>}

//             {posts.map(post => (
//                 <div key={post.id} className="post-item">
//                     {editingPostId === post.id ? (
//                         <div className="edit-form">
//                             <input
//                                 type="text"
//                                 name="title"
//                                 value={editFormData.title}
//                                 onChange={handleEditChange}
//                                 className="edit-title-input"
//                             />
//                             <textarea
//                                 name="content"
//                                 value={editFormData.content}
//                                 onChange={handleEditChange}
//                                 className="edit-content-input"
//                                 rows={5}
//                             />
//                             <div className="edit-actions">
//                                 <button 
//                                     onClick={() => handleUpdatePost(post.id)}
//                                     className="save-button"
//                                 >
//                                     Save
//                                 </button>
//                                 <button 
//                                     onClick={cancelEditing}
//                                     className="cancel-button"
//                                 >
//                                     Cancel
//                                 </button>
//                             </div>
//                         </div>
//                     ) : (
//                         <>
//                             <div className="post-header">
//                                 <h3>{post.title}</h3>
//                                 <div className="post-actions">
//                                     {(currentUserId === post.author_id || currentUserRole === 'admin') && (
//                                         <>
//                                             <button
//                                                 onClick={() => startEditing(post)}
//                                                 className="edit-button"
//                                                 title="Edit post"
//                                             >
//                                                 ✎
//                                             </button>
//                                             <button
//                                                 onClick={() => handleDeletePost(post.id, post.author_id)}
//                                                 className="delete-button"
//                                                 title={currentUserRole === 'admin' 
//                                                     ? "Delete post (admin)" 
//                                                     : "Delete your post"}
//                                             >
//                                                 ✕
//                                             </button>
//                                         </>
//                                     )}
//                                 </div>
//                             </div>
//                             <div className="post-content">
//                                 {post.content.split('\n').map((p, i) => (
//                                     <p key={i}>{p}</p>
//                                 ))}
//                             </div>
//                         </>
//                     )}
                    
//                     <div className="post-meta">
//                         <span className="author">{post.author_name}</span>
//                         <span className="separator">•</span>
//                         <span className="timestamp">
//                             {new Date(post.created_at).toLocaleDateString('en-US', {
//                                 year: 'numeric',
//                                 month: 'long',
//                                 day: 'numeric',
//                                 hour: '2-digit',
//                                 minute: '2-digit'
//                             })}
//                         </span>
//                         {currentUserRole === 'admin' && post.author_id !== currentUserId && (
//                             <span className="admin-badge">(admin action)</span>
//                         )}
//                     </div>
                    
//                     <Comments postId={post.id} />
//                 </div>
//             ))}
//         </div>
//     );
// };

// export default PostList;