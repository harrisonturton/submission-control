
// Backend API endpoints
const host        = "http://localhost";
const auth_url    = host + "/auth";
const refresh_url = host + "/refresh";

// Where we store the JWT token in localStorage
const token_key = "submission-control-token";

// Time between token refreshes. Do it every 3 minutes.
export const refresh_time = 1000 * 60 * 3;

// How long until a JWT token is invalid. 10 minutes.
export const token_timeout = 1000 * 60 * 10;

export async function refresh(token) {
	let resp = await fetch(refresh_url, postRequest({
		token: token	
	}));
	if (!resp.ok) {
		console.log(resp.statusText);
		console.log("Not OK")
		return null;
	}
	let body = await resp.json();
	if (body.token === undefined) {
		console.log("Token undefined");
		return null;
	}
	return body.token;
}

// signIn will make a request to the API. It will
// return the JWT token on success, or null on failure.
export async function signIn(email, password) {
	let resp = await fetch(auth_url, postRequest({
		email: email,
		password: password
	}));
	if (!resp.ok) {
		console.log(resp.statusText);
		console.log("Not OK")
		return null;
	}
	let body = await resp.json();
	if (body.token === undefined) {
		console.log("Token undefined");
		return null;
	}
	return body.token;
}

// store_token will save a new token to localStorage
export function storeToken(token) {
	localStorage.setItem(token_key, token);
}

export function forgetToken() {
	localStorage.removeItem(token_key);
}

// retreive_token will retreive the JWT token from localStorage
export function retreiveToken() {
	return localStorage.getItem(token_key);
}

export function hasToken() {
	return retreiveToken() !== null && retreiveToken() !== undefined;
}

// Creates a POST request with JSON data
const postRequest = data => ({
	method: "POST",
	headers: {
		"Content-Type": "application/json"	
	},
	body: JSON.stringify(data)
});
