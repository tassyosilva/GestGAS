import React, { useState, useEffect } from 'react';
import {
    Box,
    Typography,
    Paper,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    TablePagination,
    Chip,
    Button,
    TextField,
    MenuItem,
    Grid,
    CircularProgress,
    IconButton,
} from '@mui/material';
import {
    Add as AddIcon,
    Refresh as RefreshIcon,
    Search as SearchIcon,
    FilterList as FilterListIcon,
} from '@mui/icons-material';
import API_BASE_URL from '../config/api';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';

// Tipos para o componente
interface Cliente {
    id: number;
    nome: string;
    telefone: string;
}

interface Usuario {
    id: number;
    nome: string;
    perfil: string;
}

interface Pedido {
    id: number;
    cliente: Cliente | number | string; // Cliente pode ser objeto, ID ou string
    cliente_id?: number;
    cliente_nome?: string;
    atendente: Usuario;
    entregador?: Usuario;
    status: string;
    forma_pagamento: string;
    valor_total: number;
    endereco_entrega: string;
    data_entrega?: string;
    criado_em: string;
}

interface PedidoListResponse {
    pedidos: Pedido[];
    total: number;
    page: number;
    limit: number;
    pages: number;
}

// Helper para formatar moeda brasileira
const formatCurrency = (value: number): string => {
    return new Intl.NumberFormat('pt-BR', {
        style: 'currency',
        currency: 'BRL'
    }).format(value);
};

// Helper para formatar data brasileira
const formatDate = (dateString: string): string => {
    if (!dateString) return '';
    const date = new Date(dateString);
    return new Intl.DateTimeFormat('pt-BR', {
        day: '2-digit',
        month: '2-digit',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
    }).format(date);
};

// Componente para exibir o status do pedido com cores
const StatusChip = ({ status }: { status: string }) => {
    let color:
        | 'default'
        | 'primary'
        | 'secondary'
        | 'error'
        | 'info'
        | 'success'
        | 'warning' = 'default';
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
    }

    return <Chip label={label} color={color} size="small" />;
};

// Mapeamento para exibição das formas de pagamento
const formaPagamentoMap: { [key: string]: string } = {
    dinheiro: 'Dinheiro',
    cartao_credito: 'Cartão de Crédito',
    cartao_debito: 'Cartão de Débito',
    pix: 'PIX',
    fiado: 'Fiado',
};

const Pedidos = () => {
    const navigate = useNavigate();
    const [pedidos, setPedidos] = useState<Pedido[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    // Estados para paginação
    const [page, setPage] = useState(0);
    const [limit, setLimit] = useState(10);
    const [total, setTotal] = useState(0);

    // Estados para filtros
    const [filtroStatus, setFiltroStatus] = useState<string>('');
    const [filtroCliente, setFiltroCliente] = useState<string>('');
    const [filtroDataInicio, setFiltroDataInicio] = useState<string>('');
    const [filtroDataFim, setFiltroDataFim] = useState<string>('');
    const [showFilters, setShowFilters] = useState(false);

    // Função para carregar os pedidos com filtros e paginação
    const carregarPedidos = async () => {
        setLoading(true);
        setError(null);

        try {
            // Obter token de autenticação do localStorage
            const token = localStorage.getItem('token');
            if (!token) {
                setError('Não autorizado. Faça login para continuar.');
                setLoading(false);
                setPedidos([]); // Inicializar como array vazio em caso de erro
                return;
            }

            // Construir parâmetros de consulta
            const params = new URLSearchParams();
            params.append('page', String(page + 1)); // API começa em 1, MUI em 0
            params.append('limit', String(limit));

            if (filtroStatus) params.append('status', filtroStatus);
            if (filtroCliente) params.append('cliente_id', filtroCliente);
            if (filtroDataInicio) params.append('data_inicio', filtroDataInicio);
            if (filtroDataFim) params.append('data_fim', filtroDataFim);

            // Fazer requisição à API
            const response = await axios.get<PedidoListResponse>(
                `${API_BASE_URL}/pedidos?${params.toString()}`,
                {
                    headers: {
                        Authorization: `Bearer ${token}`,
                    },
                }
            );

            // Adicionar depuração para identificar a estrutura dos dados
            if (response.data && Array.isArray(response.data.pedidos) && response.data.pedidos.length > 0) {
                console.log('Estrutura do primeiro pedido:', JSON.stringify(response.data.pedidos[0], null, 2));
            }

            // Garantir que temos um array de pedidos, mesmo que a resposta seja diferente do esperado
            if (response.data && Array.isArray(response.data.pedidos)) {
                setPedidos(response.data.pedidos);
                setTotal(response.data.total || 0);
            } else if (Array.isArray(response.data)) {
                // Caso a API retorne diretamente um array
                setPedidos(response.data);
                setTotal(response.data.length);
            } else {
                console.error('Resposta inesperada da API:', response.data);
                setPedidos([]);
                setTotal(0);
            }
        } catch (err) {
            console.error('Erro ao carregar pedidos:', err);
            setError('Não foi possível carregar os pedidos. Tente novamente mais tarde.');
            setPedidos([]); // Inicializar como array vazio em caso de erro
        } finally {
            setLoading(false);
        }
    };

    // Carregar pedidos ao montar o componente ou quando os filtros/paginação mudarem
    useEffect(() => {
        carregarPedidos();
    }, [page, limit, filtroStatus, filtroCliente, filtroDataInicio, filtroDataFim]);

    // Manipuladores de eventos
    const handleChangePage = (_event: unknown, newPage: number) => {
        setPage(newPage);
    };

    const handleChangeRowsPerPage = (event: React.ChangeEvent<HTMLInputElement>) => {
        setLimit(parseInt(event.target.value, 10));
        setPage(0);
    };

    const handleAplicarFiltros = () => {
        setPage(0); // Volta para a primeira página ao aplicar novos filtros
    };

    const handleLimparFiltros = () => {
        setFiltroStatus('');
        setFiltroCliente('');
        setFiltroDataInicio('');
        setFiltroDataFim('');
        setPage(0);
    };

    const handleNovoPedido = () => {
        navigate('/pedidos/novo');
    };

    const handleVerPedido = (id: number) => {
        navigate(`/pedidos/${id}`);
    };

    // Função aprimorada para obter o nome do cliente de forma segura
    const getClienteNome = (pedido: Pedido): string => {
        // Depuração para identificar a estrutura exata dos dados
        console.log('Estrutura do pedido:', pedido);

        // Verificar todas as possibilidades comuns
        if (pedido.cliente && typeof pedido.cliente === 'object' && pedido.cliente.nome) {
            return pedido.cliente.nome;
        } else if (pedido.cliente_nome) {
            return pedido.cliente_nome;
        } else if (typeof pedido.cliente === 'string') {
            return pedido.cliente; // Se cliente for diretamente uma string do nome
        } else if (pedido.cliente_id) {
            // Se tivermos apenas o ID, podemos mostrar isso temporariamente
            return `Cliente #${pedido.cliente_id}`;
        } else if (typeof pedido.cliente === 'number') {
            // Se cliente for um ID numérico
            return `Cliente #${pedido.cliente}`;
        } else {
            return 'Cliente não informado';
        }
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
                <Typography variant="h4">Pedidos</Typography>
                <Box>
                    <Button
                        variant="contained"
                        color="primary"
                        startIcon={<AddIcon />}
                        onClick={handleNovoPedido}
                        sx={{ mr: 1 }}
                    >
                        Novo Pedido
                    </Button>
                    <Button
                        variant="outlined"
                        startIcon={<RefreshIcon />}
                        onClick={carregarPedidos}
                    >
                        Atualizar
                    </Button>
                </Box>
            </Box>

            {/* Filtros */}
            <Paper sx={{ p: 2, mb: 3 }}>
                <Box
                    display="flex"
                    justifyContent="space-between"
                    alignItems="center"
                    mb={2}
                >
                    <Typography variant="h6">Filtros</Typography>
                    <IconButton onClick={() => setShowFilters(!showFilters)}>
                        <FilterListIcon />
                    </IconButton>
                </Box>

                {showFilters && (
                    <Grid container spacing={2}>
                        <Grid item xs={12} sm={6} md={3}>
                            <TextField
                                label="Status"
                                select
                                fullWidth
                                value={filtroStatus}
                                onChange={(e) => setFiltroStatus(e.target.value)}
                                size="small"
                            >
                                <MenuItem value="">Todos</MenuItem>
                                <MenuItem value="novo">Novo</MenuItem>
                                <MenuItem value="em_preparo">Em Preparo</MenuItem>
                                <MenuItem value="em_entrega">Em Entrega</MenuItem>
                                <MenuItem value="entregue">Entregue</MenuItem>
                                <MenuItem value="finalizado">Finalizado</MenuItem>
                                <MenuItem value="cancelado">Cancelado</MenuItem>
                            </TextField>
                        </Grid>
                        <Grid item xs={12} sm={6} md={3}>
                            <TextField
                                label="ID do Cliente"
                                fullWidth
                                value={filtroCliente}
                                onChange={(e) => setFiltroCliente(e.target.value)}
                                size="small"
                                type="number"
                            />
                        </Grid>
                        <Grid item xs={12} sm={6} md={3}>
                            <TextField
                                label="Data Início"
                                type="date"
                                fullWidth
                                value={filtroDataInicio}
                                onChange={(e) => setFiltroDataInicio(e.target.value)}
                                size="small"
                                InputLabelProps={{ shrink: true }}
                            />
                        </Grid>
                        <Grid item xs={12} sm={6} md={3}>
                            <TextField
                                label="Data Fim"
                                type="date"
                                fullWidth
                                value={filtroDataFim}
                                onChange={(e) => setFiltroDataFim(e.target.value)}
                                size="small"
                                InputLabelProps={{ shrink: true }}
                            />
                        </Grid>
                        <Grid item xs={12}>
                            <Box display="flex" justifyContent="flex-end">
                                <Button
                                    variant="outlined"
                                    onClick={handleLimparFiltros}
                                    sx={{ mr: 1 }}
                                >
                                    Limpar
                                </Button>
                                <Button
                                    variant="contained"
                                    startIcon={<SearchIcon />}
                                    onClick={handleAplicarFiltros}
                                >
                                    Buscar
                                </Button>
                            </Box>
                        </Grid>
                    </Grid>
                )}
            </Paper>

            {/* Tabela de Pedidos */}
            <Paper>
                <TableContainer>
                    <Table size="small">
                        <TableHead>
                            <TableRow>
                                <TableCell>ID</TableCell>
                                <TableCell>Data</TableCell>
                                <TableCell>Cliente</TableCell>
                                <TableCell>Endereço</TableCell>
                                <TableCell>Valor</TableCell>
                                <TableCell>Pagamento</TableCell>
                                <TableCell>Status</TableCell>
                                <TableCell align="right">Ações</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {loading ? (
                                <TableRow>
                                    <TableCell colSpan={8} align="center">
                                        <CircularProgress size={24} />
                                    </TableCell>
                                </TableRow>
                            ) : error ? (
                                <TableRow>
                                    <TableCell colSpan={8} align="center" sx={{ color: 'error.main' }}>
                                        {error}
                                    </TableCell>
                                </TableRow>
                            ) : !pedidos || pedidos.length === 0 ? (
                                <TableRow>
                                    <TableCell colSpan={8} align="center">
                                        Nenhum pedido encontrado.
                                    </TableCell>
                                </TableRow>
                            ) : (
                                pedidos.map((pedido) => (
                                    <TableRow key={pedido.id} hover onClick={() => handleVerPedido(pedido.id)}>
                                        <TableCell>{pedido.id}</TableCell>
                                        <TableCell>{formatDate(pedido.criado_em)}</TableCell>
                                        <TableCell>{getClienteNome(pedido)}</TableCell>
                                        <TableCell sx={{ maxWidth: '200px', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                                            {pedido.endereco_entrega}
                                        </TableCell>
                                        <TableCell>{formatCurrency(pedido.valor_total)}</TableCell>
                                        <TableCell>{formaPagamentoMap[pedido.forma_pagamento] || pedido.forma_pagamento}</TableCell>
                                        <TableCell>
                                            <StatusChip status={pedido.status} />
                                        </TableCell>
                                        <TableCell align="right">
                                            <Button
                                                size="small"
                                                variant="text"
                                                onClick={(e) => {
                                                    e.stopPropagation();
                                                    handleVerPedido(pedido.id);
                                                }}
                                            >
                                                Detalhes
                                            </Button>
                                        </TableCell>
                                    </TableRow>
                                ))
                            )}
                        </TableBody>
                    </Table>
                </TableContainer>

                {/* Paginação */}
                <TablePagination
                    component="div"
                    count={total}
                    page={page}
                    onPageChange={handleChangePage}
                    rowsPerPage={limit}
                    onRowsPerPageChange={handleChangeRowsPerPage}
                    rowsPerPageOptions={[5, 10, 25, 50]}
                    labelRowsPerPage="Itens por página"
                    labelDisplayedRows={({ from, to, count }) => `${from}-${to} de ${count}`}
                />
            </Paper>
        </Box>
    );
};

export default Pedidos;