import React, { Component } from "react";
import PropTypes from "prop-types";
import { connect } from "react-redux";
import { Redirect } from "react-router-dom";
import { Header } from "containers";
import "./style.css";

const RedirectToLogin = () => <Redirect to="/"/>;

const _WithHeader = ({ children, className, isAuthenticated, currentCourseID }) => !isAuthenticated ? <RedirectToLogin/> : (
	<div className={"body-wrapper" + (className ? ` ${className}` : "")}>
		<Header currentCourseID={currentCourseID}/>
		{children}
	</div>
);

_WithHeader.propTypes = {
	isAuthenticated: PropTypes.bool.isRequired,
	currentCourseID: PropTypes.string.isRequired,
};

const mapStateToProps = state => ({
	isAuthenticated: state.auth.is_authenticated,
});

const WithHeader = connect(
	mapStateToProps,
	null
)(_WithHeader);

export default WithHeader;
