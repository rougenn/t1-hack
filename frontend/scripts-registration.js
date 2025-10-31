document.getElementById('registrationForm').addEventListener('submit', function(event) {
    event.preventDefault();

    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
    const confirmPassword = document.getElementById('confirmPassword').value;

    // Проверка, что пароли совпадают
    if (password !== confirmPassword) {
        document.getElementById('errorMessage').textContent = 'Passwords do not match';
        return;
    }

    // Проверка длины паролей
    if (password.length < 8) {
        document.getElementById('errorMessage').textContent = 'Password must be at least 8 characters';
        return;
    }

    const data = { email, password };  // Убираем "login", только email и password

    // Отправка запроса на сервер
    fetch('http://localhost:8090/api/admin/signup', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    })
    .then(response => response.json())
    .then(data => {
        console.log('Server response:', data);  // Отладка ответа от сервера

        // Проверка, что объект user существует в ответе
        if (data.user) {
            console.log('Registration successful');
            window.location.href = './login.html';  // Переход на страницу логина
        } else {
            document.getElementById('errorMessage').textContent = 'Registration error';
        }
    })
    .catch(error => {
        console.error('Error:', error);
        document.getElementById('errorMessage').textContent = 'Registration error';
    });
});
