import React from "react";
import { connect } from "react-redux";
import { Redirect } from "react-router-dom";
import { hasToken } from "auth";
import "./style.css";
import {
	Header,
	AssessmentItemList,
	AssessmentFeedbackList,
} from "components";

const RedirectToLogin = () => <Redirect to="/login"/>;

const Home = ({
	header,
	assignments,
	labs,
	feedback
}) => hasToken() ? RedirectToLogin() : (
	<div className="home-wrapper">
		<Header
			currentCourse={header.current_course}
			courses={header.courses}
		/>
		<div className="assignment-wrapper">
			<AssessmentItemList
				title="Upcoming Assignments"
				subtitle="View or submit"
				items={assignments}
			/>
			<AssessmentItemList
				title="Upcoming Labs"
				subtitle="View or submit"
				items={labs}
			/>
		</div>
		<div className="feedback-wrapper">
			<AssessmentItemList
				title="Assessment Feedback"
				subtitle="View or comment"
				items={feedback}
			/>
		</div>
	</div>
);

const mapStateToProps = state => ({
	header: {
		current_course: state.user.current_course,	
		courses: state.user.courses,
	},
	assignments: state.user.assignments,
	labs: state.user.labs,
	feedback: [...state.user.assignments, ...state.user.labs].map(item => ({
		title: item.title,
		due_date: item.due_date,
		feedback: item.feedback
	}))
})

const HomeScreen = connect(
	mapStateToProps,
	null
)(Home);

export default HomeScreen;
