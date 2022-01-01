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
  let wsBase = base.replace("http", "ws");
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
      let up = [];
      let down = [];
      for (const [key, value] of Object.entries(userData.users)) {
        if (value.hand) {
          up.push(key);
        } else {
          down.push(key);
        }
      }
      if (up.length !== 1 && down.length !== 1) {
        document.getElementById("result").innerHTML = `no one win`;
      } else if (up.length === 1) {
        document.getElementById("result").innerHTML = `${up[0]} wins!`;
      } else {
        document.getElementById("result").innerHTML = `${down[0]} wins!`;
      }
    }
    let doc = document.getElementById("data");
    if (userData.users[userName].hand === undefined) {
      document.getElementById("head").innerHTML = `
      <h3>choose</h3>
      <h1 style="display: inline" onclick="sendHand(true)">✋🏻</h1>
      <h1 style="display: inline" onclick="sendHand(false)">✋🏿</h1>`;
    } else {
      document.getElementById("head").innerHTML = `
      <h3>you choose</h3>
      <h1 style="display: inline"">${userData.users[userName].hand === true ? "✋🏻" : "✋🏿"}</h1>`;
    }

    doc.innerHTML = ``;
    for (const [key, value] of Object.entries(userData.users)) {
      let temp = `
      <div class="mx-4">
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
  ws.send(JSON.stringify({ name: userName, hand: data }));
};

const generateRoom = async () => {
  const res = await fetch(`${base}api/room/create`)
  const data = await res.json()
  window.location = `${base}${data.room}`;
}