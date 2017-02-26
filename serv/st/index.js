$(document).ready(function () {
		$('#login').click(function() {
				$.post('/auth',
							 JSON.stringify(
									 { 'user': $('#user')[0].value,
										 'pass': $('#password')[0].value
									 }),
							 function(data){
									 $('#resp').text(data);
							 }
							);
		});
		return false;
});
