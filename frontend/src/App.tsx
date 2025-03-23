import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { ThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import theme from './theme/theme.ts';

// Páginas
import Login from './pages/Login.tsx';
import Dashboard from './pages/Dashboard.tsx';

// Componentes de autenticação
import PrivateRoute from './components/auth/PrivateRoute.tsx';

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
            <Route path="/dashboard" element={<Dashboard />} />
            {/* Adicionar mais rotas protegidas aqui */}
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