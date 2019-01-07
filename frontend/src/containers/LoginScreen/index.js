import React, { Component } from "react";
import { Redirect } from "react-router-dom";
import { authenticate } from "api";
import "./style.css";

export default class LoginScreen extends Component {
	constructor(props) {
		super(props);
		this.state = {
			isAuthenticated: false,
			email: "",
			password: ""
		};
		this.handleSubmit = this.handleSubmit.bind(this);
		this.handleEmailChange = this.handleEmailChange.bind(this);
		this.handlePasswordChange = this.handlePasswordChange.bind(this);
	}
	async handleSubmit(event) {
		event.preventDefault();
		const { email, password } = this.state;
		const { ok, token } = await authenticate(email, password);
		this.setState({
			isAuthenticated: ok	
		})
		this.props.onAuth(ok, token)
	}
	handleEmailChange(event) {
		this.setState({
			email: event.target.value
		})
	}
	handlePasswordChange(event) {
		this.setState({
			password: event.target.value
		})
	}
	render() {
		let { isAuthenticated } = this.state;
		if (isAuthenticated) {
			return <Redirect to={{
				pathname: "/",
				state: { from: this.props.location }
			}}/>;
		}
		return (
			<div className="login-wrapper">
				<div className="image"></div>
				<div className="form-wrapper">
					<h1>Submission Control</h1>
					<p>This website is where you can view & submit assignments, labs, and give feedback on your course.</p>
					<form onSubmit={this.handleSubmit}>
						<input
							name="email"
							type="email"
							placeholder="uXXXXXXX@anu.edu.au"
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
				</div>
			</div>
		);
	}
}

/*
export default class LoginScreen extends Component {
	constructor(props) {
		super(props);
		this.state = {
			email: "",
			password: "",
		}
		this.handleChange = this.handleChange.bind(this);
		this.handleSubmit = this.handleSubmit.bind(this);
	}
	handleChange(event) {
		let { name, value } = event.target;
		this.setState({
			[name]: value
		});
	}
	async handleSubmit(event) {
		event.preventDefault();	
		let { email, password } = this.state;
		let { ok, token } = await authenticate(email, password);
		alert(`${email} ${password}`);
		this.props.onAuth(ok, token)
		console.log(JSON.stringify(this.state));
	}
	render() {
		return (
			<div className="login-wrapper">
				<div className="image"></div>
				<div className="form-wrapper">
					<h1>Submission Control</h1>
					<p>This website is where you can view & submit assignments, labs, and give feedback on your course.</p>
					<form onSubmit={this.handleSubmit}>
						<input
							name="email"
							type="email"
							placeholder="uXXXXXXX@anu.edu.au"
							value={this.state.email}
							onChange={this.handleChange}
						/>
						<input
							name="password"
							type="password"
							value={this.state.password}
							onChange={this.handleChange}
						/>
						<input type="submit"   value="Submit"/>
					</form>
				</div>
			</div>
		);
	}
}*/
