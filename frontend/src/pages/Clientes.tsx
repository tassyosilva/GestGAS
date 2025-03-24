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
    Button,
    TextField,
    InputAdornment,
    IconButton,
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    Grid,
    MenuItem,
    Chip,
    Snackbar,
    Alert,
    Tabs,
    Tab,
    CircularProgress,
    Tooltip,
} from '@mui/material';
import {
    Add as AddIcon,
    Search as SearchIcon,
    Edit as EditIcon,
    Delete as DeleteIcon,
    Phone as PhoneIcon,
    Person as PersonIcon,
    Place as PlaceIcon,
    Refresh as RefreshIcon,
    Clear as ClearIcon,
} from '@mui/icons-material';
import axios from 'axios';
import API_BASE_URL from '../config/api';

// Interfaces
interface Cliente {
    id: number;
    nome: string;
    telefone: string;
    cpf?: string;
    email?: string;
    endereco?: string;
    complemento?: string;
    bairro?: string;
    cidade?: string;
    estado?: string;
    cep?: string;
    observacoes?: string;
    canal_origem?: string;
    criado_em: string;
    atualizado_em: string;
}

interface ClienteResponse extends Cliente {
    ultimos_pedidos?: PedidoResumido[];
    total_pedidos: number;
}

interface PedidoResumido {
    id: number;
    status: string;
    forma_pagamento: string;
    valor_total: number;
    data_pedido: string;
}

interface ClienteForm {
    nome: string;
    telefone: string;
    cpf: string;
    email: string;
    endereco: string;
    complemento: string;
    bairro: string;
    cidade: string;
    estado: string;
    cep: string;
    observacoes: string;
    canal_origem: string;
}

// Componente principal
const Clientes: React.FC = () => {
    // Estados
    const [clientes, setClientes] = useState<Cliente[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    // Estados para paginação
    const [page, setPage] = useState(0);
    const [rowsPerPage, setRowsPerPage] = useState(10);
    const [total, setTotal] = useState(0);

    // Estados para filtros
    const [filtroNome, setFiltroNome] = useState('');
    const [filtroTelefone, setFiltroTelefone] = useState('');
    const [filtroCanalOrigem, setFiltroCanalOrigem] = useState('');

    // Estados para o formulário
    const [dialogOpen, setDialogOpen] = useState(false);
    const [formMode, setFormMode] = useState<'criar' | 'editar'>('criar');
    const [clienteAtual, setClienteAtual] = useState<Cliente | null>(null);
    const [formData, setFormData] = useState<ClienteForm>({
        nome: '',
        telefone: '',
        cpf: '',
        email: '',
        endereco: '',
        complemento: '',
        bairro: '',
        cidade: '',
        estado: '',
        cep: '',
        observacoes: '',
        canal_origem: 'whatsapp',
    });

    // Estado para detalhes do cliente
    const [detailsDialogOpen, setDetailsDialogOpen] = useState(false);
    const [clienteDetalhes, setClienteDetalhes] = useState<ClienteResponse | null>(null);
    const [loadingDetalhes, setLoadingDetalhes] = useState(false);

    // Estado para confirmação de exclusão
    const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
    const [clienteParaExcluir, setClienteParaExcluir] = useState<Cliente | null>(null);
    const [deletingCliente, setDeletingCliente] = useState(false);

    // Estado para snackbar de feedback
    const [snackbar, setSnackbar] = useState({
        open: false,
        message: '',
        severity: 'success' as 'success' | 'error' | 'warning' | 'info',
    });

    // Estado para abas do formulário
    const [tabValue, setTabValue] = useState(0);

    // Carregar clientes na montagem do componente e quando os filtros/paginação mudarem
    useEffect(() => {
        carregarClientes();
    }, [page, rowsPerPage, filtroNome, filtroTelefone, filtroCanalOrigem]);

    // Função para carregar a lista de clientes
    const carregarClientes = async () => {
        setLoading(true);
        setError(null);

        try {
            const token = localStorage.getItem('token');
            if (!token) {
                setError('Não autorizado. Faça login para continuar.');
                setLoading(false);
                return;
            }

            // Construir parâmetros de consulta
            const params = new URLSearchParams();
            params.append('page', String(page + 1)); // API começa em 1, MUI em 0
            params.append('limit', String(rowsPerPage));

            if (filtroNome) params.append('nome', filtroNome);
            if (filtroTelefone) params.append('telefone', filtroTelefone);
            if (filtroCanalOrigem) params.append('canal_origem', filtroCanalOrigem);

            const response = await axios.get(`${API_BASE_URL}/clientes?${params.toString()}`, {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });

            setClientes(response.data.clientes || []);
            setTotal(response.data.total || 0);
        } catch (err: any) {
            console.error('Erro ao carregar clientes:', err);
            setError('Não foi possível carregar os clientes. Tente novamente mais tarde.');
        } finally {
            setLoading(false);
        }
    };

    // Função para buscar detalhes de um cliente
    const buscarDetalhesCliente = async (id: number) => {
        setLoadingDetalhes(true);

        try {
            const token = localStorage.getItem('token');
            if (!token) {
                setError('Não autorizado. Faça login para continuar.');
                setLoadingDetalhes(false);
                return;
            }

            const response = await axios.get(`${API_BASE_URL}/clientes/${id}`, {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });

            setClienteDetalhes(response.data);
            setDetailsDialogOpen(true);
        } catch (err: any) {
            console.error('Erro ao buscar detalhes do cliente:', err);
            setSnackbar({
                open: true,
                message: 'Erro ao buscar detalhes do cliente',
                severity: 'error',
            });
        } finally {
            setLoadingDetalhes(false);
        }
    };

    // Função para abrir o formulário de criação
    const abrirFormularioCriar = () => {
        setFormMode('criar');
        setFormData({
            nome: '',
            telefone: '',
            cpf: '',
            email: '',
            endereco: '',
            complemento: '',
            bairro: '',
            cidade: '',
            estado: '',
            cep: '',
            observacoes: '',
            canal_origem: 'whatsapp',
        });
        setTabValue(0);
        setDialogOpen(true);
    };

    // Função para abrir o formulário de edição
    const abrirFormularioEditar = (cliente: Cliente) => {
        setFormMode('editar');
        setClienteAtual(cliente);
        setFormData({
            nome: cliente.nome || '',
            telefone: cliente.telefone || '',
            cpf: cliente.cpf || '',
            email: cliente.email || '',
            endereco: cliente.endereco || '',
            complemento: cliente.complemento || '',
            bairro: cliente.bairro || '',
            cidade: cliente.cidade || '',
            estado: cliente.estado || '',
            cep: cliente.cep || '',
            observacoes: cliente.observacoes || '',
            canal_origem: cliente.canal_origem || 'whatsapp',
        });
        setTabValue(0);
        setDialogOpen(true);
    };

    // Função para lidar com a mudança nos campos do formulário
    const handleFormChange = (e: React.ChangeEvent<HTMLInputElement | { name?: string; value: unknown }>) => {
        const { name, value } = e.target;
        if (name) {
            setFormData({
                ...formData,
                [name]: value,
            });
        }
    };

    // Função para salvar cliente (criar ou editar)
    const salvarCliente = async () => {
        // Validar campos obrigatórios
        if (!formData.nome || !formData.telefone) {
            setSnackbar({
                open: true,
                message: 'Nome e Telefone são campos obrigatórios',
                severity: 'error',
            });
            return;
        }

        try {
            const token = localStorage.getItem('token');
            if (!token) {
                setSnackbar({
                    open: true,
                    message: 'Não autorizado. Faça login para continuar.',
                    severity: 'error',
                });
                return;
            }

            if (formMode === 'criar') {
                // Criar novo cliente
                await axios.post(
                    `${API_BASE_URL}/clientes`,
                    formData,
                    {
                        headers: {
                            Authorization: `Bearer ${token}`,
                            'Content-Type': 'application/json',
                        },
                    }
                );

                setSnackbar({
                    open: true,
                    message: 'Cliente criado com sucesso',
                    severity: 'success',
                });
            } else if (formMode === 'editar' && clienteAtual) {
                // Editar cliente existente
                await axios.put(
                    `${API_BASE_URL}/clientes/${clienteAtual.id}`,
                    formData,
                    {
                        headers: {
                            Authorization: `Bearer ${token}`,
                            'Content-Type': 'application/json',
                        },
                    }
                );

                setSnackbar({
                    open: true,
                    message: 'Cliente atualizado com sucesso',
                    severity: 'success',
                });
            }

            // Fechar diálogo e recarregar lista
            setDialogOpen(false);
            carregarClientes();
        } catch (err: any) {
            console.error('Erro ao salvar cliente:', err);
            setSnackbar({
                open: true,
                message: err.response?.data || 'Erro ao salvar cliente',
                severity: 'error',
            });
        }
    };

    // Função para confirmar exclusão de cliente
    const confirmarExclusao = (cliente: Cliente) => {
        setClienteParaExcluir(cliente);
        setDeleteDialogOpen(true);
    };

    // Função para excluir cliente
    const excluirCliente = async () => {
        if (!clienteParaExcluir) return;

        setDeletingCliente(true);

        try {
            const token = localStorage.getItem('token');
            if (!token) {
                setSnackbar({
                    open: true,
                    message: 'Não autorizado. Faça login para continuar.',
                    severity: 'error',
                });
                setDeletingCliente(false);
                return;
            }

            await axios.delete(`${API_BASE_URL}/clientes/${clienteParaExcluir.id}`, {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });

            setSnackbar({
                open: true,
                message: 'Cliente excluído com sucesso',
                severity: 'success',
            });

            // Fechar diálogo e recarregar lista
            setDeleteDialogOpen(false);
            carregarClientes();
        } catch (err: any) {
            console.error('Erro ao excluir cliente:', err);
            setSnackbar({
                open: true,
                message: err.response?.data || 'Erro ao excluir cliente',
                severity: 'error',
            });
        } finally {
            setDeletingCliente(false);
        }
    };

    // Funções para paginação
    const handleChangePage = (_event: unknown, newPage: number) => {
        setPage(newPage);
    };

    const handleChangeRowsPerPage = (event: React.ChangeEvent<HTMLInputElement>) => {
        setRowsPerPage(parseInt(event.target.value, 10));
        setPage(0);
    };

    // Função para limpar filtros
    const limparFiltros = () => {
        setFiltroNome('');
        setFiltroTelefone('');
        setFiltroCanalOrigem('');
        setPage(0);
    };

    // Função para formatar data
    const formatarData = (dataString: string): string => {
        try {
            const data = new Date(dataString);
            return new Intl.DateTimeFormat('pt-BR', {
                day: '2-digit',
                month: '2-digit',
                year: 'numeric',
                hour: '2-digit',
                minute: '2-digit',
            }).format(data);
        } catch (e) {
            return dataString;
        }
    };

    // Componente para canal de origem
    const CanalOrigemChip = ({ canal }: { canal?: string }) => {
        if (!canal) return null;

        let color: 'default' | 'primary' | 'secondary' | 'error' | 'info' | 'success' | 'warning' = 'default';
        let label = canal;

        switch (canal) {
            case 'whatsapp':
                color = 'success';
                label = 'WhatsApp';
                break;
            case 'telefone':
                color = 'primary';
                label = 'Telefone';
                break;
            case 'presencial':
                color = 'secondary';
                label = 'Presencial';
                break;
            case 'aplicativo':
                color = 'info';
                label = 'Aplicativo';
                break;
        }

        return <Chip label={label} color={color} size="small" />;
    };

    // Renderizar os status dos pedidos
    const renderizarStatusPedido = (status: string) => {
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
        }

        return <Chip label={label} color={color} size="small" />;
    };

    // Formatar moeda
    const formatarMoeda = (valor: number): string => {
        return new Intl.NumberFormat('pt-BR', {
            style: 'currency',
            currency: 'BRL',
        }).format(valor);
    };

    return (
        <Box>
            {/* Cabeçalho */}
            <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
                <Typography variant="h4">Clientes</Typography>
                <Box>
                    <Button
                        variant="contained"
                        color="primary"
                        startIcon={<AddIcon />}
                        onClick={abrirFormularioCriar}
                        sx={{ mr: 1 }}
                    >
                        Novo Cliente
                    </Button>
                    <Button
                        variant="outlined"
                        startIcon={<RefreshIcon />}
                        onClick={carregarClientes}
                    >
                        Atualizar
                    </Button>
                </Box>
            </Box>

            {/* Filtros */}
            <Paper sx={{ p: 2, mb: 3 }}>
                <Grid container spacing={2} alignItems="center">
                    <Grid item xs={12} sm={4}>
                        <TextField
                            label="Nome"
                            value={filtroNome}
                            onChange={(e) => setFiltroNome(e.target.value)}
                            fullWidth
                            size="small"
                            InputProps={{
                                startAdornment: (
                                    <InputAdornment position="start">
                                        <PersonIcon fontSize="small" />
                                    </InputAdornment>
                                ),
                                endAdornment: filtroNome ? (
                                    <InputAdornment position="end">
                                        <IconButton
                                            size="small"
                                            onClick={() => setFiltroNome('')}
                                        >
                                            <ClearIcon fontSize="small" />
                                        </IconButton>
                                    </InputAdornment>
                                ) : null,
                            }}
                        />
                    </Grid>
                    <Grid item xs={12} sm={4}>
                        <TextField
                            label="Telefone"
                            value={filtroTelefone}
                            onChange={(e) => setFiltroTelefone(e.target.value)}
                            fullWidth
                            size="small"
                            InputProps={{
                                startAdornment: (
                                    <InputAdornment position="start">
                                        <PhoneIcon fontSize="small" />
                                    </InputAdornment>
                                ),
                                endAdornment: filtroTelefone ? (
                                    <InputAdornment position="end">
                                        <IconButton
                                            size="small"
                                            onClick={() => setFiltroTelefone('')}
                                        >
                                            <ClearIcon fontSize="small" />
                                        </IconButton>
                                    </InputAdornment>
                                ) : null,
                            }}
                        />
                    </Grid>
                    <Grid item xs={12} sm={4}>
                        <TextField
                            select
                            label="Canal de Origem"
                            value={filtroCanalOrigem}
                            onChange={(e) => setFiltroCanalOrigem(e.target.value)}
                            fullWidth
                            size="small"
                        >
                            <MenuItem value="">Todos os canais</MenuItem>
                            <MenuItem value="whatsapp">WhatsApp</MenuItem>
                            <MenuItem value="telefone">Telefone</MenuItem>
                            <MenuItem value="presencial">Presencial</MenuItem>
                            <MenuItem value="aplicativo">Aplicativo</MenuItem>
                        </TextField>
                    </Grid>
                    {(filtroNome || filtroTelefone || filtroCanalOrigem) && (
                        <Grid item xs={12}>
                            <Box display="flex" justifyContent="flex-end">
                                <Button
                                    variant="text"
                                    color="inherit"
                                    startIcon={<ClearIcon />}
                                    onClick={limparFiltros}
                                    size="small"
                                >
                                    Limpar Filtros
                                </Button>
                            </Box>
                        </Grid>
                    )}
                </Grid>
            </Paper>

            {/* Tabela de clientes */}
            <Paper>
                <TableContainer>
                    <Table size="small">
                        <TableHead>
                            <TableRow>
                                <TableCell>Nome</TableCell>
                                <TableCell>Telefone</TableCell>
                                <TableCell>Endereço</TableCell>
                                <TableCell>Canal</TableCell>
                                <TableCell>Cadastro</TableCell>
                                <TableCell align="right">Ações</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {loading ? (
                                <TableRow>
                                    <TableCell colSpan={6} align="center">
                                        <CircularProgress size={24} sx={{ my: 2 }} />
                                    </TableCell>
                                </TableRow>
                            ) : error ? (
                                <TableRow>
                                    <TableCell colSpan={6} align="center">
                                        <Typography color="error">{error}</Typography>
                                    </TableCell>
                                </TableRow>
                            ) : clientes.length === 0 ? (
                                <TableRow>
                                    <TableCell colSpan={6} align="center">
                                        <Typography>Nenhum cliente encontrado</Typography>
                                    </TableCell>
                                </TableRow>
                            ) : (
                                clientes.map((cliente) => (
                                    <TableRow key={cliente.id} hover>
                                        <TableCell
                                            sx={{ cursor: 'pointer' }}
                                            onClick={() => buscarDetalhesCliente(cliente.id)}
                                        >
                                            {cliente.nome}
                                        </TableCell>
                                        <TableCell>{cliente.telefone}</TableCell>
                                        <TableCell>
                                            {cliente.endereco ? (
                                                <Tooltip title={`${cliente.endereco}${cliente.complemento ? `, ${cliente.complemento}` : ''}, ${cliente.bairro || ''} - ${cliente.cidade || ''} / ${cliente.estado || ''}`}>
                                                    <Typography noWrap sx={{ maxWidth: 200 }}>
                                                        {cliente.endereco}
                                                    </Typography>
                                                </Tooltip>
                                            ) : (
                                                <Typography color="text.secondary" variant="body2">
                                                    Não informado
                                                </Typography>
                                            )}
                                        </TableCell>
                                        <TableCell>
                                            <CanalOrigemChip canal={cliente.canal_origem} />
                                        </TableCell>
                                        <TableCell>{formatarData(cliente.criado_em)}</TableCell>
                                        <TableCell align="right">
                                            <Tooltip title="Detalhes">
                                                <IconButton
                                                    size="small"
                                                    color="info"
                                                    onClick={() => buscarDetalhesCliente(cliente.id)}
                                                >
                                                    <SearchIcon fontSize="small" />
                                                </IconButton>
                                            </Tooltip>
                                            <Tooltip title="Editar">
                                                <IconButton
                                                    size="small"
                                                    color="primary"
                                                    onClick={() => abrirFormularioEditar(cliente)}
                                                >
                                                    <EditIcon fontSize="small" />
                                                </IconButton>
                                            </Tooltip>
                                            <Tooltip title="Excluir">
                                                <IconButton
                                                    size="small"
                                                    color="error"
                                                    onClick={() => confirmarExclusao(cliente)}
                                                >
                                                    <DeleteIcon fontSize="small" />
                                                </IconButton>
                                            </Tooltip>
                                        </TableCell>
                                    </TableRow>
                                ))
                            )}
                        </TableBody>
                    </Table>
                </TableContainer>

                {/* Paginação */}
                <TablePagination
                    rowsPerPageOptions={[5, 10, 25, 50]}
                    component="div"
                    count={total}
                    rowsPerPage={rowsPerPage}
                    page={page}
                    onPageChange={handleChangePage}
                    onRowsPerPageChange={handleChangeRowsPerPage}
                    labelRowsPerPage="Itens por página"
                    labelDisplayedRows={({ from, to, count }) => `${from}-${to} de ${count}`}
                />
            </Paper>

            {/* Diálogo para criar/editar cliente */}
            <Dialog
                open={dialogOpen}
                onClose={() => setDialogOpen(false)}
                maxWidth="md"
                fullWidth
            >
                <DialogTitle>
                    {formMode === 'criar' ? 'Novo Cliente' : 'Editar Cliente'}
                </DialogTitle>
                <DialogContent>
                    <Tabs
                        value={tabValue}
                        onChange={(_e, newValue) => setTabValue(newValue)}
                        sx={{ mb: 2 }}
                    >
                        <Tab label="Dados Básicos" />
                        <Tab label="Endereço" />
                        <Tab label="Observações" />
                    </Tabs>

                    {/* Formulário de Dados Básicos */}
                    {tabValue === 0 && (
                        <Grid container spacing={2} sx={{ mt: 1 }}>
                            <Grid item xs={12} sm={6}>
                                <TextField
                                    name="nome"
                                    label="Nome"
                                    value={formData.nome}
                                    onChange={handleFormChange}
                                    fullWidth
                                    required
                                />
                            </Grid>
                            <Grid item xs={12} sm={6}>
                                <TextField
                                    name="telefone"
                                    label="Telefone"
                                    value={formData.telefone}
                                    onChange={handleFormChange}
                                    fullWidth
                                    required
                                />
                            </Grid>
                            <Grid item xs={12} sm={6}>
                                <TextField
                                    name="cpf"
                                    label="CPF"
                                    value={formData.cpf}
                                    onChange={handleFormChange}
                                    fullWidth
                                />
                            </Grid>
                            <Grid item xs={12} sm={6}>
                                <TextField
                                    name="email"
                                    label="E-mail"
                                    type="email"
                                    value={formData.email}
                                    onChange={handleFormChange}
                                    fullWidth
                                />
                            </Grid>
                            <Grid item xs={12}>
                                <TextField
                                    name="canal_origem"
                                    label="Canal de Origem"
                                    select
                                    value={formData.canal_origem}
                                    onChange={handleFormChange}
                                    fullWidth
                                >
                                    <MenuItem value="whatsapp">WhatsApp</MenuItem>
                                    <MenuItem value="telefone">Telefone</MenuItem>
                                    <MenuItem value="presencial">Presencial</MenuItem>
                                    <MenuItem value="aplicativo">Aplicativo</MenuItem>
                                </TextField>
                            </Grid>
                        </Grid>
                    )}

                    {/* Formulário de Endereço */}
                    {tabValue === 1 && (
                        <Grid container spacing={2} sx={{ mt: 1 }}>
                            <Grid item xs={12}>
                                <TextField
                                    name="endereco"
                                    label="Endereço"
                                    value={formData.endereco}
                                    onChange={handleFormChange}
                                    fullWidth
                                />
                            </Grid>
                            <Grid item xs={12} sm={6}>
                                <TextField
                                    name="complemento"
                                    label="Complemento"
                                    value={formData.complemento}
                                    onChange={handleFormChange}
                                    fullWidth
                                />
                            </Grid>
                            <Grid item xs={12} sm={6}>
                                <TextField
                                    name="bairro"
                                    label="Bairro"
                                    value={formData.bairro}
                                    onChange={handleFormChange}
                                    fullWidth
                                />
                            </Grid>
                            <Grid item xs={12} sm={6}>
                                <TextField
                                    name="cidade"
                                    label="Cidade"
                                    value={formData.cidade}
                                    onChange={handleFormChange}
                                    fullWidth
                                />
                            </Grid>
                            <Grid item xs={12} sm={2}>
                                <TextField
                                    name="estado"
                                    label="UF"
                                    value={formData.estado}
                                    onChange={handleFormChange}
                                    fullWidth
                                />
                            </Grid>
                            <Grid item xs={12} sm={4}>
                                <TextField
                                    name="cep"
                                    label="CEP"
                                    value={formData.cep}
                                    onChange={handleFormChange}
                                    fullWidth
                                />
                            </Grid>
                        </Grid>
                    )}

                    {/* Formulário de Observações */}
                    {tabValue === 2 && (
                        <Grid container spacing={2} sx={{ mt: 1 }}>
                            <Grid item xs={12}>
                                <TextField
                                    name="observacoes"
                                    label="Observações"
                                    value={formData.observacoes}
                                    onChange={handleFormChange}
                                    fullWidth
                                    multiline
                                    rows={4}
                                />
                            </Grid>
                        </Grid>
                    )}
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setDialogOpen(false)}>Cancelar</Button>
                    <Button onClick={salvarCliente} variant="contained" color="primary">
                        Salvar
                    </Button>
                </DialogActions>
            </Dialog>

            {/* Diálogo de detalhes do cliente */}
            <Dialog
                open={detailsDialogOpen}
                onClose={() => setDetailsDialogOpen(false)}
                maxWidth="md"
                fullWidth
            >
                {loadingDetalhes ? (
                    <DialogContent>
                        <Box display="flex" justifyContent="center" p={4}>
                            <CircularProgress />
                        </Box>
                    </DialogContent>
                ) : clienteDetalhes ? (
          <>
                        <DialogTitle>
                            <Box display="flex" justifyContent="space-between" alignItems="center">
                                {clienteDetalhes.nome}
                                <CanalOrigemChip canal={clienteDetalhes.canal_origem} />
                            </Box>
                        </DialogTitle>
                        <DialogContent>
                            <Grid container spacing={3}>
                                {/* Informações básicas */}
                                <Grid item xs={12} md={6}>
                                    <Paper variant="outlined" sx={{ p: 2 }}>
                                        <Typography variant="h6" gutterBottom>
                                            Informações Básicas
                                        </Typography>
                                        <Grid container spacing={2}>
                                            <Grid item xs={12}>
                                                <Typography variant="subtitle2">Telefone:</Typography>
                                                <Typography>{clienteDetalhes.telefone}</Typography>
                                            </Grid>
                                            {clienteDetalhes.cpf && (
                                                <Grid item xs={12} sm={6}>
                                                    <Typography variant="subtitle2">CPF:</Typography>
                                                    <Typography>{clienteDetalhes.cpf}</Typography>
                                                </Grid>
                                            )}
                                            {clienteDetalhes.email && (
                                                <Grid item xs={12} sm={6}>
                                                    <Typography variant="subtitle2">E-mail:</Typography>
                                                    <Typography>{clienteDetalhes.email}</Typography>
                                                </Grid>
                                            )}
                                            <Grid item xs={12}>
                        <Typography variant="subtitle2">Cliente desde:</Typography>
                        <Typography>{formatarData(clienteDetalhes.criado_em)}</Typography>
                      </Grid>
                      {clienteDetalhes.observacoes && (
                        <Grid item xs={12}>
                          <Typography variant="subtitle2">Observações:</Typography>
                          <Typography>{clienteDetalhes.observacoes}</Typography>
                        </Grid>
                      )}
                    </Grid>
                  </Paper>
                </Grid>
                
                {/* Endereço */}
                <Grid item xs={12} md={6}>
                  <Paper variant="outlined" sx={{ p: 2 }}>
                    <Typography variant="h6" gutterBottom>
                      Endereço
                    </Typography>
                    {clienteDetalhes.endereco ? (
                      <Grid container spacing={2}>
                        <Grid item xs={12}>
                          <Typography variant="subtitle2">Endereço:</Typography>
                          <Typography>{clienteDetalhes.endereco}</Typography>
                        </Grid>
                        {clienteDetalhes.complemento && (
                          <Grid item xs={12}>
                            <Typography variant="subtitle2">Complemento:</Typography>
                            <Typography>{clienteDetalhes.complemento}</Typography>
                          </Grid>
                        )}
                        <Grid item xs={12} sm={6}>
                          <Typography variant="subtitle2">Bairro:</Typography>
                          <Typography>{clienteDetalhes.bairro || '-'}</Typography>
                        </Grid>
                        <Grid item xs={12} sm={6}>
                          <Typography variant="subtitle2">CEP:</Typography>
                          <Typography>{clienteDetalhes.cep || '-'}</Typography>
                        </Grid>
                        <Grid item xs={12} sm={6}>
                          <Typography variant="subtitle2">Cidade:</Typography>
                          <Typography>{clienteDetalhes.cidade || '-'}</Typography>
                        </Grid>
                        <Grid item xs={12} sm={6}>
                          <Typography variant="subtitle2">Estado:</Typography>
                          <Typography>{clienteDetalhes.estado || '-'}</Typography>
                        </Grid>
                      </Grid>
                    ) : (
                      <Typography color="text.secondary">
                        Nenhum endereço cadastrado
                      </Typography>
                    )}
                  </Paper>
                </Grid>
                
                {/* Últimos pedidos */}
                <Grid item xs={12}>
                  <Paper variant="outlined" sx={{ p: 2 }}>
                    <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
                      <Typography variant="h6">
                        Últimos Pedidos
                      </Typography>
                      <Chip 
                        label={`Total: ${clienteDetalhes.total_pedidos} pedido(s)`} 
                        color="primary" 
                        variant="outlined"
                      />
                    </Box>
                    
                    {clienteDetalhes.ultimos_pedidos && clienteDetalhes.ultimos_pedidos.length > 0 ? (
                      <TableContainer>
                        <Table size="small">
                          <TableHead>
                            <TableRow>
                              <TableCell>ID</TableCell>
                              <TableCell>Data</TableCell>
                              <TableCell>Valor</TableCell>
                              <TableCell>Pagamento</TableCell>
                              <TableCell>Status</TableCell>
                            </TableRow>
                          </TableHead>
                          <TableBody>
                            {clienteDetalhes.ultimos_pedidos.map((pedido) => (
                              <TableRow key={pedido.id} hover>
                                <TableCell>{pedido.id}</TableCell>
                                <TableCell>{formatarData(pedido.data_pedido)}</TableCell>
                                <TableCell>{formatarMoeda(pedido.valor_total)}</TableCell>
                                <TableCell>{pedido.forma_pagamento}</TableCell>
                                <TableCell>{renderizarStatusPedido(pedido.status)}</TableCell>
                              </TableRow>
                            ))}
                          </TableBody>
                        </Table>
                      </TableContainer>
                    ) : (
                      <Typography color="text.secondary">
                        Este cliente ainda não realizou nenhum pedido
                      </Typography>
                    )}
                  </Paper>
                </Grid>
              </Grid>
            </DialogContent>
            <DialogActions>
              <Button 
                startIcon={<EditIcon />}
                onClick={() => {
                  setDetailsDialogOpen(false);
                  abrirFormularioEditar(clienteDetalhes);
                }}
              >
                Editar Cliente
              </Button>
              <Button onClick={() => setDetailsDialogOpen(false)}>
                Fechar
              </Button>
            </DialogActions>
          </>
        ) : (
          <DialogContent>
            <Typography color="error">
              Erro ao carregar detalhes do cliente
            </Typography>
          </DialogContent>
        )}
      </Dialog>
      
      {/* Diálogo de confirmação de exclusão */}
      <Dialog
        open={deleteDialogOpen}
        onClose={() => !deletingCliente && setDeleteDialogOpen(false)}
      >
        <DialogTitle>Confirmar Exclusão</DialogTitle>
        <DialogContent>
          <Typography>
            Tem certeza que deseja excluir o cliente <strong>{clienteParaExcluir?.nome}</strong>?
          </Typography>
          <Typography color="error" variant="body2" sx={{ mt: 1 }}>
            Esta ação não poderá ser desfeita.
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button 
            onClick={() => setDeleteDialogOpen(false)} 
            disabled={deletingCliente}
          >
            Cancelar
          </Button>
          <Button 
            onClick={excluirCliente} 
            color="error" 
            variant="contained"
            disabled={deletingCliente}
          >
            {deletingCliente ? <CircularProgress size={24} /> : 'Excluir'}
          </Button>
        </DialogActions>
      </Dialog>
      
      {/* Snackbar para feedback */}
      <Snackbar
        open={snackbar.open}
        autoHideDuration={6000}
        onClose={() => setSnackbar({ ...snackbar, open: false })}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
      >
        <Alert 
          onClose={() => setSnackbar({ ...snackbar, open: false })}
          severity={snackbar.severity}
          variant="filled"
        >
          {snackbar.message}
        </Alert>
      </Snackbar>
    </Box>
  );
};

export default Clientes;