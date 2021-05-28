function queryParams() {
  var username = localStorage.getItem("username");
  var token = localStorage.getItem("token");
  return 'phone=' + username + '&token=' + token;
}