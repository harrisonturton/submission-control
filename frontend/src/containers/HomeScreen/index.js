import React, { Component } from "react";
import { connect } from "react-redux";
import PropTypes from "prop-types";
import { Redirect } from "react-router-dom";
import { Loader } from "components";

const Home = ({ courses }) => {
	let has_loaded = courses.length > 0;
	let course = courses[0];
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
	courses: PropTypes.array.isRequired
};

const mapStateToProps = state => ({
	courses: state.data.courses
});

const HomeScreen = connect(
	mapStateToProps,
	null
)(Home);

export default HomeScreen;
