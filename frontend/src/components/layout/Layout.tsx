import { useState } from 'react';
import { Outlet, useNavigate } from 'react-router-dom';
import {
    Box,
    AppBar,
    Toolbar,
    Typography,
    IconButton,
    Button,
    useMediaQuery,
    useTheme,
} from '@mui/material';
import {
    Menu as MenuIcon,
    Logout as LogoutIcon,
} from '@mui/icons-material';
import Sidebar from './Sidebar';
import { authService } from '../../services/authService';

// Largura do drawer quando aberto
const drawerWidth = 240;

const Layout = () => {
    // No mobile, começamos com a sidebar fechada
    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('md'));
    const [open, setOpen] = useState(!isMobile);

    const navigate = useNavigate();

    // Obter dados do usuário logado
    const user = authService.getUser();

    // Alternar o estado da barra lateral (apenas para mobile)
    const toggleDrawer = () => {
        if (isMobile) {
            setOpen(!open);
        }
    };

    // Lidar com o logout
    const handleLogout = () => {
        authService.logout();
        navigate('/login');
    };

    return (
        <Box sx={{ display: 'flex', height: '100vh' }}>
            {/* Barra superior */}
            <AppBar
                position="fixed"
                sx={{
                    width: { md: `calc(100% - ${drawerWidth}px)` },
                    ml: { md: `${drawerWidth}px` },
                    zIndex: (theme) => theme.zIndex.drawer + 1,
                }}
            >
                <Toolbar>
                    {isMobile && (
                        <IconButton
                            color="inherit"
                            aria-label="open drawer"
                            onClick={toggleDrawer}
                            edge="start"
                            sx={{ mr: 2 }}
                        >
                            <MenuIcon />
                        </IconButton>
                    )}

                    <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
                        Sistema de Gerenciamento para Revendedora de Gás e Água
                    </Typography>

                    <Box sx={{ display: 'flex', alignItems: 'center' }}>
                        <Typography variant="body1" sx={{ mr: 2, display: { xs: 'none', sm: 'block' } }}>
                            Olá, {user?.nome || 'Usuário'}
                        </Typography>

                        <Button
                            color="inherit"
                            startIcon={<LogoutIcon />}
                            onClick={handleLogout}
                        >
                            Sair
                        </Button>
                    </Box>
                </Toolbar>
            </AppBar>

            {/* Barra lateral */}
            <Sidebar open={open} onToggle={toggleDrawer} />

            {/* Conteúdo principal */}
            <Box
                component="main"
                sx={{
                    flexGrow: 1,
                    p: 3,
                    width: { md: `calc(100% - ${drawerWidth}px)` },
                    ml: { md: `${drawerWidth}px` },
                    mt: '64px', // Altura da barra superior
                    overflow: 'auto',
                }}
            >
                <Outlet />
            </Box>
        </Box>
    );
};

export default Layout;