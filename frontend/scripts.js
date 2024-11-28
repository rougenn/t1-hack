document.addEventListener('DOMContentLoaded', function () {
    let accessToken = localStorage.getItem('access_token');
    let refreshToken = localStorage.getItem('refresh_token');
    let assistantId = localStorage.getItem('assistant_id');  // Добавим переменную для хранения ID ассистента
    
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
        localStorage.removeItem('assistant_id');  // Удаляем ID ассистента
        window.location.href = './login.html'; // Редирект на страницу логина
    });

    // Кнопка создания ассистента
    const createAssistantBtn = document.getElementById('createAssistantBtn');
    createAssistantBtn.addEventListener('click', function () {
        const assistantName = document.getElementById('name').value;
        const modelName = document.getElementById('llm-select').value;
        const embaddingName = document.getElementById('emb-select').value;
        const chunkSize = document.getElementById('chnk-select').value;
        const filesInput = document.getElementById('document');  // Элемент для загрузки файлов
        const files = filesInput.files;

        if (!assistantName || !modelName) {
            alert("Название ассистента и модель обязательны для ввода!");
            return;
        }

        const formData = new FormData();
        formData.append('assistant_name', assistantName);
        formData.append('model_name', modelName);
        formData.append('embeddings_model_id', embaddingName);
        formData.append('chunk_size', chunkSize);

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
            assistantId = data.assistant_id;  // Сохраняем ID ассистента
            localStorage.setItem('assistant_id', assistantId);  // Сохраняем в localStorage
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

    // Обработка отправки сообщений
    const messageInput = document.getElementById('message');
    const dialog = document.getElementById('dialog');
    const sendBtn = document.getElementById('send-btn');

    sendBtn.addEventListener('click', async function() {
        if (!assistantId) {
            alert("Пожалуйста, создайте ассистента перед отправкой сообщений.");
            return;
        }

        const message = messageInput.value.trim();
        if (!message) return; // Не отправляем пустые сообщения

        // Отображаем сообщение пользователя
        const userMessageElement = document.createElement('div');
        userMessageElement.className = 'message user-message';
        userMessageElement.textContent = message;
        dialog.appendChild(userMessageElement);
        dialog.scrollTop = dialog.scrollHeight;

        messageInput.value = '';  // Очищаем поле ввода

        // Блокируем интерфейс отправки сообщений, пока идет запрос
        sendBtn.disabled = true;
        messageInput.disabled = true;

        try {
            // Отправляем сообщение
            const response = await fetch(`http://localhost:8090/api/chats/send/${assistantId}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${accessToken}`,
                },
                body: JSON.stringify({ message })
            });

            if (!response.ok) {
                if (response.status === 401) {  // 401 - Unauthorized, токен истек
                    console.log("Token expired. Trying to refresh...");
                    await refreshAccessToken();  // Обновляем токены
                    return addEventListener(formData);  // повторяем запрос после обновления токена
                } else {
                    throw new Error('Ошибка при создании ассистента. Пожалуйста, попробуйте позже.');
                }
            }

            const data = await response.json();
            const assistantMessage = data.message;

            // Отображаем сообщение ассистента
            const assistantMessageElement = document.createElement('div');
            assistantMessageElement.className = 'message bot-message';
            assistantMessageElement.textContent = assistantMessage;
            dialog.appendChild(assistantMessageElement);
            dialog.scrollTop = dialog.scrollHeight;
        } catch (error) {
            alert(error.message);
        } finally {
            // Восстанавливаем интерфейс отправки сообщений
            sendBtn.disabled = false;
            messageInput.disabled = false;
        }
    });

    // Далее идет ваш код для кастомизации, формы, сообщений и т.д.
    const fontForm = document.getElementById('customFont');
    const fontInput = document.getElementById('font');
    const windowView = document.getElementById('window-view');
    const headerWindow = document.getElementById('header-window');
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

    // Обработка отправки сообщений
    messageInput.addEventListener('keypress', function(event) {
        if (event.key === 'Enter') {
            event.preventDefault(); // Предотвращаем отправку формы
            sendMessage();
        }
    });

    async function sendMessage() {
        const userMessage = messageInput.value.trim();
        if (userMessage === '') return;

        // Создаем элемент для сообщения пользователя
        const userMessageElement = document.createElement('div');
        userMessageElement.className = 'message user-message';
        userMessageElement.textContent = userMessage;

        // Добавляем сообщение в диалог
        dialog.appendChild(userMessageElement);

        // Прокручиваем окно чата вниз
        dialog.scrollTop = dialog.scrollHeight;

        // Очистить поле ввода
        messageInput.value = '';

        // Блокируем интерфейс отправки сообщений, пока идет запрос
        sendBtn.disabled = true;
        messageInput.disabled = true;

        try {
            // Отправляем сообщение
            const response = await fetch(`http://localhost:8090/api/chats/send/${assistantId}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${accessToken}`,
                },
                body: JSON.stringify({ message: userMessage })
            });

            if (!response.ok) {
                if (response.status === 401) {  // 401 - Unauthorized, токен истек
                    console.log("Token expired. Trying to refresh...");
                    await refreshAccessToken();  // Обновляем токены
                    return sendMessage(formData);  // повторяем запрос после обновления токена
                } else {
                    throw new Error('Ошибка при создании ассистента. Пожалуйста, попробуйте позже.');
                }
            }

            const data = await response.json();
            const assistantMessage = data.message;

            // Отображаем сообщение ассистента
            const assistantMessageElement = document.createElement('div');
            assistantMessageElement.className = 'message bot-message';
            assistantMessageElement.textContent = assistantMessage;
            dialog.appendChild(assistantMessageElement);
            dialog.scrollTop = dialog.scrollHeight;
        } catch (error) {
            alert(error.message);
        } finally {
            // Восстанавливаем интерфейс отправки сообщений
            sendBtn.disabled = false;
            messageInput.disabled = false;
        }
    }

    // Функция для генерации HTML-кода
    function generateExportHTML() {
        const html = `
            <div class="window-view" id="window-view">
                <div class="header-window" id="header-window">
                    Название ассистента
                </div>
                <div class="dialog" id="dialog">
                    <!-- Здесь будут добавляться сообщения -->
                </div>
                <form action="/submit-message" method="post" class="window-ask" id="window-ask">
                    <input type="text" name="message" id="message" placeholder="Написать запрос" class="input-message">
                    <button id="applyMessage" class="base-submit"></button>
                </form>
            </div>
        `;
        return html;
    }

    // Функция для генерации CSS-кода
    function generateExportCSS() {
        const css = `
            .window-view {
                width: 100%;
                height: 60vh;
                background-color: ${windowView.style.backgroundColor};
                margin-top: 2vw;
                border-radius: 20px;
                position: relative;
                border: 1px solid var(--color-border-light);
                overflow: hidden;
            }

            .header-window {
                padding: 1vw 2vw;
                font-weight: 700;
                font-size: 18px;
                line-height: 24px;
                color: ${headerWindow.style.color};
                width: 100%;
                background-color: ${headerWindow.style.backgroundColor};
                height: auto;
                border-radius: 20px 20px 0 0;
            }
            window-ask{
                padding: 1vw 2vw;
                width: 100%;
                height: auto;
                background-color: var(--color-white) !important;
                position: absolute;
                bottom: 0;
                left: 0;
                right: 0;
                display: flex;
            }

            .dialog {
                overflow-y: auto;
                height: 300px; /* Установите нужную высоту */
                padding: 10px;
                display: flex;
                flex-direction: column;
                width: 100%; /* Убедитесь, что контейнер имеет достаточную ширину */
            }

            .dialog::after {
                content: '';
                display: table;
                clear: both;
            }

            .window-ask {
                padding: 1vw 2vw;
                width: 50%;
                height: auto;
                background-color: ${messageInput.style.backgroundColor} !important;
                position: absolute;
                bottom: 0;
                left: 0;
                right: 0;
                display: flex;
            }

            .input-message {
                width: 100%;
                background-color: ${messageInput.style.backgroundColor} !important;
            }

            .input-message::placeholder {
                font-size: 18px;
            }

            input:focus {
                outline: none;
            }

            .base-submit {
                background-image: url(./icons/submit.svg);
                width: 15px;
                height: 15px;
                border: none;
                background-color: transparent;
                margin-top: 13px;
                background-repeat: no-repeat;
            }

            /* Применение цветов к сообщениям */
            .bot-message {
                background-color: ${document.getElementById('assistentColor').value};
                color: ${document.getElementById('textColor').value};
            }

            .user-message {
                background-color: ${document.getElementById('userColor').value};
                color: ${document.getElementById('textColor').value};
            }
        `;
        return css;
    }

    // Функция для генерации JavaScript-кода
    function generateExportJS() {
        const js = `
            document.getElementById('window-ask').addEventListener('submit', async function(event) {
                event.preventDefault();
                const userMessage = document.getElementById('message').value.trim();
                if (userMessage === '') return;
    
                const userMessageElement = document.createElement('div');
                userMessageElement.className = 'message user-message';
                userMessageElement.textContent = userMessage;
                document.getElementById('dialog').appendChild(userMessageElement);
                document.getElementById('message').value = '';
                document.getElementById('dialog').scrollTop = document.getElementById('dialog').scrollHeight;
    
                try {
                    const response = await fetch('http://localhost:8090/api/chats/send/${assistantId}', {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json',
                            'Authorization': 'Bearer ${accessToken}',
                        },
                        body: JSON.stringify({ message: userMessage })
                    });
    
                    if (!response.ok) {
                        throw new Error('Ошибка при отправке сообщения.');
                    }
    
                    const data = await response.json();
                    const assistantMessage = data.message;
    
                    const assistantMessageElement = document.createElement('div');
                    assistantMessageElement.className = 'message bot-message';
                    assistantMessageElement.textContent = assistantMessage;
                    document.getElementById('dialog').appendChild(assistantMessageElement);
                    document.getElementById('dialog').scrollTop = document.getElementById('dialog').scrollHeight;
                } catch (error) {
                    console.error('Ошибка при отправке сообщения:', error);
                    const errorMessageElement = document.createElement('div');
                    errorMessageElement.className = 'message error';
                    errorMessageElement.textContent = 'Ошибка при отправке сообщения';
                    document.getElementById('dialog').appendChild(errorMessageElement);
                    document.getElementById('dialog').scrollTop = document.getElementById('dialog').scrollHeight;
                }
            });
        `;
        return js;
    }

    // Функция для скачивания экспортированного кода
    async function downloadExport() {
        try {
            // Проверяем, не истек ли токен
            const response = await fetch('http://localhost:8090/api/check-token', {
                method: 'GET',
                headers: {
                    'Authorization': `Bearer ${accessToken}`,
                },
            });

            if (!response.ok) {
                if (response.status === 401) {  // 401 - Unauthorized, токен истек
                    console.log("Token expired. Trying to refresh...");
                    await refreshAccessToken();  // Обновляем токены
                } else {
                    throw new Error('Ошибка при проверке токена.');
                }
            }

            const html = generateExportHTML();
            const css = generateExportCSS();
            const js = generateExportJS();

            const exportCode = `
                <!DOCTYPE html>
                <html lang="en">
                <head>
                    <meta charset="UTF-8">
                    <meta name="viewport" content="width=device-width, initial-scale=1.0">
                    <title>Окно знаний | YSL</title>
                    <style>${css}</style>
                </head>
                <body>
                    ${html}
                    <script>${js}</script>
                </body>
                </html>
            `;

            const blob = new Blob([exportCode], { type: 'text/html' });
            const url = URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = 'exported_dialog.html';
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
            URL.revokeObjectURL(url);
        } catch (error) {
            alert(error.message);
        }
    }

    // Привязываем функцию скачивания к кнопке экспорта
    const exportBtn = document.getElementById('exportBtn');
    exportBtn.addEventListener('click', downloadExport);
});