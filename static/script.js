var base = window.location.href;

const joinRoom = (e) => {
  const form = document.getElementById("form");
  e.preventDefault();
  const formData = new FormData(form);
  let data = {};
  for (let key of formData.keys()) {
    data[key] = formData.get(key);
  }
  window.location = `${base}${data.room}`;
};

var userName = undefined;
var ws = undefined;
var data = {};

const connectWs = (e) => {
  let wsBase = base.replace("https", "ws");
  wsBase = base.replace("http", "ws");
  wsBase = wsBase.substring(0, wsBase.length - 6)

  const form = document.getElementById("form");
  e.preventDefault();
  const formData = new FormData(form);
  let data = {};

  for (let key of formData.keys()) {
    data[key] = formData.get(key);
  }
  userName = data.name;
  ws = new WebSocket(`${wsBase}api/room/${roomId}`);

  ws.onmessage = (evt) => {
    console.log("data", evt.data);
    data = JSON.parse(evt.data);
    document.getElementById("form").innerHTML = JSON.stringify(data);
  };

  ws.onopen = () => {
    ws.send(userName);
    ws.send(JSON.stringify({ name: userName }));
  };
};
