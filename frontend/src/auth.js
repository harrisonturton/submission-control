
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
		this.retreiveFromStorage()
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
	// Check if a token has been storage in localStorage,
	// and retreive it if so.
	async retreiveFromStorage() {
		let token = localStorage.getItem(token_key)
		let hasToken = token !== undefined || token === null;
		this.token = hasToken ? token : null;
		this.is_authenticated = hasToken;
	}
	async verifySignedIn() {
		console.log("Verifying signed in...");
		if (!this.is_authenticated || this.token === null) {
			console.log("Nope!");
			return false;	
		}
		this.retreiveFromStorage();
		if (!this.is_authenticated) {
			console.log("Nope!");
			return false;	
		}
		let token = await this.refresh(this.token);
		console.log(`Maybe! ${token !== null}`);
		console.log(`Authenticated: ${this.is_authenticated}`);
		return token !== null;
	}
	async signOut() {
		this.is_authenticated = false;
		localStorage.setItem(token_key, null)
	}
	async refresh(token) {
		let resp = await fetch(refresh_url, refreshRequest(token))
		if (!resp.ok) {
			this.is_authenticated = false;
			return null;	
		}
		let body = await resp.json();
		if (body.token === undefined) {
			this.is_authenticated = false;
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
