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
} from '@mui/material';
import {
  Add as AddIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  Search as SearchIcon,
  Clear as ClearIcon,
} from '@mui/icons-material';
import { produtoService, Produto } from '../services/produtoService';

// Interface para o formulário
interface FormularioProduto {
  nome: string;
  descricao: string;
  categoria: string;
  preco: string; // String para facilitar a manipulação no formulário
}

// Categorias de produtos disponíveis
const categorias = [
  { valor: 'botija_gas', rotulo: 'Botija de Gás' },
  { valor: 'agua', rotulo: 'Água Mineral' },
  { valor: 'acessorio', rotulo: 'Acessórios' },
  { valor: 'outros', rotulo: 'Outros' },
];

const Produtos = () => {
  // Estados
  const [produtos, setProdutos] = useState<Produto[]>([]);
  const [produtosFiltrados, setProdutosFiltrados] = useState<Produto[]>([]);
  const [dialogoAberto, setDialogoAberto] = useState(false);
  const [formulario, setFormulario] = useState<FormularioProduto>({
    nome: '',
    descricao: '',
    categoria: '',
    preco: '',
  });
  const [editandoProduto, setEditandoProduto] = useState<number | null>(null);
  const [carregando, setCarregando] = useState(false);
  const [erro, setErro] = useState<string | null>(null);
  const [sucesso, setSucesso] = useState<string | null>(null);
  const [dialogoExclusaoAberto, setDialogoExclusaoAberto] = useState(false);
  const [produtoParaExcluir, setProdutoParaExcluir] = useState<number | null>(null);
  
  // Estados para paginação e busca
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [termoBusca, setTermoBusca] = useState('');

  // Buscar produtos ao carregar a página
  useEffect(() => {
    buscarProdutos();
  }, []);

  // Filtrar produtos baseado no termo de busca
  useEffect(() => {
    if (termoBusca.trim() === '') {
      setProdutosFiltrados(produtos);
    } else {
      const termoBuscaLowerCase = termoBusca.toLowerCase();
      const resultadosFiltrados = produtos.filter(produto => 
        produto.nome.toLowerCase().includes(termoBuscaLowerCase) ||
        produto.descricao.toLowerCase().includes(termoBuscaLowerCase)
      );
      setProdutosFiltrados(resultadosFiltrados);
    }
    setPage(0); // Voltar para a primeira página quando filtrar
  }, [termoBusca, produtos]);

  // Função para buscar produtos da API
  const buscarProdutos = async () => {
    setCarregando(true);
    try {
      const data = await produtoService.listarProdutos();
      setProdutos(data);
      setProdutosFiltrados(data);
    } catch (error) {
      console.error('Erro ao buscar produtos:', error);
      setErro('Não foi possível carregar os produtos. Tente novamente mais tarde.');
    } finally {
      setCarregando(false);
    }
  };

  // Funções para manipular o formulário
  const abrirFormulario = (produto?: Produto) => {
    if (produto) {
      // Edição
      setFormulario({
        nome: produto.nome,
        descricao: produto.descricao || '',
        categoria: produto.categoria,
        preco: produto.preco.toString(),
      });
      setEditandoProduto(produto.id);
    } else {
      // Novo produto
      setFormulario({
        nome: '',
        descricao: '',
        categoria: '',
        preco: '',
      });
      setEditandoProduto(null);
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

  // Handler específico para o Select de categoria
  const handleSelectChange = (e: SelectChangeEvent<string>) => {
    const { name, value } = e.target;
    setFormulario({
      ...formulario,
      [name]: value,
    });
  };

  // Handlers para paginação
  const handleChangePage = (event: unknown, newPage: number) => {
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
      setErro('O nome do produto é obrigatório');
      return false;
    }
    if (!formulario.categoria) {
      setErro('A categoria é obrigatória');
      return false;
    }
    if (!formulario.preco || isNaN(parseFloat(formulario.preco)) || parseFloat(formulario.preco) <= 0) {
      setErro('O preço deve ser um número maior que zero');
      return false;
    }
    return true;
  };

  // Salvar produto (criar ou atualizar)
  const salvarProduto = async () => {
    if (!validarFormulario()) return;

    setCarregando(true);
    setErro(null);
    
    try {
      const produtoData = {
        nome: formulario.nome,
        descricao: formulario.descricao,
        categoria: formulario.categoria,
        preco: parseFloat(formulario.preco),
      };

      if (editandoProduto) {
        // Atualizar produto existente
        await produtoService.atualizarProduto(editandoProduto, produtoData);
        setSucesso('Produto atualizado com sucesso!');
      } else {
        // Criar novo produto
        await produtoService.criarProduto(produtoData);
        setSucesso('Produto criado com sucesso!');
      }

      // Fechar o formulário e atualizar a lista
      fecharFormulario();
      buscarProdutos();
    } catch (error) {
      console.error('Erro ao salvar produto:', error);
      setErro('Ocorreu um erro ao salvar o produto. Tente novamente.');
    } finally {
      setCarregando(false);
    }
  };

  // Funções para exclusão de produto
  const confirmarExclusao = (id: number) => {
    setProdutoParaExcluir(id);
    setDialogoExclusaoAberto(true);
  };

  const fecharDialogoExclusao = () => {
    setDialogoExclusaoAberto(false);
    setProdutoParaExcluir(null);
  };

  const excluirProduto = async () => {
    if (!produtoParaExcluir) return;

    setCarregando(true);
    try {
      await produtoService.excluirProduto(produtoParaExcluir);
      
      setSucesso('Produto excluído com sucesso!');
      fecharDialogoExclusao();
      buscarProdutos();
    } catch (error) {
      console.error('Erro ao excluir produto:', error);
      setErro('Ocorreu um erro ao excluir o produto. Tente novamente.');
      fecharDialogoExclusao();
    } finally {
      setCarregando(false);
    }
  };

  // Função para obter o rótulo da categoria
  const obterRotuloCategoria = (valor: string): string => {
    const categoria = categorias.find(cat => cat.valor === valor);
    return categoria ? categoria.rotulo : valor;
  };

  // Fechar alertas
  const fecharAlerta = () => {
    setErro(null);
    setSucesso(null);
  };

  // Calcular produtos para a página atual
  const produtosPaginados = produtosFiltrados.slice(
    page * rowsPerPage,
    page * rowsPerPage + rowsPerPage
  );

  return (
    <Box>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Typography variant="h4">
          Produtos
        </Typography>
        <Button
          variant="contained"
          color="primary"
          startIcon={<AddIcon />}
          onClick={() => abrirFormulario()}
        >
          Novo Produto
        </Button>
      </Box>

      {/* Campo de busca */}
      <Paper sx={{ p: 2, mb: 2 }}>
        <TextField
          fullWidth
          placeholder="Buscar produtos pelo nome ou descrição..."
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

      {/* Tabela de produtos */}
      <Paper sx={{ width: '100%', overflow: 'hidden' }}>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Nome</TableCell>
                <TableCell>Categoria</TableCell>
                <TableCell>Preço</TableCell>
                <TableCell align="right">Ações</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {carregando && produtosFiltrados.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={4} align="center">
                    <CircularProgress size={30} />
                  </TableCell>
                </TableRow>
              ) : produtosFiltrados.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={4} align="center">
                    {termoBusca ? 'Nenhum produto encontrado para a busca' : 'Nenhum produto cadastrado'}
                  </TableCell>
                </TableRow>
              ) : produtosPaginados.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={4} align="center">
                    Não há produtos nesta página
                  </TableCell>
                </TableRow>
              ) : (
                produtosPaginados.map((produto) => (
                  <TableRow key={produto.id}>
                    <TableCell>{produto.nome}</TableCell>
                    <TableCell>{obterRotuloCategoria(produto.categoria)}</TableCell>
                    <TableCell>R$ {produto.preco.toFixed(2)}</TableCell>
                    <TableCell align="right">
                      <IconButton color="primary" onClick={() => abrirFormulario(produto)}>
                        <EditIcon />
                      </IconButton>
                      <IconButton color="error" onClick={() => confirmarExclusao(produto.id)}>
                        <DeleteIcon />
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
          count={produtosFiltrados.length}
          rowsPerPage={rowsPerPage}
          page={produtosFiltrados.length <= page * rowsPerPage && page > 0 ? page - 1 : page}
          onPageChange={handleChangePage}
          onRowsPerPageChange={handleChangeRowsPerPage}
          labelRowsPerPage="Itens por página:"
          labelDisplayedRows={({ from, to, count }) => `${from}-${to} de ${count}`}
        />
      </Paper>

      {/* Diálogo de criação/edição de produto */}
      <Dialog open={dialogoAberto} onClose={fecharFormulario} maxWidth="sm" fullWidth>
        <DialogTitle>
          {editandoProduto ? 'Editar Produto' : 'Novo Produto'}
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
              fullWidth
              label="Descrição"
              name="descricao"
              multiline
              rows={3}
              value={formulario.descricao}
              onChange={handleTextFieldChange}
            />
            
            <FormControl fullWidth margin="normal">
              <InputLabel id="categoria-label">Categoria</InputLabel>
              <Select
                labelId="categoria-label"
                name="categoria"
                value={formulario.categoria}
                label="Categoria"
                onChange={handleSelectChange}
              >
                {categorias.map((categoria) => (
                  <MenuItem key={categoria.valor} value={categoria.valor}>
                    {categoria.rotulo}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
            
            <TextField
              margin="normal"
              required
              fullWidth
              label="Preço (R$)"
              name="preco"
              type="number"
              value={formulario.preco}
              onChange={handleTextFieldChange}
              InputProps={{
                inputProps: { min: 0, step: 0.01 }
              }}
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={fecharFormulario}>Cancelar</Button>
          <Button 
            onClick={salvarProduto} 
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
            Tem certeza que deseja excluir este produto? Esta ação não pode ser desfeita.
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={fecharDialogoExclusao}>Cancelar</Button>
          <Button 
            onClick={excluirProduto} 
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

export default Produtos;