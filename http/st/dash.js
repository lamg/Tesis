$(function () {
		var jwt = getJWT();
		if (jwt === '') {
				window.location.replace('/');
		}
});

function postAccMatches() {
		//select matches ids from page
		console.log("postAccMatches");
		
		$.ajax ({
				
		});
}
