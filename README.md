# Hydra Daemon

This service powers the automated building & testing of student projects. Though it is intended to be run non-interactively, it publishes an RPC server which can be communicated with via `hydra-cli`. This allows us to manually start/remove/scale individual swarm services.

### Background

This service has some conflicting goals:

**Speed / Responsiveness**
Whether or not their project builds perfectly, generates warnings, or fails with errors, a user will want quick feedback. Ideally the server will stay responsive, even close to the deadline (and under heavy load).

**Ease-of-use**
We need to build & test projects in a variety of languages and environments. It should be really easy to introduce new languages, define new environments/dependencies, and new build instructions. Conflicting language verisons (i.e. Python 2.7 and Python 3) should not be an issue.

**Security**
We are running unverified user code! This is very dangerous, whether or not malicious code is uploaded. Suppose a student uploads a kernel exploit – how do we protect against this? What if they inadvertently change the testing environment? We need to be secure, but also log everything so that we can find the source if anything does happen.

**Reliability**
If a testing environment goes down (i.e. through malicious code or otherwise), we cannot let this influene other environments. Ideally, the failure of one node will not influence the failure of the platform as a whole. 

For example, if COMP2310 experiences an outage, this shouldn't affect COMP1130. At a finer level of granularity, `assignment-3-comp2310` outages should not effect `assignment-1-comp2310`.



### Approach

We've attempted to solve this issues through Docker containers. It gives a lot of benefits:

- Easy, contained testing environments
- Easy to restart
- A basic level of security
- Swarms can make it super performant and reliable (automatically keeps up `n` replicas)