import React, { Component } from "react";
import { BrowserRouter as Router, Switch, Route } from "react-router-dom";
import { PrivateRoute, HomeScreen, LoginContainer, CourseScreen } from "containers";

export default class App extends Component {
	render() {
		return (
			<Router>
				<Switch>
					<Route
						path="/login"
						component={LoginContainer}
					/>
					<Route
						path="/course/:course_no"
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
