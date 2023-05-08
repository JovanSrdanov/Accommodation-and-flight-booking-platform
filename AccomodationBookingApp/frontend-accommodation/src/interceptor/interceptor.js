import axios from 'axios';

const API_URL = 'http://localhost:8000/';

const interceptor = axios.create({
    baseURL: API_URL,
});

interceptor.interceptors.request.use(
    config => {
        const paseto = localStorage.getItem('paseto');
        if (paseto) {
            config.headers['Authorization'] = `Bearer ${paseto}`;
        }
        return config;
    },
    error => {
        Promise.reject(error)
    }
);


export default interceptor;