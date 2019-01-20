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
	}
	onClick = () => this.setState(prev => ({
		isExpanded: true
	}));
	onMouseLeave = () => this.setState(prev => ({
		isExpanded: false
	}));
	render() {
		let { isExpanded } = this.state;
		let { courses, currentCourse, logout } = this.props;
		console.log("Header courses: ", JSON.stringify(courses));
		return (
			<header
				className={isExpanded ? "expanded" : ""}
				onMouseLeave={this.onMouseLeave}
				style={{
					height: isExpanded ? courses.length * 40 + 10 : 15
				}}
			>
				<div className="header-bar">
					<span className="current" onClick={this.onClick}>
						{currentCourse}
						<img className="chevron" alt={""} src={chevronDown}/>
					</span>
					<button className="logout-button" onClick={() => logout()}>
						LOGOUT
					</button>
				</div>
				<ul>
					{courses.filter(course => course.name !== currentCourse).map((course, i) => (
						<li key={i} className="course-name">
							<Link to={`/course/${course.id}`}>{course.name}</Link>
						</li>
					))}
				</ul>
			</header>
		);
	}
}

_Header.propTypes = {
	currentCourse: PropTypes.string.isRequired,
	courses:       PropTypes.array.isRequired
}

const mapDispatchToProps = dispatch => ({
	logout: () => dispatch(logout())
});

const Header = connect(
	null,
	mapDispatchToProps
)(_Header);

export default Header;
