import { Box, Typography, Grid, Paper } from '@mui/material';

const Dashboard = () => {
    return (
        <Box>
            <Typography variant="h4" gutterBottom>
                Dashboard
            </Typography>

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
            </Grid>
        </Box>
    );
};

export default Dashboard;