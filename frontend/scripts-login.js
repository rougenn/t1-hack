document.getElementById('loginForm').addEventListener('submit', function(event) {
    event.preventDefault(); 
    
    // Получаем данные из формы
    const login = document.getElementById('login').value;
    const password = document.getElementById('password').value;

    // Создаем объект с данными для отправки
    const data = {
        login: login,
        password: password
    };

    // Отправляем данные на сервер с помощью fetch API
    fetch('https://ваш-сервер/api/login', {
        method: 'POST', // Метод запроса
        headers: {
            'Content-Type': 'application/json' // Указываем, что отправляем JSON
        },
        body: JSON.stringify(data) // Преобразуем объект в JSON строку
    })
    .then(response => response.json()) // Преобразуем ответ в JSON
    .then(data => {
        console.log('Успех:', data);
        if (data.success) {
            // Перенаправляем пользователя на главную страницу или другую страницу после успешной авторизации
            window.location.href = './main.html';
        } else {
            document.getElementById('errorMessage').textContent = 'Неверный логин или пароль';
        }
    })
    .catch((error) => {
        console.error('Ошибка:', error);
        document.getElementById('errorMessage').textContent = 'Ошибка авторизации';
    });
});