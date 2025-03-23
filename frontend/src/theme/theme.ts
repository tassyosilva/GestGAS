import { createTheme } from '@mui/material/styles';

// Criando o tema com as cores especificadas: azul, branco, verde e vermelho
const theme = createTheme({
    palette: {
        primary: {
            main: '#1976d2', // Azul
            light: '#42a5f5',
            dark: '#1565c0',
            contrastText: '#fff',
        },
        secondary: {
            main: '#4caf50', // Verde
            light: '#81c784',
            dark: '#388e3c',
            contrastText: '#fff',
        },
        error: {
            main: '#d32f2f', // Vermelho
            light: '#ef5350',
            dark: '#c62828',
            contrastText: '#fff',
        },
        background: {
            default: '#ffffff', // Branco
            paper: '#f5f5f5',
        },
    },
    typography: {
        fontFamily: [
            'Roboto',
            '"Helvetica Neue"',
            'Arial',
            'sans-serif',
        ].join(','),
        h4: {
            fontWeight: 600,
        },
        h5: {
            fontWeight: 500,
        },
        h6: {
            fontWeight: 500,
        },
    },
    components: {
        MuiButton: {
            styleOverrides: {
                root: {
                    textTransform: 'none',
                    borderRadius: 8,
                    padding: '8px 16px',
                },
                contained: {
                    boxShadow: 'none',
                    '&:hover': {
                        boxShadow: '0px 2px 4px -1px rgba(0,0,0,0.2)',
                    },
                },
            },
        },
        MuiCard: {
            styleOverrides: {
                root: {
                    borderRadius: 12,
                    boxShadow: '0px 2px 8px rgba(0, 0, 0, 0.05)',
                },
            },
        },
        MuiTextField: {
            styleOverrides: {
                root: {
                    '& .MuiOutlinedInput-root': {
                        borderRadius: 8,
                    },
                },
            },
        },
    },
});

export default theme;