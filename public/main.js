const messages = document.querySelector('#messages')
const send = document.querySelector('#send')

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

// ws.onmessage = function (msg) {
//     console.log(msg)
//
//     // insertMessage(JSON.parse(msg))
// };

// ws.onopen = function () {
//     ws.send('Hello Server')
//     console.log('connected')
//     const rawFile = new XMLHttpRequest();
//     rawFile.open("GET", "./logs/server.log", false);
//     rawFile.onreadystatechange = function ()
//     {
//         if(rawFile.readyState === 4)
//         {
//             if(rawFile.status === 200 || rawFile.status == 0)
//             {
//                 const allText = rawFile.responseText;
//                 insertMessage({content: allText})
//             }
//         }
//     }
//     insertMessage({content: "Connected to server"})
//
// };

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
