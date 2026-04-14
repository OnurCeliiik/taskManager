A simple API to manage tasks.
Users can create, update, and organize tasks securely with user login.
For learning how to handle databases and secure user access.

API Responsibilities:
- User Management
	- Signup and log in to get a JWT token.
	- Secure API with JWT
- Task Management
	- Create tasks (title, description, category)
	- View all tasks or filter by status/category
	- Update or delete tasks by ID
- Task Operations (CRUD)
	- Create a new task: Allow users to add a task with a title, description, and optional category(e.g, work, personal)
	- Read tasks: 
		- Returns all tasks belonging to the authenticated user
		- Filter by status (completed or pending) or category.
		- Retrieve a single task by its ID
	- Update a task
	- Delete a task

## Tech Stack
#### Language: Golang
#### Framework: Gin
#### Database: PostgreSQL
#### Authentication: JWT
#### Deployment: Docker