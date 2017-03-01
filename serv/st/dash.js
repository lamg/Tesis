var jwt = sessionStorage['Auth'];
if (jwt === undefined) {
		window.location.replace('/');
}

$(function () {
		//request user information
		
		$.ajax({
				type: 'GET',
				url: '/a/info',
				dataType: 'json',
				beforeSend: function (req) {
						req.setRequestHeader('Auth', jwt);
				},
				success: function(data, textStatus, req) {
						//TODO present data
						console.log(data);
				}
		});
		//show user information
});
