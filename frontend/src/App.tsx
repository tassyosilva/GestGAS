import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { ThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import theme from './theme/theme';

// Layout
import Layout from './components/layout/Layout';

// Páginas
import Login from './pages/Login';
import Dashboard from './pages/Dashboard';
import Produtos from './pages/Produtos';
import Pedidos from './pages/Pedidos';
import PedidoDetalhe from './pages/PedidoDetalhe';
import NovoPedido from './pages/NovoPedido';
import Estoque from './pages/Estoque';
import Clientes from './pages/Clientes';
import Entregadores from './pages/Entregadores';
import Relatorios from './pages/Relatorios';
import Configuracoes from './pages/Configuracoes';
import Usuarios from './pages/Usuarios';

// Componentes de autenticação
import PrivateRoute from './components/auth/PrivateRoute';

function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Router>
        <Routes>
          {/* Rotas públicas */}
          <Route path="/login" element={<Login />} />

          {/* Rotas protegidas */}
          <Route element={<PrivateRoute />}>
            <Route element={<Layout />}>
              <Route path="/dashboard" element={<Dashboard />} />
              <Route path="/produtos" element={<Produtos />} />
              <Route path="/pedidos" element={<Pedidos />} />
              <Route path="/pedidos/novo" element={<NovoPedido />} />
              <Route path="/pedidos/:id" element={<PedidoDetalhe />} />
              <Route path="/estoque" element={<Estoque />} />
              <Route path="/clientes" element={<Clientes />} />
              <Route path="/entregadores" element={<Entregadores />} />
              <Route path="/relatorios" element={<Relatorios />} />
              <Route path="/configuracoes" element={<Configuracoes />} />
              <Route path="/usuarios" element={<Usuarios />} />
            </Route>
          </Route>

          {/* Redirecionamento para o Dashboard se autenticado ou Login se não */}
          <Route path="/" element={<Navigate to="/dashboard" replace />} />

          {/* Rota para página não encontrada */}
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </Router>
    </ThemeProvider>
  );
}

export default App;