function getJWT() {
		var arr = document.cookie.split(';');
		var t = ['',''];
		// P ≡ 0 ≤ i ≤ arr.length ∧ (i = arr.length ∨
		//   ⟨∃i:0 ≤ i ≤ arr.length:startWith.arr[i].'auth='⟩)
		for(var i = 0; i != arr.length && t[0] !== 'auth'; i++){
				t = arr[i].split('=');
		}
		// { P }
		return t[1];
}

function logOut() {
		document.cookie = 'auth=; expires=Thu, 01 Jan 1970 00:00:01 GMT;';
		window.location.replace('/');
		// { ¬cookieHas.'auth' ∧ location = '/' }
}
