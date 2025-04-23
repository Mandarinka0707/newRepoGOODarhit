import React, { useEffect, useState } from 'react';
import axios from 'axios';
import '../MainLayout.css';

const PostList = ({ refreshTrigger }) => {
    const [posts, setPosts] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    const fetchPosts = async () => {
        try {
            setLoading(true);
            setError(null);
            
            const token = localStorage.getItem('token');
            if (!token) {
                throw new Error('Authentication token not found');
            }

            const response = await axios.get('http://localhost:8081/posts', {
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Accept': 'application/json'
                }
            });

            // Validate response structure
            if (!response.data || !Array.isArray(response.data)) {
                throw new Error('Invalid data structure received from server');
            }

            // Process posts and ensure they have valid IDs
            const processedPosts = response.data.map((post, index) => {
                // Log posts without IDs for debugging
                if (!post.id) {
                    console.warn('Post missing ID:', post);
                }
                
                return {
                    ...post,
                    id: post.id || `temp-${index}-${Date.now()}`,
                    title: post.title || 'Untitled Post',
                    content: post.content || 'No content available',
                    author_id: post.author_id || 'Unknown',
                    created_at: post.created_at || new Date().toISOString()
                };
            });

            setPosts(processedPosts);
        } catch (err) {
            console.error('Error fetching posts:', {
                error: err,
                response: err.response,
                message: err.message
            });
            
            setError(
                err.response?.data?.message || 
                err.message || 
                'Failed to load posts. Please try again.'
            );
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchPosts();
    }, [refreshTrigger]);

    if (loading) {
        return (
            <div className="loading-container">
                <div className="spinner"></div>
                <p>Loading posts...</p>
            </div>
        );
    }

    if (error) {
        return (
            <div className="error-container">
                <p className="error-message">Error: {error}</p>
                <button 
                    onClick={fetchPosts}
                    className="retry-button"
                >
                    Try Again
                </button>
            </div>
        );
    }

    return (
        <div className="posts-container">
            {posts.length > 0 ? (
                [...posts].reverse().map(post => (
                    <div key={post.id} className="post-card">
                        <h3 className="post-title">{post.title}</h3>
                        <p className="post-content">{post.content}</p>
                        <div className="post-meta">
                            <span className="post-author">Author: {post.author_id}</span>
                            <span className="post-date">
                                {new Date(post.created_at).toLocaleString()}
                            </span>
                        </div>
                    </div>
                ))
            ) : (
                <div className="empty-state">
                    <p>No posts available yet</p>
                </div>
            )}
        </div>
    );
};

export default PostList;