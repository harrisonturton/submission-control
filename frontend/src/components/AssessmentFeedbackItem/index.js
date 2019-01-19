import React from "react";
import "./style.css";

const AssessmentFeedbackItem = ({ title, due_date, feedback }) => (
	<div className="assessment-feedback-item">
		<div className="feedback-title-wrapper">
			<span className="feedback-title">{title}</span>
			<span className="feedback-due">{formatDueDate(due_date)}</span>
		</div>
		<p className="feedback">{feedback}</p>
	</div>
);

const formatDueDate = due_date => {
	let difference_in_days = new Date() - due_date;
	return `${difference_in_days} days ago`;
};

export default AssessmentFeedbackItem;
