import React, { useState } from 'react';
import axios from 'axios';
import '../MainLayout.css';

const CreatePost = ({ onPostCreated }) => {
    const [newPost, setNewPost] = useState({ 
        title: '', 
        content: '' 
    });
    const [error, setError] = useState('');

    const handleInputChange = (e) => {
        const { name, value } = e.target;
        setNewPost(prev => ({
            ...prev,
            [name]: value
        }));
    };

    const createPost = async () => {
        try {
            setError('');
            const token = localStorage.getItem('token');
            
            if (!token) {
                setError('You need to login first');
                return;
            }

            if (!newPost.title.trim() || !newPost.content.trim()) {
                setError('Title and content are required');
                return;
            }

            const response = await axios.post('http://localhost:8081/api/v1/posts', {
                title: newPost.title,
                content: newPost.content
            }, {
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                }
            });

            setNewPost({ title: '', content: '' });
            onPostCreated();
        } catch (error) {
            console.error('Full error:', error);
            const errorMessage = error.response?.data?.error || 
                               error.response?.data?.message || 
                               'Error creating post. Please try again.';
            setError(errorMessage);
        }
    };

    return (
        <div className="create-post-container">
            <h3>Create New Post</h3>
            {error && <div className="error-message">{error}</div>}
            <input
                type="text"
                name="title"
                placeholder="Title"
                value={newPost.title}
                onChange={handleInputChange}
                required
            />
            <textarea
                name="content"
                placeholder="Content"
                value={newPost.content}
                onChange={handleInputChange}
                required
            />
            <button onClick={createPost}>Create Post</button>
        </div>
    );
};

export default CreatePost;