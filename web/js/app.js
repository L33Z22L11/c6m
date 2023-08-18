var $history = $("#history")
var $toast = $("#toast");

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
        toast(JSON.parse(event.data))
      }
      // 处理错误事件
      socket.onerror = function (error) {
        // 在控制台打印错误信息
        toast('WebSocket error:', error)
        // 进行适当的错误处理逻辑
        alert("WebSocket 连接错误")
        location.reload()
      }

      // 处理关闭事件
      socket.onclose = function (event) {
        if (event.wasClean) {
          // WebSocket 关闭是有意义的
          console.log('WebSocket closed cleanly')
        } else {
          // WebSocket 关闭是意外的
          console.error('WebSocket connection closed unexpectedly')
        }
        console.log('Close code:', event.code)
        console.log('Close reason:', event.reason)
        // 进行适当的关闭处理逻辑
        alert("WebSocket 连接关闭")
        location.reload()
      }
      $('#loginpanel').css('display', 'none')
      updateList()
    },
    error: function (xhr, status, error) {
      toast(xhr.responseText)
    }
  })
}

function updateList() {
  var select = $('#msgDest')
  select.html("")
  $.ajax({
    url: "/friend/all",
    type: "get",
    success: function (response) {
      $('<option></option>')
        .attr('disabled', true)
        .text("---好友---")
        .appendTo(select)
      $.each(response, function (key, value) {
        $('<option></option>')
          .attr('value', key)
          .text(value)
          .appendTo(select)
      })
    },
    error: function (xhr, status, error) {
      toast(xhr.responseText)
    }
  })
  $.ajax({
    url: "/group/all",
    type: "get",
    success: function (response) {
      $('<option></option>')
        .attr('disabled', true)
        .text("---群---")
        .appendTo(select)
      $.each(response, function (key, value) {
        $('<option></option>')
          .attr('value', key)
          .text(value)
          .appendTo(select)
      })
      showHistory()
    },
    error: function (xhr, status, error) {
      toast(xhr.responseText)
    }
  })
}

function send() {
  var msg = {
    time: new Date().getTime(),
    src: uid,
    dest: $("#msgDest").val(),
    content: `${username}: ${$("#msgContent").val()}`,
  }
  socket.send(JSON.stringify(msg))
  if (msg.dest == $("#msgDest").val()) {
    addHistory(msg)
  }

  $("#msgContent").val("");
}

// 显示消息
function toast(msg) {
  console.log(msg);
  if (msg.src == $("#msgDest").val()) {
    addHistory(msg)
    return
  }

  $toast.prepend(`<div>${msg.content ?
    `<div class="dim dp05">${msg.src} ${new Date(msg.time)}</div>
      ${msg.content}` : msg}</div>`);
  setTimeout(function () {
    $toast.children().last().fadeOut(500, function () {
      $(this).remove();
    });
  }, 3000);

}

function showHistory() {
  var id = $("#msgDest").val()
  if (!id) return
  $.ajax({
    url: "/history",
    type: "GET",
    data: {
      id: id
    },
    success: function (response) {
      // 处理成功响应
      console.log(response);
      // 在此处进行你想要的操作，如更新页面显示等
      $history.html("")
      response.forEach(msg => addHistory(msg))
    },
    error: function (xhr, status, error) {
      // 处理错误响应
      console.log("请求失败:", xhr.responseText);
    }
  });
}

function addHistory(msg) {
  $history.prepend(`<div>${msg.content ?
    `<div class="dim dp05">${msg.src} ${new Date(msg.time)}</div>
    ${msg.content}` : msg}</div>`)
}

function upload() {
  var formData = new FormData()
  formData.append('file', $('#file-input')[0].files[0])
  formData.append('dest', $("#msgDest").val())

  $.ajax({
    url: '/upload',
    type: 'POST',
    data: formData,
    contentType: false,
    processData: false,
    success: function (response) {
      // 处理服务器响应
      toast(response)
    },
    error: function () {
      // 处理错误情况
      toast("上传文件失败")
    }
  })
}
