import React, { Component } from "react";
import { Redirect } from "react-router-dom";
import "./style.css";

export default class LoginScreen extends Component {
	// props
	//   auth: the Auth() class used for setting authentication
	constructor(props) {
		super(props);
		this.state = {
			redirectToReferrer: false,
			email: "",
			password: ""
		};
		this.handleSubmit = this.handleSubmit.bind(this);	
		this.handleEmailChange = this.handleEmailChange.bind(this);	
		this.handlePasswordChange = this.handlePasswordChange.bind(this);	
	}
	async handleSubmit(event) {
		event.preventDefault(); // Required to stop the page from refreshing
		let { email, password }	= this.state;
		let { auth } = this.props;
		let token = await auth.signIn(email, password);
		if (token === undefined) {
			return;	
		}
		console.log(`Signed in with JWT ${token}`);
		this.setState({
			redirectToReferrer: true
		});
	}
	async handleEmailChange(event) {
		this.setState({
			email: event.target.value	
		});
	}
	async handlePasswordChange(event) {
		this.setState({
			password: event.target.value	
		});
	}
	render() {
		let { redirectToReferrer } = this.state;
		if (redirectToReferrer) {
			return (
				<div>
					<p>Redirecting...</p>
					<Redirect to="/"/>
				</div>
			);
		}
		return (
			<form onSubmit={this.handleSubmit}>
				<input
					name="email"
					type="email"
					onChange={this.handleEmailChange}
					value={this.state.email}
				/>
				<input
					name="password"
					type="password"
					onChange={this.handlePasswordChange}
					value={this.state.password}
				/>
				<input type="submit" value="Submit"/>
			</form>
		);
	}
}
