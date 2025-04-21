import React, { useEffect, useState } from 'react';
import axios from 'axios';
// import 'forum-frontend/src/css/Postst.css';


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
        <PostListContainer>
            {posts && posts.length > 0 ? (
                [...posts].reverse().map(post => (
                    <PostItem key={post.id}>
                        <h3>{post.title}</h3>
                        <p>{post.content}</p>
                        <small>Created by User ID: {post.author_id}</small>
                    </PostItem>
                ))
            ) : (
                <p>No posts available.</p>
            )}
        </PostListContainer>
    );
};

export default PostList;