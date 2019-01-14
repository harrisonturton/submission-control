import React, { Component } from "react";
import { BrowserRouter as Router, Switch, Route } from "react-router-dom";
import { PrivateRoute, HomeScreen, LoginScreen } from "containers";

const CourseScreen = ({ match }) => (
	<p>{match.params.course_code}</p>
);

export default class App extends Component {
	render() {
		return (
			<Router>
				<Switch>
					<Route
						path="/login"
						component={LoginScreen}
					/>
					<Route
						path="/course/:course_code"
						component={CourseScreen}
					/>
					<PrivateRoute
						exact path="/"
						component={HomeScreen}
					/>
				</Switch>
			</Router>
		);
	}
}
