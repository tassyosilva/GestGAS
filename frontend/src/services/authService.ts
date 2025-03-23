// src/services/authService.ts
import axios from 'axios';
import API_BASE_URL from '../config/api';

// Tipos de resposta da API
export interface LoginResponse {
    id: number;
    nome: string;
    login: string;
    perfil: string;
    token: string;
}

// Configuração inicial do Axios
const api = axios.create({
    baseURL: API_BASE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
});

// Interceptor para adicionar o token em todas as requisições
api.interceptors.request.use((config) => {
    const token = localStorage.getItem('token');
    if (token && config.headers) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

// Serviço de autenticação
export const authService = {
    // Login do usuário
    async login(login: string, senha: string): Promise<LoginResponse> {
        try {
            const response = await api.post<LoginResponse>('/login', { login, senha });
            return response.data;
        } catch (error) {
            console.error('Erro ao fazer login:', error);
            throw error;
        }
    },

    // Verificar se o usuário está autenticado
    isAuthenticated(): boolean {
        return !!localStorage.getItem('token');
    },

    // Obter dados do usuário logado
    getUser() {
        const user = localStorage.getItem('user');
        return user ? JSON.parse(user) : null;
    },

    // Logout do usuário
    logout() {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        // Redirecionar para a página de login
        window.location.href = '/login';
    }
};

export default api;