import React, { Component } from "react";
import { connect } from "react-redux";
import PropTypes from "prop-types";
import { Redirect } from "react-router-dom";
import { Loader } from "components";

const Home = ({ isAuthenticated, courses }) => {
	let has_loaded = courses.length > 0;
	let course = courses[0];
	if (!isAuthenticated) {
		return <Redirect to="/login"/>;
	}
	if (!has_loaded || course.id == undefined) {
		return (
			<div style={{display: "flex", alignItems: "center", justifyContent: "center", width: "100vw", height: "100vh"}}>
				<Loader/>	
			</div>
		);
	}
	return <Redirect to={`/course/${course.id}`}/>;
};

Home.propTypes = {
	courses:         PropTypes.array.isRequired,
	isAuthenticated: PropTypes.bool.isRequired
};

const mapStateToProps = state => ({
	isAuthenticated: state.auth.is_authenticated,
	courses: state.data.courses,
});

const HomeScreen = connect(
	mapStateToProps,
	null
)(Home);

export default HomeScreen;
