import React from "react";
import "./style.css";

const Feedback = ({ feedback }) => (
	<div className="feedback">
		<div className="divider"></div>
		<p className="description">{feedback}</p>
	</div>
);

const SubmissionItem = ({ title, hasFeedback, description, dueDate, feedback }) => (
	<div className="submission-item">
		<div className="submission-header">
			<div className="submission-title-wrapper">
				<h1>{title}</h1>
				<span className="comments">All tests passed. Congratulations!</span>
			</div>
			<div className="submission-date-wrapper">
				<span className="due-date">{formatSubmittedDate(dueDate)}</span>
				<span className="timestamp">{formatTimestamp(dueDate)}</span>
			</div>
		</div>
		{description === ""
				? <p className="absent">No description</p>
				: <p className="description">{description}</p>}
		<Feedback hasFeedback={hasFeedback} feedback={feedback}/>
	</div>
);

const formatSubmittedDate = due_date => {
	let difference_in_days = new Date() - due_date;
	return `Submitted ${difference_in_days} days ago`;
};

const formatTimestamp = date => date.toLocaleDateString("en-US", {
	weekday: "long",
	month: "long",
	day: "numeric"
});


export default SubmissionItem;
