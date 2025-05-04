import React, { useEffect, useState, useCallback } from 'react';
import axios from 'axios';
import Comments from '/Users/darinautalieva/Desktop/GOProject/forum-frontend/src/components/Posts/Comment.js';
import '/Users/darinautalieva/Desktop/GOProject/forum-frontend/src/components/MainLayout.css';

const PostList = ({ refreshTrigger }) => {
    const [posts, setPosts] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    const currentUserId = parseInt(localStorage.getItem('userId'), 10);
    const currentUserRole = localStorage.getItem('userRole');
    const token = localStorage.getItem('token');

    const fetchPosts = useCallback(async () => {
        try {
            setLoading(true);
            setError(null);
    
            const response = await axios.get('http://localhost:8081/api/v1/posts', {
                headers: { 
                    'Accept': 'application/json',
                    'Authorization': `Bearer ${token}`
                }
            });
    
            // Проверка структуры ответа
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

    return (
        <div className="post-list-container">
            {loading && <div className="loading-indicator">Loading posts...</div>}
            {error && <div className="error-message">Error: {error}</div>}

            {posts.map(post => (
                <div key={post.id} className="post-item">
                    <div className="post-header">
                        <h3>{post.title}</h3>
                        {(currentUserId === post.author_id || currentUserRole === 'admin') && (
                            <button
                                onClick={() => handleDeletePost(post.id, post.author_id)}
                                className="delete-button"
                                title={currentUserRole === 'admin' 
                                    ? "Delete post (admin)" 
                                    : "Delete your post"}
                            >
                                ✕
                            </button>
                        )}
                    </div>
                    <div className="post-content">
                        {post.content.split('\n').map((p, i) => (
                            <p key={i}>{p}</p>
                        ))}
                    </div>
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
                    
                    {/* Компонент комментариев */}
                    <Comments postId={post.id} />
                </div>
            ))}
        </div>
    );
};

export default PostList;