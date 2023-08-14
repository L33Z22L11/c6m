function toast(...msg) {
  alert(msg)
}

function login() {
  $.ajax({
    url: "/login",
    type: "POST",
    data: {
      username: $("#username").val(),
      password: $("#password").val()
    },
    success: function (response) {
      uid = response.uid
      username = response.username
      // 将 token 保存到请求头的 Authorization 字段
      $.ajaxSetup({ headers: { "Authorization": "Bearer " + response.token } })

      // 建立websocket
      socket = new WebSocket(`ws://${location.host}/ws?token=Bearer ${response.token}`)
      // 接收到消息时的处理逻辑
      socket.onmessage = function (event) {
        displayMessage(JSON.parse(event.data))
      }
      $("#loginbtn").text("注销")
      $("#loginbtn").attr('onclick', 'location.reload()')
      $('#loginpanel').css('display', 'none')
      getFriendList()
    },
    error: function (xhr, status, error) {
      toast(xhr.responseText)
    }
  })
}

function getFriendList() {
  var select = $('#msgDest')
  select.html("")
  $.ajax({
    url: "/friend/all",
    type: "get",
    success: function (response) {
      $('<option></option>').attr('value', 0).text("请选择好友").appendTo(select)
      $.each(response, function (key, value) {
        $('<option></option>')
          .attr('value', key)
          .text(value)
          .appendTo(select)
      });

      $('#mySelect').change(function () {
        dest = $(this).val();
        console.log('选择的值为：', dest);
      });
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
  var msg = {
    type: $("#msgType").val(),
    time: new Date().getTime(),
    src: uid,
    dest: $("#msgDest").val(),
    content: $("#msgContent").val(),
  }
  socket.send(JSON.stringify(msg))
  displayMessage(msg)
  $("#msgContent").val("");
}

// 显示消息
function displayMessage(msg) {
  console.log(msg)
  var $chatlog = $("#chatlog")
  $chatlog.append(`<div>
    <div class="dim dp05">[${msg.type}] ${msg.src} ${new Date(msg.time)}</div>
    ${msg.content}
  </div>`)
}
