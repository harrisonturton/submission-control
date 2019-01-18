import React, { Component } from "react";
import { connect } from "react-redux";

class Course extends Component {
	render() {
		let { course_id } = this.props.match.params;
		return <h1>{course_id}</h1>;
	}
}

const CourseScreen = connect(null, null)(Course);

export default CourseScreen;
