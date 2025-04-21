import React, { useState } from 'react';
import axios from 'axios';
// import 'forum-frontend/src/css/Postst.css';



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
        <CreatePostContainer>
            <h3>Create New Post</h3>
            <Input
                type="text"
                name="title"
                placeholder="Title"
                value={newPost.title} // Исправлено: Добавлено значение
                onChange={handleInputChange}
            />
            <TextArea
                name="content"
                placeholder="Content"
                value={newPost.content}
                onChange={handleInputChange}
            />
            <Button onClick={createPost}>Create Post</Button>
        </CreatePostContainer>
    );
};

export default CreatePost;