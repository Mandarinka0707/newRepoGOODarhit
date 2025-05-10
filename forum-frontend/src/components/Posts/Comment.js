import React, { useState, useEffect, useCallback } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import '/Users/darinautalieva/Desktop/GOProject/forum-frontend/src/components/MainLayout.css';

const Comments = ({ postId }) => {
    const [comments, setComments] = useState([]);
    const [newComment, setNewComment] = useState('');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    const [submitting, setSubmitting] = useState(false);
    
    // Получаем данные пользователя из localStorage
    const token = localStorage.getItem('token');
    const userId = localStorage.getItem('userId');
    const username = localStorage.getItem('username');
    
    // Формируем объект пользователя
    const currentUser = userId ? {
        id: parseInt(userId, 10),
        username: username || `User #${userId}`
    } : null;
    
    const isAuthenticated = !!token && !!currentUser;
    const navigate = useNavigate();

    const fetchComments = useCallback(async () => {
        try {
            setLoading(true);
            setError(null);

            const config = token ? {
                headers: { 
                    'Authorization': `Bearer ${token}`
                }
            } : {};

            const response = await axios.get(
                `http://localhost:8081/api/v1/posts/${postId}/comments`,
                config
            );

            const rawData = response.data;
            let commentsArray = [];
            
            if (Array.isArray(rawData)) {
                commentsArray = rawData;
            } else if (rawData && Array.isArray(rawData.data)) {
                commentsArray = rawData.data;
            } else if (rawData && Array.isArray(rawData.comments)) {
                commentsArray = rawData.comments;
            }

            const processedComments = commentsArray.map(comment => ({
                id: parseInt(comment.id, 10),
                author_id: parseInt(comment.author_id, 10),
                post_id: parseInt(comment.post_id, 10),
                content: comment.content || '',
                author_name: comment.author_name || `User #${comment.author_id}`,
                created_at: comment.created_at || new Date().toISOString()
            }));

            setComments(processedComments);
        } catch (err) {
            setError(err.response?.data?.error || err.message);
        } finally {
            setLoading(false);
        }
    }, [postId, token]);

    const handleSubmitComment = async (e) => {
        e.preventDefault();
        
        if (!isAuthenticated || !currentUser) {
            navigate('/login');
            return;
        }

        if (!newComment.trim() || submitting) return;

        try {
            setSubmitting(true);
            const response = await axios.post(
                `http://localhost:8081/api/v1/posts/${postId}/comments`,
                { 
                    content: newComment,
                    author_id: currentUser.id
                },
                { 
                    headers: { 
                        Authorization: `Bearer ${token}`,
                        'Content-Type': 'application/json' 
                    }
                }
            );

            const newCommentWithAuthor = {
                ...response.data,
                author_name: currentUser.username,
                created_at: new Date().toISOString(),
                content: newComment,
                author_id: currentUser.id,
                post_id: postId
            };

            setComments(prev => [newCommentWithAuthor, ...prev]);
            setNewComment('');
            setError(null);
        } catch (err) {
            setError(err.response?.data?.error || err.message || 'Failed to post comment');
        } finally {
            setSubmitting(false);
        }
    };

    useEffect(() => {
        if (postId) {
            fetchComments();
        }
    }, [postId, fetchComments]);

    // Для отладки
    useEffect(() => {
        console.log('Auth status:', {
            token,
            userId,
            username,
            currentUser,
            isAuthenticated
        });
    }, []);

    return (
        <div className="comments-section">
            <h4>Comments ({comments.length})</h4>
            
            {isAuthenticated ? (
                <form onSubmit={handleSubmitComment} className="comment-form">
                    <textarea
                        value={newComment}
                        onChange={(e) => setNewComment(e.target.value)}
                        placeholder="Write a comment..."
                        rows="3"
                        disabled={submitting}
                        required
                    />
                    <button 
                        type="submit" 
                        className="submit-comment-btn"
                        disabled={submitting || !newComment.trim()}
                    >
                        {submitting ? 'Posting...' : 'Post Comment'}
                    </button>
                </form>
            ) : (
                <div className="login-prompt">
                    <p>Please <a href="/login">login</a> to post comments</p>
                </div>
            )}

            {loading && <div className="loading-indicator">Loading comments...</div>}
            {error && <div className="error-message">Error: {error}</div>}

            <div className="comments-list">
                {comments.map(comment => (
                    <div key={comment.id} className="comment-item">
                        <div className="comment-header">
                            <span className="comment-author">
                                {comment.author_name}
                            </span>
                            <span className="comment-timestamp">
                                {new Date(comment.created_at).toLocaleDateString('en-US', {
                                    hour: '2-digit',
                                    minute: '2-digit',
                                    day: 'numeric',
                                    month: 'short',
                                    year: 'numeric'
                                })}
                            </span>
                        </div>
                        <div className="comment-content">
                            {(comment.content || '').split('\n').map((line, index) => (
                                <p key={index}>{line}</p>
                            ))}
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
};

export default Comments;
// import React, { useState, useEffect, useCallback } from 'react';
// import axios from 'axios';
// import '/Users/darinautalieva/Desktop/GOProject/forum-frontend/src/components/MainLayout.css';

// const Comments = ({ postId }) => {
//     const [comments, setComments] = useState([]);
//     const [newComment, setNewComment] = useState('');
//     const [loading, setLoading] = useState(false);
//     const [error, setError] = useState(null);
//     const [submitting, setSubmitting] = useState(false);
    
//     const token = localStorage.getItem('token');
//     const currentUser = JSON.parse(localStorage.getItem('user')) || {};

//     const fetchComments = useCallback(async () => {
//         try {
//             setLoading(true);
//             setError(null);

//             const response = await axios.get(
//                 `http://localhost:8081/api/v1/posts/${postId}/comments`,
//                 {
//                     headers: { 
//                         'Authorization': `Bearer ${token}`
//                     }
//                 }
//             );

//             // Обработка разных форматов ответа
//             const rawData = response.data;
//             let commentsArray = [];
            
//             if (Array.isArray(rawData)) {
//                 commentsArray = rawData;
//             } else if (rawData && Array.isArray(rawData.data)) {
//                 commentsArray = rawData.data;
//             } else if (rawData && Array.isArray(rawData.comments)) {
//                 commentsArray = rawData.comments;
//             }

//             const processedComments = commentsArray.map(comment => ({
//                 id: parseInt(comment.id, 10),
//                 author_id: parseInt(comment.author_id, 10),
//                 post_id: parseInt(comment.post_id, 10),
//                 content: comment.content || '',
//                 author_name: comment.author_name || `User #${comment.author_id}`,
//                 created_at: comment.created_at || new Date().toISOString()
//             }));

//             setComments(processedComments);
//         } catch (err) {
//             setError(err.response?.data?.error || err.message);
//         } finally {
//             setLoading(false);
//         }
//     }, [postId, token]);

//     const handleSubmitComment = async (e) => {
//         e.preventDefault();
//         if (!newComment.trim() || submitting) return;

//         try {
//             setSubmitting(true);
//             const response = await axios.post(
//                 `http://localhost:8081/api/v1/posts/${postId}/comments`,
//                 { 
//                     content: newComment,
//                     author_id: currentUser.id
//                 },
//                 { 
//                     headers: { 
//                         Authorization: `Bearer ${token}`,
//                         'Content-Type': 'application/json' 
//                     }
//                 }
//             );

//             // Добавление нового комментария с полными данными
//             const newCommentWithAuthor = {
//                 ...response.data,
//                 author_name: currentUser.username || `User #${currentUser.id}`,
//                 created_at: new Date().toISOString(),
//                 content: newComment,
//                 author_id: currentUser.id,
//                 post_id: postId
//             };

//             setComments(prev => [newCommentWithAuthor, ...prev]);
//             setNewComment('');
//         } catch (err) {
//             setError(err.response?.data?.error || err.message);
//         } finally {
//             setSubmitting(false);
//         }
//     };

//     useEffect(() => {
//         fetchComments();
//     }, [fetchComments]);

//     return (
//         <div className="comments-section">
//             <h4>Comments ({comments.length})</h4>
            
//             <form onSubmit={handleSubmitComment} className="comment-form">
//                 <textarea
//                     value={newComment}
//                     onChange={(e) => setNewComment(e.target.value)}
//                     placeholder="Write a comment..."
//                     rows="3"
//                     disabled={submitting}
//                 />
//                 <button 
//                     type="submit" 
//                     className="submit-comment-btn"
//                     disabled={submitting || !newComment.trim()}
//                 >
//                     {submitting ? 'Posting...' : 'Post Comment'}
//                 </button>
//             </form>

//             {loading && <div className="loading-indicator">Loading comments...</div>}
//             {error && <div className="error-message">Error: {error}</div>}

//             <div className="comments-list">
//                 {comments.map(comment => (
//                     <div key={comment.id} className="comment-item">
//                         <div className="comment-header">
//                             <span className="comment-author">
//                                 {comment.author_name}
//                             </span>
//                             <span className="comment-timestamp">
//                                 {new Date(comment.created_at).toLocaleDateString('en-US', {
//                                     hour: '2-digit',
//                                     minute: '2-digit',
//                                     day: 'numeric',
//                                     month: 'short',
//                                     year: 'numeric'
//                                 })}
//                             </span>
//                         </div>
//                         <div className="comment-content">
//                             {(comment.content || '').split('\n').map((line, index) => (
//                                 <p key={index}>{line}</p>
//                             ))}
//                         </div>
//                     </div>
//                 ))}
//             </div>
//         </div>
//     );
// };

// export default Comments;