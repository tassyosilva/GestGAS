import { useState } from 'react';
import { Outlet, useNavigate } from 'react-router-dom';
import {
    AppBar,
    Box,
    CssBaseline,
    Drawer,
    IconButton,
    Toolbar,
    Typography,
    Button,
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
    const [mobileOpen, setMobileOpen] = useState(false);
    const navigate = useNavigate();

    // Obter dados do usuário logado
    const user = authService.getUser();

    const handleDrawerToggle = () => {
        setMobileOpen(!mobileOpen);
    };

    // Lidar com o logout
    const handleLogout = () => {
        authService.logout();
        navigate('/login');
    };

    return (
        <Box sx={{ display: 'flex' }}>
            <CssBaseline />

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
                    <IconButton
                        color="inherit"
                        aria-label="open drawer"
                        edge="start"
                        onClick={handleDrawerToggle}
                        sx={{ mr: 2, display: { md: 'none' } }}
                    >
                        <MenuIcon />
                    </IconButton>

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

            {/* Menu lateral */}
            <Box
                component="nav"
                sx={{ width: { md: drawerWidth }, flexShrink: { md: 0 } }}
            >
                {/* Mobile drawer */}
                <Drawer
                    variant="temporary"
                    open={mobileOpen}
                    onClose={handleDrawerToggle}
                    ModalProps={{
                        keepMounted: true, // Melhor desempenho em dispositivos móveis
                    }}
                    sx={{
                        display: { xs: 'block', md: 'none' },
                        '& .MuiDrawer-paper': { boxSizing: 'border-box', width: drawerWidth },
                    }}
                >
                    <Sidebar open={true} onToggle={handleDrawerToggle} />
                </Drawer>

                {/* Desktop drawer */}
                <Drawer
                    variant="permanent"
                    sx={{
                        display: { xs: 'none', md: 'block' },
                        '& .MuiDrawer-paper': {
                            boxSizing: 'border-box',
                            width: drawerWidth,
                            borderRight: '1px solid rgba(0, 0, 0, 0.12)',
                        },
                    }}
                    open
                >
                    <Sidebar open={true} onToggle={handleDrawerToggle} />
                </Drawer>
            </Box>

            {/* Conteúdo principal */}
            <Box
                component="main"
                sx={{
                    flexGrow: 1,
                    p: 3,
                    width: { md: `calc(100% - ${drawerWidth}px)` },
                    backgroundColor: 'background.default'
                }}
            >
                <Toolbar /> {/* Espaço para a barra de ferramentas */}
                <Outlet />
            </Box>
        </Box>
    );
};

export default Layout;