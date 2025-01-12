let ws = null;
const SHOWN_NOTIFICATIONS_KEY = 'shown_notifications';

// Получаем уже показанные нотификации
function getShownNotifications() {
    const shown = localStorage.getItem(SHOWN_NOTIFICATIONS_KEY);
    return shown ? JSON.parse(shown) : [];
}

// Добавляем нотификацию в показанные
function markNotificationAsShown(notification) {
    const shown = getShownNotifications();
    shown.push({
        title: notification.title,
        message: notification.message,
        timestamp: Date.now()
    });
    
    // Храним только последние 100 нотификаций не старше 24 часов
    const dayAgo = Date.now() - 24 * 60 * 60 * 1000;
    const filtered = shown.filter(n => n.timestamp > dayAgo).slice(-100);
    
    localStorage.setItem(SHOWN_NOTIFICATIONS_KEY, JSON.stringify(filtered));
}

// Проверяем, была ли уже показана такая нотификация
function wasNotificationShown(notification) {
    const shown = getShownNotifications();
    return shown.some(n => 
        n.title === notification.title && 
        n.message === notification.message
    );
}

function connectWebSocket() {
    if (ws !== null && ws.readyState === WebSocket.OPEN) {
        console.log('WebSocket already connected');
        return;
    }

    ws = new WebSocket(`ws://${window.location.host}/ws`);
    console.log('WebSocket connecting...');

    ws.onopen = function() {
        console.log('WebSocket connected');
    };

    ws.onmessage = function(event) {
        console.log('WebSocket message received:', event.data);
        try {
            const data = JSON.parse(event.data);
            if (data.type === 'notification' && data.payload.type !== 'system') {
                console.log('WebSocket message received:', data.payload);
                showNotification(data.payload);
            }
        } catch (e) {
            console.error('Failed to parse WebSocket message:', e);
        }
    };

    ws.onclose = function() {
        console.log('WebSocket disconnected');
        ws = null;
        if (document.visibilityState === 'visible') {
            setTimeout(connectWebSocket, 5000);
        }
    };
}

function showNotification(notification) {
    // Проверяем, показывали ли мы уже такую нотификацию
    if (wasNotificationShown(notification)) {
        return;
    }

    // Отмечаем нотификацию как показанную
    markNotificationAsShown(notification);

    const container = document.getElementById('notifications');
    const notif = document.createElement('div');
    notif.className = `p-4 rounded-lg shadow-lg transition-all transform translate-x-0 ${getBackgroundColor(notification.type)}`;
    notif.innerHTML = `
        <div class="flex items-center">
            <div class="flex-shrink-0">${getIcon(notification.type)}</div>
            <div class="ml-3">
                <h3 class="text-sm font-medium">${notification.title}</h3>
                <p class="mt-1 text-sm opacity-90">${notification.message}</p>
            </div>
        </div>
    `;
    
    container.appendChild(notif);
    
    setTimeout(() => {
        notif.style.transform = 'translateX(150%)';
        setTimeout(() => notif.remove(), 300);
    }, 5000);
}

function getBackgroundColor(type) {
    switch(type) {
        case 'success': return 'bg-green-500 text-white';
        case 'error': return 'bg-red-500 text-white';
        case 'warning': return 'bg-yellow-500 text-white';
        default: return 'bg-blue-500 text-white';
    }
}

function getIcon(type) {
    switch(type) {
        case 'success': return '<i class="fas fa-check-circle"></i>';
        case 'error': return '<i class="fas fa-exclamation-circle"></i>';
        case 'warning': return '<i class="fas fa-exclamation-triangle"></i>';
        default: return '<i class="fas fa-info-circle"></i>';
    }
}

// Подключаемся при загрузке страницы
document.addEventListener('DOMContentLoaded', connectWebSocket);

// Переподключаемся при возвращении на страницу
document.addEventListener('visibilitychange', function() {
    if (document.visibilityState === 'visible' && (!ws || ws.readyState !== WebSocket.OPEN)) {
        connectWebSocket();
    }
}); 