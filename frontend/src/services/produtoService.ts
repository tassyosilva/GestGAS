import axios from 'axios';
import API_BASE_URL from '../config/api';

// Tipos
export interface Produto {
    id: number;
    nome: string;
    descricao: string;
    categoria: string;
    preco: number;
    criado_em?: string;
    atualizado_em?: string;
}

export interface NovoProduto {
    nome: string;
    descricao: string;
    categoria: string;
    preco: number;
}

// Serviço para gerenciamento de produtos
export const produtoService = {
    // Buscar todos os produtos
    async listarProdutos(): Promise<Produto[]> {
        try {
            const token = localStorage.getItem('token');
            const response = await axios.get<Produto[]>(`${API_BASE_URL}/produtos`, {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });
            return response.data;
        } catch (error) {
            console.error('Erro ao listar produtos:', error);
            throw error;
        }
    },

    // Buscar um produto pelo ID
    async obterProduto(id: number): Promise<Produto> {
        try {
            const token = localStorage.getItem('token');
            const response = await axios.get<Produto>(`${API_BASE_URL}/produtos/${id}`, {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });
            return response.data;
        } catch (error) {
            console.error(`Erro ao obter produto ${id}:`, error);
            throw error;
        }
    },

    // Criar um novo produto
    async criarProduto(produto: NovoProduto): Promise<Produto> {
        try {
            const token = localStorage.getItem('token');
            // Adicione uma barra no final da URL
            const response = await axios.post<Produto>(`${API_BASE_URL}/produtos/`, produto, {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });
            return response.data;
        } catch (error) {
            console.error('Erro ao criar produto:', error);
            throw error;
        }
    },

    // Atualizar um produto existente
    async atualizarProduto(id: number, produto: NovoProduto): Promise<Produto> {
        try {
            const token = localStorage.getItem('token');
            const response = await axios.put<Produto>(`${API_BASE_URL}/produtos/${id}`, produto, {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });
            return response.data;
        } catch (error) {
            console.error(`Erro ao atualizar produto ${id}:`, error);
            throw error;
        }
    },

    // Excluir um produto
    async excluirProduto(id: number): Promise<void> {
        try {
            const token = localStorage.getItem('token');
            await axios.delete(`${API_BASE_URL}/produtos/${id}`, {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });
        } catch (error) {
            console.error(`Erro ao excluir produto ${id}:`, error);
            throw error;
        }
    }
};