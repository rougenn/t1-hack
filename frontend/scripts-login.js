document.getElementById('loginForm').addEventListener('submit', function(event) {
    event.preventDefault();

    const email = document.getElementById('email').value; // Получаем email
    const password = document.getElementById('password').value; // Получаем пароль

    // Отправляем данные на сервер
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
        console.log('Login response:', data); // Проверим, что приходит от сервера

        // Проверяем, что токены и email есть в ответе
        if (data.access_token && data.refresh_token && data.user && data.user.email) {
            // Сохраняем токены и email в localStorage
            localStorage.setItem('access_token', data.access_token);
            localStorage.setItem('refresh_token', data.refresh_token);
            localStorage.setItem('user_email', data.user.email);  // Сохраняем email

            console.log('Login successful');
            window.location.href = './main.html'; // Перенаправляем на другую страницу
        } else {
            // Выводим ошибку, если токенов или email нет
            document.getElementById('errorMessage').textContent = 'Invalid login or password';
        }
    })
    .catch(error => {
        // Выводим ошибку в случае сбоя запроса
        console.error('Error:', error);
        document.getElementById('errorMessage').textContent = 'Authorization error';
    });
});
