import { connect } from "react-redux";
import { makeLoginRequest } from "actions";
import Login from "containers/NewLogin/Login";

const mapStateToProps = state => ({
	is_authenticated: state.auth.is_authenticated,
	is_fetching:      state.auth.is_fetching,
});

const mapDispatchToProps = dispatch => ({
	requestLogin: (email, password) => dispatch(makeLoginRequest(email, password)),
});

const LoginContainer = connect(
	mapStateToProps,
	mapDispatchToProps,
)(Login);

export default LoginContainer;
