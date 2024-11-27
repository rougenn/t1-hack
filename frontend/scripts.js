// document.addEventListener('DOMContentLoaded', function() {
//     const fontForm = document.getElementById('customFont');
//     const fontInput = document.getElementById('font');
//     const windowView = document.getElementById('window-view');
//     const dialog = document.getElementById('dialog');
//     const headerWindow = document.getElementById('header-window');
//     const messageInput = document.getElementById('message');

//     fontForm.addEventListener('submit', function(event) {
//         event.preventDefault(); // Предотвращаем отправку формы

//         const file = fontInput.files[0];
//         if (file) {
//             const reader = new FileReader();

//             reader.onload = function(e) {
//                 const fontUrl = e.target.result;
//                 applyCustomFont(fontUrl);
//             };

//             reader.readAsDataURL(file);
//         }
//     });

//     function applyCustomFont(fontUrl) {
//         const fontFace = new FontFace('CustomFont', `url(${fontUrl})`);

//         fontFace.load().then(function(loadedFace) {
//             document.fonts.add(loadedFace);
//             windowView.style.fontFamily = 'CustomFont, sans-serif';
//             dialog.style.fontFamily = 'CustomFont, sans-serif';
//             headerWindow.style.fontFamily = 'CustomFont, sans-serif';
//             messageInput.style.fontFamily = 'CustomFont, sans-serif';
//         }).catch(function(error) {
//             console.error('Ошибка загрузки шрифта:', error);
//         });
//     }

//     // Функция для применения цвета
//     function applyColor(colorInputId, targetElementId, styleProperty) {
//         const colorInput = document.getElementById(colorInputId);
//         const targetElement = document.getElementById(targetElementId);
//         const applyButton = document.getElementById(`apply${colorInputId.charAt(0).toUpperCase() + colorInputId.slice(1)}`);

//         applyButton.addEventListener('click', function(event) {
//             event.preventDefault(); // Предотвращаем отправку формы
//             targetElement.style[styleProperty] = colorInput.value;
//         });
//     }

//     // Применение цвета фона
//     applyColor('backgroundColor', 'window-view', 'backgroundColor');

//     // Применение цвета текста
//     applyColor('textColor', 'dialog', 'color');

//     // Применение цвета ассистента
//     applyColor('assistentColor', 'header-window', 'color');

//     // Применение цвета пользователя
//     applyColor('userColor', 'window-ask', 'backgroundColor');

//     // Обработка изменения названия ассистента
//     const nameForm = document.getElementById('customName');
//     const nameInput = document.getElementById('name');

//     nameForm.addEventListener('submit', function(event) {
//         event.preventDefault(); // Предотвращаем отправку формы
//         headerWindow.textContent = nameInput.value;
//     });

//     /* chat-bot */
//     const messagesContainer = document.getElementById('dialog');
//     const userInput = document.getElementById('message');
//     const sendButton = document.getElementById('send-btn');

//     function addMessage(sender, text) {
//         const messageDiv = document.createElement('div');
//         messageDiv.classList.add(sender === 'bot' ? 'bot-message' : 'user-message');
//         messageDiv.textContent = text;
//         messagesContainer.appendChild(messageDiv);
//         messagesContainer.scrollTop = messagesContainer.scrollHeight; // Скролл вниз
//     }

//     sendButton.addEventListener('click', () => {
//         const userText = userInput.value.trim();
//         if (userText === '') return;

//         addMessage('user', userText);
//         userInput.value = '';

//         // Пример ответа бота
//         setTimeout(() => {
//             addMessage('bot', 'Это пример ответа от чат-бота!');
//         }, 500);
//     });

//     // Отправка сообщения при нажатии Enter
//     userInput.addEventListener('keydown', (event) => {
//         if (event.key === 'Enter') {
//             sendButton.click();
//         }
//     });
// });

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

    // Функция для применения цвета
    function applyColor(colorInputId, targetElementId, styleProperty) {
        const colorInput = document.getElementById(colorInputId);
        const targetElement = document.getElementById(targetElementId);
        const applyButton = document.getElementById(`apply${colorInputId.charAt(0).toUpperCase() + colorInputId.slice(1)}`);

        applyButton.addEventListener('click', function(event) {
            event.preventDefault(); // Предотвращаем отправку формы
            targetElement.style[styleProperty] = colorInput.value;
        });
    }
    // function applyColor(colorInputId, targetClassName, styleProperty) {
    //     const colorInput = document.getElementById(colorInputId);
    //     const applyButton = document.getElementById(`apply${colorInputId.charAt(0).toUpperCase() + colorInputId.slice(1)}`);
    
    //     applyButton.addEventListener('click', function(event) {
    //         event.preventDefault(); // Предотвращаем отправку формы
    //         const elements = document.querySelectorAll(`.${targetClassName}`);
    //         elements.forEach(element => {
    //             element.style[styleProperty] = colorInput.value;
    //         });
    //     });
    // } 

    // Применение цвета фона
    applyColor('backgroundColor', 'window-view', 'backgroundColor');

    // Применение цвета текста
    applyColor('textColor', 'dialog', 'color');

    // Применение цвета ассистента
    applyColor('assistentColor', 'header-window', 'color');
    // applyColor('botMessageColor', 'bot-message', 'backgroundColor');
    // applyColor('userMessageColor', 'user-message', 'backgroundColor');

    // Применение цвета пользователя
    applyColor('userColor', 'window-ask', 'backgroundColor');

    // Обработка изменения названия ассистента
    const nameForm = document.getElementById('customName');
    const nameInput = document.getElementById('name');

    nameForm.addEventListener('submit', function(event) {
        event.preventDefault(); // Предотвращаем отправку формы
        headerWindow.textContent = nameInput.value;
    });

    /* chat-bot */
    const messagesContainer = document.getElementById('dialog');
    const userInput = document.getElementById('message');
    const sendButton = document.getElementById('send-btn');

    function addMessage(sender, text) {
        const messageDiv = document.createElement('div');
        messageDiv.classList.add(sender === 'bot' ? 'bot-message' : 'user-message');
        messageDiv.textContent = text;
        messagesContainer.appendChild(messageDiv);
        messagesContainer.scrollTop = messagesContainer.scrollHeight; // Скролл вниз
    }

    // function addMessage(sender, text) {
    //     const messageDiv = document.createElement('div');
    //     messageDiv.classList.add(sender === 'bot' ? 'bot-message' : 'user-message');
    //     messageDiv.textContent = text;
    
    //     // Применяем текущие цвета
    //     const currentColor = document.getElementById(sender === 'bot' ? 'botMessageColor' : 'userMessageColor').value;
    //     messageDiv.style.backgroundColor = currentColor;
    
    //     messagesContainer.appendChild(messageDiv);
    //     messagesContainer.scrollTop = messagesContainer.scrollHeight; // Скролл вниз
    // }    

    sendButton.addEventListener('click', () => {
        const userText = userInput.value.trim();
        if (userText === '') return;

        addMessage('user', userText);
        userInput.value = '';

        // Пример ответа бота
        setTimeout(() => {
            addMessage('bot', 'Это пример ответа от чат-бота!');
        }, 500);
    });

    // Отправка сообщения при нажатии Enter
    userInput.addEventListener('keydown', (event) => {
        if (event.key === 'Enter') {
            sendButton.click();
        }
    });
});
