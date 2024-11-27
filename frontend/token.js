// Function to refresh the access token
function refreshAccessToken() {
    const refreshToken = localStorage.getItem('refresh_token');

    return fetch('http://localhost:8090/api/admin/refresh-token', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            access_token: localStorage.getItem('access_token'),
            refresh_token: refreshToken
        })
    })
    .then(response => response.json())
    .then(data => {
        if (data.access_token) {
            localStorage.setItem('access_token', data.access_token);
            console.log('Access token refreshed');
        } else {
            console.error('Token refresh failed');
            window.location.href = './login.html';
        }
    })
    .catch(error => {
        console.error('Error:', error);
        window.location.href = './login.html';
    });
}

// Function to make authenticated requests
function makeAuthenticatedRequest(url, options = {}) {
    const accessToken = localStorage.getItem('access_token');

    options.headers = {
        ...options.headers,
        'Authorization': `Bearer ${accessToken}`,
        'Content-Type': 'application/json'
    };

    return fetch(url, options).then(response => {
        if (response.status === 401) {
            // Token is invalid, refresh it
            return refreshAccessToken().then(() => {
                const newAccessToken = localStorage.getItem('access_token');
                options.headers['Authorization'] = `Bearer ${newAccessToken}`;
                return fetch(url, options);
            });
        }
        return response;
    });
}
