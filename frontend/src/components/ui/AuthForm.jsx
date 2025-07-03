import { Button, Card, Container, Typography } from '@mui/material';

export default function AuthForm({
    title,
    buttonText,
    footerText,
    footerLink,
    footerLinkText,
    children,
    onSubmit,
    disabled
}) {
    return (
        <Container
            sx={{
                display: 'flex',
                justifyContent: 'center',
                alignItems: 'center',
                minHeight: '100vh',
                p: 2,
                background: 'linear-gradient(135deg, #f5f7fa 0%, #e4ecfb 100%)',
            }}
        >
            <Card
                sx={{
                    p: 4,
                    width: '100%',
                    maxWidth: 420,
                    textAlign: 'center'
                }}
            >
                <Typography variant="h1" sx={{ mb: 3 }}>
                    {title}
                </Typography>
                {children}
                <Button
                    variant="contained"
                    color="primary"
                    onClick={onSubmit}
                    disabled={disabled}
                    fullWidth
                    sx={{ mt: 1 }}
                >
                    {buttonText}
                </Button>
                <Typography variant="body1" sx={{ mt: 2, color: 'text.secondary' }}>
                    {footerText} <a href={footerLink}>{footerLinkText}</a>
                </Typography>
            </Card>
        </Container>
    );
}