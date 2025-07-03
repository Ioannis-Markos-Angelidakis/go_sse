import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import { useAuth } from '../context/AuthContext';
import AuthForm from '../ui/AuthForm';
import TextField from '@mui/material/TextField';

export default function RegisterPage() {
    const { login } = useAuth();
    const navigate = useNavigate();
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [loading, setLoading] = useState(false);

    const handleRegister = async () => {
        setLoading(true);
        try {
            const response = await axios.post('/register', { email, password });
            login(response.data.user);
            navigate('/dashboard');
        } catch (error) {
            alert('Registration failed: ' + (error.response?.data?.error || 'Unknown error'));
        } finally {
            setLoading(false);
        }
    };

    return (
        <AuthForm
            title="Create Account"
            buttonText={loading ? 'Creating Account...' : 'Sign Up'}
            footerText="Already have an account?"
            footerLink="/login"
            footerLinkText="Sign in"
            onSubmit={handleRegister}
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