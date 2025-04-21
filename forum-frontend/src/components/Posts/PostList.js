import React, { useEffect, useState } from 'react';
import axios from 'axios';
 import '../MainLayout.css';

const PostList = () => {
    const [posts, setPosts] = useState([]);

    useEffect(() => {
        const fetchPosts = async () => {
            try {
                const response = await axios.get('http://localhost:8081/posts');
                setPosts(response.data);
            } catch (error) {
                console.error('Error fetching posts:', error);
            }
        };

        fetchPosts();
    }, []);

    return (
        <div className="post-list-container">
            {posts && posts.length > 0 ? (
                [...posts].reverse().map(post => (
                    <div key={post.id} className="post-item">
                        <h3>{post.title}</h3>
                        <p>{post.content}</p>
                        <small>Created by User ID: {post.author_id}</small>
                    </div>
                ))
            ) : (
                <p>No posts available.</p>
            )}
        </div>
    );
};

export default PostList;
