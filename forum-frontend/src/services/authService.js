// authService.js
const API_URL = 'http://localhost:50051'; // Замените на URL вашего auth_service

export const register = async (username, password) => {
  const response = await fetch(`${API_URL}/register`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ username, password }),
  });

  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.error || 'Registration failed');
  }

  return data;
};

export const login = async (username, password) => {
  const response = await fetch(`${API_URL}/login`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ username, password }),
  });

  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.error || 'Login failed');
  }
   localStorage.setItem('token', data.token);
  return data.token;
};

export const validateToken = async (token) => {
    const response = await fetch(`${API_URL}/validate`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({token})
    });

    const data = await response.json();

    if (!response.ok) {
        throw new Error(data.error || 'Validation failed');
    }

    return data;
};