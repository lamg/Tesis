$(function () {
		var jwt = getJWT();
		if (jwt === '') {
				window.location.replace('/');
		}
});

function postAccMatches() {
		//get selected AccMatches
		//post them
		//write syncronized AccMatches
		//in a decent widget
		console.log("postAccMatches");
		
		$.ajax ({
				
		});
}
