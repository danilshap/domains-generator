htmx.on("showUnreadBadge", function() {
    const badge = document.getElementById('unread-notifications-badge');
    if (badge) {
        badge.classList.remove('hidden');
    }
});

htmx.on("hideUnreadBadge", function() {
    const badge = document.getElementById('unread-notifications-badge');
    if (badge) {
        badge.classList.add('hidden');
    }
});

document.body.addEventListener('newNotification', function() {
    const badge = document.getElementById('unread-notifications-badge');
    if (badge) {
        badge.classList.remove('hidden');
    }
}); 