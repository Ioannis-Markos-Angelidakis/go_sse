import { AppBar, Toolbar, Typography, Button, Box } from '@mui/material';

export default function Header({ user, onLogout }) {
    return (
        <AppBar
            position="static"
            color="transparent"
            elevation={0}
            sx={{
                borderBottom: '1px solid',
                borderColor: 'divider',
                mb: 4,
                py: 1
            }}
        >
            <Toolbar>
                <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
                    Tasks
                </Typography>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                    <Typography variant="body2" color="text.secondary">
                        {user?.email}
                    </Typography>
                    <Button
                        color="primary"
                        onClick={onLogout}
                        sx={{ textTransform: 'none' }}
                    >
                        Sign Out
                    </Button>
                </Box>
            </Toolbar>
        </AppBar>
    );
}