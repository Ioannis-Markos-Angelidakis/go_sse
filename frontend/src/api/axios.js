import axios from 'axios';

export const baseURL = `http://${process.env.REACT_APP_HOST}:${process.env.REACT_APP_BACKEND_PORT}`;

const instance = axios.create({
    baseURL,
    withCredentials: true,
});

instance.interceptors.response.use(
    response => response,
    error => {
        if (error.response) {
            console.error(
                `Request failed with status ${error.response.status}:`,
                error.response.data
            );
        } else if (error.request) {
            console.error('No response received:', error.request);
        } else {
            console.error('Request setup error:', error.message);
        }

        if (error.response?.status === 401) {
            console.error('Authentication required');
        }

        return Promise.reject(error);
    }
);

export default instance;