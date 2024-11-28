document.addEventListener('DOMContentLoaded', function() {
    const fontForm = document.getElementById('customFont');
    const fontInput = document.getElementById('font');
    const windowView = document.getElementById('window-view');
    const dialog = document.getElementById('dialog');
    const headerWindow = document.getElementById('header-window');
    const messageInput = document.getElementById('message');
    const logoInput = document.getElementById('logo');
    const logoImage = document.querySelector('.img-example');

    // Проверка наличия кнопки
    const exportButton = document.getElementById('exportButton');
    if (!exportButton) {
        console.error('Кнопка с ID "exportButton" не найдена');
    }

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

    // Функция для применения цвета
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

        // Очищаем поле ввода
        messageInput.value = '';

        // Прокручиваем диалог вниз, чтобы видеть последнее сообщение
        dialog.scrollTop = dialog.scrollHeight;

        // Отправляем сообщение на сервер и получаем ответ от ассистента
        fetch('/api/assistant', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ message: userMessage })
        })
        .then(response => response.json())
        .then(data => {
            const assistantMessage = data.response;
            const assistantMessageElement = document.createElement('div');
            assistantMessageElement.className = 'message assistant';
            assistantMessageElement.textContent = assistantMessage;
            dialog.appendChild(assistantMessageElement);
            dialog.scrollTop = dialog.scrollHeight;
        })
        .catch(error => {
            console.error('Ошибка при отправке сообщения:', error);
            const errorMessageElement = document.createElement('div');
            errorMessageElement.className = 'message error';
            errorMessageElement.textContent = 'Ошибка при отправке сообщения';
            dialog.appendChild(errorMessageElement);
            dialog.scrollTop = dialog.scrollHeight;
        });
    }

    // Обработка загрузки изображения
    logoInput.addEventListener('change', function(event) {
        const file = event.target.files[0];
        if (file) {
            const reader = new FileReader();

            reader.onload = function(e) {
                logoImage.src = e.target.result;
            };

            reader.readAsDataURL(file);
        }
    });

    // Обработка кнопки "Экспорт"
    exportButton.addEventListener('click', function() {
        console.log('Кнопка "Экспорт" нажата'); // Отладочное сообщение
        const exportCode = generateExportCode();
        console.log('Сгенерированный код:', exportCode); // Отладочное сообщение

        // Создаем текстовый файл и запускаем его скачивание
        const blob = new Blob([exportCode], { type: 'text/plain' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = 'dialog_code.txt';
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
    });

    function generateExportCode() {
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

        const css = `
            .window-view{
                width: 100%;
                height: 60vh;
                background-color: #1C1F4A;
                /* margin-left: 20px; */
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
                color: var(--color-white);
                width: 100%;
                background-color: var(--color-main);
                height: auto;
                border-radius: 20px 20px 0 0;
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
            .window-ask{
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
            .input-message{
                width: 100%;
                background-color: var(--color-white) !important;
                }

            .input-message::placeholder{
                font-size: 18px;
            }

            input:focus{
                outline: none;
            }
            .base-submit{
                background-image: url(./icons/submit.svg);
                width: 15px;
                height: 15px;
                border: none;
                background-color: transparent;
                margin-top: 13px;
                background-repeat: no-repeat;
                }
        `;

        const js = `
            document.getElementById('window-ask').addEventListener('submit', function(event) {
                event.preventDefault();
                const userMessage = document.getElementById('message').value.trim();
                if (userMessage === '') return;

                const userMessageElement = document.createElement('div');
                userMessageElement.className = 'message user';
                userMessageElement.textContent = userMessage;
                document.getElementById('dialog').appendChild(userMessageElement);
                document.getElementById('message').value = '';
                document.getElementById('dialog').scrollTop = document.getElementById('dialog').scrollHeight;

                fetch('/submit-message', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ message: userMessage })
                })
                .then(response => response.json())
                .then(data => {
                    const assistantMessage = data.response;
                    const assistantMessageElement = document.createElement('div');
                    assistantMessageElement.className = 'message assistant';
                    assistantMessageElement.textContent = assistantMessage;
                    document.getElementById('dialog').appendChild(assistantMessageElement);
                    document.getElementById('dialog').scrollTop = document.getElementById('dialog').scrollHeight;
                })
                .catch(error => {
                    console.error('Ошибка при отправке сообщения:', error);
                    const errorMessageElement = document.createElement('div');
                    errorMessageElement.className = 'message error';
                    errorMessageElement.textContent = 'Ошибка при отправке сообщения';
                    document.getElementById('dialog').appendChild(errorMessageElement);
                    document.getElementById('dialog').scrollTop = document.getElementById('dialog').scrollHeight;
                });
            });
        `;

        return `
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
    }
});