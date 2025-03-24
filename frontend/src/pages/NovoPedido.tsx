import React, { useState, useEffect } from 'react';
import {
    Box,
    Typography,
    Paper,
    Grid,
    TextField,
    Button,
    MenuItem,
    Autocomplete,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    IconButton,
    Divider,
    CircularProgress,
    Alert,
    Snackbar,
    Checkbox,
    FormControlLabel,
} from '@mui/material';
import {
    ArrowBack as ArrowBackIcon,
    Add as AddIcon,
    Delete as DeleteIcon,
    Save as SaveIcon,
} from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';
import API_BASE_URL from '../config/api';
import axios from 'axios';

// Interfaces
interface Cliente {
    id: number;
    nome: string;
    telefone: string;
    endereco?: string;
}

interface Produto {
    id: number;
    nome: string;
    descricao?: string;
    categoria: string;
    preco: number;
}

interface ItemPedido {
    produto_id: number;
    nome_produto: string;
    quantidade: number;
    preco_unitario: number;
    subtotal: number;
    retorna_botija: boolean;
}

const formatCurrency = (value: number): string => {
    return new Intl.NumberFormat('pt-BR', {
        style: 'currency',
        currency: 'BRL'
    }).format(value);
};

const NovoPedido: React.FC = () => {
    const navigate = useNavigate();

    // Estado para o formulário do pedido
    const [formData, setFormData] = useState({
        cliente_id: 0,
        forma_pagamento: 'dinheiro',
        observacoes: '',
        endereco_entrega: '',
        canal_origem: 'whatsapp',
    });

    // Estados para clientes e produtos - inicializando com arrays vazios
    const [clientes, setClientes] = useState<Cliente[]>([]);
    const [clienteSelecionado, setClienteSelecionado] = useState<Cliente | null>(null);
    const [produtos, setProdutos] = useState<Produto[]>([]);
    const [produtoSelecionado, setProdutoSelecionado] = useState<Produto | null>(null);

    // Estado para os itens do pedido
    const [itensPedido, setItensPedido] = useState<ItemPedido[]>([]);
    const [quantidadeProduto, setQuantidadeProduto] = useState(1);
    const [retornaBotija, setRetornaBotija] = useState(false);

    // Estados para carregamento e erro
    const [carregandoClientes, setCarregandoClientes] = useState(false);
    const [carregandoProdutos, setCarregandoProdutos] = useState(false);
    const [enviandoPedido, setEnviandoPedido] = useState(false);
    const [erro, setErro] = useState<string | null>(null);

    // Estado para Snackbar
    const [snackbar, setSnackbar] = useState({
        open: false,
        message: '',
        severity: 'success' as 'success' | 'error',
    });

    // Carregar clientes ao montar o componente
    useEffect(() => {
        const buscarClientes = async () => {
            setCarregandoClientes(true);
            try {
                const token = localStorage.getItem('token');
                if (!token) {
                    setErro('Não autorizado. Faça login para continuar.');
                    setCarregandoClientes(false);
                    return;
                }

                const response = await axios.get(
                    `${API_BASE_URL}/clientes`,
                    {
                        headers: {
                            Authorization: `Bearer ${token}`,
                        },
                    }
                );

                // Extrair o array clientes do objeto de resposta paginada
                if (response.data && response.data.clientes && Array.isArray(response.data.clientes)) {
                    setClientes(response.data.clientes);
                } else if (Array.isArray(response.data)) {
                    // Caso a API retorne diretamente um array
                    setClientes(response.data);
                } else {
                    console.error('Resposta de clientes não contém um array válido:', response.data);
                    setClientes([]);
                }
            } catch (err) {
                console.error('Erro ao buscar clientes:', err);
                setErro('Não foi possível carregar os clientes. Tente novamente mais tarde.');
                setClientes([]);
            } finally {
                setCarregandoClientes(false);
            }
        };

        buscarClientes();
    }, []);

    // Carregar produtos ao montar o componente
    useEffect(() => {
        const buscarProdutos = async () => {
            setCarregandoProdutos(true);
            try {
                const token = localStorage.getItem('token');
                if (!token) {
                    setErro('Não autorizado. Faça login para continuar.');
                    setCarregandoProdutos(false);
                    return;
                }

                const response = await axios.get(
                    `${API_BASE_URL}/produtos`,
                    {
                        headers: {
                            Authorization: `Bearer ${token}`,
                        },
                    }
                );

                // Extrair o array produtos do objeto de resposta paginada
                if (response.data && response.data.produtos && Array.isArray(response.data.produtos)) {
                    setProdutos(response.data.produtos);
                } else if (Array.isArray(response.data)) {
                    // Caso a API retorne diretamente um array
                    setProdutos(response.data);
                } else {
                    console.error('Resposta de produtos não contém um array válido:', response.data);
                    setProdutos([]);
                }
            } catch (err) {
                console.error('Erro ao buscar produtos:', err);
                setErro('Não foi possível carregar os produtos. Tente novamente mais tarde.');
                setProdutos([]);
            } finally {
                setCarregandoProdutos(false);
            }
        };

        buscarProdutos();
    }, []);

    // Atualizar endereço de entrega quando o cliente é selecionado
    useEffect(() => {
        if (clienteSelecionado && clienteSelecionado.endereco) {
            setFormData(prev => ({
                ...prev,
                endereco_entrega: clienteSelecionado.endereco || '',
            }));
        }
    }, [clienteSelecionado]);

    // Função para adicionar item ao pedido
    const adicionarItem = () => {
        if (!produtoSelecionado) return;

        // Verificar se o produto já está no pedido
        const itemExistente = itensPedido.find(item => item.produto_id === produtoSelecionado.id);

        if (itemExistente) {
            // Atualizar quantidade se o produto já estiver no pedido
            setItensPedido(prevItens =>
                prevItens.map(item =>
                    item.produto_id === produtoSelecionado.id
                        ? {
                            ...item,
                            quantidade: item.quantidade + quantidadeProduto,
                            subtotal: (item.quantidade + quantidadeProduto) * item.preco_unitario,
                            retorna_botija: retornaBotija
                        }
                        : item
                )
            );
        } else {
            // Adicionar novo item ao pedido
            const novoItem: ItemPedido = {
                produto_id: produtoSelecionado.id,
                nome_produto: produtoSelecionado.nome,
                quantidade: quantidadeProduto,
                preco_unitario: produtoSelecionado.preco,
                subtotal: quantidadeProduto * produtoSelecionado.preco,
                retorna_botija: retornaBotija
            };

            setItensPedido(prevItens => [...prevItens, novoItem]);
        }

        // Resetar seleção
        setProdutoSelecionado(null);
        setQuantidadeProduto(1);
        setRetornaBotija(false);
    };

    // Função para remover item do pedido
    const removerItem = (index: number) => {
        setItensPedido(prevItens => prevItens.filter((_, i) => i !== index));
    };

    // Calcular valor total do pedido
    const valorTotal = itensPedido.reduce((total, item) => total + item.subtotal, 0);

    // Função para lidar com a mudança nos campos do formulário
    const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        const { name, value } = e.target;
        setFormData(prev => ({
            ...prev,
            [name]: value
        }));
    };

    // Função para enviar pedido
    const enviarPedido = async () => {
        // Validar dados
        if (!clienteSelecionado) {
            setSnackbar({
                open: true,
                message: 'Selecione um cliente para continuar.',
                severity: 'error',
            });
            return;
        }

        if (itensPedido.length === 0) {
            setSnackbar({
                open: true,
                message: 'Adicione pelo menos um produto ao pedido.',
                severity: 'error',
            });
            return;
        }

        if (!formData.endereco_entrega) {
            setSnackbar({
                open: true,
                message: 'Informe o endereço de entrega.',
                severity: 'error',
            });
            return;
        }

        setEnviandoPedido(true);

        try {
            const token = localStorage.getItem('token');
            if (!token) throw new Error('Não autorizado.');

            const payload = {
                cliente_id: clienteSelecionado.id,
                forma_pagamento: formData.forma_pagamento,
                observacoes: formData.observacoes,
                endereco_entrega: formData.endereco_entrega,
                canal_origem: formData.canal_origem,
                itens: itensPedido.map(item => ({
                    produto_id: item.produto_id,
                    quantidade: item.quantidade,
                    retorna_botija: item.retorna_botija
                }))
            };

            const response = await axios.post(
                `${API_BASE_URL}/pedidos/`,
                payload,
                {
                    headers: {
                        Authorization: `Bearer ${token}`,
                        'Content-Type': 'application/json',
                    },
                }
            );

            setSnackbar({
                open: true,
                message: 'Pedido criado com sucesso!',
                severity: 'success',
            });

            // Redirecionar para a página de detalhes do pedido criado
            setTimeout(() => {
                navigate(`/pedidos/${response.data.id}`);
            }, 1500);

        } catch (err: any) {
            console.error('Erro ao criar pedido:', err);

            setSnackbar({
                open: true,
                message: err.response?.data || 'Erro ao criar pedido. Tente novamente.',
                severity: 'error',
            });

            setEnviandoPedido(false);
        }
    };

    // Voltar para a lista de pedidos
    const handleVoltar = () => {
        navigate('/pedidos');
    };

    // Fechar snackbar
    const handleCloseSnackbar = () => {
        setSnackbar(prev => ({ ...prev, open: false }));
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
                <Box display="flex" alignItems="center">
                    <Button
                        startIcon={<ArrowBackIcon />}
                        onClick={handleVoltar}
                        sx={{ mr: 2 }}
                    >
                        Voltar
                    </Button>
                    <Typography variant="h4">Novo Pedido</Typography>
                </Box>

                <Button
                    variant="contained"
                    color="primary"
                    startIcon={<SaveIcon />}
                    onClick={enviarPedido}
                    disabled={enviandoPedido || itensPedido.length === 0 || !clienteSelecionado}
                >
                    {enviandoPedido ? <CircularProgress size={24} /> : 'Salvar Pedido'}
                </Button>
            </Box>

            {erro && (
                <Alert severity="error" sx={{ mb: 3 }}>
                    {erro}
                </Alert>
            )}

            <Grid container spacing={3}>
                {/* Seleção de Cliente */}
                <Grid item xs={12}>
                    <Paper sx={{ p: 3 }}>
                        <Typography variant="h6" gutterBottom>
                            Cliente
                        </Typography>

                        <Autocomplete
                            options={clientes}
                            getOptionLabel={(option) => `${option.nome} - ${option.telefone}`}
                            loading={carregandoClientes}
                            value={clienteSelecionado}
                            onChange={(_, newValue) => {
                                setClienteSelecionado(newValue);
                                if (newValue) {
                                    setFormData(prev => ({
                                        ...prev,
                                        cliente_id: newValue.id,
                                        endereco_entrega: newValue.endereco || '',
                                    }));
                                }
                            }}
                            renderInput={(params) => (
                                <TextField
                                    {...params}
                                    label="Selecione o cliente"
                                    variant="outlined"
                                    fullWidth
                                    required
                                    error={!clienteSelecionado}
                                    helperText={!clienteSelecionado ? 'Cliente é obrigatório' : ''}
                                    InputProps={{
                                        ...params.InputProps,
                                        endAdornment: (
                                            <>
                                                {carregandoClientes ? <CircularProgress color="inherit" size={20} /> : null}
                                                {params.InputProps.endAdornment}
                                            </>
                                        ),
                                    }}
                                />
                            )}
                        />
                    </Paper>
                </Grid>

                {/* Informações do Pedido */}
                <Grid item xs={12} md={6}>
                    <Paper sx={{ p: 3 }}>
                        <Typography variant="h6" gutterBottom>
                            Informações do Pedido
                        </Typography>

                        <Grid container spacing={2}>
                            <Grid item xs={12}>
                                <TextField
                                    label="Endereço de Entrega"
                                    name="endereco_entrega"
                                    value={formData.endereco_entrega}
                                    onChange={handleChange}
                                    fullWidth
                                    required
                                    multiline
                                    rows={2}
                                    error={!formData.endereco_entrega}
                                    helperText={!formData.endereco_entrega ? 'Endereço de entrega é obrigatório' : ''}
                                />
                            </Grid>

                            <Grid item xs={12} sm={6}>
                                <TextField
                                    select
                                    label="Forma de Pagamento"
                                    name="forma_pagamento"
                                    value={formData.forma_pagamento}
                                    onChange={handleChange}
                                    fullWidth
                                    required
                                >
                                    <MenuItem value="dinheiro">Dinheiro</MenuItem>
                                    <MenuItem value="cartao_credito">Cartão de Crédito</MenuItem>
                                    <MenuItem value="cartao_debito">Cartão de Débito</MenuItem>
                                    <MenuItem value="pix">PIX</MenuItem>
                                    <MenuItem value="fiado">Fiado</MenuItem>
                                </TextField>
                            </Grid>

                            <Grid item xs={12} sm={6}>
                                <TextField
                                    select
                                    label="Canal de Origem"
                                    name="canal_origem"
                                    value={formData.canal_origem}
                                    onChange={handleChange}
                                    fullWidth
                                >
                                    <MenuItem value="whatsapp">WhatsApp</MenuItem>
                                    <MenuItem value="telefone">Telefone</MenuItem>
                                    <MenuItem value="presencial">Presencial</MenuItem>
                                    <MenuItem value="aplicativo">Aplicativo</MenuItem>
                                </TextField>
                            </Grid>

                            <Grid item xs={12}>
                                <TextField
                                    label="Observações"
                                    name="observacoes"
                                    value={formData.observacoes}
                                    onChange={handleChange}
                                    fullWidth
                                    multiline
                                    rows={3}
                                />
                            </Grid>
                        </Grid>
                    </Paper>
                </Grid>

                {/* Adicionar Produtos */}
                <Grid item xs={12} md={6}>
                    <Paper sx={{ p: 3 }}>
                        <Typography variant="h6" gutterBottom>
                            Adicionar Produtos
                        </Typography>

                        <Grid container spacing={2}>
                            <Grid item xs={12}>
                                <Autocomplete
                                    options={produtos}
                                    getOptionLabel={(option) => `${option.nome} - ${formatCurrency(option.preco)}`}
                                    loading={carregandoProdutos}
                                    value={produtoSelecionado}
                                    onChange={(_, newValue) => {
                                        setProdutoSelecionado(newValue);

                                        // Se for uma botija de gás, habilitar checkbox de retorno de botija
                                        if (newValue && newValue.categoria.includes('botija_gas')) {
                                            setRetornaBotija(true);
                                        } else {
                                            setRetornaBotija(false);
                                        }
                                    }}
                                    renderInput={(params) => (
                                        <TextField
                                            {...params}
                                            label="Selecione o produto"
                                            variant="outlined"
                                            fullWidth
                                            InputProps={{
                                                ...params.InputProps,
                                                endAdornment: (
                                                    <>
                                                        {carregandoProdutos ? <CircularProgress color="inherit" size={20} /> : null}
                                                        {params.InputProps.endAdornment}
                                                    </>
                                                ),
                                            }}
                                        />
                                    )}
                                />
                            </Grid>

                            <Grid item xs={12} sm={6}>
                                <TextField
                                    label="Quantidade"
                                    type="number"
                                    value={quantidadeProduto}
                                    onChange={(e) => setQuantidadeProduto(Math.max(1, parseInt(e.target.value) || 1))}
                                    fullWidth
                                    InputProps={{ inputProps: { min: 1 } }}
                                />
                            </Grid>

                            <Grid item xs={12} sm={6}>
                                <Box display="flex" alignItems="center" height="100%">
                                    {produtoSelecionado && produtoSelecionado.categoria.includes('botija_gas') && (
                                        <FormControlLabel
                                            control={
                                                <Checkbox
                                                    checked={retornaBotija}
                                                    onChange={(e) => setRetornaBotija(e.target.checked)}
                                                />
                                            }
                                            label="Cliente vai devolver botija vazia"
                                        />
                                    )}
                                </Box>
                            </Grid>

                            <Grid item xs={12}>
                                <Button
                                    variant="contained"
                                    color="primary"
                                    startIcon={<AddIcon />}
                                    onClick={adicionarItem}
                                    disabled={!produtoSelecionado}
                                    fullWidth
                                >
                                    Adicionar ao Pedido
                                </Button>
                            </Grid>
                        </Grid>
                    </Paper>
                </Grid>

                {/* Resumo do Pedido */}
                <Grid item xs={12}>
                    <Paper>
                        <Box p={2} display="flex" justifyContent="space-between" alignItems="center">
                            <Typography variant="h6">Itens do Pedido</Typography>
                            <Typography variant="h6">
                                Total: {formatCurrency(valorTotal)}
                            </Typography>
                        </Box>
                        <Divider />
                        <TableContainer>
                            <Table>
                                <TableHead>
                                    <TableRow>
                                        <TableCell>Produto</TableCell>
                                        <TableCell align="center">Quantidade</TableCell>
                                        <TableCell align="right">Preço Unit.</TableCell>
                                        <TableCell align="right">Subtotal</TableCell>
                                        <TableCell align="center">Retorna Botija</TableCell>
                                        <TableCell align="right">Ações</TableCell>
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    {itensPedido.length === 0 ? (
                                        <TableRow>
                                            <TableCell colSpan={6} align="center">
                                                Nenhum item adicionado ao pedido.
                                            </TableCell>
                                        </TableRow>
                                    ) : (
                                        itensPedido.map((item, index) => (
                                            <TableRow key={index}>
                                                <TableCell>{item.nome_produto}</TableCell>
                                                <TableCell align="center">{item.quantidade}</TableCell>
                                                <TableCell align="right">{formatCurrency(item.preco_unitario)}</TableCell>
                                                <TableCell align="right">{formatCurrency(item.subtotal)}</TableCell>
                                                <TableCell align="center">
                                                    {item.retorna_botija ? 'Sim' : 'Não'}
                                                </TableCell>
                                                <TableCell align="right">
                                                    <IconButton
                                                        color="error"
                                                        onClick={() => removerItem(index)}
                                                    >
                                                        <DeleteIcon />
                                                    </IconButton>
                                                </TableCell>
                                            </TableRow>
                                        ))
                                    )}
                                </TableBody>
                            </Table>
                        </TableContainer>
                    </Paper>
                </Grid>
            </Grid>

            {/* Snackbar para feedback */}
            <Snackbar
                open={snackbar.open}
                autoHideDuration={6000}
                onClose={handleCloseSnackbar}
                anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
            >
                <Alert
                    onClose={handleCloseSnackbar}
                    severity={snackbar.severity}
                    sx={{ width: '100%' }}
                >
                    {snackbar.message}
                </Alert>
            </Snackbar>
        </Box>
    );
};

export default NovoPedido;