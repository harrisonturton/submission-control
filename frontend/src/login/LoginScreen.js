import React, { Component } from "react";
import "./login-screen.css";

export default class LoginScreen extends Component {
	render() {
		return (
			<div className="login-wrapper">
				<div className="login-modal">
					<form>
						<label>
							Email:
							<input type="email" name="email" placeholder="uXXXXXXX@anu.edu.au"/>
						</label>
						<label>
							Password:
							<input type="password" name="password"/>
						</label>
						<input type="submit" value="Submit"/>
					</form>
				</div>
			</div>
		);
	}
}
