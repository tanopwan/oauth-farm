function onSuccess(googleUser) {
  console.log('Google Logged in as: ' + googleUser.getBasicProfile().getName());
  
  var auth = googleUser.getAuthResponse(false)
  var idToken = auth.id_token;
  console.log('With id_token: ' + idToken);
      
  function post(path, params, method) {
    method = method || "post"; // Set method to post by default if not specified.

    // The rest of this code assumes you are not using a library.
    // It can be made less wordy if you use one.
    var form = document.createElement("form");
    form.setAttribute("method", method);
    form.setAttribute("action", path);

    for(var key in params) {
      if(params.hasOwnProperty(key)) {
        var hiddenField = document.createElement("input");
        hiddenField.setAttribute("type", "hidden");
        hiddenField.setAttribute("name", key);
        hiddenField.setAttribute("value", params[key]);

        form.appendChild(hiddenField);
      }
    }

    document.body.appendChild(form);
    form.submit();
  }
  
  post('/openid/login/google', {data: idToken})
}

function onFailure(error) {
	console.log(error);
}

function start() {
  gapi.load('auth2', function() {
    console.log("gapi.auth2.init")
    auth2 = gapi.auth2.init()

    if ($('#signinButton').length > 0) {
      // Attach signin event
      auth2.attachClickHandler('signinButton', {}, onSuccess, onFailure);
    }
  });
}