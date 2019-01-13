import React, { Component } from "react";
import { Redirect } from "react-router-dom";
import PropTypes from "prop-types";

class Login extends Component {
	// props:
	//   is_authenticated
	//   is_fetching
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
		event.preventDefault();
		let { email, password } = this.state;
		let { requestLogin } = this.props;
		requestLogin(email, password)
		console.log(`Submitting ${email} ${password}`);
	}
	handleEmailChange(event) {
		this.setState({
			email: event.target.value	
		});
	}
	handlePasswordChange(event) {
		this.setState({
			password: event.target.value	
		});
	}
	render() {
		let { is_authenticated, is_fetching } = this.props;
		if (is_authenticated) {
			return <Redirect to="/"/>;
		}
		if (is_fetching) {
			return <p>Waiting...</p>;
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

Login.propTypes = {
	is_authenticated: PropTypes.bool.isRequired,
	is_fetching:      PropTypes.bool.isRequired,
	requestLogin:     PropTypes.func.isRequired
}

export default Login;
