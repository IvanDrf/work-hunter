# Pull Requests

# Naming
Pull request name must have following structure
```bash
[direction:type]: message
```

Example
```bash
[back:auth:feat]: add password hashing
```

directions:
- front (frontend)
- back (backend)
- ml (machine learning)
 
types:
- feat (feature)
- fix (fix bug for exmaple)
- tests (tests for app)
- ref (refactor)

# Merging
- Before opening pull request you must run tests on your local machine
- Before merging, you must wait for CI/CD 
- Wait for code review
