export const environment = {
    production: true,
    apiUrl: '/api',
    get wsUrl(): string {
        if (typeof window !== 'undefined') {
            const proto = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
            return `${proto}//${window.location.host}/api`;
        }
        return '/api';
    }
};
