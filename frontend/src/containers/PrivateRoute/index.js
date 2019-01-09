import React from "react";
import { Route, Redirect } from "react-router-dom";

export default function PrivateRoute({ component: Component, auth, ...rest}) {
	console.log(`Is authenticated for private route: ${auth.is_authenticated}`);
	return (
		<Route {...rest} render={props => (
			auth.is_authenticated
			? <Component {...props}/>
			: <Redirect to={{
				pathname: "/login",
				state: { from: props.location }
			}}/>
		)} />
	);
}
