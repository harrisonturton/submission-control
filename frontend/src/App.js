import React, { Component } from "react";
import { BrowserRouter as Router, Route, Switch } from "react-router-dom";
import { PrivateRoute, LoginScreen } from "containers";
import { Header } from "components";
import { refresh } from "api";

class App extends Component {
	constructor(props) {
		super(props);
		this.state = {
			isAuthenticated: false,
			token: null
		}
		this.refreshToken = this.refreshToken.bind(this);
		this.handleLogin = this.handleLogin.bind(this);
	}
	async handleLogin(ok, token) {
		console.log(ok, token)
		this.setState({
			isAuthenticated: ok,
			token: token
		});
	}
	async refreshToken(token) {
		const { ok, new_token }	= await refresh(token);
		this.setState({
			isAuthenticated: ok,
			token: new_token
		});
	}
	render() {
		return (
			<Router>
				<Switch>
					<Route
						path="/login"
						render={() => <LoginScreen onAuth={this.handleLogin}/>}
					/>
					<PrivateRoute
						exact path="/"
						component={Home}
						auth={this.state}
					/>
					<PrivateRoute
						path="/course/:course_code"
						component={Course}
						auth={this.state}
					/>
				</Switch>
			</Router>
		);
	}
}

const Course = ({ match }) => {
	const course_code = match.params.course_code;
	if (course_code === undefined || course_code === null) {
		console.log("error, course_code is undefined!");
		return <h1>Course code is not found</h1>;
	}
	return <h1>Course: {course_code}</h1>;
};

const Home = () => (
	<Header
		currentCourse={"Concurrent & Distributed Programming"}
		courses={[
			"Computer Networks",
			"Structured Programming",
			"Programming as Problem Solving",
		]}
	/>
);
export default App;
