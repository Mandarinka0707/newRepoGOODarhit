import React, { useState } from 'react';
import axios from 'axios';
 import '../MainLayout.css';

const CreatePost = ({ onPostCreated }) => {
    const [newPost, setNewPost] = useState({ title: '', content: '' });
    const token = localStorage.getItem('token');

    const handleInputChange = (e) => {
        setNewPost({ ...newPost, [e.target.name]: e.target.value });
    };

    const createPost = async () => {
        try {
            const config = {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            };
            await axios.post('http://localhost:8081/posts', newPost, config);
            setNewPost({ title: '', content: '' });
            onPostCreated(); // Refresh the list
        } catch (error) {
            console.error('Error creating post:', error);
            alert('Error creating post. Please try again.');
        }
    };

    return (
        <div className="create-post-container">
            <h3>Create New Post</h3>
            <input
                type="text"
                name="title"
                placeholder="Title"
                value={newPost.title}
                onChange={handleInputChange}
            />
            <textarea
                name="content"
                placeholder="Content"
                value={newPost.content}
                onChange={handleInputChange}
            />
            <button onClick={createPost}>Create Post</button>
        </div>
    );
};

export default CreatePost;
