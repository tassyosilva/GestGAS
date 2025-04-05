import React, { useState, useEffect } from 'react';

import {
  Box,
  Typography,
  Paper,
  Grid,
  Chip,
  Button,
  Divider,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  CircularProgress,
  Alert,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  MenuItem,
  Snackbar,
} from '@mui/material';
import {
  ArrowBack as ArrowBackIcon,
  Check as CheckIcon,
  Cancel as CancelIcon,
  MonetizationOn as PaymentIcon,
} from '@mui/icons-material';
import { useParams, useNavigate } from 'react-router-dom';
import API_BASE_URL from '../config/api';
import axios from 'axios';

// Interfaces
interface ClienteBasico {
  id: number;
  nome: string;
  telefone: string;
}

interface UsuarioBasico {
  id: number;
  nome: string;
  perfil: string;
}

interface ItemPedido {
  id: number;
  produto_id: number;
  nome_produto: string;
  quantidade: number;
  preco_unitario: number;
  subtotal: number;
  retorna_botija: boolean;
}

interface PedidoDetalhado {
  id: number;
  cliente: ClienteBasico;
  atendente: UsuarioBasico;
  entregador?: UsuarioBasico;
  status: string;
  forma_pagamento: string;
  valor_total: number;
  observacoes?: string;
  endereco_entrega: string;
  canal_origem?: string;
  data_entrega?: string;
  motivo_cancelamento?: string;
  itens: ItemPedido[];
  criado_em: string;
  atualizado_em: string;
}

interface Entregador {
  id: number;
  nome: string;
}

const formatCurrency = (value: number): string => {
  return new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: 'BRL'
  }).format(value);
};

const formatDate = (dateString?: string): string => {
  if (!dateString) return '-';
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
interface StatusChipProps {
  status: string;
  sx?: React.CSSProperties | any;
}

const StatusChip: React.FC<StatusChipProps> = ({ status, sx }) => {
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

  return <Chip label={label} color={color} sx={sx} />;
};

// Mapeamento para exibição das formas de pagamento
const formaPagamentoMap: { [key: string]: string } = {
  dinheiro: 'Dinheiro',
  cartao_credito: 'Cartão de Crédito',
  cartao_debito: 'Cartão de Débito',
  pix: 'PIX',
  fiado: 'Fiado',
};

// Mapeamento de canais de origem
const canalOrigemMap: { [key: string]: string } = {
  whatsapp: 'WhatsApp',
  telefone: 'Telefone',
  presencial: 'Presencial',
  aplicativo: 'Aplicativo',
};

const PedidoDetalhe: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [pedido, setPedido] = useState<PedidoDetalhado | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Estados para diálogo de atualização de status
  const [dialogOpen, setDialogOpen] = useState(false);
  const [novoStatus, setNovoStatus] = useState('');
  const [entregadorId, setEntregadorId] = useState<number | ''>('');
  const [entregadores, setEntregadores] = useState<Entregador[]>([]);
  const [atualizandoStatus, setAtualizandoStatus] = useState(false);

  // Estados para diálogo de cancelamento
  const [cancelDialogOpen, setCancelDialogOpen] = useState(false);
  const [motivoCancelamento, setMotivoCancelamento] = useState('');

  // Estado para snackbar de feedback
  const [snackbar, setSnackbar] = useState({
    open: false,
    message: '',
    severity: 'success' as 'success' | 'error' | 'warning',
  });

  // Buscar dados do pedido
  const buscarPedido = async () => {
    setLoading(true);
    setError(null);

    try {
      const token = localStorage.getItem('token');
      if (!token) {
        setError('Não autorizado. Faça login para continuar.');
        setLoading(false);
        return;
      }

      const response = await axios.get<PedidoDetalhado>(
        `${API_BASE_URL}/pedidos/${id}`,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      setPedido(response.data);
    } catch (err: any) {
      console.error('Erro ao buscar pedido:', err);
      setError(err.response?.data || 'Não foi possível carregar o pedido. Tente novamente mais tarde.');
    } finally {
      setLoading(false);
    }
  };

  // Buscar lista de entregadores
  const buscarEntregadores = async () => {
    try {
      const token = localStorage.getItem('token');
      if (!token) return;

      // Supondo que haja um endpoint para listar usuários com perfil de entregador
      const response = await axios.get<Entregador[]>(
        `${API_BASE_URL}/entregadores`,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      setEntregadores(response.data);
    } catch (err) {
      console.error('Erro ao buscar entregadores:', err);
      // Alguns entregadores padrão para teste
      setEntregadores([
        { id: 1, nome: 'Roberto Entregador' },
        { id: 2, nome: 'Marcos Entregador' },
      ]);
    }
  };

  // Carregar dados ao montar o componente
  useEffect(() => {
    buscarPedido();
    buscarEntregadores();
  }, [id]);

  // Validar transição de status
  const validarTransicaoStatus = (atual: string, novo: string): boolean => {
    switch (atual) {
      case 'novo':
        return novo === 'em_preparo' || novo === 'cancelado';
      case 'em_preparo':
        return novo === 'em_entrega' || novo === 'cancelado';
      case 'em_entrega':
        return novo === 'entregue' || novo === 'cancelado';
      case 'entregue':
        return novo === 'finalizado';
      case 'cancelado':
      case 'finalizado':
        return false; // Estados finais, não podem ser alterados
      default:
        return false;
    }
  };

  // Obter status permitidos para transição
  const getStatusPermitidos = (statusAtual: string): { value: string; label: string }[] => {
    const todosStatus = [
      { value: 'novo', label: 'Novo' },
      { value: 'em_preparo', label: 'Em Preparo' },
      { value: 'em_entrega', label: 'Em Entrega' },
      { value: 'entregue', label: 'Entregue' },
      { value: 'finalizado', label: 'Finalizado' },
      { value: 'cancelado', label: 'Cancelado' },
    ];

    return todosStatus.filter(status =>
      validarTransicaoStatus(statusAtual, status.value)
    );
  };

  // Abrir diálogo para atualizar status
  const handleAbrirDialogo = () => {
    if (!pedido) return;

    // Definir o novo status para o primeiro permitido por padrão
    const statusPermitidos = getStatusPermitidos(pedido.status);
    if (statusPermitidos.length > 0) {
      setNovoStatus(statusPermitidos[0].value);
    }

    setDialogOpen(true);
  };

  // Abrir diálogo de cancelamento
  const handleAbrirDialogoCancelamento = () => {
    setCancelDialogOpen(true);
    setMotivoCancelamento('');
  };

  // Atualizar status do pedido
  const handleAtualizarStatus = async () => {
    if (!pedido || !novoStatus) return;

    setAtualizandoStatus(true);

    try {
      const token = localStorage.getItem('token');
      if (!token) throw new Error('Não autorizado.');

      const payload: any = {
        status: novoStatus,
      };

      // Adicionar entregador se estiver indo para em_entrega
      if (novoStatus === 'em_entrega' && entregadorId) {
        payload.entregador_id = entregadorId;
      }

      await axios.put(
        `${API_BASE_URL}/pedidos/${id}/status`,
        payload,
        {
          headers: {
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/json',
          },
        }
      );

      // Atualizar a tela
      await buscarPedido();

      // Mostrar mensagem de sucesso
      setSnackbar({
        open: true,
        message: 'Status atualizado com sucesso!',
        severity: 'success',
      });

      // Fechar diálogo
      setDialogOpen(false);
    } catch (err: any) {
      console.error('Erro ao atualizar status:', err);

      setSnackbar({
        open: true,
        message: err.response?.data || 'Erro ao atualizar status. Tente novamente.',
        severity: 'error',
      });
    } finally {
      setAtualizandoStatus(false);
    }
  };

  // Função para confirmar entrega e registrar botijas vazias
  const handleConfirmarEntregaCompleta = async () => {
    if (!pedido) return;
    setAtualizandoStatus(true);
    try {
      const token = localStorage.getItem('token');
      if (!token) throw new Error('Não autorizado.');

      // Primeiro confirmar a entrega
      await axios.post(
        `${API_BASE_URL}/pedidos/confirmar-entrega`,
        { pedido_id: parseInt(id as string) },
        {
          headers: {
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/json',
          },
        }
      );

      // Atualizar a tela para mostrar o pedido como entregue
      await buscarPedido();

      // Agora registrar as botijas vazias (se houver)
      // Verificar se há algum item com retorna_botija = true
      const temBotijasParaRetornar = pedido.itens.some(item => item.retorna_botija);

      if (temBotijasParaRetornar) {
        try {
          await axios.post(
            `${API_BASE_URL}/pedidos/registrar-botijas`,
            { pedido_id: parseInt(id as string) },
            {
              headers: {
                Authorization: `Bearer ${token}`,
                'Content-Type': 'application/json',
              },
            }
          );

          setSnackbar({
            open: true,
            message: 'Entrega confirmada e botijas vazias registradas com sucesso!',
            severity: 'success',
          });
        } catch (botijasErr) {
          console.error('Erro ao registrar botijas vazias:', botijasErr);
          setSnackbar({
            open: true,
            message: 'Entrega confirmada, mas houve um erro ao registrar botijas vazias.',
            severity: 'warning',
          });
        }
      } else {
        setSnackbar({
          open: true,
          message: 'Entrega confirmada com sucesso!',
          severity: 'success',
        });
      }
    } catch (err: any) {
      console.error('Erro ao confirmar entrega:', err);
      setSnackbar({
        open: true,
        message: err.response?.data || 'Erro ao confirmar entrega. Tente novamente.',
        severity: 'error',
      });
    } finally {
      setAtualizandoStatus(false);
    }
  };

  // Função para cancelar pedido e gerenciar estoque
  const handleCancelarPedidoComEstoque = async () => {
    if (!pedido || !motivoCancelamento.trim()) return;
    setAtualizandoStatus(true);
    try {
      const token = localStorage.getItem('token');
      if (!token) throw new Error('Não autorizado.');

      // Usar diretamente o endpoint de estoque, que já atualiza status e devolve estoque
      const estoquePayload = {
        pedido_id: parseInt(id as string),
        acao: 'cancelar',
        motivo_cancelamento: motivoCancelamento.trim()
      };

      await axios.post(
        `${API_BASE_URL}/pedidos/estoque`,
        estoquePayload,
        {
          headers: {
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/json',
          },
        }
      );

      // Atualizar a tela
      await buscarPedido();

      // Mostrar mensagem de sucesso
      setSnackbar({
        open: true,
        message: 'Pedido cancelado e estoque atualizado com sucesso!',
        severity: 'success',
      });

      // Fechar o diálogo de cancelamento
      setCancelDialogOpen(false);
    } catch (err: any) {
      console.error('Erro ao cancelar pedido:', err);
      setSnackbar({
        open: true,
        message: err.response?.data || 'Erro ao cancelar pedido. Tente novamente.',
        severity: 'error',
      });
    } finally {
      setAtualizandoStatus(false);
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

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
        <CircularProgress />
      </Box>
    );
  }

  if (error) {
    return (
      <Box>
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
        <Button startIcon={<ArrowBackIcon />} onClick={handleVoltar}>
          Voltar para Pedidos
        </Button>
      </Box>
    );
  }

  if (!pedido) {
    return (
      <Box>
        <Alert severity="warning">Pedido não encontrado.</Alert>
        <Button startIcon={<ArrowBackIcon />} onClick={handleVoltar} sx={{ mt: 2 }}>
          Voltar para Pedidos
        </Button>
      </Box>
    );
  }

  return (
    <Box>
      {/* Cabeçalho */}
      <Box display="flex" alignItems="center" mb={3}>
        <Button startIcon={<ArrowBackIcon />} onClick={handleVoltar} sx={{ mr: 2 }}>
          Voltar
        </Button>
        <Typography variant="h4">
          Pedido #{pedido.id}
          <StatusChip status={pedido.status} sx={{ ml: 2 }} />
        </Typography>
      </Box>

      <Grid container spacing={3}>
        {/* Informações gerais do pedido */}
        <Grid item xs={12} md={6}>
          <Paper sx={{ p: 3, height: '100%' }}>
            <Typography variant="h6" gutterBottom>
              Informações do Pedido
            </Typography>

            <Grid container spacing={2}>
              <Grid item xs={6}>
                <Typography variant="subtitle2" color="text.secondary">
                  Data de Criação
                </Typography>
                <Typography variant="body1">{formatDate(pedido.criado_em)}</Typography>
              </Grid>
              <Grid item xs={6}>
                <Typography variant="subtitle2" color="text.secondary">
                  Valor Total
                </Typography>
                <Typography variant="body1" fontWeight="bold">
                  {formatCurrency(pedido.valor_total)}
                </Typography>
              </Grid>
              <Grid item xs={6}>
                <Typography variant="subtitle2" color="text.secondary">
                  Forma de Pagamento
                </Typography>
                <Typography variant="body1">
                  {formaPagamentoMap[pedido.forma_pagamento] || pedido.forma_pagamento}
                </Typography>
              </Grid>
              <Grid item xs={6}>
                <Typography variant="subtitle2" color="text.secondary">
                  Canal de Venda
                </Typography>
                <Typography variant="body1">
                  {pedido.canal_origem ? (canalOrigemMap[pedido.canal_origem] || pedido.canal_origem) : '-'}
                </Typography>
              </Grid>
              <Grid item xs={6}>
                <Typography variant="subtitle2" color="text.secondary">
                  Data de Entrega
                </Typography>
                <Typography variant="body1">
                  {pedido.data_entrega ? formatDate(pedido.data_entrega) : 'Não entregue'}
                </Typography>
              </Grid>
              <Grid item xs={6}>
                <Typography variant="subtitle2" color="text.secondary">
                  Última Atualização
                </Typography>
                <Typography variant="body1">
                  {formatDate(pedido.atualizado_em)}
                </Typography>
              </Grid>

              {pedido.observacoes && (
                <Grid item xs={12}>
                  <Typography variant="subtitle2" color="text.secondary">
                    Observações
                  </Typography>
                  <Typography variant="body1">
                    {pedido.observacoes}
                  </Typography>
                </Grid>
              )}
            </Grid>

            <Divider sx={{ my: 2 }} />

            <Typography variant="h6" gutterBottom>
              Informações de Entrega
            </Typography>

            <Grid container spacing={2}>
              <Grid item xs={12}>
                <Typography variant="subtitle2" color="text.secondary">
                  Endereço de Entrega
                </Typography>
                <Typography variant="body1">
                  {pedido.endereco_entrega}
                </Typography>
              </Grid>
              <Grid item xs={12}>
                <Typography variant="subtitle2" color="text.secondary">
                  Entregador
                </Typography>
                <Typography variant="body1">
                  {pedido.entregador ? pedido.entregador.nome : 'Não atribuído'}
                </Typography>
              </Grid>
            </Grid>
          </Paper>
        </Grid>

        {/* Informações do cliente e responsáveis */}
        <Grid item xs={12} md={6}>
          <Paper sx={{ p: 3, height: '100%' }}>
            <Typography variant="h6" gutterBottom>
              Informações do Cliente
            </Typography>
            <Grid container spacing={2}>
              <Grid item xs={12}>
                <Typography variant="subtitle2" color="text.secondary">
                  Nome
                </Typography>
                <Typography variant="body1">{pedido.cliente.nome}</Typography>
              </Grid>
              <Grid item xs={12}>
                <Typography variant="subtitle2" color="text.secondary">
                  Telefone
                </Typography>
                <Typography variant="body1">{pedido.cliente.telefone}</Typography>
              </Grid>
            </Grid>

            <Divider sx={{ my: 2 }} />

            <Typography variant="h6" gutterBottom>
              Responsáveis
            </Typography>
            <Grid container spacing={2}>
              <Grid item xs={6}>
                <Typography variant="subtitle2" color="text.secondary">
                  Atendente
                </Typography>
                <Typography variant="body1">{pedido.atendente.nome}</Typography>
              </Grid>
              <Grid item xs={6}>
                <Typography variant="subtitle2" color="text.secondary">
                  Entregador
                </Typography>
                <Typography variant="body1">
                  {pedido.entregador ? pedido.entregador.nome : 'Não atribuído'}
                </Typography>
              </Grid>
            </Grid>

            <Divider sx={{ my: 2 }} />

            {/* Ações ou Motivo de Cancelamento */}
            <Box mt={3}>
              {pedido.status === 'cancelado' ? (
                <>
                  <Typography variant="h6" gutterBottom color="error">
                    Motivo do Cancelamento
                  </Typography>
                  <Typography variant="body1" sx={{ mt: 1 }}>
                    {pedido.motivo_cancelamento || 'Nenhum motivo registrado'}
                  </Typography>
                </>
              ) : (
                <>
                  <Typography variant="h6" gutterBottom>
                    Ações
                  </Typography>
                  <Box display="flex" flexWrap="wrap" gap={1}>
                    {pedido.status !== 'cancelado' && pedido.status !== 'finalizado' && (
                      <Button
                        variant="contained"
                        color="primary"
                        onClick={handleAbrirDialogo}
                        disabled={getStatusPermitidos(pedido.status).length === 0}
                      >
                        Atualizar Status
                      </Button>
                    )}
                    {pedido.status === 'em_entrega' && (
                      <Button
                        variant="contained"
                        color="success"
                        startIcon={<CheckIcon />}
                        onClick={handleConfirmarEntregaCompleta}
                      >
                        Confirmar Entrega
                      </Button>
                    )}
                    {pedido.status === 'entregue' && (
                      <Button
                        variant="contained"
                        color="success"
                        startIcon={<PaymentIcon />}
                        onClick={() => {
                          setNovoStatus('finalizado');
                          setDialogOpen(true);
                        }}
                      >
                        Finalizar Pedido
                      </Button>
                    )}
                    {(pedido.status === 'novo' || pedido.status === 'em_preparo' || pedido.status === 'em_entrega') && (
                      <Button
                        variant="contained"
                        color="error"
                        startIcon={<CancelIcon />}
                        onClick={handleAbrirDialogoCancelamento}
                      >
                        Cancelar Pedido
                      </Button>
                    )}
                  </Box>
                </>
              )}
            </Box>
          </Paper>
        </Grid>

        {/* Lista de itens do pedido */}
        <Grid item xs={12}>
          <Paper>
            <TableContainer>
              <Table>
                <TableHead>
                  <TableRow>
                    <TableCell>Item</TableCell>
                    <TableCell align="center">Quantidade</TableCell>
                    <TableCell align="right">Preço Unit.</TableCell>
                    <TableCell align="right">Subtotal</TableCell>
                    <TableCell align="center">Retorna Botija</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {pedido.itens.map((item) => (
                    <TableRow key={item.id}>
                      <TableCell>{item.nome_produto}</TableCell>
                      <TableCell align="center">{item.quantidade}</TableCell>
                      <TableCell align="right">{formatCurrency(item.preco_unitario)}</TableCell>
                      <TableCell align="right">{formatCurrency(item.subtotal)}</TableCell>
                      <TableCell align="center">
                        {item.retorna_botija ? (
                          <Chip label="Sim" color="success" size="small" />
                        ) : (
                          <Chip label="Não" variant="outlined" size="small" />
                        )}
                      </TableCell>
                    </TableRow>
                  ))}
                  <TableRow>
                    <TableCell colSpan={3} align="right" sx={{ fontWeight: 'bold' }}>
                      Total:
                    </TableCell>
                    <TableCell align="right" sx={{ fontWeight: 'bold' }}>
                      {formatCurrency(pedido.valor_total)}
                    </TableCell>
                    <TableCell />
                  </TableRow>
                </TableBody>
              </Table>
            </TableContainer>
          </Paper>
        </Grid>
      </Grid>

      {/* Diálogo para atualizar status */}
      <Dialog open={dialogOpen} onClose={() => setDialogOpen(false)}>
        <DialogTitle>Atualizar Status do Pedido</DialogTitle>
        <DialogContent>
          <TextField
            select
            label="Novo Status"
            value={novoStatus}
            onChange={(e) => setNovoStatus(e.target.value)}
            fullWidth
            sx={{ mt: 2, mb: 2 }}
          >
            {getStatusPermitidos(pedido.status).map((status) => (
              <MenuItem key={status.value} value={status.value}>
                {status.label}
              </MenuItem>
            ))}
          </TextField>

          {novoStatus === 'em_entrega' && (
            <TextField
              select
              label="Entregador"
              value={entregadorId}
              onChange={(e) => setEntregadorId(Number(e.target.value))}
              fullWidth
              required
              sx={{ mb: 2 }}
              error={novoStatus === 'em_entrega' && entregadorId === ''}
              helperText={
                novoStatus === 'em_entrega' && entregadorId === ''
                  ? 'Selecione um entregador'
                  : ''
              }
            >
              {entregadores.map((entregador) => (
                <MenuItem key={entregador.id} value={entregador.id}>
                  {entregador.nome}
                </MenuItem>
              ))}
            </TextField>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDialogOpen(false)}>Cancelar</Button>
          <Button
            onClick={handleAtualizarStatus}
            color="primary"
            variant="contained"
            disabled={
              atualizandoStatus ||
              (novoStatus === 'em_entrega' && !entregadorId)
            }
          >
            {atualizandoStatus ? <CircularProgress size={24} /> : 'Confirmar'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Diálogo de cancelamento */}
      <Dialog open={cancelDialogOpen} onClose={() => setCancelDialogOpen(false)}>
        <DialogTitle>Cancelar Pedido</DialogTitle>
        <DialogContent>
          <Typography variant="body1" gutterBottom>
            Informe o motivo do cancelamento:
          </Typography>
          <TextField
            autoFocus
            margin="dense"
            label="Motivo do Cancelamento"
            type="text"
            fullWidth
            multiline
            rows={4}
            value={motivoCancelamento}
            onChange={(e) => setMotivoCancelamento(e.target.value)}
            required
            error={!motivoCancelamento.trim()}
            helperText={!motivoCancelamento.trim() ? 'O motivo do cancelamento é obrigatório' : ''}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setCancelDialogOpen(false)}>Voltar</Button>
          <Button
            onClick={handleCancelarPedidoComEstoque}
            color="error"
            variant="contained"
            disabled={atualizandoStatus || !motivoCancelamento.trim()}
          >
            {atualizandoStatus ? <CircularProgress size={24} /> : 'Confirmar Cancelamento'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Snackbar para feedback */}
      <Snackbar
        open={snackbar.open}
        autoHideDuration={6000}
        onClose={handleCloseSnackbar}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
      >
        <Alert onClose={handleCloseSnackbar} severity={snackbar.severity} sx={{ width: '100%' }}>
          {snackbar.message}
        </Alert>
      </Snackbar>
    </Box>
  );
};

export default PedidoDetalhe;
