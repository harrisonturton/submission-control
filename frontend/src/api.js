
const callback_delay = 5000;

const host             = "http://localhost"
const auth_endpoint    = host + "/auth";
const refresh_endpoint = host + "/refresh";
const users_endpoint   = host + "/users";

const auth_request = (email, password) => ({
	method: "POST",
	headers: {
		"Content-Type": "application/json",
	},
	body: JSON.stringify({
		email:    email,
		password: password,
	})
})

const refresh_request = (token) => ({
	method: "POST",
	headers: {
		"Content-Type": "application/json",
	},
	body: JSON.stringify({
		token: token
	})
})

export async function authenticate(email, password, callback=null) {
	const req = auth_request(email, password)
	const res = await fetch(auth_endpoint, req)
	const body = await res.json()
	const result = {
		ok:    res.ok && body.token !== undefined,
		token: body.token
	};
	if (callback !== null && result.ok) {
		setTimeout(() => callback(result.token), callback_delay)	
	}
	return result;
}

export async function refresh(token, callback=null) {
	const req = refresh_request(token)
	const res = await fetch(refresh_endpoint, req)
	const body = await res.json()
	const result = {
		ok:    res.ok && body.token !== undefined,
		token: body.token
	};
	if (callback !== null && result.ok) {
		setTimeout(() => callback(result.token), callback_delay)	
	}
	return result;
}
