import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Box,
  Button,
  Checkbox,
  CssBaseline,
  FormControlLabel,
  TextField,
  Typography,
  Link,
  Card,
  Stack,
  FormControl,
  FormLabel,
  styled,
} from '@mui/material';
import logo from '../assets/logo.png';
import { authService } from '../services/authService';

const StyledCard = styled(Card)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'column',
  alignSelf: 'center',
  width: '100%',
  padding: theme.spacing(4),
  gap: theme.spacing(2),
  margin: 'auto',
  [theme.breakpoints.up('sm')]: {
    maxWidth: '450px',
  },
  boxShadow:
    'hsla(220, 30%, 5%, 0.05) 0px 5px 15px 0px, hsla(220, 25%, 10%, 0.05) 0px 15px 35px -5px',
}));

const LoginContainer = styled(Stack)(({ theme }) => ({
  height: '100vh',
  minHeight: '100%',
  padding: theme.spacing(2),
  [theme.breakpoints.up('sm')]: {
    padding: theme.spacing(4),
  },
  '&::before': {
    content: '""',
    display: 'block',
    position: 'absolute',
    zIndex: -1,
    inset: 0,
    backgroundImage:
      'radial-gradient(ellipse at 50% 50%, hsl(210, 100%, 97%), hsl(0, 0%, 100%))',
    backgroundRepeat: 'no-repeat',
  },
}));

const SitemarkIcon = () => (
  <Box
    component="img"
    sx={{
      height: 64,
      width: 'auto',
      alignSelf: 'center',
      mb: 2,
    }}
    alt="Logo GestGAS"
    src={logo}
  />
);

export default function Login() {
  const navigate = useNavigate();
  const [loginError, setLoginError] = useState('');
  const [formData, setFormData] = useState({
    login: '',
    senha: '',
  });
  const [loading, setLoading] = useState(false);

  // Verificar se há mensagem de erro de autenticação armazenada
  useEffect(() => {
    const authError = localStorage.getItem('authError');
    if (authError) {
      setLoginError(authError);
      localStorage.removeItem('authError');
    }
  }, []);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData({
      ...formData,
      [name]: value,
    });
  };

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setLoading(true);

    try {
      const response = await authService.login(formData.login, formData.senha);

      // Armazenar o token no localStorage
      localStorage.setItem('token', response.token);
      localStorage.setItem('user', JSON.stringify({
        id: response.id,
        nome: response.nome,
        login: response.login,
        perfil: response.perfil,
      }));

      // Redirecionar para a página principal
      navigate('/dashboard');
    } catch (error) {
      console.error('Erro de login:', error);
      setLoginError('Usuário ou senha inválidos');
    } finally {
      setLoading(false);
    }
  };

  return (
    <CssBaseline>
      <LoginContainer direction="column" justifyContent="space-between">
        <StyledCard variant="outlined">
          <SitemarkIcon />
          <Typography
            component="h1"
            variant="h4"
            sx={{ width: '100%', fontSize: 'clamp(2rem, 10vw, 2.15rem)', textAlign: 'center' }}
          >
            GestGAS
          </Typography>
          <Typography
            component="h2"
            variant="h6"
            sx={{ width: '100%', textAlign: 'center', mb: 2 }}
          >
            Sistema de Gerenciamento para Revendedora de Gás e Água
          </Typography>

          {loginError && (
            <Typography color="error" sx={{ textAlign: 'center', mb: 2 }}>
              {loginError}
            </Typography>
          )}

          <Box
            component="form"
            onSubmit={handleSubmit}
            noValidate
            sx={{
              display: 'flex',
              flexDirection: 'column',
              width: '100%',
              gap: 2,
            }}
          >
            <FormControl>
              <FormLabel htmlFor="login">Login</FormLabel>
              <TextField
                id="login"
                name="login"
                value={formData.login}
                onChange={handleChange}
                placeholder="seu.login"
                autoComplete="username"
                autoFocus
                required
                fullWidth
                variant="outlined"
              />
            </FormControl>

            <FormControl>
              <FormLabel htmlFor="senha">Senha</FormLabel>
              <TextField
                name="senha"
                value={formData.senha}
                onChange={handleChange}
                placeholder="••••••"
                type="password"
                id="senha"
                autoComplete="current-password"
                required
                fullWidth
                variant="outlined"
              />
            </FormControl>

            <FormControlLabel
              control={<Checkbox value="remember" color="primary" />}
              label="Lembrar-me"
            />

            <Button
              type="submit"
              fullWidth
              variant="contained"
              color="primary"
              disabled={loading}
            >
              {loading ? 'Entrando...' : 'Entrar'}
            </Button>

            <Link href="#" variant="body2" sx={{ alignSelf: 'center' }}>
              Esqueceu a senha?
            </Link>
          </Box>
        </StyledCard>
      </LoginContainer>
    </CssBaseline>
  );
}