// Define and initialize the auth object immediately
const auth = {
    // Use hardcoded base URL
    // API_URL: 'https://BACKEND/',
    API_URL: 'http://127.0.0.1:5000',

    APP_URL: 'https://FRONTED/', // or whatever your frontend URL is

    // Redirect functions
    redirectToAdminDashboard() {
        window.location.href = this.APP_URL + '/admin.html';
    },

    redirectToDashboard() {
        window.location.href = this.APP_URL + '/user.html';
    },

    redirectToLogin() {
        window.location.href = this.APP_URL + '/login.html';
    },

    // Handle login success
    handleLoginSuccess(userData) {
        localStorage.setItem('user', JSON.stringify(userData));
        if (userData.role === 'admin') {
            this.redirectToAdminDashboard();
        } else {
            this.redirectToDashboard();
        }
    },

    // Handle logout
    handleLogout() {
        localStorage.removeItem('user');
        this.redirectToLogin();
    },

    // Check authentication state
    checkAuthState() {
        const userJson = localStorage.getItem('user');
        const currentPath = window.location.pathname;
        
        console.log('Current path:', currentPath);
        
        if (userJson) {
            const user = JSON.parse(userJson);
            
            if (user.role === 'admin') {
                // Admin user
                if (currentPath.includes('login') || currentPath === '/' || currentPath.includes('user')) {
                    console.log('Redirecting admin to admin dashboard...');
                    this.redirectToAdminDashboard();
                }
            } else {
                // Regular user
                if (currentPath.includes('login') || currentPath === '/' || currentPath.includes('admin')) {
                    console.log('Redirecting user to dashboard...');
                    this.redirectToDashboard();
                }
            }
        } else {
            // Not logged in
            if (!currentPath.includes('login')) {
                console.log('Redirecting to login...');
                this.redirectToLogin();
            }
        }
    },

    // Handle login form submission
    

    // Get current path
    getCurrentPath() {
        const path = window.location.pathname;
        return path === '/' ? '/login' : path;
    }
};

// Make auth globally available
window.auth = auth;

// Run auth check when document loads
document.addEventListener('DOMContentLoaded', () => auth.checkAuthState());