var socket = {
  send: () => toast("尚未建立socket连接!"),
  close: () => { }
}

function toast(...msg) {
  alert(msg)
}

function register() {
  $.ajax({
    url: "/register",
    type: "POST",
    data: {
      username: $("#username").val(),
      password: $("#password").val()
    },
    success: function (response) {
      toast(JSON.stringify(response))
    },
    error: function (xhr, status, error) {
      console.log(xhr.responseText)
      toast(xhr.responseText)
    }
  })
}

function login() {
  socket.close()
  $.ajax({
    url: "/login",
    type: "POST",
    data: {
      username: $("#username").val(),
      password: $("#password").val()
    },
    success: function (response) {
      toast(JSON.stringify(response))
      // 将 token 保存到请求头的 Authorization 字段
      $.ajaxSetup({ headers: { "Authorization": "Bearer " + response.token } })

      // 建立websocket
      socket = new WebSocket(`ws://${location.host}/ws?token=Bearer ${response.token}`)
      // 接收到消息时的处理逻辑
      socket.onmessage = function (event) {
        var message = event.data
        displayMessage(message)
      }
    },
    error: function (xhr, status, error) {
      toast(xhr.responseText)
    }
  })
}

function getFriendList() {
  $.ajax({
    url: "/friend/all",
    type: "get",
    success: function (response) {
      toast(JSON.stringify(response))
    },
    error: function (xhr, status, error) {
      toast(xhr.responseText)
    }
  })
}

function getFriendReq() {
  $.ajax({
    url: "/friend/req",
    type: "get",
    success: function (response) {
      toast(JSON.stringify(response))
    },
    error: function (xhr, status, error) {
      toast(xhr.responseText)
    }
  })
}

function respFriendReq() {
  friend_uid = prompt("好友uid")
  if (!friend_uid) return
  accept = prompt("同意填1,拒绝填0")
  if (accept != '0' && accept != '1') return
  $.ajax({
    url: "/friend/req",
    type: "POST",
    data: {
      friend_uid: friend_uid,
      accept: accept
    },
    success: function (response) {
      toast(JSON.stringify(response))
    },
    error: function (xhr, status, error) {
      console.log(xhr.responseText)
      toast(xhr.responseText)
    }
  })
}

function addFriend() {
  $.ajax({
    url: "/friend/add",
    type: "POST",
    data: {
      friend_name: prompt("好友名"),
    },
    success: function (response) {
      toast(JSON.stringify(response))
    },
    error: function (xhr, status, error) {
      console.log(xhr.responseText)
      toast(xhr.responseText)
    }
  })
}

function delFriend() {
  $.ajax({
    url: "/friend/del",
    type: "POST",
    data: {
      friend_name: prompt("好友名"),
    },
    success: function (response) {
      toast(JSON.stringify(response))
    },
    error: function (xhr, status, error) {
      console.log(xhr.responseText)
      toast(xhr.responseText)
    }
  })
}

function send() {
  var msg = JSON.stringify({
    type: $("#msgType").val(),
    // src: uid,
    dest: $("#msgDest").val(),
    content: $("#msgContent").val(),
  })
  socket.send(msg)
  displayMessage(msg)
  $("#msgContent").val("");
}

// 显示消息
function displayMessage(msg) {
  var $chatLog = $("#chatLog")
  $chatLog.append("<p>" + msg + "</p>")
}
