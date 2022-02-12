const socket = new WebSocket("ws://localhost:8080/echo");
let username = document.getElementById("username");
let messageToSend = document.getElementById("messageToSend");
let messagesDiv = document.getElementById("messages");
let signinForm = document.getElementById("signinForm");
let messageForm = document.getElementById("messageForm");
let usernameBox = document.getElementById("nameBox");

socket.onopen = function () {
  username.value = "aniket";
  signin({
    preventDefault: function () { }
  });
};

socket.onmessage = function (data) {
  messages = data.data.split("|").filter(e => e.length > 0);
  messages.forEach(message => {
    const messageElement = document.createElement("div");
    messageElement.innerText = message;
    messageElement.classList.add("message")
    messagesDiv.appendChild(messageElement);
  });

  messagesDiv.scrollTop = messagesDiv.scrollHeight;
};

signinForm.addEventListener("submit", signin);
messageForm.addEventListener("submit", send);
document.getElementById("signoutButton").addEventListener("click", signout);

function signin(e) {
  e.preventDefault();
  if (username.value.length === 0) return;

  socket.send(`Init ${username.value}`);

  document.getElementById("chatbox").style.display = "block";
  document.getElementById("signin").style.display = "none";
  usernameBox.innerText = username.value;
}

function signout() {
  document.getElementById("chatbox").style.display = "none";
  document.getElementById("signin").style.display = "block";
  username.value = "";
  messagesDiv.innerHTML = "";
}

function send(e) {
  e.preventDefault();

  socket.send(messageToSend.value);
  messageToSend.value = "";
}