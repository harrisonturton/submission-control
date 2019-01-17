import React, { Component } from "react";
import { connect } from "react-redux";
import { Redirect } from "react-router-dom";
import { hasToken } from "api/auth";
import { fetchAssessmentsForCourse } from "api/api";
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
		let { token } = this.props;
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
					<button onClick={() => {
						console.log(token);
						fetchAssessmentsForCourse(1, token)
							.then(({ assignments, labs }) => {
								console.log(JSON.stringify(assignments));	
								console.log(JSON.stringify(labs));	
							});
					}}>Click me!</button>
					<AssessmentItemList
						title="Upcoming Assignments"
						subtitle=""
						items={[
							{ title: "Harmony & Money", course_code: "comp2300", due_date: new Date(), comments: "2 submissions with 8 warnings on the most recent build." },
							{ title: "Harmony & Money", course_code: "comp2300", due_date: new Date(), comments: "2 submissions with 8 warnings on the most recent build." },
						]}
					/>
					<AssessmentItemList
						title="Upcoming Labs"
						subtitle=""
						items={[
							{ title: "Structured Programming", course_code: "comp2300", due_date: new Date(), comments: "2 submissions with 8 warnings on the most recent build." },
							{ title: "Tasks", course_code: "comp2300", due_date: new Date(), comments: "2 submissions with 8 warnings on the most recent build." },
							{ title: "Protection", course_code: "comp2300", due_date: new Date(), comments: "2 submissions with 8 warnings on the most recent build." },
							{ title: "Task Lifetimes", course_code: "comp2300", due_date: new Date(), comments: "2 submissions with 8 warnings on the most recent build." },
							{ title: "Communicating Tasks", course_code: "comp2300", due_date: new Date(), comments: "2 submissions with 8 warnings on the most recent build." },
							{ title: "Distributed Server", course_code: "comp2300", due_date: new Date(), comments: "2 submissions with 8 warnings on the most recent build." },
							{ title: "Implicit Concurrency", course_code: "comp2300", due_date: new Date(), comments: "2 submissions with 8 warnings on the most recent build." },
							{ title: "Synchronised Data", course_code: "comp2300", due_date: new Date(), comments: "2 submissions with 8 warnings on the most recent build." },
						]}
					/>
				</div>
				<div className="feedback-wrapper">
					<AssessmentFeedbackList
						title="Assessment Feedback"
						subtitle=""
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

const mapStateToProps = state => ({
	token: state.auth.token
});

const mapDispatchToProps = dispatch => ({ });

const HomeScreen = connect(
	mapStateToProps,
	null
)(Home);

export default HomeScreen;
