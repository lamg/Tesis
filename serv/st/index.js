$(function () {
		var jwt = getJWT();
		// { cookieHas.'auth' ≡ jwt ≠ '' ∧ jwt ≠ undefined }
		if (jwt !== '' && jwt !== undefined) {
				window.location.replace('/dash');
		}
		// { cookieHas.'auth' ≡ location = '/dash' }
});
