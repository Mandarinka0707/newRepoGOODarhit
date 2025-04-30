import React, { useEffect, useState, useCallback } from 'react';
import axios from 'axios';
import '../MainLayout.css';

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

            const processedPosts = response.data.data.map(post => ({
                ...post,
                id: parseInt(post.id, 10),
                author_id: parseInt(post.author_id, 10),
                created_at: new Date(post.created_at).toISOString()
            }));

            setPosts(processedPosts);
        } catch (err) {
            setError(err.response?.data?.error || err.message);
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
                </div>
            ))}
        </div>
    );
};

export default PostList;

// import React, { useEffect, useState, useCallback } from 'react';
// import axios from 'axios';
// import '../MainLayout.css';

// const PostList = ({ refreshTrigger }) => {
//     const [posts, setPosts] = useState([]);
//     const [loading, setLoading] = useState(true);
//     const [error, setError] = useState(null);

//     const handleDeletePost = async (postId, authorId) => {
//         const confirmDelete = window.confirm("Are you sure you want to delete this post?");
//         if (!confirmDelete) return;
    
//         const token = localStorage.getItem('token');
//         const currentUserId = parseInt(localStorage.getItem('userId'), 10);
        
//         console.log('Trying to delete post:', postId);
//         console.log('Current user:', currentUserId, 'Post author:', authorId);
    
//         if (currentUserId !== authorId) {
//             alert("You can only delete your own posts");
//             return;
//         }
    
//         try {
//             await axios.delete(`http://localhost:8081/api/v1/posts/${postId}`, {
//                 headers: {
//                     'Authorization': `Bearer ${token}`
//                 }
//             });
    
//             setPosts(prev => prev.filter(post => post.id !== postId));
//             alert("Post deleted successfully");
//         } catch (error) {
//             console.error('Delete post error:', error);
//             alert(error.response?.data?.error || 'Failed to delete post');
//         }
//     };
    

//     const fetchPosts = useCallback(async () => {
//         try {
//             setLoading(true);
//             setError(null);

//             const response = await axios.get('http://localhost:8081/api/v1/posts', {
//                 headers: {
//                     'Accept': 'application/json'
//                 }
//             });

//             const processedPosts = response.data.data.map(post => ({
//                 ...post,
//                 author_id: parseInt(post.author_id, 10),
//                 created_at: new Date(post.created_at).toISOString()
//             }));

//             setPosts(processedPosts);
//         } catch (err) {
//             setError(err.response?.data?.error || err.message);
//         } finally {
//             setLoading(false);
//         }
//     }, []);

//     useEffect(() => {
//         fetchPosts();
//     }, [fetchPosts, refreshTrigger]);

//     const currentUserId = parseInt(localStorage.getItem('userId'), 10);

//     return (
//         <div className="post-list-container">
//             {posts.map(post => (
//                 <div key={post.id} className="post-item">
//                     <div className="post-header">
//                         <h3>{post.title}</h3>
//                         {(() => {
//     console.log('Текущий пользователь:', currentUserId);
//     console.log('Автор поста:', post.author_id);
//     return currentUserId === post.author_id ? (
//       <button
//         onClick={() => handleDeletePost(post.id, post.author_id)}
//         className="delete-button"
//         title="Delete post"
//       >
//         ✕
//       </button>
//     ) : null;
//   })()}
//                     </div>
//                     <div className="post-content">
//                         {post.content.split('\n').map((p, i) => (
//                             <p key={i}>{p}</p>
//                         ))}
//                     </div>
//                     <div className="post-meta">
//                         <span>By {post.author_name}</span>
//                         <span> • </span>
//                         <span>
//                             {new Date(post.created_at).toLocaleDateString('en-US', {
//                                 year: 'numeric',
//                                 month: 'long',
//                                 day: 'numeric',
//                                 hour: '2-digit',
//                                 minute: '2-digit'
//                             })}
//                         </span>
//                     </div>
//                 </div>
//             ))}
//         </div>
//     );
// };

// export default PostList;
// import React, { useEffect, useState, useCallback } from 'react';
// import axios from 'axios';
// import '../MainLayout.css';

// const PostList = ({ refreshTrigger }) => {
//     const [posts, setPosts] = useState([]);
//     const [loading, setLoading] = useState(true);
//     const [error, setError] = useState(null);

//     const handleDeletePost = async (postId, authorId) => {
//         const confirmDelete = window.confirm("Are you sure you want to delete this post?");
//         if (!confirmDelete) return;

//         try {
//             const token = localStorage.getItem('token');
//             const currentUserId = parseInt(localStorage.getItem('userId'));
//             if (!currentUserId || currentUserId !== authorId) {
//                 alert("You can only delete your own posts");
//                 return;
//             }

//             await axios.delete(`http://localhost:8081/api/v1/posts/${postId}`, {
//                 headers: {
//                     'Authorization': `Bearer ${token}`
//                 }
//             })

//             fetchPosts();
//             alert("Post deleted successfully");
//         } catch (error) {
//             console.error('Delete post error:', error);
//             alert(error.response?.data?.error || 'Failed to delete post');
//         }
//     };

//     const fetchPosts = useCallback(async () => {
//         let isMounted = true;

//         try {
//             setLoading(true);
//             setError(null);

//             const token = localStorage.getItem('token');
//             if (!token) {
//                 throw new Error('Authentication token not found');
//             }

//             const response = await axios.get('http://localhost:8081/api/v1/posts', {
//                 headers: {
//                     'Authorization': `Bearer ${token}`,
//                     'Accept': 'application/json'
//                 }
//             });

//             if (!isMounted) return;

//             const postsData = response.data?.data || response.data;

//             if (!postsData || !Array.isArray(postsData)) {
//                 throw new Error('Invalid data structure received from server');
//             }

//             const processedPosts = await Promise.all(
//                 postsData.map(async (post) => {
//                     let username = `User ${post.author_id}`;

//                     if (!post.author_name) {
//                         try {
//                             const userResponse = await axios.get(
//                                 `http://localhost:8081/api/users/${post.author_id}`,
//                                 { 
//                                     headers: { 
//                                         'Authorization': `Bearer ${token}` 
//                                     } 
//                                 }
//                             );
//                             username = userResponse.data.username || username;
//                         } catch (err) {
//                             console.error('Failed to fetch username:', err);
//                         }
//                     } else {
//                         username = post.author_name;
//                     }

//                     return {
//                         id: post.id,
//                         title: post.title || 'Untitled Post',
//                         content: post.content || 'No content available',
//                         author_id: post.author_id,
//                         author_name: username,
//                         created_at: post.created_at || new Date().toISOString()
//                     };
//                 })
//             );

//             if (isMounted) {
//                 setPosts(processedPosts);
//             }
//         } catch (err) {
//             if (isMounted) {
//                 console.error('Post fetch error:', {
//                     error: err,
//                     response: err.response
//                 });

//                 setError(
//                     err.response?.data?.error ||
//                     err.message ||
//                     'Failed to load posts. Please try again later.'
//                 );
//             }
//         } finally {
//             if (isMounted) {
//                 setLoading(false);
//             }
//         }

//         return () => { isMounted = false };
//     }, []);

//     useEffect(() => {
//         const abortController = new AbortController();
//         fetchPosts();
//         return () => {
//             abortController.abort();
//         };
//     }, [fetchPosts, refreshTrigger]);

//     if (loading) {
//         return (
//             <div className="loading-container">
//                 <div className="spinner"></div>
//                 <p>Loading posts...</p>
//             </div>
//         );
//     }

//     if (error) {
//         return (
//             <div className="error-container">
//                 <p className="error-message">{error}</p>
//                 <button
//                     onClick={fetchPosts}
//                     className="retry-button"
//                 >
//                     Retry
//                 </button>
//             </div>
//         );
//     }

//     return (
//         <div className="post-list-container">
//             {posts.length > 0 ? (
//                 [...posts].reverse().map((post) => (
//                     <div key={`post-${post.id}`} className="post-item">
//                         <div className="post-header">
//                             <h3>{post.title}</h3>
//                             {post.author_id === Number(localStorage.getItem('userId')) && (
//                                 <button
//                                     onClick={() => handleDeletePost(post.id, post.author_id)}
//                                     className="delete-button"
//                                     title="Delete post"
//                                 >
//                                     Delete
//                                 </button>
//                             )}
//                         </div>
//                         <div className="post-content">
//                             {post.content.split('\n').map((paragraph, i) => (
//                                 <p key={i}>{paragraph}</p>
//                             ))}
//                         </div>
//                         <div className="post-meta">
//                         <span className="author">Author: {post.author_name}</span>
//                             <span className="separator"> | </span>
//                             <span className="date">
//                                 {new Date(post.created_at).toLocaleDateString()},
//                                 {new Date(post.created_at).toLocaleTimeString()}
//                             </span>
//                         </div>
//                     </div>
//                 ))
//             ) : (
//                 <div className="no-posts-message">
//                     <p>No posts found. Create the first one!</p>
//                 </div>
//             )}
//         </div>
//     );
// };

// export default PostList;