
// Backend API endpoints
const host        = "http://localhost";
const auth_url    = host + "/auth";
const refresh_url = host + "/refresh";

// Refresh our JWT token every 5 minutes.
const refresh_delay = 1000 * 60 * 5;

// Where we store the JWT token in localStorage
const token_key = "submission-control-token";

export async function signIn(email, password) {
	let resp = await fetch("http://localhost/auth", {
		method: "POST",
		headers: {
			"Content-Type": "application/json",
		},
		body: JSON.stringify({
			email: "harrisonturton@gmail.com",
			password: "password"
		})
	});
	if (!resp.ok) {
		return null;	
	}
	let body = await resp.json();
	if (!body.token === undefined) {
		return null;
	}
	return body.token;
}

export default class Auth {
	constructor() {
		this.is_authenticated = false;
		this.token = null;
	}
	async signIn(email, password) {
		let resp = await fetch(auth_url, loginRequest(email, password))
		if (!resp.ok) {
			return null;	
		}
		let body = await resp.json();
		if (!body.token === undefined) {
			return null;
		}
		this.is_authenticated = true;
		this.token = body.token;
		localStorage.setItem(token_key, body.token)
		return body.token;
	}
	async signOut() {
		this.is_authenticated = false;
		localStorage.setItem(token_key, "")
	}
	async refresh(token) {
		let resp = await fetch(refresh_url, refreshRequest(token))
		if (!resp.ok) {
			return null;	
		}
		let body = await resp.json();
		if (!body.token === undefined) {
			return null;
		}
		this.is_authenticated = true;
		this.token = body.token;
		localStorage.setItem(token_key, body.token)
		return body.token;
	}
}

const loginRequest = (email, password) => postRequest({
	email: email,
	password: password
});

const refreshRequest = (token) => postRequest({
	token: token
});

const postRequest = (body) => ({
	method: "POST",
	headers: {
		"Content-Type": "application/json",
	},
	body: JSON.stringify(body)
});
