import { Box, Typography, Paper } from '@mui/material';

interface PagePlaceholderProps {
    title: string;
    description: string;
}

const PagePlaceholder = ({ title, description }: PagePlaceholderProps) => {
    return (
        <Box>
            <Typography variant="h4" gutterBottom>
                {title}
            </Typography>

            <Paper sx={{ p: 3, borderRadius: 2 }}>
                <Typography variant="body1">
                    {description}
                </Typography>
            </Paper>
        </Box>
    );
};

export default PagePlaceholder;