import React, { Component } from "react";
import PropTypes from "prop-types";
import { Route, Redirect } from "react-router-dom";
import { connect } from "react-redux";
import { checkTokenTimeout } from "actions/auth";

class DumbPrivateRoute extends Component {
	render() {
		let { is_authenticated, component: Component } = this.props;
		let alteredProps = this.props;
		delete alteredProps.component;
		return (
			<Route
				{...alteredProps}
				render={(props) => is_authenticated ? <Component {...props}/> : <Redirect to="/login"/>}
			/>
		);
	}
	render() {
		let { is_authenticated, component: Component } = this.props;
		checkTokenTimeout();
		if (!is_authenticated) {
			return (
				<Redirect
					to={{
						pathname: "/login",
						state: {
							from: this.props.location	
						}
					}}
				/>
			);
		}
		return (
			<Component {...this.props}/>
		);
	}
}

	/*class DumbPrivateRoute extends Component {
	render() {
		let newThis = this;
		let { is_authenticated, component: Component, path } = this.props;
		checkTokenTimeout();
		console.log("Checking timeout...");
		return (
			<Route
				path={path}
				render={(props) =>
					is_authenticated
					? <Component {...newThis.props} />
					: <Redirect to={{
						pathname: '/login',
						state: {
							from: newThis.props.location
						}
					}} />
				}
			/>
		);
	}
}*/

DumbPrivateRoute.propTypes = {
	checkTokenTimeout: PropTypes.func.isRequired
}

const mapStateToProps = state => ({
	is_authenticated: state.auth.is_authenticated
});

const mapDispatchToProps = dispatch => ({
	checkTokenTimeout: () => dispatch(checkTokenTimeout())
});

const PrivateRoute = connect(
	mapStateToProps,
	mapDispatchToProps
)(DumbPrivateRoute);

export default PrivateRoute;
