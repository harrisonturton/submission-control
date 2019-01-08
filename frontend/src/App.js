import React, { Component } from "react";
import Auth from "auth";
import { BrowserRouter as Router, Switch, Route } from "react-router-dom";
import { PrivateRoute, LoginScreen } from "containers";

const Home = () => (
	<div>
		<p>Home!</p>
	</div>
);

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
						component={Home}
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
