_2025-01-12_

> I'm considering giving up blogging because it consumes too much time that I could otherwise spend writing more code.

_2025-01-09_

> Today, I found myself spending more time than expected wrestling with build tools. Initially, I had no intention of replacing Make, but its syntax and lack of compatibility with Windows eventually got the better of me. I began exploring alternatives like Task and Just, which seemed promising, but I realized that each of these would introduce more to manage. Perhaps the simplest solution is just to stick with a Makefile, but here I am.

> After experimenting with several options, I decided to create a Go script to handle the build process for the entire project. Since the application is already written in Go, all that's needed to compile it is Go itself—no extra build tools to install. This approach also allows me to leverage the full power of a programming language, which opens up a lot of possibilities for customization and research.

> While I'm happy with this decision, there's a nagging feeling that I'm doing something unconventional. Many projects use a variety of build tools, but very few seem to rely on the project’s own programming language for the build process—though it’s somewhat common with Python projects. Regardless, I’ll stick with this approach until I encounter a problem. There are a few challenges ahead, like managing different GOOS and GOARCH environment variables, and the JS stuff, but I’m optimistic it will work out. Or at least, I hope so!

_2025-01-06_

> I decided to replace Fiber with the Go standard library's net/http package. My goal is to become more familiar with the standard library first. Similarly, I’ve chosen to forgo Chakra and focus on building my own components. I think it’s important to understand the foundational concepts that led to the creation of tools like Gin, Fiber or Chakra before diving into them, so I’m keeping things minimal until I gain more familiarity.

> I’ve also started to grasp how to serve static directories, such as Vite distributions, from Go binaries. I like how Go handles this aspect, and I start feeling confident about my tech stack choice. Go will serve as the backend, handling communication with all devices while also serving the frontend, which is essentially a React app. This setup allows users to access Neba from any device with a modern web browser. Additionally, since the backend operates like a daemon on the server, I can work on features like recursive tasks or rules/events in the future. Having a daemon will be useful.

---

_2025-01-05_

> My primary goal with Neba is to gain a deeper understanding of the software production workflow. In the past, I typically created projects based on requests or personal needs, but I’ve never released a fully functional software product. With Neba, I aim to develop something valuable and continuously improve it over time. TBH, I’m very excited about the journey ahead and hope to stay committed to this project.
