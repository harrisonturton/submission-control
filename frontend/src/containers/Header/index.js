import React, { Component } from "react";
import { Link } from "react-router-dom";
import PropTypes from "prop-types";
import { connect } from "react-redux";
import { logout } from "actions/thunks";
import chevronDown from "assets/chevron-down.png";
import "./style.css";

class _Header extends Component {
	constructor(props) {
		super(props);
		this.state = { isExpanded: false };
		this.onClick = this.onClick.bind(this);
		this.onMouseLeave = this.onMouseLeave.bind(this);
		this.renderHeaderTop = this.renderHeaderTop.bind(this);
	}
	onClick() {
		this.setState({ isExpanded: true });
	}
	onMouseLeave() {
		this.setState({ isExpanded: false });
	}
	renderHeaderTop(currentCourse) {
		let { logout } = this.props;
		return (
			<div className="header-top">
				<span className="current" onClick={this.onClick}>
					{currentCourse.name}
					<img className="chevron" alt={""} src={chevronDown}/>
				</span>
				<button className="logout-button" onClick={() => logout()}>
					LOGOUT
				</button>
			</div>
		);
	}
	renderHeaderDropdown(courses, currentCourse) {
		return (
			<ul>
				{courses.filter(course => course.name !== currentCourse).map((course, i) => (
					<li key={i} className="course-name">
						<Link to={`/course/${course.id}`}>{course.name}</Link>
					</li>
				))}
			</ul>
		);
	}
	render() {
		let { isExpanded } = this.state;
		let { courses, currentCourseID, currentCourse, logout } = this.props;
		let expandedHeight = courses.length * 40 + 10;
		let minimisedHeight = 15;
		console.log(`Rendering header for id ${currentCourseID} for ${currentCourse}`);
		return (
			<header
				className={isExpanded ? "expanded" : ""}
				onMouseLeave={this.onMouseLeave}
				style={{
					height: isExpanded ? expandedHeight : minimisedHeight
				}}
			>
				{this.renderHeaderTop(currentCourse)}
				{this.renderHeaderDropdown(courses, currentCourse)}
			</header>
		);
	}
}

_Header.propTypes = {
	currentCourseID: PropTypes.string.isRequired,
	// Calculated from the state according to currentCourseID
	currentCourse:   PropTypes.string.isRequired,
	courses:         PropTypes.array.isRequired
}

const mapStateToProps = (state, { currentCourseID }) => ({
	currentCourse: state.data.courses.find(course => course.id == currentCourseID),
	courses:       state.data.courses.filter(course => course.id != currentCourseID)
});

const mapDispatchToProps = dispatch => ({
	logout: () => dispatch(logout())
});

const Header = connect(
	mapStateToProps,
	mapDispatchToProps
)(_Header);

export default Header;
