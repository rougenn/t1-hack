document.getElementById('registrationForm').addEventListener('submit', function(event) {
    event.preventDefault();

    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
    const confirmPassword = document.getElementById('confirmPassword').value;

    if (password !== confirmPassword) {
        document.getElementById('errorMessage').textContent = 'Passwords do not match';
        return;
    }

    if (password.length < 8) {
        document.getElementById('errorMessage').textContent = 'Password must be at least 8 characters';
        return;
    }

    const data = { email, password };

    fetch('http://localhost:8090/api/admin/signup', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            console.log('Registration successful');
            window.location.href = './login.html';
        } else {
            document.getElementById('errorMessage').textContent = 'Registration error';
        }
    })
    .catch(error => {
        console.error('Error:', error);
        document.getElementById('errorMessage').textContent = 'Registration error';
    });
});
