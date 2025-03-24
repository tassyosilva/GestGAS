import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Paper,
  Grid,
  Card,
  CardContent,
  CardActions,
  Button,
  Chip,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TextField,
  MenuItem,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  CircularProgress,
  Alert,
  Snackbar,
  IconButton,
  Tabs,
  Tab,
} from '@mui/material';
import {
  Refresh as RefreshIcon,
  Add as AddIcon,
  Remove as RemoveIcon,
  Edit as EditIcon,
  Warning as WarningIcon,
} from '@mui/icons-material';
import axios from 'axios';
import API_BASE_URL from '../config/api';

// Interfaces
interface EstoqueItem {
  id: number;
  produto_id: number;
  nome_produto: string;
  categoria: string;
  quantidade: number;
  botijas_vazias?: number;
  botijas_emprestadas?: number;
  alerta_minimo?: number;
  status: string;
  atualizado_em: string;
}

interface AlertaEstoque {
  produto_id: number;
  nome_produto: string;
  quantidade: number;
  alerta_minimo: number;
  status: string;
}

interface MovimentacaoEstoqueRequest {
  tipo: string;
  quantidade: number;
  observacoes?: string;
}

// Componente principal
const Estoque: React.FC = () => {
  // Estados
  const [itensEstoque, setItensEstoque] = useState<EstoqueItem[]>([]);
  const [alertas, setAlertas] = useState<AlertaEstoque[]>([]);
  const [categoriaFiltro, setCategoriaFiltro] = useState<string>('');
  const [loading, setLoading] = useState(true);
  const [loadingAlertas, setLoadingAlertas] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Estado para diálogo de movimentação
  const [dialogoAberto, setDialogoAberto] = useState(false);
  const [produtoSelecionado, setProdutoSelecionado] = useState<EstoqueItem | null>(null);
  const [tipoMovimentacao, setTipoMovimentacao] = useState<string>('entrada');
  const [quantidade, setQuantidade] = useState<number>(1);
  const [observacoes, setObservacoes] = useState<string>('');
  const [atualizando, setAtualizando] = useState(false);

  // Estado para diálogo de alerta mínimo
  const [dialogoAlertaAberto, setDialogoAlertaAberto] = useState(false);
  const [alertaMinimo, setAlertaMinimo] = useState<number>(5);

  // Estado para diálogo de empréstimo/devolução de botijas
  const [dialogoBotijasAberto, setDialogoBotijasAberto] = useState(false);
  const [quantidadeBotijas, setQuantidadeBotijas] = useState<number>(1);
  const [tipoBotijasOperacao, setTipoBotijasOperacao] = useState<string>('emprestimo');

  // Estado para tabs
  const [tabAtiva, setTabAtiva] = useState(0);

  // Estado para snackbar
  const [snackbar, setSnackbar] = useState({
    open: false,
    message: '',
    severity: 'success' as 'success' | 'error' | 'warning' | 'info',
  });

  // Carregar dados do estoque
  const carregarEstoque = async () => {
    setLoading(true);
    setLoadingAlertas(true);
    setError(null);

    try {
      const token = localStorage.getItem('token');
      if (!token) {
        setError('Não autorizado. Faça login para continuar.');
        setLoading(false);
        setLoadingAlertas(false);
        setItensEstoque([]);
        setAlertas([]);
        return;
      }

      // Parâmetros de consulta para filtrar por categoria
      const params = new URLSearchParams();
      if (categoriaFiltro) {
        params.append('categoria', categoriaFiltro);
      }

      // Buscar itens do estoque
      const response = await axios.get<EstoqueItem[]>(
        `${API_BASE_URL}/estoque?${params.toString()}`,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      setItensEstoque(response.data || []);

      // Buscar alertas de estoque
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
      } catch (alertErr) {
        console.error('Erro ao buscar alertas de estoque:', alertErr);
        setAlertas([]);
      } finally {
        setLoadingAlertas(false);
      }
    } catch (err) {
      console.error('Erro ao carregar estoque:', err);
      setError('Não foi possível carregar os dados do estoque. Tente novamente mais tarde.');
      setItensEstoque([]);
      setAlertas([]);
    } finally {
      setLoading(false);
    }
  };

  // Carregar dados ao montar o componente
  useEffect(() => {
    carregarEstoque();
  }, [categoriaFiltro]);

  // Abrir diálogo para movimentação
  const handleAbrirDialogoMovimentacao = (produto: EstoqueItem, tipo: string = 'entrada') => {
    setProdutoSelecionado(produto);
    setTipoMovimentacao(tipo);
    setQuantidade(1);
    setObservacoes('');
    setDialogoAberto(true);
  };

  // Abrir diálogo para alerta mínimo
  const handleAbrirDialogoAlerta = (produto: EstoqueItem) => {
    setProdutoSelecionado(produto);
    setAlertaMinimo(produto.alerta_minimo || 5);
    setDialogoAlertaAberto(true);
  };

  // Abrir diálogo para empréstimo/devolução de botijas
  const handleAbrirDialogoBotijas = (produto: EstoqueItem, tipo: string) => {
    setProdutoSelecionado(produto);
    setQuantidadeBotijas(1);
    setObservacoes('');
    setTipoBotijasOperacao(tipo);
    setDialogoBotijasAberto(true);
  };

  // Realizar movimentação de estoque
  const handleRealizarMovimentacao = async () => {
    if (!produtoSelecionado) return;

    setAtualizando(true);

    try {
      const token = localStorage.getItem('token');
      if (!token) throw new Error('Não autorizado.');

      const payload: MovimentacaoEstoqueRequest = {
        tipo: tipoMovimentacao,
        quantidade: quantidade,
      };

      if (observacoes) {
        payload.observacoes = observacoes;
      }

      await axios.put(
        `${API_BASE_URL}/estoque/${produtoSelecionado.produto_id}`,
        payload,
        {
          headers: {
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/json',
          },
        }
      );

      // Atualizar a tela
      await carregarEstoque();

      // Fechar diálogo
      setDialogoAberto(false);

      // Mostrar mensagem de sucesso
      setSnackbar({
        open: true,
        message: 'Movimentação registrada com sucesso!',
        severity: 'success',
      });
    } catch (err: any) {
      console.error('Erro ao realizar movimentação:', err);

      setSnackbar({
        open: true,
        message: err.response?.data || 'Erro ao realizar movimentação. Tente novamente.',
        severity: 'error',
      });
    } finally {
      setAtualizando(false);
    }
  };

  // Atualizar alerta mínimo
  const handleAtualizarAlertaMinimo = async () => {
    if (!produtoSelecionado) return;

    setAtualizando(true);

    try {
      const token = localStorage.getItem('token');
      if (!token) throw new Error('Não autorizado.');

      await axios.put(
        `${API_BASE_URL}/estoque/${produtoSelecionado.produto_id}/alerta`,
        { alerta_minimo: alertaMinimo },
        {
          headers: {
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/json',
          },
        }
      );

      // Atualizar a tela
      await carregarEstoque();

      // Fechar diálogo
      setDialogoAlertaAberto(false);

      // Mostrar mensagem de sucesso
      setSnackbar({
        open: true,
        message: 'Alerta mínimo atualizado com sucesso!',
        severity: 'success',
      });
    } catch (err: any) {
      console.error('Erro ao atualizar alerta mínimo:', err);

      setSnackbar({
        open: true,
        message: err.response?.data || 'Erro ao atualizar alerta mínimo. Tente novamente.',
        severity: 'error',
      });
    } finally {
      setAtualizando(false);
    }
  };

  // Realizar operação de botijas
  const handleOperacaoBotijas = async () => {
    if (!produtoSelecionado) return;

    setAtualizando(true);

    try {
      const token = localStorage.getItem('token');
      if (!token) throw new Error('Não autorizado.');

      const endpoint = tipoBotijasOperacao === 'emprestimo'
        ? `${API_BASE_URL}/estoque/botijas/emprestimo`
        : `${API_BASE_URL}/estoque/botijas/devolucao`;

      const payload = {
        produto_id: produtoSelecionado.produto_id,
        quantidade: quantidadeBotijas,
        observacoes: observacoes,
      };

      await axios.post(
        endpoint,
        payload,
        {
          headers: {
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/json',
          },
        }
      );

      // Atualizar a tela
      await carregarEstoque();

      // Fechar diálogo
      setDialogoBotijasAberto(false);

      // Mostrar mensagem de sucesso
      const mensagem = tipoBotijasOperacao === 'emprestimo'
        ? 'Empréstimo de botijas registrado com sucesso!'
        : 'Devolução de botijas registrada com sucesso!';

      setSnackbar({
        open: true,
        message: mensagem,
        severity: 'success',
      });
    } catch (err: any) {
      console.error('Erro na operação de botijas:', err);

      setSnackbar({
        open: true,
        message: err.response?.data || 'Erro na operação. Tente novamente.',
        severity: 'error',
      });
    } finally {
      setAtualizando(false);
    }
  };

  // Formatadores
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

  // Helper para traduzir categorias para exibição
  const traduzirCategoria = (categoria: string): string => {
    const categorias: { [key: string]: string } = {
      'botija_gas': 'Botija de Gás',
      'agua': 'Água Mineral',
      'acessorio': 'Acessório',
    };

    return categorias[categoria] || categoria;
  };

  // Filtrar itens por categoria da botija
  const botijasGas = itensEstoque ? itensEstoque.filter(item => item.categoria.includes('botija_gas')) : [];

  // Renderizar chip de status com cores
  const renderizarStatusChip = (status: string) => {
    let color: 'default' | 'primary' | 'secondary' | 'error' | 'info' | 'success' | 'warning' = 'default';
    let label = status;

    switch (status) {
      case 'normal':
        color = 'success';
        label = 'Normal';
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
        <Typography variant="h4">Controle de Estoque</Typography>
        <Box>
          <Button
            variant="contained"
            color="primary"
            startIcon={<RefreshIcon />}
            onClick={carregarEstoque}
          >
            Atualizar
          </Button>
        </Box>
      </Box>

      {/* Alertas */}
      {alertas && alertas.length > 0 && (
        <Alert
          severity="warning"
          icon={<WarningIcon />}
          sx={{ mb: 3 }}
          action={
            <Button
              color="inherit"
              size="small"
              onClick={() => setTabAtiva(2)}
            >
              Ver Todos
            </Button>
          }
        >
          {alertas.length === 1
            ? `${alertas[0].nome_produto} está com estoque ${alertas[0].status}.`
            : `${alertas.length} produtos estão com estoque baixo ou crítico.`}
        </Alert>
      )}

      {/* Filtros e Tabs */}
      <Paper sx={{ mb: 3 }}>
        <Box display="flex" justifyContent="space-between" alignItems="center" p={2}>
          <Tabs
            value={tabAtiva}
            onChange={(_e, newValue) => setTabAtiva(newValue)}
          >
            <Tab label="Visão Geral" />
            <Tab label="Botijas de Gás" />
            <Tab label="Alertas" />
          </Tabs>

          <TextField
            select
            label="Categoria"
            value={categoriaFiltro}
            onChange={(e) => setCategoriaFiltro(e.target.value)}
            size="small"
            sx={{ width: 200 }}
          >
            <MenuItem value="">Todas</MenuItem>
            <MenuItem value="botija_gas">Botijas de Gás</MenuItem>
            <MenuItem value="agua">Água Mineral</MenuItem>
            <MenuItem value="acessorio">Acessórios</MenuItem>
          </TextField>
        </Box>
      </Paper>

      {/* Conteúdo baseado na tab ativa */}
      {tabAtiva === 0 && (
        <Box>
          {/* Visão Geral */}
          {loading ? (
            <Box display="flex" justifyContent="center" p={4}>
              <CircularProgress />
            </Box>
          ) : error ? (
            <Alert severity="error" sx={{ mb: 3 }}>
              {error}
            </Alert>
          ) : (
            <TableContainer component={Paper}>
              <Table size="small">
                <TableHead>
                  <TableRow>
                    <TableCell>Produto</TableCell>
                    <TableCell>Categoria</TableCell>
                    <TableCell align="center">Quantidade</TableCell>
                    <TableCell align="center">Alerta Mínimo</TableCell>
                    <TableCell align="center">Status</TableCell>
                    <TableCell align="right">Última Atualização</TableCell>
                    <TableCell align="right">Ações</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {!itensEstoque || itensEstoque.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={7} align="center">
                        Nenhum item encontrado.
                      </TableCell>
                    </TableRow>
                  ) : (
                    itensEstoque.map((item) => (
                      <TableRow key={item.id}>
                        <TableCell>{item.nome_produto}</TableCell>
                        <TableCell>{traduzirCategoria(item.categoria)}</TableCell>
                        <TableCell align="center">{item.quantidade}</TableCell>
                        <TableCell align="center">{item.alerta_minimo || '-'}</TableCell>
                        <TableCell align="center">
                          {renderizarStatusChip(item.status)}
                        </TableCell>
                        <TableCell align="right">{formatarData(item.atualizado_em)}</TableCell>
                        <TableCell align="right">
                          <IconButton
                            size="small"
                            color="primary"
                            onClick={() => handleAbrirDialogoMovimentacao(item, 'entrada')}
                          >
                            <AddIcon />
                          </IconButton>
                          <IconButton
                            size="small"
                            color="error"
                            onClick={() => handleAbrirDialogoMovimentacao(item, 'saida')}
                            disabled={item.quantidade <= 0}
                          >
                            <RemoveIcon />
                          </IconButton>
                          <IconButton
                            size="small"
                            color="info"
                            onClick={() => handleAbrirDialogoAlerta(item)}
                          >
                            <EditIcon />
                          </IconButton>
                        </TableCell>
                      </TableRow>
                    ))
                  )}
                </TableBody>
              </Table>
            </TableContainer>
          )}
        </Box>
      )}

      {tabAtiva === 1 && (
        <Box>
          {/* Botijas de Gás */}
          <Grid container spacing={3}>
            {botijasGas && botijasGas.map((botija) => (
              <Grid item xs={12} sm={6} md={4} lg={3} key={botija.id}>
                <Card>
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      {botija.nome_produto}
                    </Typography>
                    <Box display="flex" justifyContent="space-between" mb={1}>
                      <Typography variant="body2" color="text.secondary">
                        Cheias:
                      </Typography>
                      <Typography variant="body1" fontWeight="bold">
                        {botija.quantidade}
                      </Typography>
                    </Box>
                    <Box display="flex" justifyContent="space-between" mb={1}>
                      <Typography variant="body2" color="text.secondary">
                        Vazias:
                      </Typography>
                      <Typography variant="body1">
                        {botija.botijas_vazias || 0}
                      </Typography>
                    </Box>
                    <Box display="flex" justifyContent="space-between" mb={1}>
                      <Typography variant="body2" color="text.secondary">
                        Emprestadas:
                      </Typography>
                      <Typography variant="body1">
                        {botija.botijas_emprestadas || 0}
                      </Typography>
                    </Box>
                    <Box display="flex" justifyContent="space-between">
                      <Typography variant="body2" color="text.secondary">
                        Status:
                      </Typography>
                      {renderizarStatusChip(botija.status)}
                    </Box>
                  </CardContent>
                  <CardActions>
                    <Button
                      size="small"
                      onClick={() => handleAbrirDialogoBotijas(botija, 'emprestimo')}
                      disabled={(botija.botijas_vazias || 0) <= 0}
                    >
                      Emprestar
                    </Button>
                    <Button
                      size="small"
                      onClick={() => handleAbrirDialogoBotijas(botija, 'devolucao')}
                      disabled={(botija.botijas_emprestadas || 0) <= 0}
                    >
                      Devolver
                    </Button>
                    <Button
                      size="small"
                      onClick={() => handleAbrirDialogoMovimentacao(botija, 'entrada')}
                    >
                      Entrada
                    </Button>
                  </CardActions>
                </Card>
              </Grid>
            ))}

            {(!botijasGas || botijasGas.length === 0) && (
              <Grid item xs={12}>
                <Alert severity="info">
                  Nenhuma botija de gás encontrada.
                </Alert>
              </Grid>
            )}
          </Grid>
        </Box>
      )}

      {tabAtiva === 2 && (
        <Box>
          {/* Alertas de Estoque */}
          {loadingAlertas ? (
            <Box display="flex" justifyContent="center" p={3}>
              <CircularProgress size={24} />
            </Box>
          ) : !alertas || alertas.length === 0 ? (
            <Alert severity="success">
              Não há alertas de estoque no momento.
            </Alert>
          ) : (
            <TableContainer component={Paper}>
              <Table size="small">
                <TableHead>
                  <TableRow>
                    <TableCell>Produto</TableCell>
                    <TableCell align="center">Quantidade Atual</TableCell>
                    <TableCell align="center">Mínimo</TableCell>
                    <TableCell align="center">Status</TableCell>
                    <TableCell align="right">Ações</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {alertas.map((alerta) => (
                    <TableRow key={alerta.produto_id}>
                      <TableCell>{alerta.nome_produto}</TableCell>
                      <TableCell align="center">{alerta.quantidade}</TableCell>
                      <TableCell align="center">{alerta.alerta_minimo}</TableCell>
                      <TableCell align="center">
                        {renderizarStatusChip(alerta.status)}
                      </TableCell>
                      <TableCell align="right">
                        <Button
                          size="small"
                          variant="contained"
                          color="primary"
                          onClick={() => {
                            const item = itensEstoque ? itensEstoque.find(i => i.produto_id === alerta.produto_id) : null;
                            if (item) {
                              handleAbrirDialogoMovimentacao(item, 'entrada');
                            }
                          }}
                        >
                          Adicionar
                        </Button>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </TableContainer>
          )}
        </Box>
      )}

      {/* Diálogo para Movimentação de Estoque */}
      <Dialog open={dialogoAberto} onClose={() => setDialogoAberto(false)}>
        <DialogTitle>
          {tipoMovimentacao === 'entrada'
            ? 'Adicionar ao Estoque'
            : tipoMovimentacao === 'saida'
              ? 'Remover do Estoque'
              : 'Ajustar Estoque'}
        </DialogTitle>
        <DialogContent>
          <Box sx={{ mt: 2 }}>
            <Typography variant="subtitle1" gutterBottom>
              Produto: {produtoSelecionado?.nome_produto}
            </Typography>
            <Typography variant="body2" color="text.secondary" gutterBottom>
              Quantidade atual: {produtoSelecionado?.quantidade}
            </Typography>

            <TextField
              select
              label="Tipo de Movimentação"
              value={tipoMovimentacao}
              onChange={(e) => setTipoMovimentacao(e.target.value)}
              fullWidth
              sx={{ mt: 2, mb: 2 }}
            >
              <MenuItem value="entrada">Entrada</MenuItem>
              <MenuItem value="saida">Saída</MenuItem>
              <MenuItem value="ajuste">Ajuste Direto</MenuItem>
              {produtoSelecionado?.categoria.includes('botija_gas') && (
                <MenuItem value="botijas_vazias">Entrada de Botijas Vazias</MenuItem>
              )}
            </TextField>

            <TextField
              label={tipoMovimentacao === 'ajuste' ? 'Nova Quantidade' : 'Quantidade'}
              type="number"
              value={quantidade}
              onChange={(e) => setQuantidade(Math.max(1, parseInt(e.target.value) || 0))}
              fullWidth
              required
              InputProps={{ inputProps: { min: 1 } }}
              sx={{ mb: 2 }}
            />

            <TextField
              label="Observações"
              value={observacoes}
              onChange={(e) => setObservacoes(e.target.value)}
              fullWidth
              multiline
              rows={3}
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDialogoAberto(false)}>Cancelar</Button>
          <Button
            onClick={handleRealizarMovimentacao}
            color="primary"
            variant="contained"
            disabled={atualizando || quantidade <= 0}
          >
            {atualizando ? <CircularProgress size={24} /> : 'Confirmar'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Diálogo para Alerta Mínimo */}
      <Dialog open={dialogoAlertaAberto} onClose={() => setDialogoAlertaAberto(false)}>
        <DialogTitle>Configurar Alerta Mínimo</DialogTitle>
        <DialogContent>
          <Box sx={{ mt: 2 }}>
            <Typography variant="subtitle1" gutterBottom>
              Produto: {produtoSelecionado?.nome_produto}
            </Typography>

            <TextField
              label="Quantidade Mínima"
              type="number"
              value={alertaMinimo}
              onChange={(e) => setAlertaMinimo(Math.max(0, parseInt(e.target.value) || 0))}
              fullWidth
              required
              InputProps={{ inputProps: { min: 0 } }}
              sx={{ mt: 2 }}
              helperText="Defina 0 para desativar o alerta"
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDialogoAlertaAberto(false)}>Cancelar</Button>
          <Button
            onClick={handleAtualizarAlertaMinimo}
            color="primary"
            variant="contained"
            disabled={atualizando}
          >
            {atualizando ? <CircularProgress size={24} /> : 'Confirmar'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Diálogo para Empréstimo/Devolução de Botijas */}
      <Dialog open={dialogoBotijasAberto} onClose={() => setDialogoBotijasAberto(false)}>
        <DialogTitle>
          {tipoBotijasOperacao === 'emprestimo'
            ? 'Emprestar Botijas ao Caminhoneiro'
            : 'Registrar Devolução de Botijas'}
        </DialogTitle>
        <DialogContent>
          <Box sx={{ mt: 2 }}>
            <Typography variant="subtitle1" gutterBottom>
              Produto: {produtoSelecionado?.nome_produto}
            </Typography>

            {tipoBotijasOperacao === 'emprestimo' ? (
              <Typography variant="body2" color="text.secondary" gutterBottom>
                Botijas vazias disponíveis: {produtoSelecionado?.botijas_vazias || 0}
              </Typography>
            ) : (
              <Typography variant="body2" color="text.secondary" gutterBottom>
                Botijas emprestadas: {produtoSelecionado?.botijas_emprestadas || 0}
              </Typography>
            )}

            <TextField
              label="Quantidade"
              type="number"
              value={quantidadeBotijas}
                            onChange={(e) => setQuantidadeBotijas(Math.max(1, parseInt(e.target.value) || 0))}
                            fullWidth
                            required
                            InputProps={{
                                inputProps: {
                                min: 1,
                                max: tipoBotijasOperacao === 'emprestimo'
                                    ? produtoSelecionado?.botijas_vazias || 0
                                    : produtoSelecionado?.botijas_emprestadas || 0,
                                },
                            }}
                            sx={{ mt: 2, mb: 2 }}
                            error={
                                (tipoBotijasOperacao === 'emprestimo' && quantidadeBotijas > (produtoSelecionado?.botijas_vazias || 0)) ||
                                (tipoBotijasOperacao === 'devolucao' && quantidadeBotijas > (produtoSelecionado?.botijas_emprestadas || 0))
                            }
                            helperText={
                                (tipoBotijasOperacao === 'emprestimo' && quantidadeBotijas > (produtoSelecionado?.botijas_vazias || 0))
                                ? 'Quantidade superior às botijas vazias disponíveis'
                                : (tipoBotijasOperacao === 'devolucao' && quantidadeBotijas > (produtoSelecionado?.botijas_emprestadas || 0))
                                ? 'Quantidade superior às botijas emprestadas'
                                : ''
                            }
                        />
                        
                        <TextField
                            label="Observações"
                            value={observacoes}
                            onChange={(e) => setObservacoes(e.target.value)}
                            fullWidth
                            multiline
                            rows={3}
                        />
                    </Box>
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setDialogoBotijasAberto(false)}>Cancelar</Button>
                    <Button
                        onClick={handleOperacaoBotijas}
                        color="primary"
                        variant="contained"
                        disabled={
                            atualizando ||
                            quantidadeBotijas <= 0 ||
                            (tipoBotijasOperacao === 'emprestimo' && quantidadeBotijas > (produtoSelecionado?.botijas_vazias || 0)) ||
                            (tipoBotijasOperacao === 'devolucao' && quantidadeBotijas > (produtoSelecionado?.botijas_emprestadas || 0))
                        }
                    >
                        {atualizando ? <CircularProgress size={24} /> : 'Confirmar'}
                    </Button>
                </DialogActions>
            </Dialog>
            
            {/* Snackbar para feedback */}
            <Snackbar
                open={snackbar.open}
                autoHideDuration={6000}
                onClose={() => setSnackbar(prev => ({ ...prev, open: false }))}
                anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
            >
                <Alert
                    onClose={() => setSnackbar(prev => ({ ...prev, open: false }))}
                    severity={snackbar.severity}
                    sx={{ width: '100%' }}
                >
                    {snackbar.message}
                </Alert>
            </Snackbar>
        </Box>
    );
};

export default Estoque;