import React from "react";
import { Route, Redirect } from "react-router-dom";

export default function PrivateRoute({ component: Component, auth, ...rest}) {
	return (
		<Route {...rest} render={props => (
			auth.isAuthenticated
			? <Component {...props}/>
			: <Redirect to={{
				pathname: "/login",
				state: { from: props.location }
			}}/>
		)} />
	);
}
