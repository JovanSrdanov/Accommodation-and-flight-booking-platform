import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';
import {BrowserRouter} from 'react-router-dom';
import {createTheme, ThemeProvider} from '@mui/material/styles';
import {SnackbarProvider} from "notistack";

const darkTheme = createTheme({
    palette: {
        mode: 'dark',
    },
});

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
    <ThemeProvider theme={darkTheme}>
        <BrowserRouter>
            <SnackbarProvider maxSnack={3}>
                <App></App>
            </SnackbarProvider>
        </BrowserRouter>
    </ThemeProvider>
);

