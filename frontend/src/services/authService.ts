import axios from 'axios';
import type { LoginResponse } from '../types/auth';

const API_URL = 'http://localhost:8080';

export const authService = {
    login: async (email: string, password: string): Promise<LoginResponse> => {
        const response = await axios.post(`${API_URL}/auth/login`, {
            email,
            password,
        });
        return response.data;
    },

    register: async (name: string, email: string, password: string): Promise<LoginResponse> => {
        const response = await axios.post(`${API_URL}/auth/register`, {
            name,
            email,
            password,
        });
        return response.data;
    },

    refreshToken: async (refreshToken: string): Promise<{ access_token: string }> => {
        const response = await axios.post(`${API_URL}/auth/refresh`, {
            refresh_token: refreshToken,
        });
        return response.data;
    },

    logout: async (refreshToken: string): Promise<void> => {
        await axios.post(`${API_URL}/auth/logout`, {
            refresh_token: refreshToken,
        });
    },
};