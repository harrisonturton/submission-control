import React, { Component } from "react";
import { connect } from "react-redux";
import PropTypes from "prop-types";
import { Redirect } from "react-router-dom";

const Home = ({ courses }) => {
	let course_id = courses[0].id;	
	return <Redirect to={`/course/${course_id}`}/>
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
