# Tutorial Template
Template repo for creating tutorials.

[Use this template](https://github.com/new?owner=byu-oit&template_name=tutorial-template&template_owner=byu-oit) to create new tutorials.

## Tutorial Instructions
> **Note**
> Keep track of any issues with this tutorial and contribute to fix them by following
> the [steps below](#4-pay-it-forward) [make sure this links to the right step number]

### 0. Prerequisites
* Complete the [GitHub Tutorial](https://github.com/byu-oit/tutorial-github)
* [add more prerequisites as seen fit]

### 1. Setup
When writing new tutorials, please [use this template](https://github.com/new?owner=byu-oit&template_name=tutorial-template&template_owner=byu-oit) to create your new repository. This will help our tutorials 
be more standardized. Here are a couple standards to maintain:
* Tutorial repos should be named like this: `tutorial-<tutorial_name>`
  * ie: `tutorial-nodejs`
* When you create a new tutorial repo, go to the `Settings` tab, under `Collaborators and teams`, and give admin 
access to the `SpecOps-Developer-FTE` team and write access to the `SpecOps-Developer`, `Full Stack`, and `Cloud` 
teams.
* The tutorial should have a title that looks like this: `<tutorial_name> Tutorial`.
  * ie: `NodeJs Tutorial`
* The tutorial should have a short description after the title.
  * ie: `Tutorial for JavaScript, Node, and npm`
* For ease of access, long text chunks in the tutorial should be wrapped according to the wrap guide that shows up in 
your text editor
* The tutorial should have a `Prerequisites` step (see [step 0](#0-prerequisites))
* If the tutorial requires the user to make changes in a different branch, use [step 2](#2-clone-repository) as an 
example. If you do this, you should also have a `Clean Up` step to keep the repo from getting crowded (see 
[step 3](#3-clean-up)).
* Make sure to add the tutorial to the [Full Stack Developer Handbook](https://github.com/byu-oit/fullstack-developer-handbook/)

### 2. Clone Repository
For this tutorial we will make our changes in a branch off of this GitHub repo:
1. Clone this repository down to your local machine using git clone 
2. Create a new branch called base/<net_id> (this branch will be your working branch for the duration of this tutorial).
You won't push this branch to GitHub, but you can still make commits.

### 3. Clean Up
1. Pass off your training efforts with your mentor. Ask them any questions you might have.
2. Please delete any cloud resources created in this tutorial.

### 4. Pay it Forward

Congratulations 🙌, you've completed the GO tutorial!

1. If you found any ways to improve this tutorial (fix typos, improve text to be clearer etc.):
    1. Make a feature branch (named `feat/<something>`)
    2. Make your suggested changes, commit and push branch
    3. Make a PR to the `main` branch and assign it to be reviewed by members of the Full Stack guild
2. If you found an issue but don't know how to fix it: [make sure the following links are correct]
    1. Check this GitHub repo's [issues](https://GitHub.com/byu-oit/tutorial-template/issues) if it's already there
    2. Make a [new issue](https://GitHub.com/byu-oit/tutorial-template/issues/new/choose) and explain the problem as
       best you can (include screenshots if necessary)
