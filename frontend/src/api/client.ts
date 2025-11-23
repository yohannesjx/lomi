import axios from 'axios';
import { Platform } from 'react-native';
import { storage } from '../utils/storage';
import { useAuthStore } from '../store/authStore';

// API URL configuration
// For Telegram Mini App: Must use HTTPS and publicly accessible URL (use ngrok for local dev)
// Set EXPO_PUBLIC_API_URL environment variable to override
const getApiUrl = () => {
    // Check for environment variable (useful for ngrok URLs)
    if (process.env.EXPO_PUBLIC_API_URL) {
        return process.env.EXPO_PUBLIC_API_URL;
    }
    
    // Development URLs
    if (__DEV__) {
        if (Platform.OS === 'android') {
            return 'http://10.0.2.2:8080/api/v1'; // Android emulator
        }
        return 'http://localhost:8080/api/v1'; // Local development
    }
    
    // Production URL
    return 'https://api.lomi.social/api/v1';
};

const DEV_API_URL = getApiUrl();
const PROD_API_URL = 'https://api.lomi.social/api/v1';

export const api = axios.create({
    baseURL: __DEV__ ? DEV_API_URL : PROD_API_URL,
    headers: {
        'Content-Type': 'application/json',
    },
});

// Add interceptor to inject token
api.interceptors.request.use(
    async (config) => {
        const token = await storage.getItem('lomi_access_token');
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        } else if (__DEV__) {
            // In dev mode, log warning but don't block requests
            // This allows testing UI without authentication
            console.warn('⚠️ No auth token found. API calls will fail with 401.');
            console.warn('💡 To test with auth, log in via WelcomeScreen first.');
        }
        return config;
    },
    (error) => Promise.reject(error)
);

// Add interceptor to handle 401 (Token Expired or Missing)
api.interceptors.response.use(
    (response) => response,
    async (error) => {
        const originalRequest = error.config;
        
        if (error.response?.status === 401 && !originalRequest._retry) {
            originalRequest._retry = true;
            
            // Check if we have a refresh token
            const refreshToken = await storage.getItem('lomi_refresh_token');
            
            if (refreshToken) {
                // Try to refresh token
                try {
                    const response = await axios.post(`${api.defaults.baseURL}/auth/refresh`, {
                        refresh_token: refreshToken,
                    });
                    
                    const { access_token, refresh_token } = response.data;
                    await storage.setItem('lomi_access_token', access_token);
                    await storage.setItem('lomi_refresh_token', refresh_token);
                    
                    // Update store
                    useAuthStore.getState().setTokens({ accessToken: access_token, refreshToken: refresh_token });
                    
                    // Retry original request
                    originalRequest.headers.Authorization = `Bearer ${access_token}`;
                    return api(originalRequest);
                } catch (refreshError) {
                    // Refresh failed, clear tokens and logout
                    console.warn('Token refresh failed, logging out user');
                    await useAuthStore.getState().logout();
                    return Promise.reject(refreshError);
                }
            } else {
                // No refresh token - user needs to log in
                if (__DEV__) {
                    console.warn('⚠️ 401 Unauthorized: No authentication token found.');
                    console.warn('💡 In dev mode, you can test UI without auth, but API calls will fail.');
                    console.warn('💡 To test with real API, log in via WelcomeScreen first.');
                }
                // Don't logout if we're in dev mode and never had a token
                // This allows UI testing without breaking the app
                const hasEverLoggedIn = await storage.getItem('lomi_user');
                if (hasEverLoggedIn) {
                    await useAuthStore.getState().logout();
                }
            }
        }
        
        return Promise.reject(error);
    }
);
