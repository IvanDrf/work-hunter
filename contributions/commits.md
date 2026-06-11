# Commits
- ```main``` branch can only contain a stable app version
- commit name must have the following structure
    ``` bash
    git commit -m "[direction:type]: message"
    ```
    directions:
    - front (frontend)
    - back (backend)
    - ml (machine learning)
     
    types:
    - dev (development)
    - fix (fix bug for exmaple)
    - tests (tests for app)
    - ref (refactor)
    - feat (feature)
    - ci/cd (github workflows)

    examples:
    ```bash
    git commit -m "[back:dev]: add connection to PostgreSQL database"
    git commit -m "[back:feat]: add auth for users by jwt"
    git commit -m "[front:feat]: add user profile page"
    git commit -m "[back:tests]: add test for AuthService:RegisterUser"
    git commit -m "[front:fix]: fix logo error, now clicking on the logo redirects to the main page"
    ```

