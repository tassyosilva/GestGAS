import axios from 'axios';
import API_BASE_URL from '../config/api';

// Tipos
export interface Usuario {
    id: number;
    nome: string;
    login: string;
    cpf?: string;
    email?: string;
    perfil: string;
    criado_em: string;
    atualizado_em: string;
}

export interface NovoUsuario {
    nome: string;
    login: string;
    senha: string;
    cpf?: string;
    email?: string;
    perfil: string;
}

export interface AtualizacaoUsuario {
    nome: string;
    login?: string;
    senha?: string;
    cpf?: string;
    email?: string;
    perfil?: string;
}

// Serviço para gerenciamento de usuários
export const usuarioService = {
    // Buscar todos os usuários
    async listarUsuarios(): Promise<Usuario[]> {
        try {
            const token = localStorage.getItem('token');
            const response = await axios.get<Usuario[]>(`${API_BASE_URL}/usuarios`, {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });
            return response.data;
        } catch (error) {
            console.error('Erro ao listar usuários:', error);
            throw error;
        }
    },

    // Buscar um usuário pelo ID
    async obterUsuario(id: number): Promise<Usuario> {
        try {
            const token = localStorage.getItem('token');
            const response = await axios.get<Usuario>(`${API_BASE_URL}/usuarios/${id}`, {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });
            return response.data;
        } catch (error) {
            console.error(`Erro ao obter usuário ${id}:`, error);
            throw error;
        }
    },

    // Criar um novo usuário
    async criarUsuario(usuario: NovoUsuario): Promise<Usuario> {
        try {
            const token = localStorage.getItem('token');
            const response = await axios.post<Usuario>(`${API_BASE_URL}/usuarios`, usuario, {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });
            return response.data;
        } catch (error) {
            console.error('Erro ao criar usuário:', error);
            throw error;
        }
    },

    // Atualizar um usuário existente
    async atualizarUsuario(id: number, usuario: AtualizacaoUsuario): Promise<Usuario> {
        try {
            const token = localStorage.getItem('token');
            const response = await axios.put<Usuario>(`${API_BASE_URL}/usuarios/${id}`, usuario, {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });
            return response.data;
        } catch (error) {
            console.error(`Erro ao atualizar usuário ${id}:`, error);
            throw error;
        }
    },

    // Excluir um usuário
    async excluirUsuario(id: number): Promise<void> {
        try {
            const token = localStorage.getItem('token');
            await axios.delete(`${API_BASE_URL}/usuarios/${id}`, {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });
        } catch (error) {
            console.error(`Erro ao excluir usuário ${id}:`, error);
            throw error;
        }
    },

    // Listar entregadores
    async listarEntregadores(): Promise<Usuario[]> {
        try {
            const token = localStorage.getItem('token');
            const response = await axios.get<Usuario[]>(`${API_BASE_URL}/entregadores`, {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });
            return response.data;
        } catch (error) {
            console.error('Erro ao listar entregadores:', error);
            throw error;
        }
    }
};