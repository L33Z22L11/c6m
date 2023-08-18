function joinGroup() {
  var groupName = prompt("群组名")
  var content = prompt("申请内容")
  if (!groupName || !content) return
  $.ajax({
    url: "/group/join",
    type: "POST",
    data: {
      group_name: groupName,
      content: content,
    },
    success: function (response) {
      toast(JSON.stringify(response))
    },
    error: function (xhr, status, error) {
      console.log(xhr.responseText)
      toast(xhr.responseText)
    },
  })
}

function leaveGroup() {
  var groupName = prompt("群组名")
  if (!groupName) return
  $.ajax({
    url: "/group/leave",
    type: "POST",
    data: {
      group_name: groupName,
    },
    success: function (response) {
      toast(JSON.stringify(response))
    },
    error: function (xhr, status, error) {
      toast(xhr.responseText)
    },
  })
}

function infoGroup() {
  var groupName = prompt("群组名")
  if (!groupName) return
  $.ajax({
    url: "/group/info",
    type: "POST",
    data: {
      group_name: groupName,
    },
    success: function (response) {
      toast(JSON.stringify(response))
    },
    error: function (xhr, status, error) {
      console.log(xhr.responseText)
      toast(xhr.responseText)
    },
  })
}

function kickGroup() {
  var groupName = prompt("群组名")
  var member = prompt("成员名")
  if (!groupName || !member) return
  $.ajax({
    url: "/group/kick",
    type: "POST",
    data: {
      group_name: groupName,
      member: member,
    },
    success: function (response) {
      toast(JSON.stringify(response))
    },
    error: function (xhr, status, error) {
      console.log(xhr.responseText)
      toast(xhr.responseText)
    },
  })
}

function createGroup() {
  var groupName = prompt("群组名")
  if (!groupName) return
  $.ajax({
    url: "/group/create",
    type: "POST",
    data: {
      group_name: groupName,
    },
    success: function (response) {
      toast(JSON.stringify(response))
    },
    error: function (xhr, status, error) {
      console.log(xhr.responseText)
      toast(xhr.responseText)
    },
  })
}

function delGroup() {
  var groupName = prompt("群组名")
  if (!groupName) return
  $.ajax({
    url: "/group/del",
    type: "POST",
    data: {
      group_name: groupName,
    },
    success: function (response) {
      toast(JSON.stringify(response))
    },
    error: function (xhr, status, error) {
      console.log(xhr.responseText)
      toast(xhr.responseText)
    },
  })
}

function addGroupAdmin() {
  var groupName = prompt("群组名")
  var admin = prompt("管理员名")
  if (!groupName || !admin) return
  $.ajax({
    url: "/group/admin/add",
    type: "POST",
    data: {
      group_name: groupName,
      admin: admin,
    },
    success: function (response) {
      toast(JSON.stringify(response))
    },
    error: function (xhr, status, error) {
      console.log(xhr.responseText)
      toast(xhr.responseText)
    },
  })
}

function delGroupAdmin() {
  var groupName = prompt("群组名")
  var admin = prompt("管理员名")
  if (!groupName || !admin) return
  $.ajax({
    url: "/group/admin/del",
    type: "POST",
    data: {
      group_name: groupName,
      admin: admin,
    },
    success: function (response) {
      toast(JSON.stringify(response))
    },
    error: function (xhr, status, error) {
      console.log(xhr.responseText)
      toast(xhr.responseText)
    },
  })
}

function getGroupReq() {
  var gid = prompt("群组Gid")
  $.ajax({
    url: "/group/req?gid=" + gid,
    type: "GET",
    success: function (response) {
      toast(JSON.stringify(response))
    },
    error: function (xhr, status, error) {
      toast(xhr.responseText)
    },
  })
}

function respGroupReq() {
  var gid = prompt("群组Gid")
  if (!gid) return
  var newid = prompt("用户uid")
  if (!newid) return
  var isAccept = prompt("是否接受请求？（1-接受，0-拒绝）")
  $.ajax({
    url: "/group/req",
    type: "POST",
    data: {
      gid: gid,
      newid: newid,
      accept: isAccept,
    },
    success: function (response) {
      toast(JSON.stringify(response))
    },
    error: function (xhr, status, error) {
      toast(xhr.responseText)
    },
  })
}