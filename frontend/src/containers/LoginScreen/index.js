import { connect } from "react-redux";
import { attemptSignIn } from "actions/thunks";
import Login from "containers/LoginScreen/Login";

const mapStateToProps = state => ({
	is_authenticated: state.auth.is_authenticated,
	is_fetching:      state.auth.is_fetching,
});

const mapDispatchToProps = dispatch => ({
	requestLogin: (email, password) => dispatch(attemptSignIn(email, password))
});

const LoginContainer = connect(
	mapStateToProps,
	mapDispatchToProps,
)(Login);

export default LoginContainer;
