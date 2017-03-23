$(function () {
		var jwt = getJWT();
		// { cookieHas.'auth' ≡ jwt ≠ '' ∧ jwt ≠ undefined }
		if (jwt !== '' && jwt !== undefined) {
				window.location.replace('/dash');
		}
		// { cookieHas.'auth' ≡ location = '/dash' }
});

function sendCredentials(){
		var user = $('#user').val();
		var pass = $('#pass').val();
		$.post('/a/auth',
					 JSON.stringify({'user':user, 'pass':pass}),
					 function(data) {
							 window.location.replace('/dash');
					 }
		);
}
