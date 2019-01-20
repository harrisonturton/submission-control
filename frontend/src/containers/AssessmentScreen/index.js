import React, { Component } from "react";
import PropTypes from "prop-types";
import { Redirect } from "react-router-dom";
import { Header } from "containers";
import { connect } from "react-redux";
import { AssessmentFeedbackList, SubmissionList } from "components";
import "./style.css";

const _AssessmentScreen = ({ match, is_authenticated, courses, assessment, submissions }) => {
	let { course_id, assessment_id } = match.params;
	let filtered_submissions = submissions.filter(sub => sub.assessment_id == assessment_id);
	let assessment_name = assessment.find(ass => ass.id == assessment_id).name;
	if (!is_authenticated) {
		return <Redirect to="/login"/>;
	}
	return (
		<div className="assessment-screen">
			<Header
				currentCourse={courses.find(course => course.id == course_id).name}	
				courses={courses.filter(course => course.id != course_id)}	
			/>
			<div className="left-wrapper">
				<h1>{assessment_name}</h1>
				<p>You can make multiple submissions to a single assignment or lab. Each submission will be tested against the test suite, and you’ll see the result. Your tutors and convenors can see & give feedback on any submission you make.</p>
				<form>
					<label>Title</label>
					<input
						name="title"
						type="text"
						placeholder="Title of your submission"
					/>
					<label>Comments</label>
					<textarea
						name="comments"
						type="text"
						placeholder="Some more info about your submission"
					/>
					<label>Zipped submission folder</label>
					<input type="file"/>
					<input type="submit" value="Submit"/>
				</form>
			</div>
			<div className="right-wrapper">
				<SubmissionList
					title="Your Submissions"
					subtitle=""
					courseID={course_id}
					assessmentID={assessment_id}
					submissions={filtered_submissions}
				/>
			</div>
		</div>
	);
};

_AssessmentScreen.propTypes = {
	is_authenticated: PropTypes.bool.isRequired,
	courses:          PropTypes.array.isRequired,
	assessment:       PropTypes.array.isRequired,
	submissions:      PropTypes.array.isRequired,
};

const mapStateToProps = state => ({
	is_authenticated: state.auth.is_authenticated,
	courses: state.data.courses,
	assessment: [ ...state.data.assessment.assignments, ...state.data.assessment.labs ],
	submissions: state.data.submissions,
});

const AssessmentScreen = connect(
	mapStateToProps,
	null
)(_AssessmentScreen);

export default AssessmentScreen;
