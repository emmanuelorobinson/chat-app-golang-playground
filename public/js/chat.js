const socket = io.connect('http://127.0.0.1:8000/socket.io');

//connect to golang socket io server
socket.on('chat', function() {
    console.log('connected to server');
});

const getJ = async function (url, callback) {
    const data = await fetch(url);
    const json = await data.json();
    callback(json);
};

getJ('http://127.0.0.1:8000/jsons', function (json) {
    console.log(json);
});