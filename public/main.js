const messages = document.querySelector('#messages')

const socket = io("ws://localhost:8080",{
    path: "/socket.io",
    transports: ['websocket']
})

socket.on('connect', () => {
    console.log('connected')
})

socket.on('connect', () => {
    socket.on('some', (message) => {
        console.log(message)
        insertMessage({content:message})
    })
})

/**
 * Insert a message into the UI
 * @param {Message that will be displayed in the UI} messageObj
 */
function insertMessage(messageObj) {
    // Create a div object which will hold the message
    const message = document.createElement('div')

    // Set the attribute of the message div
    message.setAttribute('class', 'chat-message')
    console.log(" content: " +" messageObj.content")
    message.textContent = messageObj.content

    // Append the message to our chat div
    messages.appendChild(message)

    // Insert the message as the first message of our chat
    messages.insertBefore(message, messages.firstChild)
}
