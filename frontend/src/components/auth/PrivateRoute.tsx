import { Navigate, Outlet } from 'react-router-dom';
import { authService } from '../../services/authService.ts';

type PrivateRouteProps = {
    redirectPath?: string;
    requiredRoles?: string[];
};

/**
 * Componente de rota privada que verifica se o usuário está autenticado
 * e opcionalmente verifica se possui as permissões necessárias
 */
const PrivateRoute = ({
    redirectPath = '/login',
    requiredRoles,
}: PrivateRouteProps) => {
    const isAuthenticated = authService.isAuthenticated();
    const user = authService.getUser();

    // Se não estiver autenticado, redirecionar para o login
    if (!isAuthenticated) {
        return <Navigate to={redirectPath} replace />;
    }

    // Se tiver requisitos de perfil e o usuário não tiver permissão
    if (requiredRoles && requiredRoles.length > 0 && user) {
        const hasRequiredRole = requiredRoles.includes(user.perfil);

        if (!hasRequiredRole) {
            // Armazenar mensagem de erro e redirecionar para o login
            localStorage.setItem('authError', 'Você não tem permissão para acessar esta página');
            return <Navigate to={redirectPath} replace />;
        }
    }

    // Se autenticado e tiver as permissões, renderizar a rota protegida
    return <Outlet />;
};

export default PrivateRoute;