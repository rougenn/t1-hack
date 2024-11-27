document.getElementById('registrationForm').addEventListener('submit', function(event) {
    event.preventDefault(); 
    
    // Получаем данные из формы
    const login = document.getElementById('login').value;
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
    const confirmPassword = document.getElementById('confirmPassword').value;

    // Проверяем, совпадают ли пароли
    if (password !== confirmPassword) {
        document.getElementById('errorMessage').textContent = 'Пароли не совпадают';
        return;
    } else {
        document.getElementById('errorMessage').textContent = '';
    }

    // Проверяем, что пароль имеет минимальную длину 8 символов
    if (password.length < 8) {
        document.getElementById('errorMessage').textContent = 'Пароль должен быть не менее 8 символов';
        return;
    } else {
        document.getElementById('errorMessage').textContent = '';
    }

    // Создаем объект с данными для отправки
    const data = {
        email: email,
        password: password
    };

    // Отправляем данные на сервер с помощью fetch API
    fetch('https://ваш-сервер/api/register', {
        method: 'POST', // Метод запроса
        headers: {
            'Content-Type': 'application/json' // Указываем, что отправляем JSON
        },
        body: JSON.stringify(data) // Преобразуем объект в JSON строку
    })
    .then(response => response.json()) // Преобразуем ответ в JSON
    .then(data => {
        console.log('Успех:', data);
        // Перенаправляем пользователя на страницу для авторизации
        window.location.href = './login.html';
    })
    .catch((error) => {
        console.error('Ошибка:', error);
        document.getElementById('errorMessage').textContent = 'Ошибка регистрации';
    });
});