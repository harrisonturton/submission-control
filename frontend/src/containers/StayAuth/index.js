import React, { Component } from "react";
import { PrivateRoute } from "containers";
import { authenticate, refresh } from "api";

export default class StayAuth extends Component {
	// props:
	//   username: the username for authenticating
	//   password: the password for authenticating
	//   children: the components to render when authenticated
	constructor(props) {
		super(props);
		this.state = {
			isAuthenticated: false,
			token: null
		};
		this.authenticate   = this.authenticate.bind(this);
		this.reauthenticate = this.reauthenticate.bind(this);
	}
	async authenticate(email, password) {
		const { ok, token }	= await authenticate(
			this.props.username,
			this.props.password,
		);
		if (ok) {
			setTimeout(() => this.reauthenticate(token), 5000);
		}
		this.setState({
			isAuthenticated: ok,
			token: token
		});
	}
	async reauthenticate(token) {
		const { ok, token }	= await refresh(token);
		if (ok) {
			setTimeout(() => this.reauthenticate(token), 5000);
		}
		this.setState({
			isAuthenticated: ok,
			token: token
		});
	}
	render() {
		let { isAuthenticated }	= this.state;
	}
}

/*
export class StayAuthenticated extends Component {
	constructor(props) {
		super(props);
		this.state = {
			isAuthenticated: false,
			token: null
		};
		this.authenticate   = this.authenticate.bind(this);
		this.reauthenticate = this.reauthenticate.bind(this);
	}
	async authenticate(email, password) {
		const { ok, token }	= await authenticate(
			"harrisonturton@gmail.com",
			"password"
		);
		if (ok && this.props.stayAuthenticated) {
			setTimeout(() => this.reauthenticate(token), 5000);
		}
		this.setState({
			isAuthenticated: ok,
			token: token
		});
	}
	async reauthenticate(token) {
		const { ok, token }	= await refresh(token);
		if (ok && this.props.stayAuthenticated) {
			setTimeout(() => this.reauthenticate(token), 5000);
		}
		this.setState({
			isAuthenticated: ok,
			token: token
		});
	}
	render() {
		let { isAuthenticated } = this.state;
		return isAuthenticated ? this.props.AuthComponent : this.props.UnauthComponent;
	}
}*/
