import React from "react";
import { Link } from "react-router-dom";
import "./style.css";

const AssessmentItem = ({ title, due_date, comments }) => (
	<div className="assessment-item">
		<div className="title-wrapper">
			<span className="assessment-title">{title}</span>
			<span className="comments">{comments}</span>
		</div>
		<div className="date-wrapper">
			<span className="due-date">{formatDueDate(due_date)}</span>
			<span className="timestamp">{formatTimestamp(due_date)}</span>
		</div>
	</div>
);

const formatDueDate = due_date => {
	let difference_in_days = new Date() - due_date;
	return `Due in ${difference_in_days} days`;
};

const formatTimestamp = date => date.toLocaleDateString("en-US", {
	weekday: "long",
	month: "long",
	day: "numeric"
});

export default AssessmentItem;
