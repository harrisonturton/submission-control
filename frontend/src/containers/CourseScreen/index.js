import React, { Component } from "react";
import { connect } from "react-redux";
import PropTypes from "prop-types";
import { Redirect } from "react-router-dom";
import { StudentCourseScreen, ConvenorCourseScreen } from "containers";
import { Loader } from "components";

const _CourseScreen = ({ role, match }) => {
	let { course_id } = match.params;
	console.log(`COURSE SCREEN HAS ROLE ${role}`);
	switch (role) {
		case "student":
			return <StudentCourseScreen courseID={course_id}/>;
		case "convenor":
		case "admin":
		case "tutor":
			return <ConvenorCourseScreen courseID={course_id}/>;
		default:
			return <h1>Unknown role {role}. Please email harrison.turton@anu.edu.au</h1>;
	}
};

const getRole = (courses, course_id) => {
	let course = courses.find(course => course.id == course_id)
	if (!course) {
		return "student";
	}
	return course.role;
}

const mapStateToProps = (state, { match }) => ({
	role: getRole(state.data.courses, match.params.course_id)
});

const CourseScreen = connect(
	mapStateToProps,
	null
)(_CourseScreen);

export default CourseScreen;
