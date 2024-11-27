document.getElementById('loginForm').addEventListener('submit', function(event) {
    event.preventDefault();

    const email = document.getElementById('login').value;
    const password = document.getElementById('password').value;

    const data = { email, password };

    fetch('http://localhost:8090/api/admin/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    })
    .then(response => response.json())
    .then(data => {
        if (data.access_token && data.refresh_token) {
            localStorage.setItem('access_token', data.access_token);
            localStorage.setItem('refresh_token', data.refresh_token);

            console.log('Login successful');
            window.location.href = './main.html';
        } else {
            document.getElementById('errorMessage').textContent = 'Invalid login or password';
        }
    })
    .catch(error => {
        console.error('Error:', error);
        document.getElementById('errorMessage').textContent = 'Authorization error';
    });
});
