window.onload = function () {
  var conn;
  var currentRoom = null;
  var ROOMS = ["general", "dashboard", "customers"];
  var MAX_MESSAGES = 500;

  var msg = document.getElementById("msg");
  var msgbtn = document.getElementById("msgbtn");
  var log = document.getElementById("log");
  var messages = document.getElementById("messages");

  var NAME = sessionStorage.getItem("name");
  if (!NAME) {
    NAME = "User" + Math.floor(Math.random() * 10000);
    sessionStorage.setItem("name", NAME);
  }

  function setActive(room) {
    ROOMS.forEach(function (r) {
      var li = document.getElementById(r);
      if (li) li.querySelector("a").classList.toggle("is-active", r === room);
    });
  }

  function connectToRoom(room) {
    if (room === currentRoom) return;
    if (conn) {
      conn.onclose = null;
      conn.close(1000, "switch");
    }

    currentRoom = room;
    messages.innerHTML = "";
    document.getElementById("online").innerText = "";
    setActive(room);

    var proto = location.protocol === "https:" ? "wss://" : "ws://";
    conn = new WebSocket(
      proto +
        document.location.host +
        "/ws?name=" +
        encodeURIComponent(NAME) +
        "&room=" +
        encodeURIComponent(room),
    );
    conn.onclose = function (evt) {
      renderSystemText("Connection closed.");
      msg.disabled = true;
      msgbtn.disabled = true;
    };
    conn.onmessage = function (evt) {
      var chunks = evt.data.split("\n");
      for (var i = 0; i < chunks.length; i++) {
        var m = JSON.parse(chunks[i]);

        switch (m.Type) {
          case "Online":
            updateOnline(m);
            break;
          case "System":
            renderSystem(m);
            break;
          default:
            renderMessage(m);
        }
      }
    };
  }

  ROOMS.forEach(function (r) {
    document.getElementById(r).addEventListener("click", function () {
      connectToRoom(r);
    });
  });

  function appendLog(item, forceScroll) {
    var atBottom = log.scrollHeight - log.scrollTop - log.clientHeight < 80;
    messages.appendChild(item);

    while (messages.children.length > MAX_MESSAGES) {
      messages.firstElementChild.remove();
    }

    if (atBottom || forceScroll) {
      log.scrollTop = log.scrollHeight;
    }
  }

  function updateOnline(n) {
    document.getElementById("online").innerText = n.Online;
  }

  function renderSystem(m) {
    var p = document.createElement("p");
    p.classList.add("has-text-centered", "has-text-grey", "is-size-7", "my-3");
    p.innerText = m.Text;
    appendLog(p);
  }

  function renderSystemText(m) {
    var p = document.createElement("p");
    p.classList.add("has-text-centered", "has-text-grey", "is-size-7", "my-3");
    p.innerText = m;
    appendLog(p);
  }

  function renderMessage(m) {
    var mine = m.User === NAME;

    var row = document.createElement("div");
    row.classList.add("msg-row");
    if (mine) row.classList.add("is-mine");

    var item = document.createElement("article");
    item.classList.add("message", mine ? "is-info" : "is-dark");

    var header = document.createElement("div");
    header.classList.add("message-header");

    var name = document.createElement("p");
    name.innerText = mine ? m.User : m.User;

    var time = document.createElement("span");
    time.classList.add("is-size-7");
    time.innerText = new Date(m.Time).toLocaleTimeString("ru-RU", {
      hour: "2-digit",
      minute: "2-digit",
    });

    header.appendChild(name);
    header.appendChild(time);

    var text = document.createElement("div");
    text.classList.add("message-body");
    text.innerText = m.Text;

    item.appendChild(header);
    item.appendChild(text);
    row.appendChild(item);
    appendLog(row, mine);
  }

  document.getElementById("form").onsubmit = function () {
    if (!conn) {
      return false;
    }
    if (!msg.value) {
      return false;
    }
    msg.value = msg.value.trim();
    if (msg.value.length >= 1) {
      conn.send(msg.value);
      msg.value = "";
    }
    return false;
  };

  if (window["WebSocket"]) {
    connectToRoom("general");
  } else {
    var item = document.createElement("div");
    item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
    appendLog(item);
  }
};
