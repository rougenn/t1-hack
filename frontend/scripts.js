document.addEventListener('DOMContentLoaded', function () {
    let accessToken = localStorage.getItem('access_token');
    let refreshToken = localStorage.getItem('refresh_token');
    
    // Если токен отсутствует, перенаправляем на страницу регистрации
    if (!accessToken) {
        window.location.href = './login.html';
        return;
    }

    // Отображаем email пользователя
    const userLogin = document.getElementById('userLogin');
    const logoutBtn = document.getElementById('logoutBtn');
    const userEmail = localStorage.getItem('user_email');
    if (userEmail) {
        userLogin.textContent = userEmail;
    }

    logoutBtn.addEventListener('click', function () {
        // Удаляем токены и email при выходе
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        localStorage.removeItem('user_email');
        window.location.href = './login.html'; // Редирект на страницу логина
    });

    // Кнопка создания ассистента
    const createAssistantBtn = document.getElementById('createAssistantBtn');
    createAssistantBtn.addEventListener('click', function () {
        const assistantName = document.getElementById('name').value;
        const modelName = document.getElementById('llm-select').value;
        const filesInput = document.getElementById('document');  // Элемент для загрузки файлов
        const files = filesInput.files;

        if (!assistantName || !modelName) {
            alert("Название ассистента и модель обязательны для ввода!");
            return;
        }

        const formData = new FormData();
        formData.append('assistant_name', assistantName);
        formData.append('model_name', modelName);

        // Добавляем файлы в formData, если они выбраны
        if (files.length > 0) {
            for (let i = 0; i < files.length; i++) {
                formData.append('files[]', files[i]);  // Добавляем каждый файл в formData
            }
        }

        sendCreateAssistantRequest(formData);
    });

    async function sendCreateAssistantRequest(formData) {
        try {
            const response = await fetch('http://localhost:8090/api/admin/create-chat-assistant', {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${accessToken}`,
                    // Content-Type не указываем вручную, так как FormData сам его установит
                },
                body: formData
            });

            if (!response.ok) {
                if (response.status === 401) {  // 401 - Unauthorized, токен истек
                    console.log("Token expired. Trying to refresh...");
                    await refreshAccessToken();  // Обновляем токены
                    return sendCreateAssistantRequest(formData);  // повторяем запрос после обновления токена
                } else {
                    throw new Error('Ошибка при создании ассистента. Пожалуйста, попробуйте позже.');
                }
            }

            const data = await response.json();
            console.log('Assistant created:', data);
            alert('Ассистент успешно создан!');
        } catch (error) {
            // Если ошибка, выводим сообщение об ошибке на странице
            showError(error.message);
        }
    }

    // Функция для обновления access_token с использованием refresh_token
    async function refreshAccessToken() {
        const response = await fetch('http://localhost:8090/api/admin/refresh-token', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                access_token: accessToken,  // Передаем токен в теле запроса, а не в заголовке
                refresh_token: refreshToken
            })
        });

        const data = await response.json();

        if (data.access_token && data.refresh_token) {
            // Обновляем токены в localStorage и в переменных
            localStorage.setItem('access_token', data.access_token);
            localStorage.setItem('refresh_token', data.refresh_token);

            // Обновляем переменные в коде
            accessToken = data.access_token;
            refreshToken = data.refresh_token;

            console.log('Tokens refreshed');
        } else {
            console.error('Failed to refresh token');
            alert('Ошибка при обновлении токенов!');
            window.location.href = './login.html'; // Перенаправляем на страницу логина, если refresh_token тоже невалиден
        }
    }

    // Функция для отображения ошибок на сайте
    function showError(errorMessage) {
        const errorContainer = document.getElementById('error-container');
        if (errorContainer) {
            errorContainer.textContent = errorMessage;
            errorContainer.style.display = 'block';  // Показываем контейнер с ошибкой
        }
    }

    // Далее идет ваш код для кастомизации, формы, сообщений и т.д.
    const fontForm = document.getElementById('customFont');
    const fontInput = document.getElementById('font');
    const windowView = document.getElementById('window-view');
    const dialog = document.getElementById('dialog');
    const headerWindow = document.getElementById('header-window');
    const messageInput = document.getElementById('message');
    const logoInput = document.getElementById('logo');
    const logoImage = document.querySelector('.img-example');

    fontForm.addEventListener('submit', function(event) {
        event.preventDefault(); // Предотвращаем отправку формы

        const file = fontInput.files[0];
        if (file) {
            const reader = new FileReader();

            reader.onload = function(e) {
                const fontUrl = e.target.result;
                applyCustomFont(fontUrl);
            };

            reader.readAsDataURL(file);
        }
    });

    function applyCustomFont(fontUrl) {
        const fontFace = new FontFace('CustomFont', `url(${fontUrl})`);

        fontFace.load().then(function(loadedFace) {
            document.fonts.add(loadedFace);
            windowView.style.fontFamily = 'CustomFont, sans-serif';
            dialog.style.fontFamily = 'CustomFont, sans-serif';
            headerWindow.style.fontFamily = 'CustomFont, sans-serif';
            messageInput.style.fontFamily = 'CustomFont, sans-serif';
        }).catch(function(error) {
            console.error('Ошибка загрузки шрифта:', error);
        });
    }

    function applyColor(colorInputId, targetElementId, styleProperty) {
        const colorInput = document.getElementById(colorInputId);
        const targetElement = document.getElementById(targetElementId);
        const applyButton = document.getElementById(`apply${colorInputId.charAt(0).toUpperCase() + colorInputId.slice(1)}`);

        applyButton.addEventListener('click', function(event) {
            event.preventDefault(); // Предотвращаем отправку формы
            targetElement.style[styleProperty] = colorInput.value;

            // Применяем цвет к элементам в окне диалога и шапке
            if (targetElementId === 'window-view') {
                dialog.style[styleProperty] = colorInput.value;
                // Убеждаемся, что цвет фона поля ввода сообщения не изменяется
                messageInput.style.backgroundColor = 'white';
            }
            if (targetElementId === 'dialog') {
                headerWindow.style[styleProperty] = colorInput.value;
            }
            if (targetElementId === 'window-ask') {
                headerWindow.style[styleProperty] = colorInput.value;
            }
        });
    }

    // Применение цвета фона
    applyColor('backgroundColor', 'window-view', 'backgroundColor');

    // Применение цвета текста
    applyColor('textColor', 'dialog', 'color');

    // Применение цвета ассистента
    applyColor('assistentColor', 'header-window', 'color');

    // Применение цвета пользователя
    applyColor('userColor', 'window-ask', 'backgroundColor');

    // Обработка изменения названия ассистента
    const nameForm = document.getElementById('customName');
    const nameInput = document.getElementById('name');

    nameForm.addEventListener('submit', function(event) {
        event.preventDefault(); // Предотвращаем отправку формы
        headerWindow.textContent = nameInput.value;
    });

    // Обработка отправки сообщений
    messageInput.addEventListener('keypress', function(event) {
        if (event.key === 'Enter') {
            event.preventDefault(); // Предотвращаем отправку формы
            sendMessage();
        }
    });

    function sendMessage() {
        const userMessage = messageInput.value.trim();
        if (userMessage === '') return;

        // Создаем элемент для сообщения пользователя
        const userMessageElement = document.createElement('div');
        userMessageElement.className = 'message user';
        userMessageElement.textContent = userMessage;

        // Добавляем сообщение в диалог
        dialog.appendChild(userMessageElement);

        // Прокручиваем окно чата вниз
        dialog.scrollTop = dialog.scrollHeight;

        // Очистить поле ввода
        messageInput.value = '';

        // Здесь можно добавить логику для отправки сообщения на сервер, если необходимо
    }
});
