import React, { Component } from "react";
import PropTypes from "prop-types";
import { connect } from "react-redux";
import { WithHeader } from "containers";
import { TutorialItem } from "components";

const dummyTutors = [
	{ firstname: "Harry", lastname: "Turton", uid: "u6386433", tutorials: "abc" },
	{ firstname: "John", lastname: "Smith", uid: "u7262488", tutorials: "bc" },
	{ firstname: "Avril", lastname: "Lavigne", uid: "u726534", tutorials: "ad" },
];

const dummyStudents = [
	{ firstname: "Harry", lastname: "Turton", uid: "u6386433", tutorials: "abc" },
	{ firstname: "John", lastname: "Smith", uid: "u7262488", tutorials: "bc" },
	{ firstname: "John", lastname: "Smith", uid: "u7262488", tutorials: "bc" },
	{ firstname: "Avril", lastname: "Lavigne", uid: "u726534", tutorials: "ad" },
	{ firstname: "Harry", lastname: "Turton", uid: "u6386433", tutorials: "abc" },
	{ firstname: "John", lastname: "Smith", uid: "u7262488", tutorials: "bc" },
	{ firstname: "John", lastname: "Smith", uid: "u7262488", tutorials: "bc" },
	{ firstname: "Avril", lastname: "Lavigne", uid: "u726534", tutorials: "ad" },
	{ firstname: "Avril", lastname: "Lavigne", uid: "u726534", tutorials: "ad" },
];

const ManageTutorialsScreen = ({ match, students }) => {
	let { course_id } = match.params;
	return (
		<WithHeader className="column-parent" currentCourseID={course_id}>
			<div className="column-left">
				<h1>Manage Tutorial Groups</h1>
			</div>
			<div className="column-right">
				<TutorialItem
					name="Tutorial A"
					tutors={dummyTutors}
					students={dummyStudents}
				/>
				<TutorialItem
					name="Tutorial B"
					tutors={dummyTutors}
					students={dummyStudents}
				/>
				<TutorialItem
					name="Tutorial C"
					tutors={dummyTutors}
					students={dummyStudents}
				/>
			</div>
		</WithHeader>
	);
};

export default ManageTutorialsScreen;
