var jwt = sessionStorage['Auth'];
if (jwt !== undefined) {
		setDash(jwt);
}

function setDash(jwt) {
		//TODO redirect to dash
		window.location.replace('/dash');
		console.log(jwt);
}
