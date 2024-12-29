// Add token to all HTMX requests
document.body.addEventListener('htmx:configRequest', function(evt) {
    const token = localStorage.getItem('token');
    if (token) {
        evt.detail.headers['Authorization'] = 'Bearer ' + token;
    }
});

// Handle unauthorized responses
document.body.addEventListener('htmx:responseError', function(evt) {
    if (evt.detail.xhr.status === 401) {
        localStorage.removeItem('token');
        window.location.href = '/login';
    }
}); 