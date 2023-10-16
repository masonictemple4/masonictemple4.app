# Masonictemple4.app:
A go cli and backend for the masonictemple4.com blog. This includes tools that parse and upload 
blogs from markdown files.


Files are stored in GCP storage bucket.

DB is a hosted postgres instance.

### How it works?
Create a blog using markdown. At the top of the markdown file include your frontmatter with the `---` indicator. This allows you to specify the author, title, description, tags, etc.. and is used when parsing the document to fill out the specifics for the post.

**Example:** ~/path-to-blogs/example.md

```markdown
---
title: "Hello, world :)"
subtitle: "my first blog post"
description: "Welcome to my first blog!"
authors:
    - username: "masonictemple4"
      profilepicture: "url to image, you can leave this out if the user already has one."
tags:
    - "linux"
    - "wifi"
    - "os"
    - "bcm43602"
    - "broadcom"
    - "firmware"
---

# Hello, world!
A really interesting introduction to my blog

....
```

Upload the file:

`$ masonictemple4.app blog ~/path-to-blogs/example.md`

Update the blog:  
`$ masonictemple4.app blog update -i 1 ~/path-to-blogs/example.md`

Where the `-i` flag represents the blog id.


### TODOS:
- [X] Fix connection to the Cloud SQL instance
- [ ] CI/CD
- [ ] Refactor 
    - [ ] Make sure to use the utils where necessary instead of full functionality.
    - [ ] Redo the repository. 
    - [ ] Cleanup tests. 
- [ ] Implement real authentication.
- [ ] Setup CLI to manage remote data.


