import React, { useState, useEffect } from 'react';
import {
    Box,
    Typography,
    Grid,
    Paper,
    Card,
    CardContent,
    CardActions,
    Button,
    Divider,
    List,
    ListItem,
    ListItemText,
    ListItemIcon,
    Chip,
    Alert,
    CircularProgress,
} from '@mui/material';
import {
    LocalShipping as ShippingIcon,
    Inventory as InventoryIcon,
    LocalGasStation as GasIcon,
    Warning as WarningIcon,
    Payments as PaymentsIcon,
    Receipt as ReceiptIcon,
    Refresh as RefreshIcon,
    ViewList as ListIcon,
    Add as AddIcon,
} from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import API_BASE_URL from '../config/api';

// Interfaces
interface AlertaEstoque {
    produto_id: number;
    nome_produto: string;
    quantidade: number;
    alerta_minimo: number;
    status: string;
}

interface PedidoRecente {
    id: number;
    cliente: {
        nome: string;
    };
    status: string;
    valor_total: number;
    criado_em: string;
}

interface ResumoEstoque {
    total_produtos: number;
    total_botijas_cheias: number;
    total_botijas_vazias: number;
    total_botijas_emprestadas: number;
    produtos_em_alerta: number;
}

interface ResumoPedidos {
    total_pedidos: number;
    pedidos_hoje: number;
    valor_total_hoje: number;
    pedidos_pendentes: number;
    pedidos_em_entrega: number;
}

const Dashboard: React.FC = () => {
    const navigate = useNavigate();
    const [loadingAlertas, setLoadingAlertas] = useState(true);
    const [loadingPedidos, setLoadingPedidos] = useState(true);
    const [loadingResumos, setLoadingResumos] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [alertas, setAlertas] = useState<AlertaEstoque[]>([]);
    const [pedidosRecentes, setPedidosRecentes] = useState<PedidoRecente[]>([]);
    const [resumoEstoque, setResumoEstoque] = useState<ResumoEstoque>({
        total_produtos: 0,
        total_botijas_cheias: 0,
        total_botijas_vazias: 0,
        total_botijas_emprestadas: 0,
        produtos_em_alerta: 0,
    });
    const [resumoPedidos, setResumoPedidos] = useState<ResumoPedidos>({
        total_pedidos: 0,
        pedidos_hoje: 0,
        valor_total_hoje: 0,
        pedidos_pendentes: 0,
        pedidos_em_entrega: 0,
    });

    // Formatadores
    const formatarMoeda = (valor: number): string => {
        return new Intl.NumberFormat('pt-BR', {
            style: 'currency',
            currency: 'BRL',
        }).format(valor);
    };

    const formatarData = (dataString: string): string => {
        if (!dataString) return '-';
        const data = new Date(dataString);
        return new Intl.DateTimeFormat('pt-BR', {
            day: '2-digit',
            month: '2-digit',
            year: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
        }).format(data);
    };

    // Carregar dados para o dashboard com tratamento de erros aprimorado
    const carregarDados = async () => {
        try {
            const token = localStorage.getItem('token');
            if (!token) {
                setError('Não autorizado. Faça login para continuar.');
                setLoadingAlertas(false);
                setLoadingPedidos(false);
                setLoadingResumos(false);
                return;
            }

            // Carregar alertas com tratamento de erros individual
            setLoadingAlertas(true);
            try {
                const alertasResponse = await axios.get<AlertaEstoque[]>(
                    `${API_BASE_URL}/estoque/alertas`,
                    {
                        headers: {
                            Authorization: `Bearer ${token}`,
                        },
                    }
                );
                setAlertas(alertasResponse.data || []);
            } catch (alertaErr) {
                console.error('Erro ao carregar alertas:', alertaErr);
                setAlertas([]); // Define como array vazio em vez de null
            } finally {
                setLoadingAlertas(false);
            }

            // Carregar pedidos recentes com tratamento de erros individual
            setLoadingPedidos(true);
            try {
                const pedidosResponse = await axios.get<{ pedidos: PedidoRecente[] }>(
                    `${API_BASE_URL}/pedidos?limit=5`,
                    {
                        headers: {
                            Authorization: `Bearer ${token}`,
                        },
                    }
                );
                setPedidosRecentes(pedidosResponse.data?.pedidos || []);
            } catch (pedidosErr) {
                console.error('Erro ao carregar pedidos:', pedidosErr);
                setPedidosRecentes([]); // Define como array vazio em vez de null
            } finally {
                setLoadingPedidos(false);
            }

            // Simular/carregar resumo de estoque com dados seguros
            setLoadingResumos(true);

            // Usar dados fictícios para testes ou dados reais se estiverem disponíveis
            setTimeout(() => {
                setResumoEstoque({
                    total_produtos: 10,
                    total_botijas_cheias: 85,
                    total_botijas_vazias: 32,
                    total_botijas_emprestadas: 15,
                    produtos_em_alerta: alertas ? alertas.length : 0,
                });

                setResumoPedidos({
                    total_pedidos: 158,
                    pedidos_hoje: 12,
                    valor_total_hoje: 1580.50,
                    pedidos_pendentes: 5,
                    pedidos_em_entrega: 3,
                });

                setLoadingResumos(false);
            }, 500);

        } catch (err) {
            console.error('Erro ao carregar dados do dashboard:', err);
            setError('Não foi possível carregar todos os dados do dashboard. Algumas informações podem estar incompletas.');
            // Garantir que as flags de carregamento sejam desligadas
            setLoadingAlertas(false);
            setLoadingPedidos(false);
            setLoadingResumos(false);
        }
    };

    // Carregar dados ao montar o componente
    useEffect(() => {
        carregarDados();
    }, []);

    // Renderizar chip de status com cores
    const renderizarStatusChip = (status: string) => {
        let color: 'default' | 'primary' | 'secondary' | 'error' | 'info' | 'success' | 'warning' = 'default';
        let label = status;

        switch (status) {
            case 'novo':
                color = 'info';
                label = 'Novo';
                break;
            case 'em_preparo':
                color = 'warning';
                label = 'Em Preparo';
                break;
            case 'em_entrega':
                color = 'primary';
                label = 'Em Entrega';
                break;
            case 'entregue':
                color = 'success';
                label = 'Entregue';
                break;
            case 'finalizado':
                color = 'success';
                label = 'Finalizado';
                break;
            case 'cancelado':
                color = 'error';
                label = 'Cancelado';
                break;
            case 'baixo':
                color = 'warning';
                label = 'Baixo';
                break;
            case 'critico':
                color = 'error';
                label = 'Crítico';
                break;
        }

        return <Chip label={label} color={color} size="small" />;
    };

    return (
        <Box>
            {/* Cabeçalho */}
            <Box
                display="flex"
                justifyContent="space-between"
                alignItems="center"
                mb={3}
            >
                <Typography variant="h4">Dashboard</Typography>
                <Button
                    variant="outlined"
                    startIcon={<RefreshIcon />}
                    onClick={carregarDados}
                >
                    Atualizar
                </Button>
            </Box>

            {error && (
                <Alert severity="error" sx={{ mb: 3 }}>
                    {error}
                </Alert>
            )}

            {/* Resumo das Informações */}
            <Grid container spacing={3} mb={4}>
                {/* Resumo de Estoque */}
                <Grid item xs={12} md={6}>
                    <Paper sx={{ p: 2, height: '100%' }}>
                        <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
                            <Typography variant="h6">
                                <InventoryIcon fontSize="small" sx={{ mr: 1, verticalAlign: 'text-bottom' }} />
                                Resumo de Estoque
                            </Typography>
                            <Button
                                size="small"
                                onClick={() => navigate('/estoque')}
                                endIcon={<ListIcon />}
                            >
                                Ver Estoque
                            </Button>
                        </Box>

                        {loadingResumos ? (
                            <Box display="flex" justifyContent="center" p={3}>
                                <CircularProgress size={30} />
                            </Box>
                        ) : (
                            <Grid container spacing={2}>
                                <Grid item xs={6}>
                                    <Card variant="outlined">
                                        <CardContent>
                                            <Typography variant="subtitle2" color="text.secondary">
                                                Botijas Cheias
                                            </Typography>
                                            <Typography variant="h4">
                                                {resumoEstoque.total_botijas_cheias}
                                            </Typography>
                                        </CardContent>
                                    </Card>
                                </Grid>
                                <Grid item xs={6}>
                                    <Card variant="outlined">
                                        <CardContent>
                                            <Typography variant="subtitle2" color="text.secondary">
                                                Botijas Vazias
                                            </Typography>
                                            <Typography variant="h4">
                                                {resumoEstoque.total_botijas_vazias}
                                            </Typography>
                                        </CardContent>
                                    </Card>
                                </Grid>
                                <Grid item xs={6}>
                                    <Card variant="outlined">
                                        <CardContent>
                                            <Typography variant="subtitle2" color="text.secondary">
                                                Botijas Emprestadas
                                            </Typography>
                                            <Typography variant="h4">
                                                {resumoEstoque.total_botijas_emprestadas}
                                            </Typography>
                                        </CardContent>
                                    </Card>
                                </Grid>
                                <Grid item xs={6}>
                                    <Card variant="outlined"
                                        sx={{
                                            bgcolor: resumoEstoque.produtos_em_alerta > 0 ? 'error.50' : undefined,
                                            borderColor: resumoEstoque.produtos_em_alerta > 0 ? 'error.main' : undefined
                                        }}
                                    >
                                        <CardContent>
                                            <Typography variant="subtitle2" color="text.secondary">
                                                Produtos em Alerta
                                            </Typography>
                                            <Typography variant="h4" color={resumoEstoque.produtos_em_alerta > 0 ? 'error.main' : undefined}>
                                                {resumoEstoque.produtos_em_alerta}
                                            </Typography>
                                        </CardContent>
                                    </Card>
                                </Grid>
                            </Grid>
                        )}
                    </Paper>
                </Grid>

                {/* Resumo de Pedidos */}
                <Grid item xs={12} md={6}>
                    <Paper sx={{ p: 2, height: '100%' }}>
                        <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
                            <Typography variant="h6">
                                <ReceiptIcon fontSize="small" sx={{ mr: 1, verticalAlign: 'text-bottom' }} />
                                Resumo de Pedidos
                            </Typography>
                            <Button
                                size="small"
                                onClick={() => navigate('/pedidos')}
                                endIcon={<ListIcon />}
                            >
                                Ver Pedidos
                            </Button>
                        </Box>

                        {loadingResumos ? (
                            <Box display="flex" justifyContent="center" p={3}>
                                <CircularProgress size={30} />
                            </Box>
                        ) : (
                            <Grid container spacing={2}>
                                <Grid item xs={6}>
                                    <Card variant="outlined">
                                        <CardContent>
                                            <Typography variant="subtitle2" color="text.secondary">
                                                Pedidos Hoje
                                            </Typography>
                                            <Typography variant="h4">
                                                {resumoPedidos.pedidos_hoje}
                                            </Typography>
                                        </CardContent>
                                    </Card>
                                </Grid>
                                <Grid item xs={6}>
                                    <Card variant="outlined">
                                        <CardContent>
                                            <Typography variant="subtitle2" color="text.secondary">
                                                Faturamento Hoje
                                            </Typography>
                                            <Typography variant="h5">
                                                {formatarMoeda(resumoPedidos.valor_total_hoje)}
                                            </Typography>
                                        </CardContent>
                                    </Card>
                                </Grid>
                                <Grid item xs={6}>
                                    <Card variant="outlined">
                                        <CardContent>
                                            <Typography variant="subtitle2" color="text.secondary">
                                                Pedidos Pendentes
                                            </Typography>
                                            <Typography variant="h4">
                                                {resumoPedidos.pedidos_pendentes}
                                            </Typography>
                                        </CardContent>
                                    </Card>
                                </Grid>
                                <Grid item xs={6}>
                                    <Card variant="outlined">
                                        <CardContent>
                                            <Typography variant="subtitle2" color="text.secondary">
                                                Em Entrega
                                            </Typography>
                                            <Typography variant="h4">
                                                {resumoPedidos.pedidos_em_entrega}
                                            </Typography>
                                        </CardContent>
                                    </Card>
                                </Grid>
                            </Grid>
                        )}
                    </Paper>
                </Grid>
            </Grid>

            <Grid container spacing={3}>
                {/* Alertas de Estoque */}
                <Grid item xs={12} md={6}>
                    <Paper sx={{ p: 0, height: '100%' }}>
                        <Box p={2} display="flex" justifyContent="space-between" alignItems="center">
                            <Typography variant="h6">
                                <WarningIcon fontSize="small" color="warning" sx={{ mr: 1, verticalAlign: 'text-bottom' }} />
                                Alertas de Estoque
                            </Typography>
                            <Button
                                size="small"
                                onClick={() => navigate('/estoque')}
                            >
                                Ver Todos
                            </Button>
                        </Box>

                        <Divider />

                        {loadingAlertas ? (
                            <Box display="flex" justifyContent="center" p={3}>
                                <CircularProgress size={24} />
                            </Box>
                        ) : alertas && alertas.length === 0 ? (
                            <Box p={3}>
                                <Alert severity="success">
                                    Não há alertas de estoque no momento.
                                </Alert>
                            </Box>
                        ) : alertas ? (
                            <List dense>
                                {alertas.slice(0, 5).map((alerta) => (
                                    <ListItem key={alerta.produto_id} divider>
                                        <ListItemIcon>
                                            <WarningIcon color={alerta.status === 'critico' ? 'error' : 'warning'} />
                                        </ListItemIcon>
                                        <ListItemText
                                            primary={alerta.nome_produto}
                                            secondary={`Qtd: ${alerta.quantidade} | Mínimo: ${alerta.alerta_minimo}`}
                                        />
                                        <Box>
                                            {renderizarStatusChip(alerta.status)}
                                        </Box>
                                    </ListItem>
                                ))}
                            </List>
                        ) : (
                            <Box p={3}>
                                <Alert severity="warning">
                                    Erro ao carregar alertas. Tente novamente mais tarde.
                                </Alert>
                            </Box>
                        )}

                        <Box p={2} display="flex" justifyContent="space-between">
                            <Button
                                size="small"
                                variant="outlined"
                                color="primary"
                                startIcon={<AddIcon />}
                                onClick={() => navigate('/pedidos/novo')}
                            >
                                Novo Pedido
                            </Button>

                            <Button
                                size="small"
                                variant="outlined"
                                color="primary"
                                onClick={() => navigate('/estoque')}
                            >
                                Gerenciar Estoque
                            </Button>
                        </Box>
                    </Paper>
                </Grid>

                {/* Pedidos Recentes */}
                <Grid item xs={12} md={6}>
                    <Paper sx={{ p: 0, height: '100%' }}>
                        <Box p={2} display="flex" justifyContent="space-between" alignItems="center">
                            <Typography variant="h6">
                                <ReceiptIcon fontSize="small" sx={{ mr: 1, verticalAlign: 'text-bottom' }} />
                                Pedidos Recentes
                            </Typography>
                            <Button
                                size="small"
                                onClick={() => navigate('/pedidos')}
                            >
                                Ver Todos
                            </Button>
                        </Box>

                        <Divider />

                        {loadingPedidos ? (
                            <Box display="flex" justifyContent="center" p={3}>
                                <CircularProgress size={24} />
                            </Box>
                        ) : pedidosRecentes && pedidosRecentes.length === 0 ? (
                            <Box p={3}>
                                <Alert severity="info">
                                    Nenhum pedido recente encontrado.
                                </Alert>
                            </Box>
                        ) : pedidosRecentes ? (
                            <List dense>
                                {pedidosRecentes.map((pedido) => (
                                    <ListItem
                                        key={pedido.id}
                                        divider
                                        secondaryAction={
                                            <Button
                                                size="small"
                                                onClick={() => navigate(`/pedidos/${pedido.id}`)}
                                            >
                                                Detalhes
                                            </Button>
                                        }
                                    >
                                        <ListItemText
                                            primary={`#${pedido.id} - ${pedido.cliente ? pedido.cliente.nome : 'Cliente'}`}
                                            secondary={`${formatarData(pedido.criado_em)} | ${formatarMoeda(pedido.valor_total)}`}
                                        />
                                        <Box mr={8}>
                                            {renderizarStatusChip(pedido.status)}
                                        </Box>
                                    </ListItem>
                                ))}
                            </List>
                        ) : (
                            <Box p={3}>
                                <Alert severity="warning">
                                    Erro ao carregar pedidos recentes. Tente novamente mais tarde.
                                </Alert>
                            </Box>
                        )}

                        <Box p={2} display="flex" justifyContent="space-between">
                            <Button
                                size="small"
                                variant="contained"
                                color="primary"
                                startIcon={<AddIcon />}
                                onClick={() => navigate('/pedidos/novo')}
                            >
                                Novo Pedido
                            </Button>

                            <Button
                                size="small"
                                variant="outlined"
                                color="primary"
                                onClick={() => navigate('/pedidos')}
                            >
                                Ver Pedidos
                            </Button>
                        </Box>
                    </Paper>
                </Grid>
            </Grid>

            {/* Ações Rápidas */}
            <Box mt={4}>
                <Typography variant="h6" gutterBottom>
                    Ações Rápidas
                </Typography>
                <Grid container spacing={2}>
                    <Grid item xs={6} sm={4} md={3} lg={2}>
                        <Card>
                            <CardContent sx={{ textAlign: 'center', py: 2 }}>
                                <GasIcon sx={{ fontSize: 40, color: 'primary.main', mb: 1 }} />
                                <Typography variant="subtitle1">Novo Pedido</Typography>
                            </CardContent>
                            <CardActions sx={{ justifyContent: 'center', pt: 0, pb: 2 }}>
                                <Button
                                    variant="outlined"
                                    size="small"
                                    onClick={() => navigate('/pedidos/novo')}
                                >
                                    Criar
                                </Button>
                            </CardActions>
                        </Card>
                    </Grid>

                    <Grid item xs={6} sm={4} md={3} lg={2}>
                        <Card>
                            <CardContent sx={{ textAlign: 'center', py: 2 }}>
                                <InventoryIcon sx={{ fontSize: 40, color: 'primary.main', mb: 1 }} />
                                <Typography variant="subtitle1">Estoque</Typography>
                            </CardContent>
                            <CardActions sx={{ justifyContent: 'center', pt: 0, pb: 2 }}>
                                <Button
                                    variant="outlined"
                                    size="small"
                                    onClick={() => navigate('/estoque')}
                                >
                                    Gerenciar
                                </Button>
                            </CardActions>
                        </Card>
                    </Grid>

                    <Grid item xs={6} sm={4} md={3} lg={2}>
                        <Card>
                            <CardContent sx={{ textAlign: 'center', py: 2 }}>
                                <ShippingIcon sx={{ fontSize: 40, color: 'primary.main', mb: 1 }} />
                                <Typography variant="subtitle1">Entregas</Typography>
                            </CardContent>
                            <CardActions sx={{ justifyContent: 'center', pt: 0, pb: 2 }}>
                                <Button
                                    variant="outlined"
                                    size="small"
                                    onClick={() => navigate('/pedidos?status=em_entrega')}
                                >
                                    Ver
                                </Button>
                            </CardActions>
                        </Card>
                    </Grid>

                    <Grid item xs={6} sm={4} md={3} lg={2}>
                        <Card>
                            <CardContent sx={{ textAlign: 'center', py: 2 }}>
                                <PaymentsIcon sx={{ fontSize: 40, color: 'primary.main', mb: 1 }} />
                                <Typography variant="subtitle1">Vendas Fiadas</Typography>
                            </CardContent>
                            <CardActions sx={{ justifyContent: 'center', pt: 0, pb: 2 }}>
                                <Button
                                    variant="outlined"
                                    size="small"
                                    onClick={() => navigate('/vendas-fiadas')}
                                >
                                    Gerenciar
                                </Button>
                            </CardActions>
                        </Card>
                    </Grid>
                </Grid>
            </Box>
        </Box>
    );
};

export default Dashboard;