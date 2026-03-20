import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
export default defineConfig(function (_a) {
    var command = _a.command;
    return ({
        base: command === 'build' ? './' : '/',
        plugins: [react()],
        server: {
            proxy: {
                '/api/v1': {
                    target: 'http://127.0.0.1:10000',
                    changeOrigin: true,
                    ws: true,
                },
                '/healthz': {
                    target: 'http://127.0.0.1:10000',
                    changeOrigin: true,
                },
            },
        },
    });
});
