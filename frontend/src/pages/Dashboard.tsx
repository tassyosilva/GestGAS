import { Box, Typography, Container, Grid, Paper, Button } from '@mui/material';
import { Logout as LogoutIcon } from '@mui/icons-material';
import { authService } from '../services/authService.ts';

const Dashboard = () => {
    const user = authService.getUser();

    const handleLogout = () => {
        authService.logout();
    };

    return (
        <Box sx={{ display: 'flex', flexDirection: 'column', minHeight: '100vh' }}>
            {/* Header */}
            <Box sx={{ bgcolor: 'primary.main', color: 'white', py: 2, px: 3, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <Typography variant="h5" component="h1">
                    GestGAS - Dashboard
                </Typography>
                <Box sx={{ display: 'flex', alignItems: 'center' }}>
                    <Typography variant="body1" sx={{ mr: 2 }}>
                        Olá, {user?.nome || 'Usuário'}
                    </Typography>
                    <Button
                        variant="outlined"
                        color="inherit"
                        startIcon={<LogoutIcon />}
                        onClick={handleLogout}
                    >
                        Sair
                    </Button>
                </Box>
            </Box>

            {/* Main content */}
            <Container maxWidth="lg" sx={{ mt: 4, mb: 4, flexGrow: 1 }}>
                <Grid container spacing={3}>
                    {/* Welcome Message */}
                    <Grid item xs={12}>
                        <Paper
                            sx={{
                                p: 3,
                                display: 'flex',
                                flexDirection: 'column',
                                borderRadius: 2
                            }}
                        >
                            <Typography variant="h5" gutterBottom>
                                Bem-vindo ao Sistema de Gerenciamento de Revendedora de Gás e Água
                            </Typography>
                            <Typography variant="body1">
                                Este é o dashboard inicial do sistema. Aqui você terá acesso a todas as funcionalidades
                                de gerenciamento de produtos, pedidos, entregas e controle de estoque.
                            </Typography>
                        </Paper>
                    </Grid>

                    {/* Quick Stats */}
                    <Grid item xs={12} md={4}>
                        <Paper
                            sx={{
                                p: 3,
                                display: 'flex',
                                flexDirection: 'column',
                                height: 200,
                                borderRadius: 2,
                                bgcolor: 'info.light',
                                color: 'white'
                            }}
                        >
                            <Typography component="h2" variant="h6" color="inherit" gutterBottom>
                                Pedidos Pendentes
                            </Typography>
                            <Typography component="p" variant="h3">
                                12
                            </Typography>
                            <Typography variant="body2" sx={{ mt: 1 }}>
                                Pedidos aguardando entrega
                            </Typography>
                        </Paper>
                    </Grid>

                    <Grid item xs={12} md={4}>
                        <Paper
                            sx={{
                                p: 3,
                                display: 'flex',
                                flexDirection: 'column',
                                height: 200,
                                borderRadius: 2,
                                bgcolor: 'success.light',
                                color: 'white'
                            }}
                        >
                            <Typography component="h2" variant="h6" color="inherit" gutterBottom>
                                Vendas Hoje
                            </Typography>
                            <Typography component="p" variant="h3">
                                R$ 2.450,00
                            </Typography>
                            <Typography variant="body2" sx={{ mt: 1 }}>
                                8 vendas realizadas
                            </Typography>
                        </Paper>
                    </Grid>

                    <Grid item xs={12} md={4}>
                        <Paper
                            sx={{
                                p: 3,
                                display: 'flex',
                                flexDirection: 'column',
                                height: 200,
                                borderRadius: 2,
                                bgcolor: 'warning.light',
                                color: 'white'
                            }}
                        >
                            <Typography component="h2" variant="h6" color="inherit" gutterBottom>
                                Estoque Crítico
                            </Typography>
                            <Typography component="p" variant="h3">
                                3
                            </Typography>
                            <Typography variant="body2" sx={{ mt: 1 }}>
                                Produtos abaixo do mínimo
                            </Typography>
                        </Paper>
                    </Grid>

                    {/* Menu de Acesso Rápido */}
                    <Grid item xs={12}>
                        <Paper
                            sx={{
                                p: 3,
                                display: 'flex',
                                flexDirection: 'column',
                                borderRadius: 2,
                                mt: 2
                            }}
                        >
                            <Typography component="h2" variant="h6" gutterBottom>
                                Acesso Rápido
                            </Typography>
                            <Grid container spacing={2} sx={{ mt: 1 }}>
                                <Grid item xs={12} sm={6} md={3}>
                                    <Button variant="contained" color="primary" fullWidth>
                                        Novo Pedido
                                    </Button>
                                </Grid>
                                <Grid item xs={12} sm={6} md={3}>
                                    <Button variant="contained" color="secondary" fullWidth>
                                        Gerenciar Estoque
                                    </Button>
                                </Grid>
                                <Grid item xs={12} sm={6} md={3}>
                                    <Button variant="contained" color="info" fullWidth>
                                        Relatórios
                                    </Button>
                                </Grid>
                                <Grid item xs={12} sm={6} md={3}>
                                    <Button variant="contained" color="warning" fullWidth>
                                        Cadastrar Cliente
                                    </Button>
                                </Grid>
                            </Grid>
                        </Paper>
                    </Grid>
                </Grid>
            </Container>

            {/* Footer */}
            <Box
                component="footer"
                sx={{
                    py: 2,
                    px: 2,
                    mt: 'auto',
                    backgroundColor: (theme) => theme.palette.grey[200]
                }}
            >
                <Typography variant="body2" color="text.secondary" align="center">
                    GestGAS - Sistema de Gerenciamento para Revendedora de Gás e Água © {new Date().getFullYear()}
                </Typography>
            </Box>
        </Box>
    );
};

export default Dashboard;