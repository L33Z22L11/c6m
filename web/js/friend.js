function addFriend() {
  var friend_name = prompt("好友名")
  if (!friend_name) return
  $.ajax({
    url: "/friend/add",
    type: "POST",
    data: {
      friend_name: friend_name
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
  var friend_name = prompt("好友名")
  if (!friend_name) return
  $.ajax({
    url: "/friend/del",
    type: "POST",
    data: {
      friend_name: friend_name,
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
  var friend_uid = prompt("好友uid")
  if (!friend_uid) return
  var accept = prompt("同意填1,拒绝填0")
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

