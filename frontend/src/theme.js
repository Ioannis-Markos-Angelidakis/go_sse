import { createTheme } from '@mui/material/styles';

const theme = createTheme({
    palette: {
        primary: {
            main: '#0a84ff',
            dark: '#0066cc',
        },
        secondary: {
            main: '#86868b',
        },
        background: {
            default: '#f5f5f7',
            paper: '#ffffff',
        },
        text: {
            primary: '#1d1d1f',
            secondary: '#86868b',
        },
        success: {
            main: '#30d158',
        },
        error: {
            main: '#ff453a',
        },
    },
    typography: {
        fontFamily: [
            '-apple-system',
            'BlinkMacSystemFont',
            '"SF Pro"',
            '"Helvetica Neue"',
            'sans-serif',
        ].join(','),
        h1: {
            fontWeight: 600,
            fontSize: '1.8rem',
        },
        h2: {
            fontWeight: 500,
            fontSize: '1.5rem',
            marginBottom: '1rem',
        },
        body1: {
            fontSize: '1rem',
            lineHeight: 1.5,
        },
    },
    shape: {
        borderRadius: 14,
    },
    components: {
        MuiButton: {
            styleOverrides: {
                root: {
                    textTransform: 'none',
                    fontWeight: 500,
                    padding: '12px 24px',
                    boxShadow: '0 4px 10px rgba(10, 132, 255, 0.2)',
                    transition: 'all 0.3s cubic-bezier(0.25, 0.1, 0.25, 1)',
                    '&:hover': {
                        transform: 'translateY(-2px)',
                        boxShadow: '0 6px 15px rgba(10, 132, 255, 0.25)',
                    },
                },
                contained: {
                    boxShadow: '0 4px 10px rgba(10, 132, 255, 0.2)',
                },
            },
        },
        MuiCard: {
            styleOverrides: {
                root: {
                    boxShadow: '0 4px 20px rgba(0, 0, 0, 0.05)',
                    border: '1px solid rgba(0, 0, 0, 0.03)',
                    transition: 'all 0.3s cubic-bezier(0.25, 0.1, 0.25, 1)',
                    '&:hover': {
                        boxShadow: '0 8px 25px rgba(0, 0, 0, 0.08)',
                    },
                },
            },
        },
        MuiTextField: {
            styleOverrides: {
                root: {
                    '& .MuiOutlinedInput-root': {
                        borderRadius: 10,
                        backgroundColor: 'rgba(245, 245, 247, 0.6)',
                    },
                },
            },
        },
    },
});

export default theme;