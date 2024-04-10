# Holos Front End

The holos front end is a React + TypeScript single page app with strict type
checking enabled via eslint and prettier.

Frameworks were evaluated, but all were eschewed because they have too many
unnecessary features and make too many assumptions which are distracting for
velocity.  Instead, the app is built as a combination of:

 1. PatternFly for a consistent design and component library.
 2. React Router for client side url routing.
 3. SupaBase for authentication (and only auth!)
 4. Buf to generate TypeScript types from proto bufs.
 5. ConnectRPC + TanStack Query for the client side query rpc.

This stack provides a well integrated, strongly typed, front and backend
service that runs well on top of the Holos Platform.
