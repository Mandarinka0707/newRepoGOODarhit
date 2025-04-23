import React, { useEffect, useState, useCallback } from 'react';
import axios from 'axios';
import '../MainLayout.css';

const PostList = ({ refreshTrigger }) => {
    const [posts, setPosts] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    const fetchPosts = useCallback(async () => {
        let isMounted = true;
    
        try {
            setLoading(true);
            setError(null);

            const token = localStorage.getItem('token');
            if (!token) {
                throw new Error('Authentication token not found');
            }

            // 1. Get posts
            const response = await axios.get('http://localhost:8081/posts', {
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Accept': 'application/json'
                }
            });

            if (!isMounted) return;

            // 2. Check response structure
            const postsData = response.data?.data || response.data;
            
            if (!postsData || !Array.isArray(postsData)) {
                throw new Error('Invalid data structure received from server');
            }

            // 3. Process posts
            const processedPosts = await Promise.all(
                postsData.map(async (post) => {
                    let username = `User ${post.author_id}`;
                    
                    // 4. If no author name, fetch it
                    if (!post.author_name) {
                        try {
                            const userResponse = await axios.get(
                                `http://localhost:8081/api/users/${post.author_id}`,
                                { headers: { 'Authorization': `Bearer ${token}` } }
                            );
                            username = userResponse.data.username || username;
                        } catch (err) {
                            console.error('Failed to fetch username:', err);
                        }
                    } else {
                        username = post.author_name;
                    }

                    return {
                        id: post.id,
                        title: post.title || 'Untitled Post',
                        content: post.content || 'No content available',
                        author_id: post.author_id,
                        author_name: username,
                        created_at: post.created_at || new Date().toISOString()
                    };
                })
            );

            if (isMounted) {
                setPosts(processedPosts);
            }
        } catch (err) {
            if (isMounted) {
                console.error('Post fetch error:', {
                    error: err,
                    response: err.response
                });

                setError(
                    err.response?.data?.error || 
                    err.response?.data?.message || 
                    err.message || 
                    'Failed to load posts. Please try again later.'
                );
            }
        } finally {
            if (isMounted) {
                setLoading(false);
            }
        }
    }, []);

    useEffect(() => {
        const abortController = new AbortController();
        fetchPosts();
        return () => abortController.abort();
    }, [fetchPosts, refreshTrigger]);

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
                <p className="error-message">{error}</p>
                <button
                    onClick={fetchPosts}
                    className="retry-button"
                >
                    Retry
                </button>
            </div>
        );
    }

    return (
        <div className="post-list-container">
            {posts.length > 0 ? (
                [...posts].reverse().map((post) => (
                    <div key={`post-${post.id}`} className="post-item">
                        <h3>{post.title}</h3>
                        <div className="post-content">
                            {post.content.split('\n').map((paragraph, i) => (
                                <p key={i}>{paragraph}</p>
                            ))}
                        </div>
                        <div className="post-meta">
                            <span className="author">Author: {post.author_name}</span>
                            <span className="separator"> | </span>
                            <span className="date">
                                {new Date(post.created_at).toLocaleDateString()}, 
                                {new Date(post.created_at).toLocaleTimeString()}
                            </span>
                        </div>
                    </div>
                ))
            ) : (
                <div className="no-posts-message">
                    <p>No posts found. Create the first one!</p>
                </div>
            )}
        </div>
    );
};

export default PostList;