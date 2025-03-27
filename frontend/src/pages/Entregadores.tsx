import { useState, useEffect } from 'react';
import {
    Box,
    Typography,
    Paper,
    Button,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    TablePagination,
    IconButton,
    TextField,
    CircularProgress,
    Snackbar,
    Alert,
    InputAdornment,
    Card,
    CardContent,
    Grid,
} from '@mui/material';
import {
    Add as AddIcon,
    Search as SearchIcon,
    Clear as ClearIcon,
    Person as PersonIcon,
} from '@mui/icons-material';
import { usuarioService, Usuario } from '../services/usuarioService';
import { Link } from 'react-router-dom';

const Entregadores = () => {
    // Estados - iniciados com arrays vazios para garantir que nunca sejam null
    const [entregadores, setEntregadores] = useState<Usuario[]>([]);
    const [entregadoresFiltrados, setEntregadoresFiltrados] = useState<Usuario[]>([]);
    const [carregando, setCarregando] = useState(false);
    const [erro, setErro] = useState<string | null>(null);

    // Estados para paginação e busca
    const [page, setPage] = useState(0);
    const [rowsPerPage, setRowsPerPage] = useState(10);
    const [termoBusca, setTermoBusca] = useState('');

    // Buscar entregadores ao carregar a página
    useEffect(() => {
        buscarEntregadores();
    }, []);

    // Filtrar entregadores baseado no termo de busca
    useEffect(() => {
        if (!entregadores) return; // Proteção adicional

        if (termoBusca.trim() === '') {
            setEntregadoresFiltrados(entregadores);
        } else {
            const termoBuscaLowerCase = termoBusca.toLowerCase();
            const resultadosFiltrados = entregadores.filter(entregador =>
                entregador.nome.toLowerCase().includes(termoBuscaLowerCase) ||
                (entregador.email && entregador.email.toLowerCase().includes(termoBuscaLowerCase))
            );
            setEntregadoresFiltrados(resultadosFiltrados);
        }
        setPage(0);
    }, [termoBusca, entregadores]);

    // Função para buscar entregadores da API
    const buscarEntregadores = async () => {
        setCarregando(true);
        try {
            const data = await usuarioService.listarEntregadores();
            setEntregadores(data || []); // Garantir que sempre seja um array
            setEntregadoresFiltrados(data || []); // Garantir que sempre seja um array
        } catch (error) {
            console.error('Erro ao buscar entregadores:', error);
            setErro('Não foi possível carregar os entregadores. Tente novamente mais tarde.');
            // Em caso de erro, garantir que os estados sejam arrays vazios e não null
            setEntregadores([]);
            setEntregadoresFiltrados([]);
        } finally {
            setCarregando(false);
        }
    };

    // Handlers para paginação
    const handleChangePage = (_event: unknown, newPage: number) => {
        setPage(newPage);
    };

    const handleChangeRowsPerPage = (event: React.ChangeEvent<HTMLInputElement>) => {
        setRowsPerPage(parseInt(event.target.value, 10));
        setPage(0);
    };

    // Handler para busca
    const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setTermoBusca(e.target.value);
    };

    const limparBusca = () => {
        setTermoBusca('');
    };

    // Fechar alertas
    const fecharAlerta = () => {
        setErro(null);
    };

    // Garantir que entregadoresFiltrados nunca seja null antes de usar slice
    const entregadoresPaginados = entregadoresFiltrados?.slice(
        page * rowsPerPage,
        page * rowsPerPage + rowsPerPage
    ) || [];

    return (
        <Box>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
                <Typography variant="h4">
                    Entregadores
                </Typography>
                <Button
                    variant="contained"
                    color="primary"
                    startIcon={<AddIcon />}
                    component={Link}
                    to="/usuarios"
                >
                    Cadastrar Entregador
                </Button>
            </Box>

            {/* Estatísticas */}
            <Box sx={{ mb: 3 }}>
                <Card>
                    <CardContent>
                        <Typography variant="h6" gutterBottom>
                            Estatísticas de Entregadores
                        </Typography>
                        <Grid container spacing={2}>
                            <Grid item xs={12} sm={6} md={4}>
                                <Paper sx={{ p: 2, textAlign: 'center' }}>
                                    <Typography variant="body2" color="text.secondary">
                                        Total de Entregadores
                                    </Typography>
                                    <Typography variant="h4">
                                        {carregando ? <CircularProgress size={20} /> : entregadores?.length || 0}
                                    </Typography>
                                </Paper>
                            </Grid>
                            {/* Outras estatísticas poderiam ser adicionadas aqui */}
                        </Grid>
                    </CardContent>
                </Card>
            </Box>

            {/* Campo de busca */}
            <Paper sx={{ p: 2, mb: 2 }}>
                <TextField
                    fullWidth
                    placeholder="Buscar entregadores pelo nome ou email..."
                    value={termoBusca}
                    onChange={handleSearchChange}
                    InputProps={{
                        startAdornment: (
                            <InputAdornment position="start">
                                <SearchIcon />
                            </InputAdornment>
                        ),
                        endAdornment: termoBusca && (
                            <InputAdornment position="end">
                                <IconButton onClick={limparBusca} size="small">
                                    <ClearIcon />
                                </IconButton>
                            </InputAdornment>
                        ),
                    }}
                />
            </Paper>

            {/* Tabela de entregadores */}
            <Paper sx={{ width: '100%', overflow: 'hidden' }}>
                <TableContainer>
                    <Table>
                        <TableHead>
                            <TableRow>
                                <TableCell>Nome</TableCell>
                                <TableCell>Login</TableCell>
                                <TableCell>Email</TableCell>
                                <TableCell>CPF</TableCell>
                                <TableCell align="center">Ações</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {carregando && entregadoresFiltrados.length === 0 ? (
                                <TableRow>
                                    <TableCell colSpan={5} align="center">
                                        <CircularProgress size={30} />
                                    </TableCell>
                                </TableRow>
                            ) : entregadoresFiltrados.length === 0 ? (
                                <TableRow>
                                    <TableCell colSpan={5} align="center">
                                        {termoBusca ? 'Nenhum entregador encontrado para a busca' : 'Nenhum entregador cadastrado'}
                                    </TableCell>
                                </TableRow>
                            ) : entregadoresPaginados.length === 0 ? (
                                <TableRow>
                                    <TableCell colSpan={5} align="center">
                                        Não há entregadores nesta página
                                    </TableCell>
                                </TableRow>
                            ) : (
                                entregadoresPaginados.map((entregador) => (
                                    <TableRow key={entregador.id}>
                                        <TableCell>{entregador.nome}</TableCell>
                                        <TableCell>{entregador.login}</TableCell>
                                        <TableCell>{entregador.email || '-'}</TableCell>
                                        <TableCell>{entregador.cpf || '-'}</TableCell>
                                        <TableCell align="center">
                                            <IconButton
                                                color="primary"
                                                component={Link}
                                                to={`/usuarios?id=${entregador.id}`}
                                                title="Ver detalhes"
                                            >
                                                <PersonIcon />
                                            </IconButton>
                                        </TableCell>
                                    </TableRow>
                                ))
                            )}
                        </TableBody>
                    </Table>
                </TableContainer>
                {/* Controles de paginação */}
                <TablePagination
                    rowsPerPageOptions={[5, 10, 25, 50]}
                    component="div"
                    count={entregadoresFiltrados?.length || 0}
                    rowsPerPage={rowsPerPage}
                    page={entregadoresFiltrados.length <= page * rowsPerPage && page > 0 ? page - 1 : page}
                    onPageChange={handleChangePage}
                    onRowsPerPageChange={handleChangeRowsPerPage}
                    labelRowsPerPage="Itens por página:"
                    labelDisplayedRows={({ from, to, count }) => `${from}-${to} de ${count}`}
                />
            </Paper>

            {/* Alertas de erro */}
            <Snackbar open={!!erro} autoHideDuration={6000} onClose={fecharAlerta}>
                <Alert onClose={fecharAlerta} severity="error" sx={{ width: '100%' }}>
                    {erro}
                </Alert>
            </Snackbar>
        </Box>
    );
};

export default Entregadores;