import React, { Component } from "react";
import PropTypes from "prop-types";
import { Redirect } from "react-router-dom";
import { Header } from "containers";
import { connect } from "react-redux";
import "./style.css";

const _SubmissionScreen = ({ match, is_authenticated, submissions }) => {
	let { course_id, assessment_id, submission_id } = match.params;
	let submission = submissions.find(sub => sub.id == submission_id);
	if (!is_authenticated) {
		return <Redirect to="/login"/>;
	}
	return (
		<div className="submission-screen">
			<h1>{submission.title}</h1>
			<p>{submission.description}</p>
			<p>{submission.feedback}</p>
		</div>
	);
};

_SubmissionScreen.propTypes = {
	is_authenticated: PropTypes.bool.isRequired,
	submissions:      PropTypes.array.isRequired,
};

const mapStateToProps = state => ({
	is_authenticated: state.auth.is_authenticated,
	submissions: state.data.submissions,
});

const SubmissionScreen = connect(
	mapStateToProps,
	null
)(_SubmissionScreen);

export default SubmissionScreen;
