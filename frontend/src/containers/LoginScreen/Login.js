import React, { Component } from "react";
import { Redirect } from "react-router-dom";
import { Loader } from "components";
import PropTypes from "prop-types";
import "./style.css";

class Login extends Component {
	// props:
	//   is_authenticated
	//   is_fetching
	constructor(props) {
		super(props);
		this.state = {
			redirectToReferrer: false,
			uid: "",
			password: ""
		};
		this.handleSubmit = this.handleSubmit.bind(this);
		this.handleUIDChange = this.handleUIDChange.bind(this);	
		this.handlePasswordChange = this.handlePasswordChange.bind(this);	
	}
	async handleSubmit(event) {
		event.preventDefault();
		let { uid, password } = this.state;
		let { requestLogin } = this.props;
		requestLogin(uid, password)
		console.log(`Submitting ${uid} ${password}`);
	}
	handleUIDChange(event) {
		this.setState({
			uid: event.target.value	
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
		return (
			<div className="login-wrapper">
				<div className="login-splash-image"></div>
				<div className="login-panel-wrapper">
					<h1>Submission Control</h1>
					<p>This website is where you can view & submit assignments, and give feedback on lectures.</p>
					{is_fetching
							? <Loader/>
							: <form onSubmit={this.handleSubmit}>
						<label>UID</label>
						<input
							name="uid"
							type="text"
							onChange={this.handleUIDChange}
							value={this.state.uid}
						/>
						<label>Password</label>
						<input
							name="password"
							type="password"
							onChange={this.handlePasswordChange}
							value={this.state.password}
						/>
						<input type="submit" value="Login"/>
					</form>}
				</div>
			</div>
		);
	}
}

Login.propTypes = {
	is_authenticated: PropTypes.bool.isRequired,
	is_fetching:      PropTypes.bool.isRequired,
	requestLogin:     PropTypes.func.isRequired
}

export default Login;
