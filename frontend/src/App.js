import React, { Component } from "react";
import Auth from "auth";
import { BrowserRouter as Router, Switch, Route } from "react-router-dom";
import { PrivateRoute, HomeScreen, LoginScreen } from "containers";

const Course = ({ match }) => (
	<div>
		<p>Course {match.params.course_no}</p>
	</div>
);

export default class App extends Component {
	constructor(props) {
		super(props);
		this.state = { auth: new Auth() };
	}
	async componentDidMount() {
		let { auth } = this.state;
		let token = await auth.refresh(auth.token);
		if (token === null) {
			console.log("Could not verify...")	
		}
		this.forceUpdate();
	}
	render() {
		let { auth } = this.state;
		return (
			<Router>
				<Switch>
					<Route
						path="/login"
						render={() => <LoginScreen auth={auth}/>}
					/>
					<PrivateRoute
						exact path="/"
						component={HomeScreen}
						auth={auth}
					/>
					<PrivateRoute
						path="/course/:course_no"
						component={Course}
						auth={auth}
					/>
				</Switch>
			</Router>
		);
	}
}
