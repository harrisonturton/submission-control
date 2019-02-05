import React, { Component } from "react";
import PropTypes from "prop-types";
import { connect } from "react-redux";
import { WithHeader } from "containers";

// Expects students to have the following form:
// { firstname, lastname, uid, tutorials }
const StudentTable = ({ students }) => (
	<table>
		<tr>
			<th>Name</th>
			<th>UID</th>
			<th>Tutorials</th>
		</tr>
		{students.map(student => (
			<tr>
				<td>{student.firstname + " " + student.lastname}</td>
				<td>{student.uid}</td>
				<td>{student.tutorials.split("").join(" ")}</td>
			</tr>
		))}
	</table>
);

const mapStateToProps = state => ({
	
});

const _ManageStudentsScreen = ({ match, students }) => {
	let { course_id } = match.params;
	return (
		<WithHeader className="column-parent" currentCourseID={course_id}>
			<div className="column-left">
				<h1>Manage Students</h1>
			</div>
			<div className="column-right">
				<StudentTable
					students={[
						{ firstname: "Harry", lastname: "Turton", uid: "u6386433", tutorials: "abc" },
						{ firstname: "John", lastname: "Smith", uid: "u7262488", tutorials: "bc" },
						{ firstname: "Avril", lastname: "Lavigne", uid: "u726534", tutorials: "ad" },
					]}
				/>
			</div>
		</WithHeader>
	);
};

_ManageStudentsScreen.propTypes = {
	students: PropTypes.array.isRequired
};

const ManageStudentsScreen = connect(
	mapStateToProps,
	null
)(_ManageStudentsScreen);

export default ManageStudentsScreen;
