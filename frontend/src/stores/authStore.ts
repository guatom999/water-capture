import { create } from "zustand";
import { persist } from "zustand/middleware";
import type { User , LoginResponse} from "../types/auth";

interface AuthStore {
    // State
    user: User | null;
    accessToken: string | null;
    refreshToken: string | null;
    isAuthenticated: boolean;
    
    // Actions
    login: (response: LoginResponse) => void;
    logout: () => void;
    setAccessToken: (token: string) => void;
}
export const useAuthStore = create<AuthStore>()(
    persist(
        (set) => ({
            // Initial state
            user: null,
            accessToken: null,
            refreshToken: null,
            isAuthenticated: false,
            // Login - เก็บ token และ user info
            login: (response: LoginResponse) => set({
                user: response.user,
                accessToken: response.access_token,
                refreshToken: response.refresh_token,
                isAuthenticated: true,
            }),
            // Logout - clear ทุกอย่าง
            logout: () => set({
                user: null,
                accessToken: null,
                refreshToken: null,
                isAuthenticated: false,
            }),
            // Update access token (หลัง refresh)
            setAccessToken: (token: string) => set({
                accessToken: token,
            }),
        }),
        {
            name: 'auth-storage', // key ใน localStorage
            partialize: (state) => ({
                // เลือกเก็บเฉพาะ refreshToken + user (ไม่เก็บ accessToken)
                refreshToken: state.refreshToken,
                user: state.user,
                isAuthenticated: state.isAuthenticated,
            }),
        }
    )
);