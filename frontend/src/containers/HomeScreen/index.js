import React, { Component } from "react";
import { connect } from "react-redux";
import { Link, Redirect } from "react-router-dom";
import { hasToken } from "auth";
import "./style.css";

import { 
	Header,
	AssessmentItemList,
	AssessmentFeedbackList
} from "components";

class Home extends Component {
	render() {
		let is_authenticated = hasToken();
		if (!is_authenticated) {
			return <Redirect to="/login"/>;	
		}
		return (
			<div className="home-wrapper">
				<Header
					currentCourse="Concurrent & Distributed Programming"
					courses={[
						"Computer Networks",
						"Structured Programming",
						"Programming as Problem Solving",
					]}
				/>
				<div className="assessment-wrapper">
					<AssessmentItemList
						title="Upcoming Assignments"
						subtitle="View or submit"
						items={[
							{ title: "Harmony & Money", course_code: "comp2300", due_date: new Date(), comments: "2 submissions with 8 warnings on the most recent build." },
							{ title: "Harmony & Money", course_code: "comp2300", due_date: new Date(), comments: "2 submissions with 8 warnings on the most recent build." },
						]}
					/>
					<AssessmentItemList
						title="Upcoming Labs"
						subtitle="View or submit"
						items={[
							{ title: "Harmony & Money", course_code: "comp2300", due_date: new Date(), comments: "2 submissions with 8 warnings on the most recent build." },
							{ title: "Harmony & Money", course_code: "comp2300", due_date: new Date(), comments: "2 submissions with 8 warnings on the most recent build." },
						]}
					/>
				</div>
				<div className="feedback-wrapper">
					<AssessmentFeedbackList
						title="Assessment Feedback"
						subtitle="View or comment"
						items={[
							{
								title: "Distributed Server",
								dute_date: new Date(),
								feedback: "Although I really like what you did with the solution, your code style could be a bit better! Remember to name your variables in a readable way, and comment the tricky bits that aren’t immediately apparent to the reader. Comment why, not what. We know what the code does — we can read it!"
							},
							{
								title: "Synchronized Data",
								due_date: new Date(),
								feedback: "Good stuff, the commenting is effective."
							},
							{
								title: "Harmony & Money",
								due_date: new Date(),
								feedback: "Looks good, but I would've liked to see the flight controls refactored into a separate package. Make sure to test thoroughly!"
							},
						]}
					/>
				</div>
			</div>
		);
	}
}

const HomeScreen = connect(
	null,
	null
)(Home);

export default HomeScreen;
