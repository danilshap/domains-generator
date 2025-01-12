let ws = null;

function connectWebSocket() {
    if (ws !== null && ws.readyState === WebSocket.OPEN) {
        return;
    }

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    ws = new WebSocket(`${protocol}//${window.location.host}/ws`);

    ws.onopen = function() {
        console.log('WebSocket connected');
    };

    ws.onmessage = function(event) {
        console.log('WebSocket message received:', event.data);
        try {
            const data = JSON.parse(event.data);
            if (data.type === 'notification' && data.payload) {
                showNotification(data.payload);
                document.body.dispatchEvent(new Event('newNotification'));
            }
        } catch (e) {
            console.error('Failed to parse WebSocket message:', e);
        }
    };

    ws.onclose = function() {
        ws = null;
        setTimeout(connectWebSocket, 1000);
    };

    ws.onerror = function(error) {
        console.error('WebSocket error:', error);
    };
}

function showNotification(notification) {
    const toast = document.createElement('div');
    toast.className = 'fixed top-4 right-4 bg-white shadow-lg rounded-lg p-4 max-w-sm w-full';
    toast.innerHTML = `
        <div class="flex items-center">
            <div class="flex-shrink-0">
                <i class="fa-solid fa-info-circle text-blue-500"></i>
            </div>
            <div class="ml-3">
                <p class="text-sm font-medium text-gray-900">${notification.title}</p>
                <p class="text-sm text-gray-500">${notification.message}</p>
            </div>
        </div>
    `;

    document.body.appendChild(toast);

    // Удаляем уведомление через 5 секунд
    setTimeout(() => {
        toast.remove();
    }, 5000);
}

// Подключаемся при загрузке страницы
document.addEventListener('DOMContentLoaded', connectWebSocket);

// Переподключаемся при возвращении на страницу
document.addEventListener('visibilitychange', function() {
    if (document.visibilityState === 'visible' && (!ws || ws.readyState !== WebSocket.OPEN)) {
        connectWebSocket();
    }
}); 