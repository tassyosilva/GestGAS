import { useNavigate, useLocation } from 'react-router-dom';
import {
    Box,
    Drawer,
    List,
    ListItem,
    ListItemButton,
    ListItemIcon,
    ListItemText,
    Divider,
    IconButton,
    Typography,
    useTheme,
    useMediaQuery,
} from '@mui/material';
import {
    Dashboard as DashboardIcon,
    ShoppingCart as PedidosIcon,
    Inventory as ProdutosIcon,
    People as ClientesIcon,
    LocalShipping as EntregadoresIcon,
    Insights as RelatoriosIcon,
    Settings as ConfiguracoesIcon,
    Menu as MenuIcon,
} from '@mui/icons-material';
import logo from '../../assets/logo.png';

// Definição dos itens do menu
const menuItems = [
    { text: 'Dashboard', icon: <DashboardIcon />, path: '/dashboard' },
    { text: 'Pedidos', icon: <PedidosIcon />, path: '/pedidos' },
    { text: 'Produtos', icon: <ProdutosIcon />, path: '/produtos' },
    { text: 'Clientes', icon: <ClientesIcon />, path: '/clientes' },
    { text: 'Entregadores', icon: <EntregadoresIcon />, path: '/entregadores' },
    { text: 'Relatórios', icon: <RelatoriosIcon />, path: '/relatorios' },
    { text: 'Configurações', icon: <ConfiguracoesIcon />, path: '/configuracoes' },
];

// Largura do drawer quando aberto
const drawerWidth = 240;

interface SidebarProps {
    open: boolean;
    onToggle: () => void;
}

const Sidebar = ({ open, onToggle }: SidebarProps) => {
    const theme = useTheme();
    const navigate = useNavigate();
    const location = useLocation();
    const isMobile = useMediaQuery(theme.breakpoints.down('md'));

    // Lidar com o clique em um item do menu
    const handleMenuItemClick = (path: string) => {
        navigate(path);
        if (isMobile) {
            onToggle(); // Fechar o drawer em dispositivos móveis
        }
    };

    return (
        <Drawer
            variant={isMobile ? "temporary" : "permanent"}
            open={open}
            onClose={onToggle}
            sx={{
                width: drawerWidth,
                flexShrink: 0,
                '& .MuiDrawer-paper': {
                    width: drawerWidth,
                    boxSizing: 'border-box',
                },
            }}
        >
            <Box sx={{ display: 'flex', alignItems: 'center', padding: theme.spacing(2) }}>
                <Box sx={{ display: 'flex', alignItems: 'center' }}>
                    <Box
                        component="img"
                        src={logo}
                        alt="GestGAS"
                        sx={{
                            width: 40,
                            height: 40,
                            marginRight: 1
                        }}
                    />
                    <Typography
                        variant="h6"
                        component="div"
                        sx={{ fontWeight: 'bold' }}
                    >
                        GestGAS
                    </Typography>
                </Box>
                {isMobile && (
                    <IconButton onClick={onToggle} sx={{ ml: 'auto' }}>
                        <MenuIcon />
                    </IconButton>
                )}
            </Box>

            <Divider />

            {/* Lista de itens do menu */}
            <List sx={{ mt: 2 }}>
                {menuItems.map((item) => (
                    <ListItem key={item.text} disablePadding sx={{ display: 'block' }}>
                        <ListItemButton
                            sx={{
                                minHeight: 48,
                                px: 2.5,
                                backgroundColor: location.pathname === item.path ?
                                    theme.palette.action.selected : 'transparent',
                            }}
                            onClick={() => handleMenuItemClick(item.path)}
                        >
                            <ListItemIcon
                                sx={{
                                    minWidth: 0,
                                    mr: 3,
                                    justifyContent: 'center',
                                    color: location.pathname === item.path ?
                                        theme.palette.primary.main : 'inherit',
                                }}
                            >
                                {item.icon}
                            </ListItemIcon>
                            <ListItemText
                                primary={item.text}
                                sx={{
                                    color: location.pathname === item.path ?
                                        theme.palette.primary.main : 'inherit',
                                }}
                            />
                        </ListItemButton>
                    </ListItem>
                ))}
            </List>
        </Drawer>
    );
};

export default Sidebar;