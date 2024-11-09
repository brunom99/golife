let socket;
let tryToConnecting = false;
let bubbles = {};

function gID(id) {
    return document.getElementById(id);
}

function showError(err) {
    gID("error").innerHTML = err;
}

async function fetchAPI(apiName) {
    try {
        const response = await fetch("api/" + apiName);
        if (!response.ok) {
            throw new Error("response is not ok");
        }
        return await response.json();
    } catch (error) {
        showError(apiName + ": " + error);
    }
}

function onload() {
    lastActivity();
}

async function lastActivity() {
    const activity = await fetchAPI("activity");
    const divActivity = gID("activity");
    if (activity) {
        divActivity.innerHTML = 'last server activity: ' + new Date(activity.last_activity).toTimeString().slice(0, 8);
    } else {
        divActivity.innerHTML = '';
    }
    setTimeout(lastActivity, 2000);
}

async function connect() {
    if (tryToConnecting) return;
    tryToConnecting = true;
    // disconnect first
    disconnect();
    // script websocket
    socket = new WebSocket("ws://" + location.host + "/ws");
    socket.onmessage = (event) => {
        onMessage(JSON.parse(event.data));
    };
    socket.onopen = () => {
        tryToConnecting = false;
    };
}

function disconnect() {
    if (socket) socket.close();
    bubbles = {};
    gID("grid").innerHTML = "";
    gID("info").innerHTML = "";
    showError('');
}

function initGrid(size) {
    const divGrid = gID("grid");
    divGrid.style.setProperty("--grid-rows", size);
    divGrid.style.setProperty("--grid-cols", size);
    for (let row = 0; row < size; row++) {
        for (let col = 0; col < size; col++) {
            let divCell = document.createElement("div");
            divCell.id = "cell_" + row + "_" + col;
            divGrid.appendChild(divCell).className = "grid-item";
        }
    }
}

function onMessage(msg) {
    const bubble = msg.bubble;
    // info client
    gID("info").innerHTML = "seed: " + msg.info.seed + " | total bubbles: " + msg.info.total_bubbles;
    // no bubble in msg -> msg config ?
    if (!bubble || !bubble.id || bubble.id.length === 0) {
        // init grid
        if (msg.info) {
            // create grid
            initGrid(msg.info.grid_size);
        }
        return;
    }
    // bubble id & pos
    const bubbleID = bubble.id;
    const pos = bubble.pos;
    // bubble info
    let bubbleInfo = bubbles[bubbleID];
    let samePosition = false;
    let sameInvisibility = false;
    // bubble not in dictionary ?
    if (!bubbleInfo) {
        // bubble is finish ?
        if (bubble.is_finish) {
            // ignore message
            return;
        }
        // create div bubble
        const divBubble = document.createElement("div");
        divBubble.id = "bubble_" + bubbleID;
        divBubble.className = "bubble bubble-" + bubble.rarity;
        // create bubble info
        bubbleInfo = {
            bubble: bubble,
            div: divBubble,
        }
    } else {
        // bubble is finish ?
        if (bubble.is_finish) {
            if (bubbleInfo.divParent) bubbleInfo.divParent.removeChild(bubbleInfo.div);
            delete bubbles[bubbleID];
            return;
        }
        // bubble position has change ?
        samePosition = pos.column === bubbleInfo.bubble.pos.column && pos.row === bubbleInfo.bubble.pos.row;
        sameInvisibility = bubble.is_invisible === bubbleInfo.bubble.is_invisible;
        bubbleInfo.bubble = bubble;
    }
    // move bubble
    if (!samePosition) {
        // new div parent
        const divCell = gID("cell_" + bubble.pos.row + "_" + bubble.pos.column);
        divCell.appendChild(bubbleInfo.div);
        // save div parent (to remove child if bubble is finish)
        bubbleInfo.divParent = divCell;
    }
    // change invisibility status ?
    if (!sameInvisibility) {
        if (bubble.is_invisible) bubbleInfo.div.classList.add("invisible");
        if (!bubble.is_invisible) bubbleInfo.div.classList.remove("invisible");
    }
    // set bubble info
    bubbles[bubbleID] = bubbleInfo;
}