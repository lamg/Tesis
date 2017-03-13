$(function () {
		var jwt = getJWT();
		if (jwt === '') {
				window.location.replace('/');
		}
});
