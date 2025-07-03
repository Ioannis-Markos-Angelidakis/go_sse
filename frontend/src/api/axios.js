import axios from 'axios';

const instance = axios.create({
    baseURL: process.env.REACT_APP_API_URL || 'http://192.168.0.13:5000',
    withCredentials: true,
});

// Add response interceptor to handle 401 errors
instance.interceptors.response.use(
    response => response,
    error => {
        if (error.response?.status === 401) {
            // Handle unauthorized errors globally
            console.error('Authentication required');
            // Optionally: redirect to login or refresh token
        }
        return Promise.reject(error);
    }
);

export default instance;