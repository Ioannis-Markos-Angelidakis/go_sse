import { createContext, useState, useEffect, useContext } from 'react';
import axios from '../../api/axios';

export const AuthContext = createContext();

export function useAuth() {
    const context = useContext(AuthContext);

    if (context === undefined) {
        throw new Error('useAuth must be used within an AuthProvider');
    }

    return context;
}

export function AuthProvider({ children }) {
    const [user, setUser] = useState(() => {
        const savedUser = sessionStorage.getItem('user');
        return savedUser ? JSON.parse(savedUser) : null;
    });

    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchUser = async () => {
            try {
                const response = await axios.get('http://192.168.0.13:5000/me');
                const userData = response.data.user;
                setUser(userData);
                sessionStorage.setItem('user', JSON.stringify(userData));
            } catch (error) {
                console.error('Failed to fetch user:', error);
                sessionStorage.removeItem('user');
            } finally {
                setLoading(false);
            }
        };

        if (!user) fetchUser();
        else setLoading(false);
    }, []);

    const login = (userData) => {
        setUser(userData);
        sessionStorage.setItem('user', JSON.stringify(userData));
    };

    const logout = () => {
        axios.post('/logout')
            .finally(() => {
                setUser(null);
                sessionStorage.removeItem('user');
            });
    };

    const value = {
        user,
        login,
        logout,
        loading
    };

    return (
        <AuthContext.Provider value={value}>
            {children}
        </AuthContext.Provider>
    );
}