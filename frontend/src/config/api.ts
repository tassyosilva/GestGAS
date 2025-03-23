// Define a URL base da API baseada no ambiente
const API_BASE_URL = window.location.origin.includes('localhost')
    ? 'http://localhost:8080/api'
    : `${window.location.origin}/api`;

export default API_BASE_URL;