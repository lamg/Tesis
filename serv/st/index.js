var jwt = sessionStorage['Auth'];
if (jwt !== undefined) {
		setDash(jwt);
}

$(function(){
		if (jwt === undefined) {
				setLogScreen();
		}
});

function setLogScreen() {
		var jwt;
		$('#login').click(function() {
				var r;
				$.ajax({
						type:'POST',
						url:'/a/auth',
						data: JSON.stringify(
								{ 'user': $('#user')[0].value,
									'pass': $('#password')[0].value
								}),
						success: function(data, textStatus, request){
								jwt = request.getResponseHeader('Auth');
								$('#resp').text(data);
								if (jwt != null) {
										sessionStorage['Auth'] = jwt;
										//ready to set dashboard
										setDash(jwt);
								}
						},
						error: function (request, textStatus, err){
								console.log('Error making POST');
						}
				});
				//posted
		});
}

function setDash(jwt) {
		//TODO redirect to dash
		window.location.replace('/dash');
		console.log(jwt);
}
