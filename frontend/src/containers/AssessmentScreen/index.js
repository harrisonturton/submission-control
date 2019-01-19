import React, { Component } from "react";
import PropTypes from "prop-types";
import { connect } from "react-redux";
import { Header } from "containers";
import { AssessmentFeedbackList } from "components";
import "./style.css";

const _AssessmentScreen = ({ match, courses, assessment, submissions }) => {
	let { course_id, assessment_id } = match.params;
	let filtered_submissions = submissions.filter(sub => sub.assessment_id == assessment_id);
	let assessment_name = assessment.find(ass => ass.id == assessment_id).name;
	return (
		<div className="assessment-screen">
			<Header
				currentCourse={courses.find(course => course.id == course_id).name}	
				courses={courses.filter(course => course.id != course_id)}	
			/>
			<div className="left-wrapper">
				<h1>{assessment_name}</h1>
			</div>
			<div className="right-wrapper">
				{filtered_submissions.map(sub => (
					<div>
						<h1>{sub.name}</h1>	
						<h2>{sub.description}</h2>	
						<h2>{sub.feedback}</h2>	
					</div>
				))}
			</div>
		</div>
	);
};

_AssessmentScreen.propTypes = {
	courses:     PropTypes.array.isRequired,
	assessment:  PropTypes.array.isRequired,
	submissions: PropTypes.array.isRequired,
};

const mapStateToProps = state => ({
	courses: state.data.courses,
	assessment: [ ...state.data.assessment.assignments, ...state.data.assessment.labs ],
	submissions: state.data.submissions,
});

const AssessmentScreen = connect(
	mapStateToProps,
	null
)(_AssessmentScreen);

export default AssessmentScreen;
