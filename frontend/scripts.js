document.addEventListener('DOMContentLoaded', function() {
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

        // Очищаем поле ввода
        messageInput.value = '';

        // Прокручиваем диалог вниз, чтобы видеть последнее сообщение
        dialog.scrollTop = dialog.scrollHeight;

        // Здесь можно добавить логику для отправки сообщения на сервер и получения ответа от ассистента
        // Например, с помощью fetch или axios
        // Пример:
        // fetch('/api/assistant', {
        //     method: 'POST',
        //     headers: {
        //         'Content-Type': 'application/json'
        //     },
        //     body: JSON.stringify({ message: userMessage })
        // })
        // .then(response => response.json())
        // .then(data => {
        //     const assistantMessage = data.response;
        //     const assistantMessageElement = document.createElement('div');
        //     assistantMessageElement.className = 'message assistant';
        //     assistantMessageElement.textContent = assistantMessage;
        //     dialog.appendChild(assistantMessageElement);
        //     dialog.scrollTop = dialog.scrollHeight;
        // });
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
});