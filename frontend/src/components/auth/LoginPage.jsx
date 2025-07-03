import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from '../../api/axios';
import { useAuth } from '../context/AuthContext';
import AuthForm from '../ui/AuthForm';
import TextField from '@mui/material/TextField';

export default function LoginPage() {
    const { login } = useAuth();
    const navigate = useNavigate();
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [loading, setLoading] = useState(false);

    const handleLogin = async () => {
        setLoading(true);
        try {
            const response = await axios.post('/login', { email, password });
            login(response.data.user);
            navigate('/dashboard');
        } catch (error) {
            alert('Login failed: ' + (error.response?.data?.error || 'Unknown error'));
        } finally {
            setLoading(false);
        }
    };

    return (
        <AuthForm
            title="Sign In"
            buttonText={loading ? 'Signing In...' : 'Sign In'}
            footerText="Don't have an account?"
            footerLink="/register"
            footerLinkText="Sign up"
            onSubmit={handleLogin}
            disabled={loading}
        >
            <TextField
                type="email"
                label="Email Address"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                fullWidth
                sx={{ mb: 2 }}
            />
            <TextField
                type="password"
                label="Password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                fullWidth
                sx={{ mb: 2 }}
            />
        </AuthForm>
    );
}