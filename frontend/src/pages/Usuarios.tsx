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
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    TextField,
    MenuItem,
    FormControl,
    InputLabel,
    Select,
    SelectChangeEvent,
    CircularProgress,
    Snackbar,
    Alert,
    InputAdornment,
    Chip,
} from '@mui/material';
import {
    Add as AddIcon,
    Edit as EditIcon,
    Delete as DeleteIcon,
    Search as SearchIcon,
    Clear as ClearIcon,
} from '@mui/icons-material';
import { usuarioService, Usuario, NovoUsuario, AtualizacaoUsuario } from '../services/usuarioService';
import { authService } from '../services/authService';

// Interface para o formulário
interface FormularioUsuario {
    nome: string;
    login: string;
    senha: string;
    confirmaSenha: string;
    cpf: string;
    email: string;
    perfil: string;
}

// Perfis disponíveis no sistema
const perfis = [
    { valor: 'admin', rotulo: 'Administrador' },
    { valor: 'gerente', rotulo: 'Gerente' },
    { valor: 'atendente', rotulo: 'Atendente' },
    { valor: 'entregador', rotulo: 'Entregador' },
];

const Usuarios = () => {
    // Estados
    const [usuarios, setUsuarios] = useState<Usuario[]>([]);
    const [usuariosFiltrados, setUsuariosFiltrados] = useState<Usuario[]>([]);
    const [dialogoAberto, setDialogoAberto] = useState(false);
    const [formulario, setFormulario] = useState<FormularioUsuario>({
        nome: '',
        login: '',
        senha: '',
        confirmaSenha: '',
        cpf: '',
        email: '',
        perfil: '',
    });
    const [editandoUsuario, setEditandoUsuario] = useState<number | null>(null);
    const [carregando, setCarregando] = useState(false);
    const [erro, setErro] = useState<string | null>(null);
    const [sucesso, setSucesso] = useState<string | null>(null);
    const [dialogoExclusaoAberto, setDialogoExclusaoAberto] = useState(false);
    const [usuarioParaExcluir, setUsuarioParaExcluir] = useState<number | null>(null);

    // Estados para paginação e busca
    const [page, setPage] = useState(0);
    const [rowsPerPage, setRowsPerPage] = useState(10);
    const [termoBusca, setTermoBusca] = useState('');

    // Obter usuário atual
    const usuarioAtual = authService.getUser();
    const isAdmin = usuarioAtual?.perfil === 'admin';

    // Buscar usuários ao carregar a página
    useEffect(() => {
        buscarUsuarios();
    }, []);

    // Filtrar usuários baseado no termo de busca
    useEffect(() => {
        if (termoBusca.trim() === '') {
            setUsuariosFiltrados(usuarios);
        } else {
            const termoBuscaLowerCase = termoBusca.toLowerCase();
            const resultadosFiltrados = usuarios.filter(usuario =>
                usuario.nome.toLowerCase().includes(termoBuscaLowerCase) ||
                usuario.login.toLowerCase().includes(termoBuscaLowerCase) ||
                (usuario.email && usuario.email.toLowerCase().includes(termoBuscaLowerCase))
            );
            setUsuariosFiltrados(resultadosFiltrados);
        }
        setPage(0); // Voltar para a primeira página quando filtrar
    }, [termoBusca, usuarios]);

    // Função para buscar usuários da API
    const buscarUsuarios = async () => {
        setCarregando(true);
        try {
            const data = await usuarioService.listarUsuarios();
            setUsuarios(data);
            setUsuariosFiltrados(data);
        } catch (error) {
            console.error('Erro ao buscar usuários:', error);
            setErro('Não foi possível carregar os usuários. Tente novamente mais tarde.');
        } finally {
            setCarregando(false);
        }
    };

    // Funções para manipular o formulário
    const abrirFormulario = (usuario?: Usuario) => {
        if (usuario) {
            // Edição - não preenchemos a senha para edição
            setFormulario({
                nome: usuario.nome,
                login: usuario.login,
                senha: '',
                confirmaSenha: '',
                cpf: usuario.cpf || '',
                email: usuario.email || '',
                perfil: usuario.perfil,
            });
            setEditandoUsuario(usuario.id);
        } else {
            // Novo usuário
            setFormulario({
                nome: '',
                login: '',
                senha: '',
                confirmaSenha: '',
                cpf: '',
                email: '',
                perfil: '',
            });
            setEditandoUsuario(null);
        }
        setDialogoAberto(true);
    };

    const fecharFormulario = () => {
        setDialogoAberto(false);
        setErro(null);
    };

    // Handler para campos de texto
    const handleTextFieldChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setFormulario({
            ...formulario,
            [name]: value,
        });
    };

    // Handler específico para o Select de perfil
    const handleSelectChange = (e: SelectChangeEvent<string>) => {
        const { name, value } = e.target;
        setFormulario({
            ...formulario,
            [name]: value,
        });
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

    // Validação do formulário
    const validarFormulario = (): boolean => {
        if (!formulario.nome.trim()) {
            setErro('O nome do usuário é obrigatório');
            return false;
        }
        if (!formulario.login.trim()) {
            setErro('O login é obrigatório');
            return false;
        }
        if (!editandoUsuario && !formulario.senha) {
            setErro('A senha é obrigatória para novos usuários');
            return false;
        }
        if (formulario.senha && formulario.senha.length < 6) {
            setErro('A senha deve ter pelo menos 6 caracteres');
            return false;
        }
        if (formulario.senha && formulario.senha !== formulario.confirmaSenha) {
            setErro('As senhas não conferem');
            return false;
        }
        if (!formulario.perfil) {
            setErro('O perfil é obrigatório');
            return false;
        }
        return true;
    };

    // Salvar usuário (criar ou atualizar)
    const salvarUsuario = async () => {
        if (!validarFormulario()) return;

        setCarregando(true);
        setErro(null);

        try {
            if (editandoUsuario) {
                // Atualizar usuário existente
                const usuarioData: AtualizacaoUsuario = {
                    nome: formulario.nome,
                    login: formulario.login,
                    cpf: formulario.cpf || undefined,
                    email: formulario.email || undefined,
                    perfil: formulario.perfil,
                };

                // Adicionar senha apenas se foi preenchida
                if (formulario.senha) {
                    usuarioData.senha = formulario.senha;
                }

                await usuarioService.atualizarUsuario(editandoUsuario, usuarioData);
                setSucesso('Usuário atualizado com sucesso!');
            } else {
                // Criar novo usuário
                const usuarioData: NovoUsuario = {
                    nome: formulario.nome,
                    login: formulario.login,
                    senha: formulario.senha,
                    cpf: formulario.cpf || undefined,
                    email: formulario.email || undefined,
                    perfil: formulario.perfil,
                };

                await usuarioService.criarUsuario(usuarioData);
                setSucesso('Usuário criado com sucesso!');
            }

            // Fechar o formulário e atualizar a lista
            fecharFormulario();
            buscarUsuarios();
        } catch (error: any) {
            console.error('Erro ao salvar usuário:', error);
            setErro(error.response?.data || 'Ocorreu um erro ao salvar o usuário. Tente novamente.');
        } finally {
            setCarregando(false);
        }
    };

    // Funções para exclusão de usuário
    const confirmarExclusao = (id: number) => {
        setUsuarioParaExcluir(id);
        setDialogoExclusaoAberto(true);
    };

    const fecharDialogoExclusao = () => {
        setDialogoExclusaoAberto(false);
        setUsuarioParaExcluir(null);
    };

    const excluirUsuario = async () => {
        if (!usuarioParaExcluir) return;

        setCarregando(true);
        try {
            await usuarioService.excluirUsuario(usuarioParaExcluir);

            setSucesso('Usuário excluído com sucesso!');
            fecharDialogoExclusao();
            buscarUsuarios();
        } catch (error: any) {
            console.error('Erro ao excluir usuário:', error);
            setErro(error.response?.data || 'Ocorreu um erro ao excluir o usuário. Tente novamente.');
            fecharDialogoExclusao();
        } finally {
            setCarregando(false);
        }
    };

    // Função para obter o rótulo do perfil
    const obterRotuloPerfil = (valor: string): string => {
        const perfil = perfis.find(p => p.valor === valor);
        return perfil ? perfil.rotulo : valor;
    };

    // Fechar alertas
    const fecharAlerta = () => {
        setErro(null);
        setSucesso(null);
    };

    // Calcular usuários para a página atual
    const usuariosPaginados = usuariosFiltrados.slice(
        page * rowsPerPage,
        page * rowsPerPage + rowsPerPage
    );

    return (
        <Box>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
                <Typography variant="h4">
                    Usuários
                </Typography>
                {isAdmin && (
                    <Button
                        variant="contained"
                        color="primary"
                        startIcon={<AddIcon />}
                        onClick={() => abrirFormulario()}
                    >
                        Novo Usuário
                    </Button>
                )}
            </Box>

            {/* Campo de busca */}
            <Paper sx={{ p: 2, mb: 2 }}>
                <TextField
                    fullWidth
                    placeholder="Buscar usuários pelo nome, login ou email..."
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

            {/* Tabela de usuários */}
            <Paper sx={{ width: '100%', overflow: 'hidden' }}>
                <TableContainer>
                    <Table>
                        <TableHead>
                            <TableRow>
                                <TableCell>Nome</TableCell>
                                <TableCell>Login</TableCell>
                                <TableCell>Email</TableCell>
                                <TableCell>Perfil</TableCell>
                                <TableCell align="right">Ações</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {carregando && usuariosFiltrados.length === 0 ? (
                                <TableRow>
                                    <TableCell colSpan={5} align="center">
                                        <CircularProgress size={30} />
                                    </TableCell>
                                </TableRow>
                            ) : usuariosFiltrados.length === 0 ? (
                                <TableRow>
                                    <TableCell colSpan={5} align="center">
                                        {termoBusca ? 'Nenhum usuário encontrado para a busca' : 'Nenhum usuário cadastrado'}
                                    </TableCell>
                                </TableRow>
                            ) : usuariosPaginados.length === 0 ? (
                                <TableRow>
                                    <TableCell colSpan={5} align="center">
                                        Não há usuários nesta página
                                    </TableCell>
                                </TableRow>
                            ) : (
                                usuariosPaginados.map((usuario) => (
                                    <TableRow key={usuario.id}>
                                        <TableCell>{usuario.nome}</TableCell>
                                        <TableCell>{usuario.login}</TableCell>
                                        <TableCell>{usuario.email || '-'}</TableCell>
                                        <TableCell>
                                            <Chip
                                                label={obterRotuloPerfil(usuario.perfil)}
                                                color={
                                                    usuario.perfil === 'admin' ? 'error' :
                                                        usuario.perfil === 'gerente' ? 'warning' :
                                                            usuario.perfil === 'atendente' ? 'success' :
                                                                'primary'
                                                }
                                                size="small"
                                            />
                                        </TableCell>
                                        <TableCell align="right">
                                            <IconButton
                                                color="primary"
                                                onClick={() => abrirFormulario(usuario)}
                                                disabled={!isAdmin && usuario.id !== usuarioAtual?.id}
                                            >
                                                <EditIcon />
                                            </IconButton>
                                            {isAdmin && usuario.id !== usuarioAtual?.id && (
                                                <IconButton color="error" onClick={() => confirmarExclusao(usuario.id)}>
                                                    <DeleteIcon />
                                                </IconButton>
                                            )}
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
                    count={usuariosFiltrados.length}
                    rowsPerPage={rowsPerPage}
                    page={usuariosFiltrados.length <= page * rowsPerPage && page > 0 ? page - 1 : page}
                    onPageChange={handleChangePage}
                    onRowsPerPageChange={handleChangeRowsPerPage}
                    labelRowsPerPage="Itens por página:"
                    labelDisplayedRows={({ from, to, count }) => `${from}-${to} de ${count}`}
                />
            </Paper>

            {/* Diálogo de criação/edição de usuário */}
            <Dialog open={dialogoAberto} onClose={fecharFormulario} maxWidth="sm" fullWidth>
                <DialogTitle>
                    {editandoUsuario ? 'Editar Usuário' : 'Novo Usuário'}
                </DialogTitle>
                <DialogContent>
                    <Box component="form" sx={{ mt: 1 }}>
                        <TextField
                            margin="normal"
                            required
                            fullWidth
                            label="Nome"
                            name="nome"
                            value={formulario.nome}
                            onChange={handleTextFieldChange}
                        />

                        <TextField
                            margin="normal"
                            required
                            fullWidth
                            label="Login"
                            name="login"
                            value={formulario.login}
                            onChange={handleTextFieldChange}
                        />

                        <TextField
                            margin="normal"
                            required={!editandoUsuario}
                            fullWidth
                            label="Senha"
                            name="senha"
                            type="password"
                            value={formulario.senha}
                            onChange={handleTextFieldChange}
                            helperText={editandoUsuario ? "Deixe em branco para manter a senha atual" : ""}
                        />

                        <TextField
                            margin="normal"
                            required={!editandoUsuario}
                            fullWidth
                            label="Confirmar Senha"
                            name="confirmaSenha"
                            type="password"
                            value={formulario.confirmaSenha}
                            onChange={handleTextFieldChange}
                        />

                        <TextField
                            margin="normal"
                            fullWidth
                            label="CPF"
                            name="cpf"
                            value={formulario.cpf}
                            onChange={handleTextFieldChange}
                        />

                        <TextField
                            margin="normal"
                            fullWidth
                            label="Email"
                            name="email"
                            type="email"
                            value={formulario.email}
                            onChange={handleTextFieldChange}
                        />

                        <FormControl fullWidth margin="normal">
                            <InputLabel id="perfil-label">Perfil</InputLabel>
                            <Select
                                labelId="perfil-label"
                                name="perfil"
                                value={formulario.perfil}
                                label="Perfil"
                                onChange={handleSelectChange}
                                disabled={!isAdmin}
                            >
                                {perfis.map((perfil) => (
                                    <MenuItem key={perfil.valor} value={perfil.valor}>
                                        {perfil.rotulo}
                                    </MenuItem>
                                ))}
                            </Select>
                        </FormControl>
                    </Box>
                </DialogContent>
                <DialogActions>
                    <Button onClick={fecharFormulario}>Cancelar</Button>
                    <Button
                        onClick={salvarUsuario}
                        variant="contained"
                        color="primary"
                        disabled={carregando}
                    >
                        {carregando ? <CircularProgress size={24} /> : 'Salvar'}
                    </Button>
                </DialogActions>
            </Dialog>

            {/* Diálogo de confirmação de exclusão */}
            <Dialog
                open={dialogoExclusaoAberto}
                onClose={fecharDialogoExclusao}
            >
                <DialogTitle>Confirmar Exclusão</DialogTitle>
                <DialogContent>
                    <Typography>
                        Tem certeza que deseja excluir este usuário? Esta ação não pode ser desfeita.
                    </Typography>
                </DialogContent>
                <DialogActions>
                    <Button onClick={fecharDialogoExclusao}>Cancelar</Button>
                    <Button
                        onClick={excluirUsuario}
                        color="error"
                        variant="contained"
                        disabled={carregando}
                    >
                        {carregando ? <CircularProgress size={24} /> : 'Excluir'}
                    </Button>
                </DialogActions>
            </Dialog>

            {/* Alertas de sucesso e erro */}
            <Snackbar open={!!erro} autoHideDuration={6000} onClose={fecharAlerta}>
                <Alert onClose={fecharAlerta} severity="error" sx={{ width: '100%' }}>
                    {erro}
                </Alert>
            </Snackbar>

            <Snackbar open={!!sucesso} autoHideDuration={6000} onClose={fecharAlerta}>
                <Alert onClose={fecharAlerta} severity="success" sx={{ width: '100%' }}>
                    {sucesso}
                </Alert>
            </Snackbar>
        </Box>
    );
};

export default Usuarios;