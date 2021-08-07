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
var userData = {};

const connectWs = (e) => {
  let wsBase = base.replace("https", "ws");
  wsBase = base.replace("http", "ws");
  wsBase = wsBase.substring(0, wsBase.length - 6);

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
    userData = JSON.parse(evt.data);
    userData = JSON.parse(userData.Payload);
    if (userData.status === "all") {
    }
    let doc = document.getElementById("data");
    document.getElementById("head").innerHTML = `
    <h3>choose</h3>
    <h1 style="display: inline" onclick="sendHand(true)">✋🏻</h1>
    <h1 style="display: inline" onclick="sendHand(false)">✋🏿</h1>`;
    doc.innerHTML = ``;
    for (const [key, value] of Object.entries(userData.users)) {
      let temp = `
      <div>
        <p>${key}</p>
        ${
          userData.status === "all"
            ? `<h1>${value.hand === true ? "✋🏻" : "✋🏿"}</h1>`
            : `<h1>${value.hand === undefined ? "💤" : "✋"}</h1>`
        }
      </div>`;

      doc.innerHTML += temp;
    }
  };

  ws.onopen = () => {
    ws.send(userName);
    ws.send(JSON.stringify({ name: userName }));
    document.getElementById("form").innerHTML = "";
  };
};

const sendHand = (data) => {
  ws.send(JSON.stringify({ name: userName, name: data }));
};
