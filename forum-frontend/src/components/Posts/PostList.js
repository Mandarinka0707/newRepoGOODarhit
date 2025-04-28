import React, { useEffect, useState, useCallback } from 'react';
import axios from 'axios';
import '../MainLayout.css';

const PostList = ({ refreshTrigger }) => {
    const [posts, setPosts] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    const handleDeletePost = async (postId, authorId) => {
        const confirmDelete = window.confirm("Are you sure you want to delete this post?");
        if (!confirmDelete) return;

        try {
            const token = localStorage.getItem('token');
            const currentUserId = parseInt(localStorage.getItem('userId'));
            if (!currentUserId || currentUserId !== authorId) {
                alert("You can only delete your own posts");
                return;
            }

            await axios.delete(`http://localhost:8081/api/v1/posts/${postId}`, {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });

            fetchPosts();
            alert("Post deleted successfully");
        } catch (error) {
            console.error('Delete post error:', error);
            alert(error.response?.data?.error || 'Failed to delete post');
        }
    };

    const fetchPosts = useCallback(async () => {
        let isMounted = true;

        try {
            setLoading(true);
            setError(null);

            const token = localStorage.getItem('token');
            if (!token) {
                throw new Error('Authentication token not found');
            }

            const response = await axios.get('http://localhost:8081/api/v1/posts', {
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Accept': 'application/json'
                }
            });

            if (!isMounted) return;

            const postsData = response.data?.data || response.data;

            if (!postsData || !Array.isArray(postsData)) {
                throw new Error('Invalid data structure received from server');
            }

            const processedPosts = await Promise.all(
                postsData.map(async (post) => {
                    let username = `User ${post.author_id}`;

                    if (!post.author_name) {
                        try {
                            const userResponse = await axios.get(
                                `http://localhost:8081/api/users/${post.author_id}`,
                                { 
                                    headers: { 
                                        'Authorization': `Bearer ${token}` 
                                    } 
                                }
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
                    err.message ||
                    'Failed to load posts. Please try again later.'
                );
            }
        } finally {
            if (isMounted) {
                setLoading(false);
            }
        }

        return () => { isMounted = false };
    }, []);

    useEffect(() => {
        const abortController = new AbortController();
        fetchPosts();
        return () => {
            abortController.abort();
        };
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
                        <div className="post-header">
                            <h3>{post.title}</h3>
                            {post.author_id === Number(localStorage.getItem('userId')) && (
                                <button
                                    onClick={() => handleDeletePost(post.id, post.author_id)}
                                    className="delete-button"
                                    title="Delete post"
                                >
                                    Delete
                                </button>
                            )}
                        </div>
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