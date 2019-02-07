import React from "react";
import "./style.css";

const joinNames = students =>
	students
		.map(student => student.firstname + " " + student.lastname + ",")
		.join(" ")
		.slice(0, -1);

const TutorialItem = ({ name, tutors, students }) => (
	<div className="tutorial-item">
		<h2 className="tutorial-name">{name}</h2>
		<div className="tutors">
			<span className="tutor">{joinNames(tutors)}</span>
		</div>
		<div className="students">
			<span className="student">{joinNames(students)}</span>
		</div>
	</div>
);

export default TutorialItem;
